.PHONY: build test lint clean plugin

build:
	go build -o bin/logcheck ./cmd/logcheck

test:
	go test ./... -v -count=1 -race

vet:
	go vet ./...

plugin:
	go build -buildmode=plugin -o logcheck.so ./plugin/