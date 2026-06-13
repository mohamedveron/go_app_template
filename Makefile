NAME = $(notdir $(PWD))
SOURCE := $(shell git rev-parse --show-toplevel)

VERSION = $(shell printf "%s.%s" \
		$$(git rev-list --count HEAD) \
		$$(git rev-parse --short HEAD))

BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

.PHONY: generate-api-specs
generate-api-specs:
	@ echo "Bundling open api specs version 1"
	@ redocly bundle $(SOURCE)/api/contracts/v1/api-specs.yaml --output $(SOURCE)/api/contracts/v1/app-api-bundled.yaml
	@ redocly bundle $(SOURCE)/api/contracts/v1/api-specs.yaml --output $(SOURCE)/api/contracts/v1/app-api-bundled-json.json
	@ echo "Generating open api server interface golang"
	@ go tool oapi-codegen -config $(SOURCE)/api/contracts/server-generation-cfg.yaml -o $(SOURCE)/internal/transport/http/api_server_gen/v1/app-server-interface.gen.go $(SOURCE)/api/contracts/v1/app-api-bundled.yaml

build:  $(OUTPUT)
	CGO_ENABLED=0 GOOS=linux go build -o bin/app \
		-ldflags "-X main.version=$(VERSION)" \
		-gcflags "-trimpath $(GOPATH)/src" \
		./cmd/main.go

test: generate
	@echo :: run tests
	go test -v ./test

run:
	@echo :: start http server at port 9090
	go run ./cmd/main.go

format:
	@echo "Formatting code..."
	@gofmt -s -w .

.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run

all: generate format build lint test run


$(OUTPUT):
	mkdir -p $(OUTPUT)