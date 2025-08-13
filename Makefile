PROJECTROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
LINKERFLAGS = -X main.Version=`git describe --tags --always --dirty` -X main.BuildTimestamp=`date -u '+%Y-%m-%d_%I:%M:%S_UTC'`

all: clean build

.PHONY: clean
clean:
	@echo Running clean job...
	cd $(PROJECTROOT) && rm -rf bin/
	go fmt ./...

dep:
	@echo Running dep job...
	go mod tidy

generate:
	@echo Running generate job...
	go generate ./...

build: clean dep generate
	@echo Running build job...
	mkdir -p bin
	go build -o bin/ ./...

run:
	@echo Running run job...
	go run .
