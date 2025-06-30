package client

import (
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// Test that all client implementations satisfy the Client interface
func TestClientInterface(t *testing.T) {
	var _ Client = &LocalClient{}
	var _ Client = &RESTClient{}
	var _ Client = &GRPCClient{}
}

func TestClientImplementations(t *testing.T) {
	// Test that we can create all client types
	clients := map[string]Client{
		"local": NewLocalClient(storage.NewMemoryStorage()),
		"rest":  NewRESTClient("http://localhost:8080"),
	}
	
	// Add gRPC client
	grpcClient, err := NewGRPCClient("localhost:9090")
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	clients["grpc"] = grpcClient
	
	for name, client := range clients {
		t.Run(name, func(t *testing.T) {
			if client == nil {
				t.Errorf("Client %s should not be nil", name)
			}
		})
	}
}

func TestClientCRUDInterface(t *testing.T) {
	// Test that all clients implement the same interface methods
	storage := storage.NewMemoryStorage()
	client := NewLocalClient(storage)
	
	p := &types.Persona{
		Name:   "Interface Test",
		Topic:  "Testing",
		Prompt: "You are a testing expert.",
	}
	
	// Test Create method signature
	err := client.Create(p)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// Test Get method signature
	_, err = client.Get(p.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	
	// Test List method signature
	_, err = client.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	
	// Test Update method signature
	err = client.Update(p.ID, *p)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	
	// Test Delete method signature
	err = client.Delete(p.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}
