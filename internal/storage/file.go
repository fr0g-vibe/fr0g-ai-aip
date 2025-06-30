package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// FileStorage implements file-based storage for personas
type FileStorage struct {
	dataDir string
	mu      sync.RWMutex
}

// NewFileStorage creates a new file storage instance
func NewFileStorage(dataDir string) (*FileStorage, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}
	
	return &FileStorage{
		dataDir: dataDir,
	}, nil
}

func (f *FileStorage) Create(p *types.Persona) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	if p == nil {
		return fmt.Errorf("persona cannot be nil")
	}
	if p.Name == "" {
		return fmt.Errorf("persona name is required")
	}
	if p.Topic == "" {
		return fmt.Errorf("persona topic is required")
	}
	if p.Prompt == "" {
		return fmt.Errorf("persona prompt is required")
	}
	
	p.ID = f.generateID()
	return f.writePersona(*p)
}

func (f *FileStorage) Get(id string) (types.Persona, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	return f.readPersona(id)
}

func (f *FileStorage) List() ([]types.Persona, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	files, err := os.ReadDir(f.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %v", err)
	}
	
	var personas []types.Persona
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			id := file.Name()[:len(file.Name())-5] // Remove .json extension
			if p, err := f.readPersona(id); err == nil {
				personas = append(personas, p)
			}
		}
	}
	
	return personas, nil
}

func (f *FileStorage) Update(id string, p types.Persona) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	// Check if persona exists
	if _, err := f.readPersona(id); err != nil {
		return fmt.Errorf("persona not found: %s", id)
	}
	
	p.ID = id
	return f.writePersona(p)
}

func (f *FileStorage) Delete(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	filePath := filepath.Join(f.dataDir, id+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("persona not found: %s", id)
	}
	
	return os.Remove(filePath)
}

func (f *FileStorage) generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (f *FileStorage) readPersona(id string) (types.Persona, error) {
	filePath := filepath.Join(f.dataDir, id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return types.Persona{}, fmt.Errorf("persona not found: %s", id)
		}
		return types.Persona{}, fmt.Errorf("failed to read persona file: %v", err)
	}
	
	var p types.Persona
	if err := json.Unmarshal(data, &p); err != nil {
		return types.Persona{}, fmt.Errorf("failed to parse persona data: %v", err)
	}
	
	return p, nil
}

func (f *FileStorage) writePersona(p types.Persona) error {
	filePath := filepath.Join(f.dataDir, p.ID+".json")
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal persona data: %v", err)
	}
	
	return os.WriteFile(filePath, data, 0644)
}
