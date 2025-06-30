package types

import (
	"time"
)

// Identity represents a persona-based identity with additional identifying attributes
type Identity struct {
	ID          string            `json:"id"`
	PersonaID   string            `json:"persona_id"`  // Reference to the base persona
	Name        string            `json:"name"`        // Identity-specific name
	Description string            `json:"description"` // Identity description
	Attributes  map[string]string `json:"attributes"`  // Identity-specific attributes
	Preferences map[string]string `json:"preferences"` // Identity preferences
	Background  string            `json:"background"`  // Personal background story
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	IsActive    bool              `json:"is_active"`
	Tags        []string          `json:"tags"` // Identity tags for categorization
}

// IdentityWithPersona combines an identity with its base persona
type IdentityWithPersona struct {
	Identity Identity `json:"identity"`
	Persona  Persona  `json:"persona"`
}

// IdentityFilter represents filters for listing identities
type IdentityFilter struct {
	PersonaID string   `json:"persona_id,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	IsActive  *bool    `json:"is_active,omitempty"`
	Search    string   `json:"search,omitempty"`
}
