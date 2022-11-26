.PHONY: default
default: build

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: install-vet
install-vet:
	go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/kisielk/errcheck@latest

.PHONY: vet
vet:
	go vet ./...
	exhaustive ./...
	ineffassign ./...
	errcheck ./...

.PHONY: upgrade-deps
upgrade-deps:
	go get golang.org/x/tools
	go mod tidy
