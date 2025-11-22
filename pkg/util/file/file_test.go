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
	assert.Equal(t, "ee0fce81001f723072c97d89fd0cb19c5ff17790", SHA1("../main.go"))
	assert.Equal(t, "da34fc90790bcaf0fc8ac48bb56307675d6aee151000386046b113540c31e99d", SHA256("../main.go"))
}

func Test_File_Read(t *testing.T) {
	content := Read("../main.go")
	expected := `package main

import "github.com/JamesClonk/plato/cmd"

func main() {
	cmd.Execute()
}
`
	assert.Equal(t, expected, content)
}

func Test_File_ReadLines(t *testing.T) {
	lines, err := ReadLines("../main.go")
	assert.NoError(t, err)
	assert.Equal(t, 7, len(lines))
	assert.Equal(t, `func main() {`, lines[4])

	expected := `func main() {
	cmd.Execute()
}`
	assert.Equal(t, expected, strings.Join(lines[4:6], "\n"))
}
