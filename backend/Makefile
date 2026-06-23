.PHONY: dev test lint build

dev:
	go run main.go

test:
	go test ./...

lint:
	go vet ./...

build:
	go build -o bin/server main.go
