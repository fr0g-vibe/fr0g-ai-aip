#!/bin/bash

# Setup script for fr0g-ai-aip

echo "Setting up fr0g-ai-aip development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed. Please install Protocol Buffers compiler."
    echo "On Ubuntu/Debian: sudo apt-get install protobuf-compiler"
    echo "On macOS: brew install protobuf"
    echo "On Arch Linux: sudo pacman -S protobuf"
    exit 1
fi

# Install Go protobuf tools
echo "Installing Go protobuf tools..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Add GOPATH/bin to PATH if not already there
GOPATH_BIN=$(go env GOPATH)/bin
if [[ ":$PATH:" != *":$GOPATH_BIN:"* ]]; then
    echo "Adding $GOPATH_BIN to PATH for this session..."
    export PATH="$GOPATH_BIN:$PATH"
    echo "Note: Add 'export PATH=\"\$(go env GOPATH)/bin:\$PATH\"' to your shell profile for permanent effect"
fi

# Generate protobuf files
echo "Generating protobuf files..."
make proto

# Build the application
echo "Building application..."
make build

echo "Setup complete! You can now run:"
echo "  ./bin/fr0g-ai-aip -help"
echo "  ./bin/fr0g-ai-aip -server"
echo "  ./bin/fr0g-ai-aip -grpc"
