package types

import (
	"time"
	pb "github.com/fr0g-vibe/fr0g-ai-aip/internal/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Identity represents a persona-based identity with rich attributes
type Identity struct {
	Id             string                 `json:"id"`
	PersonaId      string                 `json:"persona_id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Attributes     map[string]string      `json:"attributes,omitempty"`     // Legacy
	Preferences    map[string]string      `json:"preferences,omitempty"`    // Legacy
	Background     string                 `json:"background"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	IsActive       bool                   `json:"is_active"`
	Tags           []string               `json:"tags"`
	RichAttributes *RichAttributes        `json:"rich_attributes,omitempty"`
}

// Use protobuf types for rich attributes
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

// ProtoToIdentity converts protobuf Identity to internal Identity
func ProtoToIdentity(pb *pb.Identity) *Identity {
	if pb == nil {
		return nil
	}
	
	var createdAt, updatedAt time.Time
	if pb.CreatedAt != nil {
		createdAt = pb.CreatedAt.AsTime()
	}
	if pb.UpdatedAt != nil {
		updatedAt = pb.UpdatedAt.AsTime()
	}
	
	return &Identity{
		Id:             pb.Id,
		PersonaId:      pb.PersonaId,
		Name:           pb.Name,
		Description:    pb.Description,
		Attributes:     pb.Attributes,
		Preferences:    pb.Preferences,
		Background:     pb.Background,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		IsActive:       pb.IsActive,
		Tags:           pb.Tags,
		RichAttributes: pb.RichAttributes,
	}
}

// IdentityToProto converts internal Identity to protobuf Identity
func IdentityToProto(i *Identity) *pb.Identity {
	if i == nil {
		return nil
	}
	
	var createdAt, updatedAt *timestamppb.Timestamp
	if !i.CreatedAt.IsZero() {
		createdAt = timestamppb.New(i.CreatedAt)
	}
	if !i.UpdatedAt.IsZero() {
		updatedAt = timestamppb.New(i.UpdatedAt)
	}
	
	return &pb.Identity{
		Id:             i.Id,
		PersonaId:      i.PersonaId,
		Name:           i.Name,
		Description:    i.Description,
		Attributes:     i.Attributes,
		Preferences:    i.Preferences,
		Background:     i.Background,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		IsActive:       i.IsActive,
		Tags:           i.Tags,
		RichAttributes: i.RichAttributes,
	}
}

// NewIdentity creates a new identity with default values
func NewIdentity() *Identity {
	now := time.Now()
	return &Identity{
		RichAttributes: &RichAttributes{},
		Tags:           []string{},
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
