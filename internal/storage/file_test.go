package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestFileStorage_New(t *testing.T) {
	tmpDir := t.TempDir()
	
	storage, err := NewFileStorage(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}
	
	if storage == nil {
		t.Error("Expected storage instance")
	}
	
	// Check directory was created
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Expected data directory to be created")
	}
}

func TestFileStorage_CreateAndPersist(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)
	
	p := &types.Persona{
		Name:   "File Test",
		Topic:  "File Testing",
		Prompt: "Test prompt",
	}
	
	err := storage.Create(p)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// Check file exists
	filePath := filepath.Join(tmpDir, p.ID+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Expected persona file to be created")
	}
	
	// Create new storage instance to test persistence
	storage2, _ := NewFileStorage(tmpDir)
	retrieved, err := storage2.Get(p.ID)
	if err != nil {
		t.Fatalf("Failed to get persisted persona: %v", err)
	}
	
	if retrieved.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, retrieved.Name)
	}
}

func TestFileStorage_List(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)
	
	p1 := &types.Persona{Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"}
	p2 := &types.Persona{Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"}
	
	storage.Create(p1)
	storage.Create(p2)
	
	list, err := storage.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	
	if len(list) != 2 {
		t.Errorf("Expected 2 personas, got %d", len(list))
	}
}

func TestFileStorage_DeleteFile(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)
	
	p := &types.Persona{Name: "Delete Test", Topic: "Deleting", Prompt: "Test"}
	storage.Create(p)
	
	filePath := filepath.Join(tmpDir, p.ID+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("File should exist before delete")
	}
	
	err := storage.Delete(p.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("File should not exist after delete")
	}
}

func TestFileStorage_CreateValidation(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)
	
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

func TestFileStorage_GetNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)
	
	_, err := storage.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent persona")
	}
}

func TestFileStorage_UpdateNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)
	
	p := types.Persona{
		Name:   "Test",
		Topic:  "Test",
		Prompt: "Test",
	}
	
	err := storage.Update("nonexistent", p)
	if err == nil {
		t.Error("Expected error for nonexistent persona")
	}
}

func TestFileStorage_DeleteNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)
	
	err := storage.Delete("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent persona")
	}
}

func TestFileStorage_ListWithCorruptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)
	
	// Create a corrupted JSON file
	corruptedFile := filepath.Join(tmpDir, "corrupted.json")
	os.WriteFile(corruptedFile, []byte("invalid json"), 0644)
	
	// List should still work, just skip the corrupted file
	personas, err := storage.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	
	if len(personas) != 0 {
		t.Errorf("Expected 0 personas (corrupted file should be skipped), got %d", len(personas))
	}
}
