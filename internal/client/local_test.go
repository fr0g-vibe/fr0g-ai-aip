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
	p := &types.Persona{
		Name:   "CRUD Test",
		Topic:  "Testing",
		Prompt: "Test prompt",
		Context: map[string]string{
			"test": "value",
		},
		RAG: []string{"doc1", "doc2"},
	}
	if err := client.Create(p); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	if p.ID == "" {
		t.Error("Expected ID to be generated")
	}
	
	// Read
	retrieved, err := client.Get(p.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, retrieved.Name)
	}
	if len(retrieved.Context) != len(p.Context) {
		t.Errorf("Expected context length %d, got %d", len(p.Context), len(retrieved.Context))
	}
	if len(retrieved.RAG) != len(p.RAG) {
		t.Errorf("Expected RAG length %d, got %d", len(p.RAG), len(retrieved.RAG))
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
	retrieved.Name = "Updated CRUD Test"
	retrieved.Context["updated"] = "true"
	if err := client.Update(p.ID, retrieved); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	
	// Verify update
	updated, err := client.Get(p.ID)
	if err != nil {
		t.Fatalf("Get after update failed: %v", err)
	}
	if updated.Name != "Updated CRUD Test" {
		t.Errorf("Expected updated name 'Updated CRUD Test', got %s", updated.Name)
	}
	if updated.Context["updated"] != "true" {
		t.Error("Expected context to be updated")
	}
	
	// Delete
	if err := client.Delete(p.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	
	_, err = client.Get(p.ID)
	if err == nil {
		t.Error("Expected error after delete")
	}
	
	// Verify list is empty after delete
	finalList, err := client.List()
	if err != nil {
		t.Fatalf("Final list failed: %v", err)
	}
	if len(finalList) != 0 {
		t.Errorf("Expected 0 personas after delete, got %d", len(finalList))
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
