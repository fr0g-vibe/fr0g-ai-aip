package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// PersonaServer implements the gRPC PersonaService
type PersonaServer struct {
	// Embed the unimplemented server for forward compatibility
	UnimplementedPersonaServiceServer
}

// CreatePersona implements the CreatePersona RPC
func (s *PersonaServer) CreatePersona(ctx context.Context, req *CreatePersonaRequest) (*CreatePersonaResponse, error) {
	if req.Persona == nil {
		return nil, fmt.Errorf("persona is required")
	}

	p := &types.Persona{
		Name:    req.Persona.Name,
		Topic:   req.Persona.Topic,
		Prompt:  req.Persona.Prompt,
		Context: req.Persona.Context,
		RAG:     req.Persona.Rag,
	}

	if err := persona.CreatePersona(p); err != nil {
		return nil, err
	}

	return &CreatePersonaResponse{
		Persona: &Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}, nil
}

// GetPersona implements the GetPersona RPC
func (s *PersonaServer) GetPersona(ctx context.Context, req *GetPersonaRequest) (*GetPersonaResponse, error) {
	p, err := persona.GetPersona(req.Id)
	if err != nil {
		return nil, err
	}

	return &GetPersonaResponse{
		Persona: &Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}, nil
}

// ListPersonas implements the ListPersonas RPC
func (s *PersonaServer) ListPersonas(ctx context.Context, req *ListPersonasRequest) (*ListPersonasResponse, error) {
	personas := persona.ListPersonas()
	
	var grpcPersonas []*Persona
	for _, p := range personas {
		grpcPersonas = append(grpcPersonas, &Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		})
	}

	return &ListPersonasResponse{
		Personas: grpcPersonas,
	}, nil
}

// UpdatePersona implements the UpdatePersona RPC
func (s *PersonaServer) UpdatePersona(ctx context.Context, req *UpdatePersonaRequest) (*UpdatePersonaResponse, error) {
	if req.Persona == nil {
		return nil, fmt.Errorf("persona is required")
	}

	p := types.Persona{
		ID:      req.Id,
		Name:    req.Persona.Name,
		Topic:   req.Persona.Topic,
		Prompt:  req.Persona.Prompt,
		Context: req.Persona.Context,
		RAG:     req.Persona.Rag,
	}

	if err := persona.UpdatePersona(req.Id, p); err != nil {
		return nil, err
	}

	return &UpdatePersonaResponse{
		Persona: &Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}, nil
}

// DeletePersona implements the DeletePersona RPC
func (s *PersonaServer) DeletePersona(ctx context.Context, req *DeletePersonaRequest) (*DeletePersonaResponse, error) {
	if err := persona.DeletePersona(req.Id); err != nil {
		return nil, err
	}

	return &DeletePersonaResponse{}, nil
}

// StartGRPCServer starts the gRPC server on the specified port
func StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	
	// Register the persona service
	RegisterPersonaServiceServer(s, &PersonaServer{})
	
	// Enable reflection for debugging with tools like grpcurl
	reflection.Register(s)

	fmt.Printf("gRPC server listening on port %s\n", port)
	return s.Serve(lis)
}

// Placeholder types - these will be replaced by generated protobuf code
type UnimplementedPersonaServiceServer struct{}

func (UnimplementedPersonaServiceServer) CreatePersona(context.Context, *CreatePersonaRequest) (*CreatePersonaResponse, error) {
	return nil, fmt.Errorf("method CreatePersona not implemented")
}

func (UnimplementedPersonaServiceServer) GetPersona(context.Context, *GetPersonaRequest) (*GetPersonaResponse, error) {
	return nil, fmt.Errorf("method GetPersona not implemented")
}

func (UnimplementedPersonaServiceServer) ListPersonas(context.Context, *ListPersonasRequest) (*ListPersonasResponse, error) {
	return nil, fmt.Errorf("method ListPersonas not implemented")
}

func (UnimplementedPersonaServiceServer) UpdatePersona(context.Context, *UpdatePersonaRequest) (*UpdatePersonaResponse, error) {
	return nil, fmt.Errorf("method UpdatePersona not implemented")
}

func (UnimplementedPersonaServiceServer) DeletePersona(context.Context, *DeletePersonaRequest) (*DeletePersonaResponse, error) {
	return nil, fmt.Errorf("method DeletePersona not implemented")
}

// Placeholder types - these will be replaced by generated protobuf code
type Persona struct {
	Id      string
	Name    string
	Topic   string
	Prompt  string
	Context map[string]string
	Rag     []string
}

type CreatePersonaRequest struct {
	Persona *Persona
}

type CreatePersonaResponse struct {
	Persona *Persona
}

type GetPersonaRequest struct {
	Id string
}

type GetPersonaResponse struct {
	Persona *Persona
}

type ListPersonasRequest struct{}

type ListPersonasResponse struct {
	Personas []*Persona
}

type UpdatePersonaRequest struct {
	Id      string
	Persona *Persona
}

type UpdatePersonaResponse struct {
	Persona *Persona
}

type DeletePersonaRequest struct {
	Id string
}

type DeletePersonaResponse struct{}

// Placeholder function - this will be replaced by generated protobuf code
func RegisterPersonaServiceServer(s *grpc.Server, srv interface{}) {
	// This is a placeholder - the real implementation will be generated
}
