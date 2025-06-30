package client

import (
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestLocalClient_Create(t *testing.T) {
	storage := storage.NewMemoryStorage()
	client := NewLocalClient(storage)
	
	p := &types.Persona{
		Name:   "Local Test",
		Topic:  "Local Testing",
		Prompt: "You are a local testing expert.",
	}
	
	err := client.Create(p)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	if p.ID == "" {
		t.Error("Expected ID to be generated")
	}
}

func TestLocalClient_CRUD(t *testing.T) {
	storage := storage.NewMemoryStorage()
	client := NewLocalClient(storage)
	
	// Create
	p := &types.Persona{Name: "CRUD Test", Topic: "Testing", Prompt: "Test"}
	if err := client.Create(p); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// Read
	retrieved, err := client.Get(p.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, retrieved.Name)
	}
	
	// List
	list, err := client.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 persona, got %d", len(list))
	}
	
	// Update
	retrieved.Name = "Updated"
	if err := client.Update(p.ID, retrieved); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	
	// Delete
	if err := client.Delete(p.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	
	_, err = client.Get(p.ID)
	if err == nil {
		t.Error("Expected error after delete")
	}
}

func TestLocalClient_GetNotFound(t *testing.T) {
	storage := storage.NewMemoryStorage()
	client := NewLocalClient(storage)
	
	_, err := client.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent persona")
	}
}
