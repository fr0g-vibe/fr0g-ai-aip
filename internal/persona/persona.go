// Package persona provides AI persona management functionality.
//
// This package implements the core persona service that handles creation,
// retrieval, updating, and deletion of AI subject matter experts. Each
// persona consists of a name, topic, system prompt, and optional context
// and RAG (Retrieval-Augmented Generation) documents.
//
// The service supports both in-memory and file-based storage backends,
// with comprehensive validation and error handling. All operations are
// safe for concurrent use.
//
// Example usage:
//
//	storage := storage.NewMemoryStorage()
//	service := persona.NewService(storage)
//	
//	p := &types.Persona{
//		Name:   "Go Expert",
//		Topic:  "Golang Programming",
//		Prompt: "You are an expert Go programmer with deep knowledge of best practices.",
//		Context: map[string]string{
//			"experience": "10 years",
//			"specialty":  "backend development",
//		},
//	}
//	
//	err := service.CreatePersona(p)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The package also provides identity management functionality, allowing
// creation of persona instances with rich demographic and behavioral
// attributes for community simulation and analysis.
package persona

import (
	"fmt"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/middleware"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Persona is a type alias for backward compatibility.
// Use types.Persona directly in new code.
type Persona = types.Persona

// Service provides persona and identity management operations.
//
// The service acts as the main interface for all persona and identity
// operations, providing validation, error handling, and business logic
// on top of the storage layer. All methods are safe for concurrent use.
//
// The service supports:
//   - CRUD operations for personas
//   - CRUD operations for identities with rich attributes
//   - Validation of all input data
//   - Automatic timestamp management
//   - Reference integrity checking
type Service struct {
	storage storage.Storage
}

// NewService creates a new persona service with the given storage backend.
//
// The storage backend must implement the storage.Storage interface and
// be safe for concurrent use. The service will use this backend for all
// data persistence operations.
//
// Example:
//
//	memStorage := storage.NewMemoryStorage()
//	service := persona.NewService(memStorage)
//
//	fileStorage, err := storage.NewFileStorage("./data")
//	if err != nil {
//		log.Fatal(err)
//	}
//	service := persona.NewService(fileStorage)
func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// CreatePersona creates a new AI persona with validation.
//
// The persona must have a non-empty name, topic, and prompt. The function
// will generate a unique ID and store the persona using the configured
// storage backend. Input data is automatically sanitized and validated.
//
// The persona's name, topic, and prompt are required fields and will be
// validated for length and content. Context and RAG fields are optional
// but will be validated if provided.
//
// Returns an error if:
//   - persona is nil
//   - required fields are empty or contain only whitespace
//   - field values exceed maximum length limits
//   - storage operation fails
//
// Example:
//
//	p := &types.Persona{
//		Name:   "Security Expert",
//		Topic:  "Cybersecurity", 
//		Prompt: "You are a cybersecurity expert with extensive knowledge...",
//		Context: map[string]string{
//			"domain": "enterprise security",
//			"experience": "15 years",
//		},
//	}
//	err := service.CreatePersona(p)
//	if err != nil {
//		log.Printf("Failed to create persona: %v", err)
//	}
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

// GetPersona retrieves a persona by ID.
//
// Returns the complete persona data including name, topic, prompt,
// context, and RAG documents. The ID must be a valid persona identifier
// as returned by CreatePersona or ListPersonas.
//
// Returns an error if the persona is not found or if there's a storage error.
//
// Example:
//
//	persona, err := service.GetPersona("abc123")
//	if err != nil {
//		log.Printf("Persona not found: %v", err)
//		return
//	}
//	fmt.Printf("Found persona: %s (%s)", persona.Name, persona.Topic)
func (s *Service) GetPersona(id string) (types.Persona, error) {
	return s.storage.Get(id)
}

// ListPersonas returns all personas.
//
// Returns a slice containing all stored personas. The slice will be empty
// if no personas exist. The order of personas in the slice is not guaranteed.
//
// Returns an error only if there's a storage error. An empty result set
// is not considered an error.
//
// Example:
//
//	personas, err := service.ListPersonas()
//	if err != nil {
//		log.Printf("Failed to list personas: %v", err)
//		return
//	}
//	fmt.Printf("Found %d personas", len(personas))
//	for _, p := range personas {
//		fmt.Printf("- %s: %s", p.Name, p.Topic)
//	}
func (s *Service) ListPersonas() ([]types.Persona, error) {
	return s.storage.List()
}

// DeletePersona removes a persona by ID.
//
// Permanently deletes the specified persona from storage. This operation
// cannot be undone. The persona ID must be valid and the persona must exist.
//
// Note: This does not automatically delete any identities that reference
// this persona. Consider checking for dependent identities before deletion.
//
// Returns an error if the persona is not found or if there's a storage error.
//
// Example:
//
//	err := service.DeletePersona("abc123")
//	if err != nil {
//		log.Printf("Failed to delete persona: %v", err)
//	}
func (s *Service) DeletePersona(id string) error {
	return s.storage.Delete(id)
}

// UpdatePersona updates an existing persona with validation.
//
// Updates the persona with the provided data. All fields in the persona
// struct will replace the existing values. The persona must exist and
// the new data must pass validation.
//
// Input data is automatically sanitized and validated using the same
// rules as CreatePersona. The persona ID cannot be changed.
//
// Returns an error if:
//   - the persona is not found
//   - validation fails
//   - storage operation fails
//
// Example:
//
//	persona.Name = "Updated Security Expert"
//	persona.Context["updated"] = "true"
//	err := service.UpdatePersona("abc123", persona)
//	if err != nil {
//		log.Printf("Failed to update persona: %v", err)
//	}
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

// GetStorage returns the underlying storage interface
func (s *Service) GetStorage() storage.Storage {
	return s.storage
}

// CreateIdentity creates a new identity with validation.
//
// An identity represents an instance of a persona with specific demographic
// and behavioral attributes. The identity must reference a valid persona
// and include a name. Rich attributes, tags, and other fields are optional.
//
// The function automatically:
//   - Validates the referenced persona exists
//   - Sets creation and update timestamps
//   - Initializes default values for optional fields
//   - Generates a unique identity ID
//
// Rich attributes can include demographic information like age, gender,
// location, political leaning, education level, and interests. These
// attributes are used for community generation and analysis.
//
// Returns an error if:
//   - identity is nil
//   - referenced persona does not exist
//   - required fields are missing
//   - storage operation fails
//
// Example:
//
//	identity := &types.Identity{
//		PersonaId:   "abc123",
//		Name:        "Alice Johnson",
//		Description: "Senior cybersecurity analyst",
//		RichAttributes: &types.RichAttributes{
//			Age:     32,
//			Gender:  "female",
//			Education: "master",
//		},
//		Tags: []string{"security", "analyst"},
//	}
//	err := service.CreateIdentity(identity)
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

// GetCommunity retrieves a community by ID
func (s *Service) GetCommunity(id string) (types.Community, error) {
	return s.storage.GetCommunity(id)
}

// ListCommunities returns communities with optional filtering
func (s *Service) ListCommunities(filter *types.CommunityFilter) ([]types.Community, error) {
	return s.storage.ListCommunities(filter)
}

// UpdateCommunity updates an existing community
func (s *Service) UpdateCommunity(id string, community types.Community) error {
	return s.storage.UpdateCommunity(id, community)
}

// DeleteCommunity removes a community
func (s *Service) DeleteCommunity(id string) error {
	return s.storage.DeleteCommunity(id)
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
