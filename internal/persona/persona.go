package persona

import (
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/middleware"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// Type alias for backward compatibility
type Persona = types.Persona

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

// CreatePersona creates a new persona with validation
func (s *Service) CreatePersona(p *types.Persona) error {
	// Sanitize input
	middleware.SanitizePersona(p)
	
	// Validate input
	if err := middleware.ValidatePersona(p); err != nil {
		return err
	}
	
	// Create persona
	return s.storage.Create(p)
}

// GetPersona retrieves a persona by ID
func (s *Service) GetPersona(id string) (types.Persona, error) {
	return s.storage.Get(id)
}

// ListPersonas returns all personas
func (s *Service) ListPersonas() ([]types.Persona, error) {
	return s.storage.List()
}

// DeletePersona removes a persona by ID
func (s *Service) DeletePersona(id string) error {
	return s.storage.Delete(id)
}

// UpdatePersona updates an existing persona with validation
func (s *Service) UpdatePersona(id string, p types.Persona) error {
	// Sanitize input
	middleware.SanitizePersona(&p)
	
	// Validate input
	if err := middleware.ValidatePersona(&p); err != nil {
		return err
	}
	
	// Update persona
	return s.storage.Update(id, p)
}

// Global service instance for backward compatibility
var defaultService *Service

// SetDefaultService sets the default service instance
func SetDefaultService(service *Service) {
	defaultService = service
}

// Legacy functions for backward compatibility
func CreatePersona(p *types.Persona) error {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.CreatePersona(p)
}

func GetPersona(id string) (types.Persona, error) {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.GetPersona(id)
}

func ListPersonas() []types.Persona {
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

func UpdatePersona(id string, p types.Persona) error {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.UpdatePersona(id, p)
}
