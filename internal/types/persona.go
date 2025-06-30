package types

import (
	"time"
	pb "github.com/fr0g-vibe/fr0g-ai-aip/internal/grpc/pb"
)

// Persona represents an AI persona with specific expertise
type Persona struct {
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	Topic   string            `json:"topic"`
	Prompt  string            `json:"prompt"`
	Context map[string]string `json:"context"`
	Rag     []string          `json:"rag"`
	
	// Additional fields not in proto
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProtoToPersona converts protobuf Persona to internal Persona
func ProtoToPersona(pb *pb.Persona) *Persona {
	if pb == nil {
		return nil
	}
	return &Persona{
		Id:      pb.Id,
		Name:    pb.Name,
		Topic:   pb.Topic,
		Prompt:  pb.Prompt,
		Context: pb.Context,
		Rag:     pb.Rag,
	}
}

// PersonaToProto converts internal Persona to protobuf Persona
func PersonaToProto(p *Persona) *pb.Persona {
	if p == nil {
		return nil
	}
	return &pb.Persona{
		Id:      p.Id,
		Name:    p.Name,
		Topic:   p.Topic,
		Prompt:  p.Prompt,
		Context: p.Context,
		Rag:     p.Rag,
	}
}
