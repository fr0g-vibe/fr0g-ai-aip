.PHONY: build test clean run-server run-grpc run-both run-cli proto help

# Build the application
build: proto-if-needed
	@echo "Building application..."
	go build -o bin/fr0g-ai-aip ./cmd/fr0g-ai-aip

# Build with local gRPC support (no external dependencies)
build-with-grpc: proto-if-needed build
	@echo "gRPC support built using local JSON-over-HTTP implementation"

# Generate protobuf code (force regeneration)
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

# Generate protobuf code only if files don't exist
proto-if-needed:
	@if [ ! -f "internal/grpc/pb/persona.pb.go" ] || [ ! -f "internal/grpc/pb/persona_grpc.pb.go" ]; then \
		echo "Protobuf files missing, generating..."; \
		$(MAKE) proto; \
	else \
		echo "Protobuf files already exist, skipping generation"; \
	fi

# Run tests
test: proto-if-needed
	go test ./...

# Run tests with coverage
test-coverage: proto-if-needed
	go test -cover ./...

# Run tests with detailed coverage report
test-coverage-detailed: proto-if-needed
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with verbose output
test-verbose: proto-if-needed
	go test -v ./...

# Run tests for specific package
test-storage: proto-if-needed
	go test -v ./internal/storage/...

test-persona: proto-if-needed
	go test -v ./internal/persona/...

test-client: proto-if-needed
	go test -v ./internal/client/...

test-grpc: proto-if-needed
	go test -v ./internal/grpc/...

test-main: proto-if-needed
	go test -v ./cmd/fr0g-ai-aip/...

# Run tests with race detection
test-race: proto-if-needed
	go test -race ./...

# Run benchmarks
test-bench: proto-if-needed
	go test -bench=. ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf internal/grpc/pb/*.pb.go

# Run HTTP REST API server
run-server: proto-if-needed
	go run ./cmd/fr0g-ai-aip -server

# Run gRPC server
run-grpc: proto-if-needed
	go run ./cmd/fr0g-ai-aip -grpc

# Run both HTTP and gRPC servers
run-both: proto-if-needed
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
	@echo "  proto              - Force generate protobuf code"
	@echo "  proto-if-needed    - Generate protobuf code only if missing"
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
