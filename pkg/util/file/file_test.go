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
	assert.Equal(t, "e478c6c055f8d780b2c3e760144b445e2ebef73f", SHA1("../main.go"))
	assert.Equal(t, "01627f3c6730f54ae8da947408b77ec6350624eadf1b01fbc800695a7365bd93", SHA256("../main.go"))
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
	assert.Equal(t, expected, strings.Join(lines[4:7], "\n"))
}
