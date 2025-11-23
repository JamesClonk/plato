.DEFAULT_GOAL := help
SHELL := /bin/bash
APP = plato
COMMIT_SHA = $(shell git rev-parse --short HEAD)

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: trust
## trust: trust all mise config files
trust:
	mise trust --all || true

.PHONY: dev
## dev: runs main.go with the golang race detector
dev: trust
	source .env*; go run -race main.go

.PHONY: run
## run: runs binary
run: build
	source .env*; ./plato

.PHONY: build
## build: builds the application
build: clean
	@echo "Building binary ..."
	go build -o ${APP}

.PHONY: install
## install: installs the application
install: clean
	@echo "Installing binary ..."
	go install

.PHONY: clean
## clean: cleans up binary files
clean: trust
	@echo "Cleaning up ..."
	go clean

.PHONY: test
## test: runs go test with the race detector
test: trust
	@echo "{{{ .ssh.public_key -}}}" | PLATO_WORKING_DIR=_fixtures/combined SOPS_AGE_KEY_FILE=../age.key go run -race main.go template | grep 'AAAAC3NzaC1lZDI1NTE5AAAAIMmROMaltL5DRYCYNHW9BE2BTZWUqZ3dISyzZopyrRjg' || (echo "failed basic stdin/stdout test" && exit 1)
	SOPS_AGE_KEY_FILE="$$PWD/_fixtures/age.key" GOARCH=amd64 GOOS=linux TZ="" go test -tags=faketime -v -race ./...
	#GOARCH=amd64 GOOS=linux go test -v -race ./...

.PHONY: vendor
## vendor: cleans up binary files
vendor:
	@echo "Updating vendor ..."
	go mod tidy
	go mod vendor

.PHONY: release-test
## release-test: run a test release with goreleaser
release-test:
	@echo "Testing goreleaser ..."
	goreleaser release --clean --skip=announce,publish,validate

.PHONY: release
## release: publish a release with goreleaser
release:
	@echo "Running goreleaser ..."
	goreleaser release --clean

# ======================================================================================================================
# individual commands for testing
# ======================================================================================================================
.PHONY: plato-combined-stdin
## plato-combined-stdin: renders a template provided via STDIN into STDOUT, using an encrypted plato.yaml
plato-combined-stdin:
	@echo "{{{ .ssh.public_key -}}}" | PLATO_WORKING_DIR=_fixtures/combined SOPS_AGE_KEY_FILE=../age.key go run -race main.go template

.PHONY: plato-combined-stdout
## plato-combined-stdout: renders a template file into STDOUT, using an encrypted plato.yaml
plato-combined-stdout:
	@PLATO_WORKING_DIR=_fixtures/combined SOPS_AGE_KEY_FILE=../age.key go run -race main.go template ssh_private_key

.PHONY: plato-combined-output
## plato-combined-output: renders a template file into an output file, using an encrypted plato.yaml
plato-combined-output:
	@PLATO_WORKING_DIR=_fixtures/combined SOPS_AGE_KEY_FILE=../age.key go run -race main.go template ssh_private_key tmp/ssh_private_key

.PHONY: plato-file-stdout
## plato-file-stdout: renders a template file into STDOUT
plato-file-stdout:
	@PLATO_WORKING_DIR=_fixtures SOPS_AGE_KEY_FILE=age.key go run -race main.go template input/minio.yaml

.PHONY: plato-file-output
## plato-file-output: renders a template file into an output file
plato-file-output:
	@PLATO_WORKING_DIR=_fixtures SOPS_AGE_KEY_FILE=age.key go run -race main.go template input/minio.yaml tmp/test/minio.yaml

.PHONY: plato-render
## plato-render: renders all templates
plato-render:
	@PLATO_WORKING_DIR=_fixtures SOPS_AGE_KEY_FILE=age.key go run -race main.go render

.PHONY: plato-store-secrets
## plato-store-secrets: stores secrets back into SOPS file
plato-store-secrets: plato-render
	@PLATO_WORKING_DIR=_fixtures SOPS_AGE_KEY_FILE=age.key go run -race main.go store-secrets

.PHONY: plato-version
## plato-version: displays PLATO version
plato-version:
	@PLATO_WORKING_DIR=_fixtures SOPS_AGE_KEY_FILE=age.key go run -race main.go version

.PHONY: vault-run
## vault-run: runs vault locally in dev mode
vault-run:
	@source .env*; vault server -dev -dev-root-token-id="root"

.PHONY: vault-status
## vault-status: show vault status
vault-status:
	VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=root vault status

.PHONY: vault-prepare
## vault-prepare: prepare vault for SOPS
vault-prepare:
	VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=root vault secrets enable -path=sops transit || true
	VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=root vault write sops/keys/plato type=rsa-4096 || true

.PHONY: vault-update-secrets
## vault-update-secrets: add/update vault key to SOPS encrypted file
vault-update-secrets:
	@source .env*; cd _fixtures; PLATO_WORKING_DIR=_fixtures SOPS_AGE_KEY_FILE=age.key VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=root sops -r -i --rm-hc-vault-transit http://127.0.0.1:8200/v1/sops/keys/plato secrets.yaml || true
	@source .env*; cd _fixtures; PLATO_WORKING_DIR=_fixtures SOPS_AGE_KEY_FILE=age.key VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=root sops -r -i --add-hc-vault-transit http://127.0.0.1:8200/v1/sops/keys/plato secrets.yaml

.PHONY: vault-show-secrets
## vault-show-secrets: open SOPS encrypted file via vault
vault-show-secrets:
	@source .env*; cd _fixtures; VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=root sops secrets.yaml

.PHONY: age-show-secrets
## age-show-secrets: open SOPS encrypted file via age.key
age-show-secrets:
	@source .env*; cd _fixtures; SOPS_AGE_KEY_FILE=age.key sops secrets.yaml
