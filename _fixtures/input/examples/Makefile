.DEFAULT_GOAL := help
SHELL := /bin/bash

# e.g. darwin, windows
export BIN_OS ?= linux
# e.g. arm64
export BIN_ARCH ?= amd64

ROOT_DIR = $(realpath .)
TERRAFORM_DIR = ${ROOT_DIR}/terraform
TERRAFORM_OUTPUT = ${TERRAFORM_DIR}/output.json

# container registry
REGISTRY_HOST={{{ .registry.hostname }}}:{{{ .registry.port }}}
REGISTRY_USERNAME={{{ .registry.username }}}
REGISTRY_PASSWORD={{{ .registry.password }}}

# ======================================================================================================================
.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
# ======================================================================================================================

# ======================================================================================================================
.PHONY: registry-login
registry-login:
	@docker login ${REGISTRY_HOST} --username ${REGISTRY_USERNAME} --password ${REGISTRY_PASSWORD}
	@helm registry login ${REGISTRY_HOST} --username ${REGISTRY_USERNAME} --password ${REGISTRY_PASSWORD}
# ======================================================================================================================
