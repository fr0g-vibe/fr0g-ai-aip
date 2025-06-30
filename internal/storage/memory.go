package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// MemoryStorage implements in-memory storage for personas
type MemoryStorage struct {
	personas map[string]types.Persona
	mu       sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		personas: make(map[string]types.Persona),
	}
}

// generateID creates a random ID for a persona
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (m *MemoryStorage) Create(p *types.Persona) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
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
	
	p.ID = generateID()
	m.personas[p.ID] = *p
	return nil
}

func (m *MemoryStorage) Get(id string) (types.Persona, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	p, exists := m.personas[id]
	if !exists {
		return types.Persona{}, fmt.Errorf("persona not found: %s", id)
	}
	return p, nil
}

func (m *MemoryStorage) List() ([]types.Persona, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	result := make([]types.Persona, 0, len(m.personas))
	for _, p := range m.personas {
		result = append(result, p)
	}
	return result, nil
}

func (m *MemoryStorage) Update(id string, p types.Persona) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.personas[id]; !exists {
		return fmt.Errorf("persona not found: %s", id)
	}
	
	p.ID = id
	m.personas[id] = p
	return nil
}

func (m *MemoryStorage) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.personas[id]; !exists {
		return fmt.Errorf("persona not found: %s", id)
	}
	delete(m.personas, id)
	return nil
}
