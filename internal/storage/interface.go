package storage

import "github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"

// Storage defines the interface for persona storage backends
type Storage interface {
	Create(p *persona.Persona) error
	Get(id string) (persona.Persona, error)
	List() ([]persona.Persona, error)
	Update(id string, p persona.Persona) error
	Delete(id string) error
}
