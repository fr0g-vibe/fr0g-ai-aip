package storage

import "github.com/fr0g-vibe/fr0g-ai-aip/internal/types"

// Storage defines the interface for persona storage backends
type Storage interface {
	Create(p *types.Persona) error
	Get(id string) (types.Persona, error)
	List() ([]types.Persona, error)
	Update(id string, p types.Persona) error
	Delete(id string) error
}
