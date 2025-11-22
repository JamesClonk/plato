package file

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/JamesClonk/plato/pkg/util/color"
	"github.com/JamesClonk/plato/pkg/util/command"
	"github.com/JamesClonk/plato/pkg/util/log"
)

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func Delete(filename string) error {
	return os.Remove(filename)
}

func SHA1(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf(color.Red("could not open file [%s]: %v", color.Magenta(filename), err))
	}
	defer f.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatalf(color.Red("could not read file [%s]: %v", color.Magenta(f.Name()), err))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func SHA256(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf(color.Red("could not open file [%s]: %v", color.Magenta(filename), err))
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatalf(color.Red("could not read file [%s]: %v", color.Magenta(f.Name()), err))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func Checksum(file, expectedSHA string) {
	if err := ChecksumError(file, expectedSHA); err != nil {
		log.Fatalf("%v", err)
	}
}

func ChecksumError(file, expectedSHA string) error {
	switch len(expectedSHA) {
	case 40:
		sha1 := SHA1(file)
		if sha1 != expectedSHA {
			return fmt.Errorf("sha1 checksum mismatch [%s]! actual:[%s], expected:[%s]",
				color.Cyan(file), color.Red(sha1), color.Green(expectedSHA))
		}
	case 64:
		sha256 := SHA256(file)
		if sha256 != expectedSHA {
			return fmt.Errorf("sha256 checksum mismatch for [%s]! actual:[%s], expected:[%s]",
				color.Cyan(file), color.Red(sha256), color.Green(expectedSHA))
		}
	default:
		return fmt.Errorf(color.Red("valid sha checksum missing for [%s]", file))
	}
	return nil
}

func Read(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf(color.Red("could not read filename [%s]: %v", color.Magenta(filename), err))
	}
	return string(data)
}

func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func Contains(filename, content string) bool {
	return strings.Contains(Read(filename), content)
}

func Touch(filename string) {
	if _, err := os.Create(filename); err != nil {
		log.Fatalf(color.Red("could not create file [%s]: %v", color.Magenta(filename), err))
	}
}

func Write(filename, content string) {
	if err := os.WriteFile(filename, []byte(content), 0664); err != nil {
		log.Errorf(color.Red("could not write to file [%s]:\n%s", color.Magenta(filename), content))
		log.Fatalf(color.Red("%v", err))
	}
}

func Prepend(filename, content string) {
	data := Read(filename)

	info, err := os.Stat(filename)
	if err != nil {
		log.Fatalf(color.Red("could not stat filename [%s]: %v", color.Magenta(filename), err))
	}
	mode := info.Mode()

	if err := os.WriteFile(filename, []byte(content+data), mode); err != nil {
		log.Errorf(color.Red("could not write to file [%s]:\n%s", color.Magenta(filename), content))
		log.Fatalf(color.Red("%v", err))
	}
}

func CopyTo(file, folder string) {
	command.Run([]string{"cp", "-f", file, folder + "/."})
}

func CopySymlink(source, dest string) {
	link, err := os.Readlink(source)
	if err != nil {
		log.Errorf("could not read symlink [%s]", color.Magenta(source))
		log.Fatalf(color.Red("%v", err))
	}
	if err := os.Symlink(link, dest); err != nil {
		log.Errorf("could not copy symlink [%s] to [%s]", color.Magenta(source), color.Magenta(dest))
		log.Fatalf(color.Red("%v", err))
	}
}
