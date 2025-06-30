package storage

import (
	"fmt"
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestMemoryStorage_Create(t *testing.T) {
	storage := NewMemoryStorage()
	
	p := &types.Persona{
		Name:   "Test Expert",
		Topic:  "Testing",
		Prompt: "You are a testing expert.",
		Context: map[string]string{
			"domain": "testing",
		},
		RAG: []string{"test-doc1", "test-doc2"},
	}
	
	err := storage.Create(p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}
	
	if p.ID == "" {
		t.Error("Expected persona ID to be generated")
	}
	
	// Verify the persona was stored correctly
	retrieved, err := storage.Get(p.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve created persona: %v", err)
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
}

func TestMemoryStorage_CreateValidation(t *testing.T) {
	storage := NewMemoryStorage()
	
	tests := []struct {
		name    string
		persona *types.Persona
		wantErr bool
	}{
		{"nil persona", nil, true},
		{"missing name", &types.Persona{Topic: "Test", Prompt: "Test"}, true},
		{"empty name", &types.Persona{Name: "", Topic: "Test", Prompt: "Test"}, true},
		{"missing topic", &types.Persona{Name: "Test", Prompt: "Test"}, true},
		{"empty topic", &types.Persona{Name: "Test", Topic: "", Prompt: "Test"}, true},
		{"missing prompt", &types.Persona{Name: "Test", Topic: "Test"}, true},
		{"empty prompt", &types.Persona{Name: "Test", Topic: "Test", Prompt: ""}, true},
		{"valid persona", &types.Persona{Name: "Test", Topic: "Test", Prompt: "Test"}, false},
		{"valid with context", &types.Persona{Name: "Test", Topic: "Test", Prompt: "Test", Context: map[string]string{"key": "value"}}, false},
		{"valid with RAG", &types.Persona{Name: "Test", Topic: "Test", Prompt: "Test", RAG: []string{"doc1"}}, false},
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

func BenchmarkMemoryStorage_Create(b *testing.B) {
	storage := NewMemoryStorage()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := &types.Persona{
			Name:   "Benchmark Test",
			Topic:  "Benchmarking",
			Prompt: "You are a benchmarking expert.",
		}
		storage.Create(p)
	}
}

func BenchmarkMemoryStorage_Get(b *testing.B) {
	storage := NewMemoryStorage()
	
	// Create test data
	personas := make([]*types.Persona, 1000)
	for i := 0; i < 1000; i++ {
		p := &types.Persona{
			Name:   "Test Expert",
			Topic:  "Testing",
			Prompt: "You are a testing expert.",
		}
		storage.Create(p)
		personas[i] = p
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.Get(personas[i%1000].ID)
	}
}

func TestMemoryStorage_ConcurrentAccess(t *testing.T) {
	storage := NewMemoryStorage()
	
	// Test concurrent writes
	numGoroutines := 20
	done := make(chan bool, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			p := &types.Persona{
				Name:   fmt.Sprintf("Concurrent Test %d", id),
				Topic:  "Concurrency",
				Prompt: "You are a concurrency expert.",
			}
			
			// Create
			err := storage.Create(p)
			if err != nil {
				t.Errorf("Concurrent create failed: %v", err)
				return
			}
			
			// Read
			_, err = storage.Get(p.ID)
			if err != nil {
				t.Errorf("Concurrent get failed: %v", err)
				return
			}
			
			// Update
			p.Name = fmt.Sprintf("Updated Concurrent Test %d", id)
			err = storage.Update(p.ID, *p)
			if err != nil {
				t.Errorf("Concurrent update failed: %v", err)
				return
			}
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// Verify all personas were created
	personas, err := storage.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(personas) != numGoroutines {
		t.Errorf("Expected %d personas, got %d", numGoroutines, len(personas))
	}
}
