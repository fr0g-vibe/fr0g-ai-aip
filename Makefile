.PHONY: build test clean run-server run-grpc run-both run-cli proto help

# Build the application
build:
	@echo "Building application..."
	go build -o bin/fr0g-ai-aip ./cmd/fr0g-ai-aip

# Generate protobuf code
proto:
	@echo "Generating protobuf files..."
	@mkdir -p internal/grpc/pb
	@if command -v protoc >/dev/null 2>&1; then \
		if PATH="$(shell go env GOPATH)/bin:$(PATH)" command -v protoc-gen-go >/dev/null 2>&1 && PATH="$(shell go env GOPATH)/bin:$(PATH)" command -v protoc-gen-go-grpc >/dev/null 2>&1; then \
			PATH="$(shell go env GOPATH)/bin:$(PATH)" protoc --go_out=internal/grpc/pb --go_opt=paths=source_relative \
				--go-grpc_out=internal/grpc/pb --go-grpc_opt=paths=source_relative \
				proto/persona.proto && \
			echo "Protobuf generation complete" && \
			echo "Updating imports..." && \
			sed -i 's|package grpc|package pb|g' internal/grpc/pb/persona.pb.go internal/grpc/pb/persona_grpc.pb.go 2>/dev/null || true && \
			echo "Generated files:" && \
			ls -la internal/grpc/pb/; \
		else \
			echo "protoc-gen-go or protoc-gen-go-grpc not found. Run 'make install-proto-tools' first."; \
			exit 1; \
		fi; \
	else \
		echo "protoc not found. Install Protocol Buffers compiler first."; \
		echo "On Ubuntu/Debian: sudo apt install protobuf-compiler"; \
		echo "On macOS: brew install protobuf"; \
		exit 1; \
	fi

# Build with protobuf support
build-with-grpc: proto
	@echo "Building application with gRPC support..."
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
	rm -rf internal/grpc/pb/

# Run HTTP REST API server
run-server:
	go run ./cmd/fr0g-ai-aip -server

# Run gRPC server
run-grpc:
	go run ./cmd/fr0g-ai-aip -grpc

# Run both HTTP and gRPC servers
run-both:
	go run ./cmd/fr0g-ai-aip -server -grpc

# Run CLI help
run-cli:
	go run ./cmd/fr0g-ai-aip -help

# Install dependencies
deps:
	go mod tidy

# Install protobuf tools (optional - only needed for full gRPC implementation)
install-proto-tools:
	@echo "Installing protobuf tools (this will add external dependencies)..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Make sure $(shell go env GOPATH)/bin is in your PATH"
	@echo "Note: You'll also need to add gRPC dependencies to go.mod for full gRPC support"

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build the application (no external deps)"
	@echo "  build-with-grpc    - Build with full gRPC support"
	@echo "  proto              - Generate protobuf code (optional)"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  clean              - Clean build artifacts"
	@echo "  run-server         - Run HTTP REST API server"
	@echo "  run-grpc           - Run gRPC server placeholder"
	@echo "  run-both           - Run both HTTP and gRPC servers"
	@echo "  run-cli            - Show CLI help"
	@echo "  deps               - Install/update dependencies"
	@echo "  install-proto-tools - Install protobuf generation tools"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo "  help               - Show this help"
