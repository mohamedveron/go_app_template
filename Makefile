NAME = $(notdir $(PWD))

VERSION = $(shell printf "%s.%s" \
		$$(git rev-list --count HEAD) \
		$$(git rev-parse --short HEAD))

BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

generate:
	@echo :: getting generator
	go get -v -d

	@echo :: generating code


build:  $(OUTPUT)
	CGO_ENABLED=0 GOOS=linux go build -o bin/app \
		-ldflags "-X main.version=$(VERSION)" \
		-gcflags "-trimpath $(GOPATH)/src"

test: generate
	@echo :: run tests
	go test -v ./test

run:
	@echo :: start http server at port 50051
	go run main.go

all: generate build test run


$(OUTPUT):
	mkdir -p $(OUTPUT)