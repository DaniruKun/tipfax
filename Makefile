.PHONY: build run clean

BIN_NAME := tipfax-server

# Build the application
build:
	go build -o bin/${BIN_NAME} cmd/server/main.go

# Run the application
run: build
	./bin/${BIN_NAME}

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy
