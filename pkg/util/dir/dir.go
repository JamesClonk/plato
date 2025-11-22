package dir

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/JamesClonk/plato/pkg/util/color"
	"github.com/JamesClonk/plato/pkg/util/file"
	"github.com/JamesClonk/plato/pkg/util/log"
)

type File struct {
	Name     string
	Filename string
	File     string
	Content  string
}

func ExpandTilde(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("could not get current OS user: %s", color.Red("%v", err))
	}
	return usr.HomeDir + path[1:]
}

func Exists(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

func Create(path string, mode ...os.FileMode) {
	if !Exists(path) {
		if len(mode) == 0 {
			mode = append(mode, 0750)
		}
		if err := os.MkdirAll(path, mode[0]); err != nil {
			log.Fatalf("%s", color.Red("could not create directory [%s]: %v", color.Magenta(path), err.Error()))
		}
		log.Debugf("created directory [%s]", color.Magenta(path))
	}
}

func Remove(path string) {
	if Exists(path) {
		if err := os.RemoveAll(path); err != nil {
			log.Fatalf("%s", color.Red("could not delete directory [%s]: %v", color.Magenta(path), err.Error()))
		}
		log.Debugf("removed directory [%s]", color.Magenta(path))
	}
}

func Files(path string) []File {
	files := make([]File, 0)

	fs, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("%s", color.Red("could not read directory [%s]: %v", color.Magenta(path), err.Error()))
	}

	for _, f := range fs {
		if !f.IsDir() {
			fp := filepath.Join(path, f.Name())
			files = append(files, File{
				Name:     strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())),
				Filename: f.Name(),
				File:     fp,
				Content:  strings.Trim(file.Read(fp), "\n"),
			})
		}
	}
	return files
}

func Dirs(path string) []string {
	dirs := make([]string, 0)

	fs, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("%s", color.Red("could not read directory [%s]: %v", color.Magenta(path), err.Error()))
	}

	for _, f := range fs {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}
	return dirs
}
