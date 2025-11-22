package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JamesClonk/plato/pkg/config"
	"github.com/JamesClonk/plato/pkg/util/color"
	"github.com/JamesClonk/plato/pkg/util/command"
	"github.com/JamesClonk/plato/pkg/util/dir"
	"github.com/JamesClonk/plato/pkg/util/file"
	"github.com/JamesClonk/plato/pkg/util/log"
	"github.com/spf13/viper"
)

func StoreGeneratedSecrets() {
	log.Infof("storing secrets back into [%s] ...", color.Magenta("secrets.yaml"))

	// re-encrypt all former .sops_enc files back to their original location
	err := filepath.Walk(config.DirSource(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".sops_enc" {
			baseFilename := strings.TrimPrefix(path, config.DirSource()+string(os.PathSeparator))
			renderedFilename := strings.TrimSuffix(filepath.Join(config.DirTarget(), baseFilename), ".sops_enc")

			// check if the rendered file exists
			if !file.Exists(renderedFilename) {
				// doesnt' exist, we can't re-encrypt and re-store it.. obviously!
				return nil
			}

			// check if content has changed, no need to re-encrypt the file back otherwise (avoids unnecessary git spam)
			decrypted, err := command.ExecOutput([]string{"sops", "-d", path})
			if err != nil {
				log.Errorf("could not decrypt file [%s]", color.Magenta(path))
				return err
			}
			if decrypted == file.Read(renderedFilename) {
				// content matches, don't re-encrypt!
				return nil
			}

			log.Debugf("encrypt file [%s] into [%s]", color.Magenta(renderedFilename), color.Magenta(path))
			data, err := command.ExecOutput([]string{"sops", "-e", "--input-type", "binary", renderedFilename})
			if err != nil {
				log.Errorf("could not encrypt file [%s]", color.Magenta(renderedFilename))
				return err
			}
			file.Write(path, data)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("could not re-encrypt *.sops_enc files from [%s]: %v", color.Magenta(config.DirSource()), err)
	}

	// go through all */secrets files
	if dir.Exists(config.DirGeneratedSecrets()) {
		err = filepath.Walk(config.DirGeneratedSecrets(), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// skip directories
			if info.IsDir() {
				return nil
			}
			return processFile(path, info) // store file contents back into secrets.yaml
		})
		if err != nil {
			log.Fatalf("could not work through [%s]: %v", color.Magenta(config.DirGeneratedSecrets()), err)
		}
	}

	// delete temporary .secrets-updated marker file to remove gitrepo taint
	_ = os.Remove(".secrets-updated")
}

func processFile(path string, info os.FileInfo) error {
	// exclude files we obviously didn't template and/or want to store in secrets.yaml
	ext := filepath.Ext(path)
	if ext == ".md" ||
		ext == ".txt" ||
		ext == ".zip" ||
		ext == ".tar.gz" ||
		ext == ".tgz" {
		return nil
	}

	filename := filepath.Base(path)
	data := file.Read(path)

	// check if data has actually changed, no need to store it back otherwise (avoids unnecessary git spam)
	if data == viper.GetString(filename) {
		// content matches, don't store!
		return nil
	}

	// turn data into single-line JSON compatible string (needed for SOPS, it can only take JSON input)
	data = strings.Replace(data, "\r", "", -1) // have to remove windows garbage if present
	data = strings.Replace(data, "\n", `\n`, -1)
	data = strings.Replace(data, `"`, `\"`, -1)
	data = strings.Trim(data, "\n")

	// prepare command value, for format see https://github.com/getsops/sops#set-a-sub-part-in-a-document-tree
	pathParts := strings.SplitN(filename, ".", -1)
	var value string
	for _, part := range pathParts {
		value = fmt.Sprintf("%s[\"%s\"]", value, part)
	}
	value = fmt.Sprintf("%s \"%s\"", value, data)

	// set value in-place
	log.Debugf("store secret [%s]", color.Magenta(filename))
	if err := command.Exec([]string{"sops", "--set", value, "secrets.yaml"}); err != nil {
		log.Errorf("could not store secret [%s]", color.Magenta(filename))
		return err
	}
	return nil
}
