package storage

import "github.com/fr0g-vibe/fr0g-ai-aip/internal/types"

// Storage defines the interface for persona storage backends
type Storage interface {
	// Persona operations
	Create(p *types.Persona) error
	Get(id string) (types.Persona, error)
	List() ([]types.Persona, error)
	Update(id string, p types.Persona) error
	Delete(id string) error

	// Identity operations
	CreateIdentity(i *types.Identity) error
	GetIdentity(id string) (types.Identity, error)
	ListIdentities(filter *types.IdentityFilter) ([]types.Identity, error)
	UpdateIdentity(id string, i types.Identity) error
	DeleteIdentity(id string) error
	GetIdentityWithPersona(id string) (types.IdentityWithPersona, error)
}
