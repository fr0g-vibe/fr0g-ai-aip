package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestMemoryStorage_LargeDataset(t *testing.T) {
	storage := NewMemoryStorage()

	// Create a large number of personas
	numPersonas := 1000
	for i := 0; i < numPersonas; i++ {
		p := &types.Persona{
			Name:   "Expert " + string(rune(i)),
			Topic:  "Topic " + string(rune(i)),
			Prompt: "You are expert number " + string(rune(i)),
		}

		err := storage.Create(p)
		if err != nil {
			t.Fatalf("Failed to create persona %d: %v", i, err)
		}
	}

	// Verify all personas were created
	personas, err := storage.List()
	if err != nil {
		t.Fatalf("Failed to list personas: %v", err)
	}

	if len(personas) != numPersonas {
		t.Errorf("Expected %d personas, got %d", numPersonas, len(personas))
	}
}

func TestMemoryStorage_UnicodeContent(t *testing.T) {
	storage := NewMemoryStorage()

	p := &types.Persona{
		Name:   "Unicode Expert 测试专家",
		Topic:  "Unicode Testing 测试",
		Prompt: "You are a Unicode testing expert. 你是一个Unicode测试专家。",
		Context: map[string]string{
			"language": "中文",
			"encoding": "UTF-8",
		},
		Rag: []string{
			"Unicode best practices",
			"UTF-8 编码最佳实践",
		},
	}

	err := storage.Create(p)
	if err != nil {
		t.Fatalf("Failed to create Unicode persona: %v", err)
	}

	retrieved, err := storage.Get(p.Id)
	if err != nil {
		t.Fatalf("Failed to get Unicode persona: %v", err)
	}

	if retrieved.Name != p.Name {
		t.Errorf("Unicode name not preserved: expected %s, got %s", p.Name, retrieved.Name)
	}
}

func TestFileStorage_SpecialCharactersInPath(t *testing.T) {
	// Create a directory with special characters
	tmpDir := t.TempDir()
	specialDir := filepath.Join(tmpDir, "test-dir with spaces & symbols")

	storage, err := NewFileStorage(specialDir)
	if err != nil {
		t.Fatalf("Failed to create file storage with special path: %v", err)
	}

	p := &types.Persona{
		Name:   "Special Path Test",
		Topic:  "Path Testing",
		Prompt: "Testing special characters in paths",
	}

	err = storage.Create(p)
	if err != nil {
		t.Fatalf("Failed to create persona in special path: %v", err)
	}

	// Verify file was created in the main directory
	files, err := os.ReadDir(specialDir)
	if err != nil {
		t.Fatalf("Failed to read special directory: %v", err)
	}

	// Count only JSON files (persona files)
	jsonFiles := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			jsonFiles++
		}
	}

	if jsonFiles != 1 {
		t.Errorf("Expected 1 JSON file in special directory, got %d", jsonFiles)
	}
}

func TestFileStorage_VeryLongFilenames(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)

	// Create persona with very long content that might affect filename
	longContent := strings.Repeat("a", 1000)
	p := &types.Persona{
		Name:   "Long Content Test",
		Topic:  "Long Content",
		Prompt: longContent,
	}

	err := storage.Create(p)
	if err != nil {
		t.Fatalf("Failed to create persona with long content: %v", err)
	}

	// Verify we can retrieve it
	retrieved, err := storage.Get(p.Id)
	if err != nil {
		t.Fatalf("Failed to get persona with long content: %v", err)
	}

	if retrieved.Prompt != longContent {
		t.Error("Long content not preserved correctly")
	}
}

func TestFileStorage_ReadOnlyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Make directory read-only
	err := os.Chmod(tmpDir, 0444)
	if err != nil {
		t.Skipf("Cannot make directory read-only: %v", err)
	}

	// Restore permissions for cleanup
	defer os.Chmod(tmpDir, 0755)

	storage, err := NewFileStorage(tmpDir)
	if err == nil {
		t.Error("Expected error when creating file storage in read-only directory")
		return
	}

	// If NewFileStorage failed as expected, the test passes
	// If it somehow succeeded, we should test that Create fails
	if storage != nil {
		p := &types.Persona{
			Name:   "Read Only Test",
			Topic:  "Read Only",
			Prompt: "Testing read-only directory",
		}

		err = storage.Create(p)
		if err == nil {
			t.Error("Expected error when creating in read-only directory")
		}
	}
}

func TestFileStorage_DiskSpaceSimulation(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)

	// Create a very large persona to simulate disk space issues
	largeContent := strings.Repeat("x", 10*1024*1024) // 10MB
	p := &types.Persona{
		Name:   "Large Content Test",
		Topic:  "Large Content",
		Prompt: largeContent,
	}

	// This should work on most systems, but might fail on very constrained environments
	err := storage.Create(p)
	if err != nil {
		t.Logf("Large content creation failed (expected on constrained systems): %v", err)
		return
	}

	// If creation succeeded, verify retrieval
	retrieved, err := storage.Get(p.Id)
	if err != nil {
		t.Fatalf("Failed to get large persona: %v", err)
	}

	if len(retrieved.Prompt) != len(largeContent) {
		t.Errorf("Large content size mismatch: expected %d, got %d", len(largeContent), len(retrieved.Prompt))
	}
}

func TestMemoryStorage_NilPointerHandling(t *testing.T) {
	storage := NewMemoryStorage()

	// Test with nil persona
	err := storage.Create(nil)
	if err == nil {
		t.Error("Expected error when creating nil persona")
	}

	// Test with empty persona (missing required fields)
	p := &types.Persona{} // Empty persona with required fields missing
	err = storage.Create(p)
	if err == nil {
		t.Error("Expected error when creating persona with missing required fields")
	}
}

func TestFileStorage_ConcurrentFileAccess(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewFileStorage(tmpDir)

	// Create initial persona
	p := &types.Persona{
		Name:   "Concurrent Test",
		Topic:  "Concurrency",
		Prompt: "Testing concurrent access",
	}
	storage.Create(p)

	// Test concurrent reads
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			for j := 0; j < 100; j++ {
				_, err := storage.Get(p.Id)
				if err != nil {
					t.Errorf("Concurrent read failed: %v", err)
					return
				}
			}
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
