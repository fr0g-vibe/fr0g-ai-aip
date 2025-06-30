package client

import (
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestNewGRPCClient(t *testing.T) {
	client, err := NewGRPCClient("localhost:9090")
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	
	if client == nil {
		t.Error("Expected client to be created")
	}
	
	// Close the connection
	client.Close()
}

func TestGRPCClient_Interface(t *testing.T) {
	// Test that GRPCClient implements the Client interface
	var _ Client = &GRPCClient{}
}

func TestGRPCClient_Methods(t *testing.T) {
	// Test that all required methods exist (will fail at runtime if no server)
	client, err := NewGRPCClient("localhost:9090")
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer client.Close()
	
	// Test method signatures exist
	p := &types.Persona{
		Name:   "Test",
		Topic:  "Testing", 
		Prompt: "Test prompt",
	}
	
	// These will fail at runtime without a server, but we're just testing the interface
	_ = client.Create(p)
	_, _ = client.Get("test-id")
	_, _ = client.List()
	_ = client.Update("test-id", *p)
	_ = client.Delete("test-id")
}
