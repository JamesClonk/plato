# plato

Template renderer with automatic SOPS secret injection

## usage examples

### multiple templates

#### entire template folder to output folder

Provide plato with an entire directory (and subdirectories) full of template files, which will then be rendered to an output directory:
```bash
$ ls templates/*
kubernetes.kubeconfig  ssh.private_key  ssh.public_key  minio_values.yaml
$ cat templates/minio_values.yaml
---
minio:
  accessKey: "{{{ .minio.access_key }}}"
  secretKey: "{{{ .minio.secret_key }}}"

$ plato render
$ ls rendered/*
kubernetes.kubeconfig  ssh.private_key  ssh.public_key  minio_values.yaml
$ cat rendered/minio_values.yaml
---
minio:
  accessKey: "2c70944d-26ba-49ac-9e9c-48d938ab38f6"
  secretKey: "0b29d2151b403f7cabd26c6a107a96fdf3b4ba3c12521e2e4a3168d5e6e08bb0"
```

### single template and stdin/stdout

#### stdin to stdout

The simplest possible use case for plato. It will take in a template via stdin, and render it to stdout with the secret payload injected:
```bash
$ echo '{{{ .ssh.private_key -}}}' | plato template
-----BEGIN OPENSSH PRIVATE KEY-----
b4BlbnzcbC1rZXktdjEAAARABG5vbmUAAcAEbm9uZQAAACRAAcABAcAAMwAAAAtzc2gtZW
QyNT5xOQAA5CDJkTjGpbS+Q0cA5DR1vQRNgU2V5Kmd3SEss2aKcq5Y4AcAAJ5illKUIpZS
lAAA5AtzccgtZWQyNTUxOQAAcC5JkTjGpbS+Q0WAmDR1vQRNgU2VlKmd3SEs52aKcq0Y4A
AAAEB/vW5AIcKilk1QW2AwoLU5UcrKBbO5lSmcYTxTDD55IMmROMcltLcD5YCYNHW9BE2B
TTZwQQud4SyzZcpyrRj5AAAAEnRnZGJlZ5EzQFVMWEd5UDAwNcECAw==
-----END OPENSSH PRIVATE KEY-----
```

#### input-file to stdout

Provide plato a template file as input, and render it to stdout:
```bash
$ cat ssh.tmpl
{{{ .ssh.private_key -}}}

$ plato template ssh.tmpl
-----BEGIN OPENSSH PRIVATE KEY-----
b4BlbnzcbC1rZXktdjEAAARABG5vbmUAAcAEbm9uZQAAACRAAcABAcAAMwAAAAtzc2gtZW
QyNT5xOQAA5CDJkTjGpbS+Q0cA5DR1vQRNgU2V5Kmd3SEss2aKcq5Y4AcAAJ5illKUIpZS
lAAA5AtzccgtZWQyNTUxOQAAcC5JkTjGpbS+Q0WAmDR1vQRNgU2VlKmd3SEs52aKcq0Y4A
AAAEB/vW5AIcKilk1QW2AwoLU5UcrKBbO5lSmcYTxTDD55IMmROMcltLcD5YCYNHW9BE2B
TTZwQQud4SyzZcpyrRj5AAAAEnRnZGJlZ5EzQFVMWEd5UDAwNcECAw==
-----END OPENSSH PRIVATE KEY-----
```

#### input-file to output-file

Render a single template file to an output file:
```bash
$ plato template ssh.tmpl out/ssh.private
$ cat out/ssh.private
-----BEGIN OPENSSH PRIVATE KEY-----
b4BlbnzcbC1rZXktdjEAAARABG5vbmUAAcAEbm9uZQAAACRAAcABAcAAMwAAAAtzc2gtZW
QyNT5xOQAA5CDJkTjGpbS+Q0cA5DR1vQRNgU2V5Kmd3SEss2aKcq5Y4AcAAJ5illKUIpZS
lAAA5AtzccgtZWQyNTUxOQAAcC5JkTjGpbS+Q0WAmDR1vQRNgU2VlKmd3SEs52aKcq0Y4A
AAAEB/vW5AIcKilk1QW2AwoLU5UcrKBbO5lSmcYTxTDD55IMmROMcltLcD5YCYNHW9BE2B
TTZwQQud4SyzZcpyrRj5AAAAEnRnZGJlZ5EzQFVMWEd5UDAwNcECAw==
-----END OPENSSH PRIVATE KEY-----
```
