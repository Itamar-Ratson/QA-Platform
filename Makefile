.PHONY: build run apply destroy test clean

build:
	go build -o bin/qa-test-app cmd/main.go

run:
	go run cmd/main.go $(ARGS)

apply:
	go run cmd/main.go -apply

destroy:
	go run cmd/main.go -destroy

test:
	go test ./...

clean:
	rm -rf bin/
	go clean

deps:
	go mod tidy
	go mod download
