.PHONY: default build test install-vet vet

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
	exhaustive -checking-strategy name ./...
	ineffassign ./...
	errcheck ./...
