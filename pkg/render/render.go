package render

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/JamesClonk/plato/pkg/config"
	"github.com/JamesClonk/plato/pkg/util/color"
	"github.com/JamesClonk/plato/pkg/util/command"
	"github.com/JamesClonk/plato/pkg/util/dir"
	"github.com/JamesClonk/plato/pkg/util/file"
	"github.com/JamesClonk/plato/pkg/util/log"
	"github.com/Masterminds/semver/v3"
	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/viper"
	"github.com/tredoe/osutil/user/crypt/sha512_crypt"
	"gopkg.in/yaml.v3"
)

func RenderFile(inputFile, outputFile string) {
	log.Infof("rendering template [%s] to [%s]...", color.Magenta(inputFile), color.Cyan(outputFile))

	// parse and render template file
	baseFilename := filepath.Base(inputFile)
	if err := writeFile(baseFilename, filepath.Dir(inputFile), outputFile, viper.AllSettings()); err != nil {
		log.Fatalf("could not render template file [%s]: %v", color.Magenta(baseFilename), err)
	}
}

func RenderTemplates(removeTerraformFiles, removeAllDirectories bool) {
	log.Infof("preparing to render templates ...")

	// fail if temporary .secrets-updated marker file / gitrepo taint exists
	if file.Exists(".secrets-updated") {
		log.Fatalf("[%s] marker file exists, git repository is tainted, abort!", color.Magenta(".secrets-updated"))
	}

	// cleanup directories
	if removeAllDirectories {
		dir.Remove(config.DirTarget())
	} else if dir.Exists(config.DirTarget()) {
		// delete only files
		err := filepath.Walk(config.DirTarget(), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// skip directories
			if info.IsDir() {
				return nil
			}
			// check if we want to also delete terraform init files
			if !removeTerraformFiles && strings.Contains(path, ".terraform"+string(os.PathSeparator)) {
				return nil
			}
			// delete file
			return file.Delete(path)
		})
		if err != nil {
			log.Fatalf("could not cleanup existing rendered files: %v", err)
		}
	}

	// ensure directory exists
	dir.Create(config.DirTarget())

	// go through all files
	err := filepath.Walk(config.DirSource(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip directories
		if info.IsDir() {
			return nil
		}
		return processFile(path, info)
	})
	if err != nil {
		log.Fatalf("could not render template files: %v", err)
	}
}

var funcMap = template.FuncMap{
	"PLATO":       platoHeader,
	"IPofCIDR":    ipOfCIDR,
	"MKPasswd":    mkpasswd,
	"ToYAML":      toYaml,
	"SemverCheck": semverCheck,
}

func platoHeader() string {
	// empty function for PLATO headers in template files
	return ""
}

func semverCheck(version string, constraint string, message ...string) bool {
	v, err := semver.NewVersion(strings.TrimPrefix(version, "v"))
	if err != nil {
		log.Fatalf("invalid semver version [%s]: %v", version, err)
	}
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		log.Fatalf("invalid semver constraint [%s]: %v", constraint, err)
	}

	check := c.Check(v)
	if check && len(message) > 0 {
		for _, m := range message {
			log.Errorf("semver check passed: %s", color.Red(m))
		}
	}
	return check
}

func ipOfCIDR(cidr string, pos int) string {
	_, c, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatalf("could not parse CIDR [%s]: %v", cidr, err)
	}

	for i := 0; i < pos; i++ {
		c.IP[3]++
	}
	return c.IP.String()
}

func mkpasswd(password string) string {
	c := sha512_crypt.New()
	s := sha512_crypt.GetSalt()

	salt := s.GenerateWRounds(s.SaltLenMax, 8192)
	hash, err := c.Generate([]byte(password), salt)
	if err != nil {
		log.Fatalf("could not generate a hashed password with salt [%s]: %v", salt, err)
	}

	return hash
}

func toYaml(object any, indent int) string {
	// out, err := yaml.Marshal(object)
	buf := bytes.Buffer{}
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(indent)
	err := enc.Encode(object)
	if err != nil {
		log.Fatalf("could not encode yaml: %v", err)
	}
	return strings.TrimSpace(buf.String())
}

func processFile(path string, info os.FileInfo) error {
	baseFilename := strings.TrimPrefix(path, config.DirSource()+string(os.PathSeparator))
	renderedFilename := filepath.Join(config.DirTarget(), baseFilename)

	// ensure path exists to write file to
	if err := os.MkdirAll(filepath.Dir(renderedFilename), 0700); err != nil { // use mode 0700, since we are likely rendering sensitive data
		return err
	}

	// begin .symlink marker handling
	// if its a .symlink marker file, then instead of templating/copying over the file its meant for,
	// we will instead create a symlink to the original file.
	// this allows us to deal with dynamically generated and/or changing state files that should to be checked back into git,
	// like *.tfstate, etc..
	if filepath.Ext(path) == ".symlink" {
		renderedFilename = strings.TrimSuffix(renderedFilename, ".symlink")
		path = strings.TrimSuffix(path, ".symlink")

		relativePath, err := filepath.Rel(filepath.Dir(renderedFilename), path)
		if err != nil {
			return fmt.Errorf("could not calculate path of symlink [%s]: %v", color.Magenta(renderedFilename), err)
		}
		if err := os.Symlink(relativePath, renderedFilename); err != nil {
			return fmt.Errorf("could not create symlink [%s]: %v", color.Magenta(renderedFilename), err)
		}
		return nil
	}
	// check if current file has a .symlink marker companion
	// if so we skip these files, we don't want to template/copy them over, we create symlinks for them (see above)
	if file.Exists(path + ".symlink") {
		return nil
	}
	// end of .symlink marker handling

	// if its a normal symlink then we copy it unmodified as-is
	if !info.Mode().IsRegular() && info.Mode()&fs.ModeSymlink != 0 {
		file.CopySymlink(path, renderedFilename)
		log.Debugf("copied symlink from [%s] to [%s]", color.Magenta(path), color.Magenta(renderedFilename))
		return nil
	}

	// decrypt .sops_enc files on the fly, write decrypted content to target
	if filepath.Ext(path) == ".sops_enc" {
		renderedFilename = strings.TrimSuffix(renderedFilename, ".sops_enc")

		log.Debugf("decrypt file [%s] into [%s]", color.Magenta(path), color.Magenta(renderedFilename))
		data, err := command.ExecOutput([]string{"sops", "-d", path})
		if err != nil {
			log.Errorf("could not decrypt file [%s]", color.Magenta(path))
			return err
		}
		file.Write(renderedFilename, data)
		return nil
	}

	// parse and render template
	if err := writeFile(baseFilename, config.DirSource(), renderedFilename, viper.AllSettings()); err != nil {
		return fmt.Errorf("could not render [%s]: %v", color.Magenta(baseFilename), err)
	}
	return nil
}

func writeFile(baseFilename, sourcePath, targetFile string, data interface{}) error {
	// ensure path exists to write file to
	if err := os.MkdirAll(filepath.Dir(targetFile), 0700); err != nil { // use mode 0700, since we are likely rendering sensitive data
		return err
	}

	f, err := os.Create(targetFile)
	if err != nil {
		log.Errorf("could not create file [%s]", color.Magenta(targetFile))
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	funcMap["filepath"] = func() string {
		return baseFilename
	}
	tmpl := template.New(baseFilename).Funcs(funcMap).Funcs(sprig.FuncMap()).Delims(config.DelimiterLeft(), config.DelimiterRight()).Option("missingkey=error")

	// parse template
	tmpl, err = tmpl.Parse(file.Read(filepath.Join(sourcePath, baseFilename)))
	if err != nil {
		log.Errorf("could not parse template [%s]", color.Magenta(baseFilename))
		return err
	}

	// use template, write output file
	if err := tmpl.Execute(w, data); err != nil {
		log.Errorf("could not execute template [%s]", color.Magenta(baseFilename))
		return err
	}

	// chmod +x *.sh
	if filepath.Ext(targetFile) == ".sh" {
		if err := command.Exec([]string{"chmod", "+x", targetFile}); err != nil {
			log.Errorf("could not chmod+x [%s]", color.Magenta(targetFile))
			return err
		}
	}
	// safety chmod for secrets
	if strings.Contains(baseFilename, "secrets") && targetFile != "/dev/stdout" {
		if err := command.Exec([]string{"chmod", "go-rwx", targetFile}); err != nil {
			log.Errorf("could not chmod go-rwx [%s]", color.Magenta(targetFile))
			return err
		}
	}
	return nil
}
