package store

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/JamesClonk/plato/pkg/config"
	"github.com/JamesClonk/plato/pkg/util/command"
	"github.com/JamesClonk/plato/pkg/util/dir"
	"github.com/JamesClonk/plato/pkg/util/file"
	"github.com/JamesClonk/plato/pkg/util/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Chdir("../../_fixtures") // change to fixtures dir
	log.Initialize()
	dir.Remove(config.DirTarget())
	dir.Create(config.DirTarget())

	os.Setenv("SOPS_AGE_KEY_FILE", "age.key")
	config.InitConfig()
	config.InitSecrets()

	seed, _ := strconv.ParseInt(strings.Trim(command.RunOutput([]string{"date", `+%s%N`}), "\n"), 10, 64)
	rand.Seed(seed)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Test_StoreGeneratedSecrets(t *testing.T) {
	targetFile := "input/infrastructure/example.json.sops_enc"
	_ = os.Remove(targetFile)
	assert.False(t, file.Exists(targetFile))

	file.Touch(targetFile)
	assert.True(t, file.Exists(targetFile))
	err := command.Exec([]string{"sops", "-e", "-i", "--input-type", "binary", targetFile})
	assert.NoError(t, err)

	sourceFile := "output/infrastructure/example.json"
	content := `{"test": {
	"debug": "on",
	"window": {
		"title": "Sample Konfabulator Widget",
		"name": "main_window",
		"width": 500,
		"height": 500
	}
}}`
	dir.Create(filepath.Dir(sourceFile))
	file.Write(sourceFile, content)

	StoreGeneratedSecrets()
	assert.True(t, file.Exists(targetFile))
	data := file.Read(targetFile)
	assert.True(t, strings.Contains(data, `"recipient": "age1yapc0k0tfz8cketuldrjq3vyuzne4587zmf3d2ejypaftg95yvrs8r44yh",`))
	assert.True(t, strings.Contains(data, `"sops": {`))
	assert.True(t, strings.Contains(data, `"enc": "-----BEGIN AGE ENCRYPTED FILE-----`))
	assert.True(t, strings.Contains(data, `"data": "ENC[AES256_GCM,data:`))

	// read and compare encrypted file to original json content
	decryptedData, err := command.ExecOutput([]string{"sops", "-d", targetFile})
	assert.NoError(t, err)
	assert.Equal(t, content, decryptedData)
}

func Test_processFile_store_secrets(t *testing.T) {
	testFile := "output/secrets/test.myconfig"
	_ = os.Remove(testFile)
	dir.Create(filepath.Dir(testFile))

	// reset viper test property, to be sure
	viper.Set("test.myconfig", "TEST")

	// generate randomized multiline test string
	randomName := randStringRunes(16)
	data := `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: deadbeef-beefdead
    server: https://my.super.kubernetes.cluster:6443
  name: default
contexts:
- context:
    cluster: default
    user: ` + randomName + `
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: default
  user:
    client-certificate-data: deadbeef-beefdead
`
	file.Write(testFile, data)
	assert.True(t, file.Exists(testFile))

	info, err := os.Lstat(testFile)
	assert.NoError(t, err)

	// processFile() should store back into secrets.yaml
	err = processFile(testFile, info)
	assert.NoError(t, err)

	// read secrets
	decryptedSecrets, err := command.ExecOutput([]string{"sops", "-d", "secrets.yaml"})
	assert.NoError(t, err)
	assert.True(t, strings.Contains(decryptedSecrets, "server: https://my.super.kubernetes.cluster:6443"))
	assert.True(t, strings.Contains(decryptedSecrets, "client-certificate-data: deadbeef-beefdead\n"))
	assert.True(t, strings.Contains(decryptedSecrets, "user: "+randomName+"\n"))

	// reset viper, to be absolutely sure the payload got written correctly
	viper.Reset()
	viper.SetConfigType("yaml")
	config.InitSecrets()
	assert.Equal(t, data, viper.GetString("test.myconfig"))
}

func Test_processFile_filetype_exclusion(t *testing.T) {
	testFile := "output/secrets/test.md"
	_ = os.Remove(testFile)
	dir.Create(filepath.Dir(testFile))

	data := `This is a test file!
:)
`
	file.Write(testFile, data)

	info, err := os.Lstat(testFile)
	assert.NoError(t, err)

	// processFile() should not store a markdown file back into secrets.yaml
	err = processFile(testFile, info)
	assert.NoError(t, err)
	assert.True(t, file.Exists(testFile))

	// read secrets
	decryptedSecrets, err := command.ExecOutput([]string{"sops", "-d", "secrets.yaml"})
	assert.NoError(t, err)
	assert.True(t, !strings.Contains(decryptedSecrets, "This is a test file"))

	// reset viper, to be absolutely sure the payload did not get written back
	viper.Reset()
	viper.SetConfigType("yaml")
	config.InitSecrets()
	assert.Equal(t, "", viper.GetString("test.md"))
}
