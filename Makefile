SHELL := bash # the shell used internally by Make

GOBIN ?= $(shell which go)

.PHONY: all build lint

ifeq ($(OS),Windows_NT)     # is Windows_NT on XP, 2000, 7, Vista, 10...
 detected_OS := Windows
else
 detected_OS := $(strip $(shell uname))
endif

all: build

deps: lint-install

build:
	${GOBIN} build -tags="${BUILD_TAGS}" $(BUILD_FLAGS) -o build/komainu ./cmd/komainu

vendor:
	${GOBIN} mod tidy

lint-install:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		bash -s -- -b $(shell ${GOBIN} env GOPATH)/bin v1.50.0

lint:
	@echo "lint"
	@golangci-lint --exclude=SA1019 run ./... --deadline=5m

generate:
	${GOBIN} generate ./api/generate.go
	${GOBIN} generate ./pkg/persistence/sqlcipher/sql/doc.go
	${GOBIN} generate ./pkg/persistence/sqlite/sql/doc.go

tests:
	${GOBIN} run test/*.go