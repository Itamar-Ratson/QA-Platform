.PHONY: build run test clean

build:
	go build -o bin/qa-test-app cmd/main.go

run:
	go run cmd/main.go

test:
	go test ./...

clean:
	rm -rf bin/
	go clean

deps:
	go mod tidy
	go mod download
