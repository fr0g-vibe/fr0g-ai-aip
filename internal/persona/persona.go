package persona

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

// Persona represents an AI persona with specific expertise
type Persona struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Topic   string            `json:"topic"`
	Prompt  string            `json:"prompt"`
	Context map[string]string `json:"context,omitempty"`
	RAG     []string          `json:"rag,omitempty"`
}

// In-memory storage for personas (replace with persistent storage later)
var (
	personas = make(map[string]Persona)
	mu       sync.RWMutex
)

// generateID creates a random ID for a persona
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CreatePersona creates a new persona
func CreatePersona(p Persona) error {
	mu.Lock()
	defer mu.Unlock()
	
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
	personas[p.ID] = p
	return nil
}

// GetPersona retrieves a persona by ID
func GetPersona(id string) (Persona, error) {
	mu.RLock()
	defer mu.RUnlock()
	
	p, exists := personas[id]
	if !exists {
		return Persona{}, fmt.Errorf("persona not found: %s", id)
	}
	return p, nil
}

// ListPersonas returns all personas
func ListPersonas() []Persona {
	mu.RLock()
	defer mu.RUnlock()
	
	result := make([]Persona, 0, len(personas))
	for _, p := range personas {
		result = append(result, p)
	}
	return result
}

// DeletePersona removes a persona by ID
func DeletePersona(id string) error {
	mu.Lock()
	defer mu.Unlock()
	
	if _, exists := personas[id]; !exists {
		return fmt.Errorf("persona not found: %s", id)
	}
	delete(personas, id)
	return nil
}

// UpdatePersona updates an existing persona
func UpdatePersona(id string, p Persona) error {
	mu.Lock()
	defer mu.Unlock()
	
	if _, exists := personas[id]; !exists {
		return fmt.Errorf("persona not found: %s", id)
	}
	
	p.ID = id
	personas[id] = p
	return nil
}
