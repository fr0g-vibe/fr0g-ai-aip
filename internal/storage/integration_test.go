package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestStorageIntegration(t *testing.T) {
	// Test both memory and file storage with the same operations
	storages := map[string]Storage{
		"memory": NewMemoryStorage(),
	}
	
	// Add file storage
	tmpDir := t.TempDir()
	fileStorage, err := NewFileStorage(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}
	storages["file"] = fileStorage
	
	for name, storage := range storages {
		t.Run(name, func(t *testing.T) {
			testStorageOperations(t, storage)
		})
	}
}

func testStorageOperations(t *testing.T, storage Storage) {
	// Test complete CRUD workflow
	personas := []*types.Persona{
		{Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"},
		{Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"},
		{Name: "Expert 3", Topic: "Topic 3", Prompt: "Prompt 3"},
	}
	
	// Create multiple personas
	for _, p := range personas {
		err := storage.Create(p)
		if err != nil {
			t.Fatalf("Failed to create persona %s: %v", p.Name, err)
		}
		if p.ID == "" {
			t.Errorf("Expected ID to be generated for persona %s", p.Name)
		}
	}
	
	// List and verify count
	list, err := storage.List()
	if err != nil {
		t.Fatalf("Failed to list personas: %v", err)
	}
	if len(list) != len(personas) {
		t.Errorf("Expected %d personas, got %d", len(personas), len(list))
	}
	
	// Get each persona and verify
	for _, p := range personas {
		retrieved, err := storage.Get(p.ID)
		if err != nil {
			t.Fatalf("Failed to get persona %s: %v", p.ID, err)
		}
		if retrieved.Name != p.Name {
			t.Errorf("Expected name %s, got %s", p.Name, retrieved.Name)
		}
	}
	
	// Update first persona
	personas[0].Name = "Updated Expert 1"
	err = storage.Update(personas[0].ID, *personas[0])
	if err != nil {
		t.Fatalf("Failed to update persona: %v", err)
	}
	
	// Verify update
	updated, err := storage.Get(personas[0].ID)
	if err != nil {
		t.Fatalf("Failed to get updated persona: %v", err)
	}
	if updated.Name != "Updated Expert 1" {
		t.Errorf("Expected updated name 'Updated Expert 1', got %s", updated.Name)
	}
	
	// Delete second persona
	err = storage.Delete(personas[1].ID)
	if err != nil {
		t.Fatalf("Failed to delete persona: %v", err)
	}
	
	// Verify deletion
	_, err = storage.Get(personas[1].ID)
	if err == nil {
		t.Error("Expected error when getting deleted persona")
	}
	
	// Verify final count
	finalList, err := storage.List()
	if err != nil {
		t.Fatalf("Failed to list final personas: %v", err)
	}
	if len(finalList) != len(personas)-1 {
		t.Errorf("Expected %d personas after deletion, got %d", len(personas)-1, len(finalList))
	}
}

func TestFileStorageCorruption(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)
	
	// Create a valid persona
	p := &types.Persona{Name: "Valid", Topic: "Valid", Prompt: "Valid"}
	storage.Create(p)
	
	// Create a corrupted file manually
	corruptedPath := filepath.Join(tmpDir, "corrupted.json")
	os.WriteFile(corruptedPath, []byte("{invalid json"), 0644)
	
	// Create a non-JSON file
	nonJSONPath := filepath.Join(tmpDir, "notjson.txt")
	os.WriteFile(nonJSONPath, []byte("not json"), 0644)
	
	// List should still work and return only valid personas
	list, err := storage.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	
	// Should only have the one valid persona
	if len(list) != 1 {
		t.Errorf("Expected 1 valid persona, got %d", len(list))
	}
	
	if list[0].Name != "Valid" {
		t.Errorf("Expected valid persona name 'Valid', got %s", list[0].Name)
	}
}

func TestMemoryStorageConcurrency(t *testing.T) {
	storage := NewMemoryStorage()
	
	// Test concurrent operations
	numGoroutines := 50
	numOperations := 10
	
	done := make(chan bool, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			for j := 0; j < numOperations; j++ {
				p := &types.Persona{
					Name:   fmt.Sprintf("Concurrent %d-%d", id, j),
					Topic:  "Concurrency",
					Prompt: "Concurrent test",
				}
				
				// Create
				if err := storage.Create(p); err != nil {
					t.Errorf("Concurrent create failed: %v", err)
					return
				}
				
				// Read
				if _, err := storage.Get(p.ID); err != nil {
					t.Errorf("Concurrent get failed: %v", err)
					return
				}
				
				// Update
				p.Name = fmt.Sprintf("Updated %d-%d", id, j)
				if err := storage.Update(p.ID, *p); err != nil {
					t.Errorf("Concurrent update failed: %v", err)
					return
				}
			}
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// Verify final state
	list, err := storage.List()
	if err != nil {
		t.Fatalf("Final list failed: %v", err)
	}
	
	expectedCount := numGoroutines * numOperations
	if len(list) != expectedCount {
		t.Errorf("Expected %d personas after concurrent operations, got %d", expectedCount, len(list))
	}
}
