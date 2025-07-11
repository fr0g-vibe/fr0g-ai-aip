package client

import (
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// LocalClient implements local storage client for persona and identity service
type LocalClient struct {
	storage storage.Storage
}

// NewLocalClient creates a new local client with the given storage backend
func NewLocalClient(storage storage.Storage) *LocalClient {
	return &LocalClient{
		storage: storage,
	}
}

// Persona operations
func (l *LocalClient) Create(p *types.Persona) error {
	return l.storage.Create(p)
}

func (l *LocalClient) Get(id string) (types.Persona, error) {
	return l.storage.Get(id)
}

func (l *LocalClient) List() ([]types.Persona, error) {
	return l.storage.List()
}

func (l *LocalClient) Update(id string, p types.Persona) error {
	return l.storage.Update(id, p)
}

func (l *LocalClient) Delete(id string) error {
	return l.storage.Delete(id)
}

// Identity operations
func (l *LocalClient) CreateIdentity(i *types.Identity) error {
	return l.storage.CreateIdentity(i)
}

func (l *LocalClient) GetIdentity(id string) (types.Identity, error) {
	return l.storage.GetIdentity(id)
}

func (l *LocalClient) ListIdentities(filter *types.IdentityFilter) ([]types.Identity, error) {
	return l.storage.ListIdentities(filter)
}

func (l *LocalClient) UpdateIdentity(id string, i types.Identity) error {
	return l.storage.UpdateIdentity(id, i)
}

func (l *LocalClient) DeleteIdentity(id string) error {
	return l.storage.DeleteIdentity(id)
}

func (l *LocalClient) GetIdentityWithPersona(id string) (types.IdentityWithPersona, error) {
	return l.storage.GetIdentityWithPersona(id)
}

func (l *LocalClient) Close() error {
	// Local client doesn't need cleanup
	return nil
}
