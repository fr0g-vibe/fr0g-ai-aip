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

# Run tests with verbose coverage for specific package
test-coverage-verbose-grpc: proto-if-needed
	go test -v -coverprofile=grpc_coverage.out ./internal/grpc/
	go tool cover -func=grpc_coverage.out
	go tool cover -html=grpc_coverage.out -o grpc_coverage.html
	@echo "gRPC coverage report generated: grpc_coverage.html"

# Run tests with verbose output
test-verbose: proto-if-needed
	go test -v ./...

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

# Generate documentation
docs: proto-if-needed
	@echo "Generating documentation..."
	@mkdir -p docs/generated
	go doc -all ./... > docs/generated/godoc.txt
	@echo "Documentation generated in docs/generated/"

# Serve documentation locally
docs-serve: docs
	@echo "Starting documentation server on http://localhost:6060"
	@echo "Visit http://localhost:6060/pkg/github.com/fr0g-vibe/fr0g-ai-aip/ for package docs"
	godoc -http=:6060

# Generate OpenAPI documentation
docs-openapi:
	@echo "OpenAPI specification available at docs/OPENAPI_SPEC.yaml"
	@echo "View with: swagger-ui-serve docs/OPENAPI_SPEC.yaml"

# Generate MCP documentation
docs-mcp:
	@echo "MCP integration documentation available at docs/MCP_INTEGRATION.md"

# Show help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Building:"
	@echo "  build              - Build the application (no external deps)"
	@echo "  build-with-grpc    - Build with full gRPC support"
	@echo "  clean              - Clean build artifacts"
	@echo ""
	@echo "Protocol Buffers:"
	@echo "  proto              - Force generate protobuf code"
	@echo "  proto-if-needed    - Generate protobuf code only if missing"
	@echo ""
	@echo "Testing:"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  test-coverage-detailed - Generate HTML coverage report"
	@echo "  test-verbose       - Run tests with verbose output"
	@echo "  test-race          - Run tests with race detection"
	@echo "  test-bench         - Run benchmarks"
	@echo ""
	@echo "Running:"
	@echo "  run-server         - Run HTTP REST API server"
	@echo "  run-grpc           - Run gRPC server"
	@echo "  run-both           - Run both HTTP and gRPC servers"
	@echo "  run-cli            - Show CLI help"
	@echo ""
	@echo "Documentation:"
	@echo "  docs               - Generate all documentation"
	@echo "  docs-serve         - Serve documentation locally (port 6060)"
	@echo "  docs-openapi       - Show OpenAPI documentation info"
	@echo "  docs-mcp           - Show MCP integration documentation info"
	@echo ""
	@echo "Development:"
	@echo "  deps               - Install/update dependencies"
	@echo "  install-proto-tools - Install protobuf generation tools"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code (requires golangci-lint)"
	@echo ""
	@echo "Help:"
	@echo "  help               - Show this help"
