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
- API-driven architecture for integration flexibility

## Setup

```bash
# Install protobuf tools
make install-proto-tools

# Ensure GOPATH/bin is in your PATH
export PATH="$(go env GOPATH)/bin:$PATH"

# Generate protobuf files and build
make build
```

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

## Contributing

Please follow Go best practices for code documentation and CLI design.
