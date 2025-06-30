package persona

import (
	"fmt"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/middleware"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Type alias for backward compatibility
type Persona = types.Persona

// Service provides persona and identity management operations
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
	if p == nil {
		return fmt.Errorf("persona cannot be nil")
	}

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

// CreateIdentity creates a new identity with validation
func (s *Service) CreateIdentity(i *types.Identity) error {
	if i == nil {
		return fmt.Errorf("identity cannot be nil")
	}

	// Validate that the referenced persona exists
	if _, err := s.storage.Get(i.PersonaId); err != nil {
		return fmt.Errorf("referenced persona not found: %v", err)
	}

	// Set timestamps
	now := time.Now()
	i.CreatedAt = timestamppb.New(now)
	i.UpdatedAt = timestamppb.New(now)

	// Set default values
	if i.RichAttributes == nil {
		i.RichAttributes = &types.RichAttributes{}
	}
	if i.Tags == nil {
		i.Tags = []string{}
	}
	if !i.IsActive {
		i.IsActive = true // Default to active
	}

	// Create identity
	return s.storage.CreateIdentity(i)
}

// GetIdentity retrieves an identity by ID
func (s *Service) GetIdentity(id string) (types.Identity, error) {
	return s.storage.GetIdentity(id)
}

// ListIdentities returns identities with optional filtering
func (s *Service) ListIdentities(filter *types.IdentityFilter) ([]types.Identity, error) {
	return s.storage.ListIdentities(filter)
}

// UpdateIdentity updates an existing identity with validation
func (s *Service) UpdateIdentity(id string, i types.Identity) error {
	// Validate that the referenced persona exists
	if _, err := s.storage.Get(i.PersonaId); err != nil {
		return fmt.Errorf("referenced persona not found: %v", err)
	}

	// Update timestamp
	i.UpdatedAt = timestamppb.New(time.Now())

	// Update identity
	return s.storage.UpdateIdentity(id, i)
}

// DeleteIdentity removes an identity by ID
func (s *Service) DeleteIdentity(id string) error {
	return s.storage.DeleteIdentity(id)
}

// GetIdentityWithPersona retrieves an identity with its associated persona
func (s *Service) GetIdentityWithPersona(id string) (types.IdentityWithPersona, error) {
	return s.storage.GetIdentityWithPersona(id)
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
