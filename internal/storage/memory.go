package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
)

// MemoryStorage implements in-memory storage for personas
type MemoryStorage struct {
	personas map[string]persona.Persona
	mu       sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		personas: make(map[string]persona.Persona),
	}
}

// generateID creates a random ID for a persona
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (m *MemoryStorage) Create(p *persona.Persona) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
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

func (m *MemoryStorage) Get(id string) (persona.Persona, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	p, exists := m.personas[id]
	if !exists {
		return persona.Persona{}, fmt.Errorf("persona not found: %s", id)
	}
	return p, nil
}

func (m *MemoryStorage) List() ([]persona.Persona, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	result := make([]persona.Persona, 0, len(m.personas))
	for _, p := range m.personas {
		result = append(result, p)
	}
	return result, nil
}

func (m *MemoryStorage) Update(id string, p persona.Persona) error {
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
