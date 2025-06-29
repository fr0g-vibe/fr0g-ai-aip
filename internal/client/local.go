package client

import (
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
)

// LocalClient implements local storage client for persona service
type LocalClient struct {
	storage storage.Storage
}

// NewLocalClient creates a new local client with the given storage backend
func NewLocalClient(storage storage.Storage) *LocalClient {
	return &LocalClient{
		storage: storage,
	}
}

func (l *LocalClient) Create(p *persona.Persona) error {
	return l.storage.Create(p)
}

func (l *LocalClient) Get(id string) (persona.Persona, error) {
	return l.storage.Get(id)
}

func (l *LocalClient) List() ([]persona.Persona, error) {
	return l.storage.List()
}

func (l *LocalClient) Update(id string, p persona.Persona) error {
	return l.storage.Update(id, p)
}

func (l *LocalClient) Delete(id string) error {
	return l.storage.Delete(id)
}
