.PHONY: build test clean run-server run-grpc run-both run-cli proto help

# Build the application
build: proto
	@echo "Building application..."
	go build -o bin/fr0g-ai-aip ./cmd/fr0g-ai-aip

# Build with local gRPC support (no external dependencies)
build-with-grpc: proto build
	@echo "gRPC support built using local JSON-over-HTTP implementation"

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@if [ ! -f "internal/grpc/proto/persona.proto" ]; then \
		echo "Error: internal/grpc/proto/persona.proto not found"; \
		echo "Please create the proto file first"; \
		exit 1; \
	fi
	@echo "Using protoc: $(shell which protoc)"
	@echo "Using protoc-gen-go: $(shell which protoc-gen-go || echo "$(shell go env GOPATH)/bin/protoc-gen-go")"
	@echo "Using protoc-gen-go-grpc: $(shell which protoc-gen-go-grpc || echo "$(shell go env GOPATH)/bin/protoc-gen-go-grpc")"
	@mkdir -p internal/grpc/pb
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc \
		--proto_path=internal/grpc/proto \
		--go_out=internal/grpc/pb --go_opt=paths=source_relative \
		--go-grpc_out=internal/grpc/pb --go-grpc_opt=paths=source_relative \
		internal/grpc/proto/persona.proto
	@echo "Protobuf code generated successfully in internal/grpc/pb/"

# Run tests
test: proto
	go test ./...

# Run tests with coverage
test-coverage: proto
	go test -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf internal/grpc/pb/

# Run HTTP REST API server
run-server: proto
	go run ./cmd/fr0g-ai-aip -server

# Run gRPC server
run-grpc: proto
	go run ./cmd/fr0g-ai-aip -grpc

# Run both HTTP and gRPC servers
run-both: proto
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
