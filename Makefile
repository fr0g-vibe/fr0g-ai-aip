.PHONY: build test clean run-server run-cli help

# Build the application
build:
	go build -o bin/fr0g-ai-aip ./cmd/fr0g-ai-aip

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run in server mode
run-server:
	go run ./cmd/fr0g-ai-aip -server

# Run CLI help
run-cli:
	go run ./cmd/fr0g-ai-aip -help

# Install dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  run-server    - Run in server mode"
	@echo "  run-cli       - Show CLI help"
	@echo "  deps          - Install/update dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  help          - Show this help"
