# fr0g-ai-aip

AI Personas - A customizable AI subject matter expert system

## Overview

fr0g-ai-aip is a collection of customizable AI "personas" that function as on-demand subject matter experts. Each persona consists of a chatbot system prompt with accompanying RAG (Retrieval-Augmented Generation) and context for a specific AI identity or domain expertise.

## Purpose

This system provides specialized AI personas that can be instantiated as subject matter experts on specific topics or perspectives. These personas are designed to be used via MCP (Model Context Protocol) to provide knowledge and perspective when making decisions or taking actions.

## Architecture

- **API-based**: RESTful API interface for persona management and interaction
- **Golang**: Written entirely in Go for performance and reliability
- **CLI-first**: All management operations through Go CLI tools
- **No GUI**: Headless operation, no web UI or graphical interfaces
- **MCP Integration**: Designed for use with Model Context Protocol

## Technical Requirements

- Go 1.21 or higher
- Protocol Buffers compiler (protoc) for gRPC functionality
- gRPC and protobuf dependencies (automatically managed)

## Setup

```bash
# Install protobuf tools
make install-proto-tools
export PATH="$(go env GOPATH)/bin:$PATH"

# Build with full gRPC support
make build-with-grpc

# Or basic build (HTTP REST only)
make build
```

## Dependency Philosophy

This project supports multiple client/server modes:
- **Core functionality**: Uses Go standard library where possible
- **HTTP REST API**: Built with `net/http` (standard library)
- **gRPC API**: Uses google.golang.org/grpc and google.golang.org/protobuf
- **JSON handling**: Built with `encoding/json` (standard library)
- **File storage**: Built with `os` and `path/filepath` (standard library)

## Documentation

- General project documentation: This README.md
- CLI documentation: Generated using Go best practices
- API documentation: Embedded in Go code using standard conventions

## Getting Started

```bash
# Build the project
make build

# Run the CLI help
./bin/fr0g-ai-aip -help

# CLI with local storage (default: in-memory)
./bin/fr0g-ai-aip create -name "Go Expert" -topic "Golang Programming" -prompt "You are an expert Go programmer with deep knowledge of best practices, performance optimization, and modern Go development."

# CLI with file storage
FR0G_STORAGE_TYPE=file FR0G_DATA_DIR=./personas ./bin/fr0g-ai-aip create -name "Security Expert" -topic "Cybersecurity" -prompt "You are a cybersecurity expert."

# CLI using REST API client (requires server running)
FR0G_CLIENT_TYPE=rest FR0G_SERVER_URL=http://localhost:8080 ./bin/fr0g-ai-aip list

# CLI using gRPC client (requires gRPC server running)
FR0G_CLIENT_TYPE=grpc FR0G_SERVER_URL=localhost:9090 ./bin/fr0g-ai-aip list

# List all personas
./bin/fr0g-ai-aip list

# Get a specific persona
./bin/fr0g-ai-aip get <persona-id>

# Update a persona
./bin/fr0g-ai-aip update <persona-id> -name "Updated Name" -topic "Updated Topic"

# Delete a persona
./bin/fr0g-ai-aip delete <persona-id>

# Start HTTP REST API server with in-memory storage
./bin/fr0g-ai-aip -server

# Start HTTP REST API server with file storage
./bin/fr0g-ai-aip -server -storage file -data-dir ./server-data

# Start HTTP REST API server on custom port
./bin/fr0g-ai-aip -server -port 9090

# Start gRPC server
./bin/fr0g-ai-aip -grpc

# CLI using gRPC client
FR0G_CLIENT_TYPE=grpc FR0G_SERVER_URL=localhost:9090 ./bin/fr0g-ai-aip list

# Start both HTTP and gRPC servers
./bin/fr0g-ai-aip -server -grpc -storage file -data-dir ./shared-data
```

## Configuration

The CLI can be configured via environment variables:

- `FR0G_CLIENT_TYPE`: Client type (`local`, `rest`, `grpc`) - default: `local`
- `FR0G_STORAGE_TYPE`: Storage type (`memory`, `file`) - default: `memory` (only for local client)
- `FR0G_DATA_DIR`: Data directory for file storage - default: `./data`
- `FR0G_SERVER_URL`: Server URL for REST client - default: `http://localhost:8080`

Server mode supports command-line flags:

- `-storage`: Storage type (`memory`, `file`) - default: `memory`
- `-data-dir`: Data directory for file storage - default: `./data`
- `-port`: HTTP server port - default: `8080`
- `-grpc-port`: gRPC server port - default: `9090`

## API Usage

### HTTP REST API

```bash
# Health check
curl http://localhost:8080/health

# List all personas
curl http://localhost:8080/personas

# Create a persona
curl -X POST http://localhost:8080/personas \
  -H "Content-Type: application/json" \
  -d '{"name":"Security Expert","topic":"Cybersecurity","prompt":"You are a cybersecurity expert with extensive knowledge of threat analysis, security best practices, and incident response."}'

# Get a specific persona
curl http://localhost:8080/personas/<persona-id>

# Delete a persona
curl -X DELETE http://localhost:8080/personas/<persona-id>
```

### gRPC API

The gRPC service runs on port 9090 by default and provides the same functionality as the REST API with better performance and type safety.

## Testing

The project maintains comprehensive test coverage across all packages:

- **API**: 100% coverage
- **gRPC**: 96.9% coverage  
- **Storage**: 94.1% coverage
- **Client**: 90.7% coverage
- **Persona**: 82.6% coverage
- **CLI**: 85.6% coverage

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Generate detailed HTML coverage report
make test-coverage-detailed

# Run gRPC-specific coverage with detailed report
make test-coverage-verbose-grpc

# Run tests with verbose output
make test-verbose

# Run tests with race detection
make test-race

# Run benchmarks
make test-bench
```

### Test Coverage

The test suite includes:

- **Unit tests** for all core functionality
- **Integration tests** across storage implementations
- **Concurrent operation tests** for thread safety
- **Error handling tests** for edge cases
- **Network failure simulation** for REST client
- **gRPC server validation** with comprehensive scenarios
- **JSON marshaling/unmarshaling** with special characters
- **File corruption handling** for file storage
- **Performance benchmarks** for critical paths

### Test Organization

- `*_test.go` files contain unit tests for each package
- `integration_test.go` files test cross-package functionality
- Mock servers and storage implementations for isolated testing
- Comprehensive validation of error conditions and edge cases

## Development

### Available Make Targets

```bash
# Building
make build              # Build the application (no external deps)
make build-with-grpc    # Build with full gRPC support
make clean              # Clean build artifacts

# Protocol Buffers
make proto              # Force generate protobuf code
make proto-if-needed    # Generate protobuf code only if missing

# Testing
make test               # Run tests
make test-coverage      # Run tests with coverage
make test-coverage-detailed # Generate HTML coverage report
make test-coverage-verbose-grpc # Detailed gRPC coverage
make test-verbose       # Run tests with verbose output
make test-race          # Run tests with race detection
make test-bench         # Run benchmarks

# Running
make run-server         # Run HTTP REST API server
make run-grpc           # Run gRPC server
make run-both           # Run both HTTP and gRPC servers
make run-cli            # Show CLI help

# Dependencies and Tools
make deps               # Install/update dependencies
make install-proto-tools # Install protobuf generation tools
make fmt                # Format code

# Help
make help               # Show all available targets
```

### Code Quality

```bash
# Format code
make fmt

# Install dependencies
make deps

# Note: Linting requires external tools not included in this Makefile
# You can install golangci-lint separately if needed
```

## Contributing

Please follow Go best practices for code documentation and CLI design.
