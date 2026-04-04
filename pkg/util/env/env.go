package env

import (
	"bufio"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/JamesClonk/plato/pkg/util/color"
	"github.com/JamesClonk/plato/pkg/util/log"
)

func Set(key string, value string) error {
	return os.Setenv(key, value)
}

func Get(key string, nvl string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return nvl
	}
	return value
}

func MustGet(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatal("%s", color.Fail("environment variable [%s] is missing!", key))
	}
	return value
}

func SourceFile(filename string) {
	sourceFile := func(file *os.File) {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "export") || strings.Contains(line, "=") {
				line = strings.TrimPrefix(line, "export ")
				keyValue := strings.SplitN(line, "=", 2)
				os.Setenv(strings.TrimSpace(keyValue[0]), strings.Trim(strings.TrimSpace(keyValue[1]), `"'`))
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(color.Red("%v", err))
		}
	}

	// source file from etc directory first (lowest priority)
	etcFile, err := os.Open(path.Join("/etc", strings.TrimPrefix(filename, ".")))
	if err == nil {
		sourceFile(etcFile)
	}

	// then try to source file from home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatal("%s", color.Red("%v", err))
	}
	homeFile, err := os.Open(path.Join(usr.HomeDir, filename))
	if err == nil {
		sourceFile(homeFile)
	}

	// finally source file from current working directory (highest priority)
	workFile, err := os.Open(filename)
	if err == nil {
		sourceFile(workFile)
	}
}
