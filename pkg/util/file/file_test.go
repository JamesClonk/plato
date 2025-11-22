package file

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Chdir("../../../_fixtures") // change to fixtures dir
}

func Test_File_SHA1(t *testing.T) {
	assert.Equal(t, "4565f4184adba56983d389849b58365085324db4", SHA1("../main.go"))
	assert.Equal(t, "2ced1d4d19270d86a1a2db3323d7c088a2d13fedce946bf712a986731e66d23d", SHA256("../main.go"))
}

func Test_File_Read(t *testing.T) {
	content := Read("../main.go")
	expected := `package main

import "github.com/JamesClonk/plato/cmd"

var (
	version = "0.0.0"
	commit  = "-"
	date    = "now"
)

func main() {
	cmd.Execute(version, commit, date)
}
`
	assert.Equal(t, expected, content)
}

func Test_File_ReadLines(t *testing.T) {
	lines, err := ReadLines("../main.go")
	assert.NoError(t, err)
	assert.Equal(t, 13, len(lines))
	assert.Equal(t, `var (`, lines[4])

	expected := `func main() {
	cmd.Execute(version, commit, date)
}`
	assert.Equal(t, expected, strings.Join(lines[10:13], "\n"))
}
