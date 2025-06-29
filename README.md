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
- No external dependencies for core functionality
- API-driven architecture for integration flexibility

## Documentation

- General project documentation: This README.md
- CLI documentation: Generated using Go best practices
- API documentation: Embedded in Go code using standard conventions

## Getting Started

```bash
# Build the project
go build ./cmd/fr0g-ai-aip

# Run the CLI help
./fr0g-ai-aip --help
```

## Contributing

Please follow Go best practices for code documentation and CLI design.
