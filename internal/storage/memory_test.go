package storage

import (
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestMemoryStorage_Create(t *testing.T) {
	storage := NewMemoryStorage()
	
	p := &types.Persona{
		Name:   "Test Expert",
		Topic:  "Testing",
		Prompt: "You are a testing expert.",
	}
	
	err := storage.Create(p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}
	
	if p.ID == "" {
		t.Error("Expected persona ID to be generated")
	}
}

func TestMemoryStorage_CreateValidation(t *testing.T) {
	storage := NewMemoryStorage()
	
	tests := []struct {
		name    string
		persona *types.Persona
		wantErr bool
	}{
		{"missing name", &types.Persona{Topic: "Test", Prompt: "Test"}, true},
		{"missing topic", &types.Persona{Name: "Test", Prompt: "Test"}, true},
		{"missing prompt", &types.Persona{Name: "Test", Topic: "Test"}, true},
		{"valid persona", &types.Persona{Name: "Test", Topic: "Test", Prompt: "Test"}, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Create(tt.persona)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryStorage_GetNotFound(t *testing.T) {
	storage := NewMemoryStorage()
	
	_, err := storage.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent persona")
	}
}

func TestMemoryStorage_CRUD(t *testing.T) {
	storage := NewMemoryStorage()
	
	// Create
	p := &types.Persona{Name: "CRUD Test", Topic: "Testing", Prompt: "Test prompt"}
	if err := storage.Create(p); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// Read
	retrieved, err := storage.Get(p.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, retrieved.Name)
	}
	
	// Update
	retrieved.Name = "Updated Name"
	if err := storage.Update(p.ID, retrieved); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	
	updated, _ := storage.Get(p.ID)
	if updated.Name != "Updated Name" {
		t.Errorf("Expected updated name, got %s", updated.Name)
	}
	
	// Delete
	if err := storage.Delete(p.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	
	_, err = storage.Get(p.ID)
	if err == nil {
		t.Error("Expected error after delete")
	}
}
