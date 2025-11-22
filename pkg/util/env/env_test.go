package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	os.Chdir("../../../") // change to plato-src root
}

func Test_Env_Get(t *testing.T) {
	assert.NotEqual(t, "foo", Get("PATH", "foo"))
	assert.Equal(t, "bar", Get("foo", "bar"))
}

func Test_Env_SourceFile(t *testing.T) {
	assert.Nil(t, Set("MY_KEY", "foobar"))
	assert.Equal(t, "foobar", Get("MY_KEY", "foobar"))
	SourceFile("_fixtures/env_file")
	assert.Equal(t, "my_value", Get("MY_KEY", "foobar"))
}
