NAME = $(notdir $(PWD))
SOURCE := $(shell git rev-parse --show-toplevel)

VERSION = $(shell printf "%s.%s" \
		$$(git rev-list --count HEAD) \
		$$(git rev-parse --short HEAD))

BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

generate:
	$(call chdir,$(SOURCE))
	@echo "Bundling open api specs"
	swagger-cli bundle $(SOURCE)/cmd/server/contracts/api-specs.yaml --type yaml --outfile $(SOURCE)/cmd/server/contracts/api-bundled.yaml
	@echo "Generating bundled open api golang"
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	@oapi-codegen --config $(SOURCE)/cmd/server/contracts/cfg.yaml -o $(SOURCE)/go-services/$(SERVICE)/cmd/server/http $(SOURCE)/cmd/dsp-server/contracts/api-bundled.yaml


build:  $(OUTPUT)
	CGO_ENABLED=0 GOOS=linux go build -o bin/app \
		-ldflags "-X main.version=$(VERSION)" \
		-gcflags "-trimpath $(GOPATH)/src"

test: generate
	@echo :: run tests
	go test -v ./test

run:
	@echo :: start http server at port 9090
	go run ./cmd/main.go

all: generate build test run


$(OUTPUT):
	mkdir -p $(OUTPUT)