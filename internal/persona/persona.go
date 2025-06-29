package persona

import (
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
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

// Service provides persona management operations
type Service struct {
	storage storage.Storage
}

// NewService creates a new persona service with the given storage backend
func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// CreatePersona creates a new persona
func (s *Service) CreatePersona(p *Persona) error {
	return s.storage.Create(p)
}

// GetPersona retrieves a persona by ID
func (s *Service) GetPersona(id string) (Persona, error) {
	return s.storage.Get(id)
}

// ListPersonas returns all personas
func (s *Service) ListPersonas() ([]Persona, error) {
	return s.storage.List()
}

// DeletePersona removes a persona by ID
func (s *Service) DeletePersona(id string) error {
	return s.storage.Delete(id)
}

// UpdatePersona updates an existing persona
func (s *Service) UpdatePersona(id string, p Persona) error {
	return s.storage.Update(id, p)
}

// Global service instance for backward compatibility
var defaultService *Service

// SetDefaultService sets the default service instance
func SetDefaultService(service *Service) {
	defaultService = service
}

// Legacy functions for backward compatibility
func CreatePersona(p *Persona) error {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.CreatePersona(p)
}

func GetPersona(id string) (Persona, error) {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.GetPersona(id)
}

func ListPersonas() []Persona {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	personas, _ := defaultService.ListPersonas()
	return personas
}

func DeletePersona(id string) error {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.DeletePersona(id)
}

func UpdatePersona(id string, p Persona) error {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.UpdatePersona(id, p)
}
