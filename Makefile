.PHONY: default build test install-vet vet upgrade-deps

default: build

build:
	go build ./...

test:
	go test -cover ./...

install-vet:
	go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/kisielk/errcheck@latest

vet:
	go vet ./...
	exhaustive ./...
	ineffassign ./...
	errcheck ./...

upgrade-deps:
	go get golang.org/x/tools
	go mod tidy
