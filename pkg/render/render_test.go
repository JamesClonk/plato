package render

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/JamesClonk/plato/pkg/config"
	"github.com/JamesClonk/plato/pkg/util/dir"
	"github.com/JamesClonk/plato/pkg/util/file"
	"github.com/JamesClonk/plato/pkg/util/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Chdir("../../_fixtures") // change to fixtures dir
	log.Initialize()
	os.Setenv("SOPS_AGE_KEY_FILE", "age.key")
	config.InitConfig()
	config.LoadSecrets()
	dir.Remove(config.DirTarget())
	dir.Create(config.DirTarget())
}

func Test_RenderFile(t *testing.T) {
	sourceFile := "input/minio.yaml"
	targetFile := "tmp/test/minio.yaml"
	targetFolder := "tmp/test"

	_ = os.RemoveAll(targetFolder)
	assert.False(t, file.Exists(targetFile))
	assert.False(t, dir.Exists(targetFolder))

	RenderFile(sourceFile, targetFile)
	assert.True(t, file.Exists(targetFile))
	assert.True(t, dir.Exists(targetFolder))

	data := file.Read(targetFile)
	assert.Equal(t, `---
minio:
  accessKey: "2c70944d-26ba-49ac-9e9c-48d938ab38f6"
  secretKey: "0b29d2151b403f7cabd26c6a107a96fdf3b4ba3c12521e2e4a3168d5e6e08bb0"
`, data)
}

func Test_RenderTemplates(t *testing.T) {
	targetFileA := "output/infrastructure/terraform/.terraform.lock.hcl"
	targetFileB := "output/infrastructure/terraform/.terraform/some-file"
	targetFileC := "output/infrastructure/terraform/.terraform/providers/another-file"
	targetFolderA := "output/infrastructure/terraform/.terraform/providers"
	targetFolderB := "output/infrastructure/terraform/fake-folder"

	assert.False(t, file.Exists(targetFileA))
	assert.False(t, file.Exists(targetFileB))
	assert.False(t, file.Exists(targetFileC))
	assert.False(t, dir.Exists(targetFolderA))
	assert.False(t, dir.Exists(targetFolderB))

	RenderTemplates(true, true)
	assert.True(t, file.Exists(targetFileA))
	assert.False(t, file.Exists(targetFileB))
	assert.False(t, file.Exists(targetFileC))
	assert.False(t, dir.Exists(targetFolderA))
	assert.False(t, dir.Exists(targetFolderB))

	// create .terraform data
	dir.Create(targetFolderA)
	file.Touch(targetFileB)
	file.Touch(targetFileC)
	RenderTemplates(false, false) // don't cleanup anything
	assert.True(t, file.Exists(targetFileA))
	assert.True(t, file.Exists(targetFileB))
	assert.True(t, file.Exists(targetFileC))
	assert.True(t, dir.Exists(targetFolderA))
	assert.False(t, dir.Exists(targetFolderB))

	// create fake folder
	dir.Create(targetFolderB)
	RenderTemplates(true, false) // cleanup only .terraform
	assert.True(t, file.Exists(targetFileA))
	assert.False(t, file.Exists(targetFileB))
	assert.False(t, file.Exists(targetFileC))
	assert.True(t, dir.Exists(targetFolderA)) // only cleans files within .terraform/, not the dir itself
	assert.True(t, dir.Exists(targetFolderB))

	RenderTemplates(true, true) // cleanup everything
	assert.True(t, file.Exists(targetFileA))
	assert.False(t, file.Exists(targetFileB))
	assert.False(t, file.Exists(targetFileC))
	assert.False(t, dir.Exists(targetFolderA))
	assert.False(t, dir.Exists(targetFolderB))

	data := file.Read(filepath.Join(config.DirTarget(), "infrastructure/yaml.yaml"))
	assert.Equal(t, `---
apiVersion: super.cluster.io/v1
kind: CustomResource
metadata:
  kubernetes:
    kubeconfig: |-
      apiVersion: v1
      clusters:
      - cluster:
          certificate-authority-data: deadbeef-beefdead
          server: https://my.super.kubernetes.cluster:6443
        name: default
      contexts:
      - context:
          cluster: default
          user: default
        name: default
      current-context: default
      kind: Config
      preferences: {}
      users:
      - name: default
        user:
          client-certificate-data: deadbeef-beefdead
    metallb:
      bgp_config:
        my_asn: 65477
        password: 6GUD4rh3QIejbsD7yTF5n7zEUb7ofYNp
        peer_asn: 4777444999
    server: https://my.super.kubernetes.cluster:6443`, data)
}

func Test_writeFile_with_custom_funcmaps(t *testing.T) {
	kubernetes := make(map[any]any)
	kubernetes["server"] = "https://my.super.kubernetes.cluster:6443"
	payload := make(map[any]any)
	payload["kubernetes"] = kubernetes
	payload["cidr"] = "100.106.160.64/26"

	filename := "infrastructure/terraform/settings/22_folder_test.yaml"
	err := writeFile(filename, config.DirSource(), filepath.Join(config.DirTarget(), filename), payload)
	assert.NoError(t, err)

	data := file.Read(filepath.Join(config.DirTarget(), filename))
	assert.Equal(t, `---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  labels:
    my.super.kubernetes.cluster/component: test-ingress
    my.super.kubernetes.cluster/subdomain: settings
  name: test-ingress
spec:
  rules:
  - host: settings.my.super.kubernetes.cluster
  tls:
  - hosts:
    - settings.my.super.kubernetes.cluster
`, data)

	filename = "infrastructure/terraform/cidr.yaml"
	err = writeFile(filename, config.DirSource(), filepath.Join(config.DirTarget(), filename), payload)
	assert.NoError(t, err)

	data = file.Read(filepath.Join(config.DirTarget(), filename))
	assert.Equal(t, `---
ip: 100.106.160.77
`, data)
}

func Test_processFile_with_Symlink(t *testing.T) {
	// processFile first for the target of the symlink, so that it gets copied too.
	// otherwise the symlink will point to void and further checks will fail.
	source := "input/secrets/kubernetes.kubeconfig"
	target := "output/secrets/kubernetes.kubeconfig"
	_ = os.Remove(target)

	info, err := os.Lstat(source)
	assert.NoError(t, err)

	err = processFile(source, info)
	assert.NoError(t, err)
	assert.True(t, file.Exists(target))

	// now process and check the symlink if it gets copied over as expected
	source = "input/infrastructure/terraform/kubernetes.kubeconfig"
	target = "output/infrastructure/terraform/kubernetes.kubeconfig"
	_ = os.Remove(target)

	info, err = os.Lstat(source)
	assert.NoError(t, err)

	err = processFile(source, info)
	assert.NoError(t, err)
	assert.True(t, file.Exists(target))

	lines, err := file.ReadLines(target)
	assert.NoError(t, err)
	assert.Equal(t, `apiVersion: v1`, lines[0])
	assert.Equal(t, `clusters:`, lines[1])

	info, err = os.Lstat(target)
	assert.True(t, !info.Mode().IsRegular())
	assert.True(t, info.Mode()&fs.ModeSymlink != 0)

	link, err := os.Readlink(target)
	assert.NoError(t, err)
	assert.Equal(t, `../../secrets/kubernetes.kubeconfig`, link)
}

func Test_processFile_with_Symlink_Marker(t *testing.T) {
	source := "input/infrastructure/terraform/terraform.tfstate"
	target := "output/infrastructure/terraform/terraform.tfstate"
	_ = os.Remove(target)

	info, err := os.Lstat(source)
	assert.NoError(t, err)

	err = processFile(source, info)
	assert.NoError(t, err)
	assert.True(t, !file.Exists(target)) // files that have a .symlink companion should be ignored

	source = "input/infrastructure/terraform/terraform.tfstate.symlink"
	target = "output/infrastructure/terraform/terraform.tfstate"
	_ = os.Remove(target)

	info, err = os.Lstat(source)
	assert.NoError(t, err)

	err = processFile(source, info)
	assert.NoError(t, err)
	assert.True(t, file.Exists(target))
	assert.Equal(t, "dev-state\n", file.Read(target))

	info, err = os.Lstat(target)
	assert.True(t, !info.Mode().IsRegular())
	assert.True(t, info.Mode()&fs.ModeSymlink != 0)

	link, err := os.Readlink(target)
	assert.NoError(t, err)
	assert.Equal(t, `../../../input/infrastructure/terraform/terraform.tfstate`, link)
}

func Test_processFile_with_SOPS_encrypted_file(t *testing.T) {
	source := "input/infrastructure/terraform/terraform.tfstate.backup.sops_enc"
	target := "output/infrastructure/terraform/terraform.tfstate.backup"
	_ = os.Remove(target)

	info, err := os.Lstat(source)
	assert.NoError(t, err)

	err = processFile(source, info)
	assert.NoError(t, err)
	assert.True(t, file.Exists(target))
	assert.Equal(t, `{"widget": {
	"debug": "on",
	"window": {
		"title": "Sample Konfabulator Widget",
		"name": "main_window",
		"width": 500,
		"height": 500
	},
	"image": {
		"src": "Images/Sun.png",
		"name": "sun1",
		"hOffset": 250,
		"vOffset": 250,
		"alignment": "center"
	},
	"text": {
		"data": "Click Here",
		"size": 36,
		"style": "bold",
		"name": "text1",
		"hOffset": 250,
		"vOffset": 100,
		"alignment": "center",
		"onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;"
	}
}}
`, file.Read(target))
}
