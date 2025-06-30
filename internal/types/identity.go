package types

import (
	pb "github.com/fr0g-vibe/fr0g-ai-aip/internal/grpc/pb"
)

// Use protobuf types directly to eliminate duplication
type Identity = pb.Identity
type RichAttributes = pb.RichAttributes
type Demographics = pb.Demographics
type Psychographics = pb.Psychographics
type LifeHistory = pb.LifeHistory
type CulturalReligious = pb.CulturalReligious
type PoliticalSocial = pb.PoliticalSocial
type Health = pb.Health
type Preferences = pb.Preferences
type BehavioralTendencies = pb.BehavioralTendencies
type CurrentContext = pb.CurrentContext
type Location = pb.Location
type Personality = pb.Personality
type LifeEvent = pb.LifeEvent
type Education = pb.Education
type Career = pb.Career
type AgeRange = pb.AgeRange

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

	// New filters for rich attributes
	AgeRange         *AgeRange    `json:"age_range,omitempty"`
	Location         *Location    `json:"location,omitempty"`
	PoliticalLeaning string       `json:"political_leaning,omitempty"`
	Education        string       `json:"education,omitempty"`
	Occupation       string       `json:"occupation,omitempty"`
	Personality      *Personality `json:"personality,omitempty"`
}

// Helper function for creating new identities
func NewIdentity() *Identity {
	return &Identity{
		RichAttributes: &RichAttributes{},
		Tags:           []string{},
		IsActive:       true,
	}
}
