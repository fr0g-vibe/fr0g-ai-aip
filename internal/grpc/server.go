package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
	"github.com/fr0g-ai/fr0g-ai-aip/internal/persona"
)

// Server implements the PersonaService gRPC server
type Server struct {
	pb.UnimplementedPersonaServiceServer
}

// NewServer creates a new gRPC server instance
func NewServer() *Server {
	return &Server{}
}

// StartGRPCServer starts the gRPC server on the specified port
func StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	pb.RegisterPersonaServiceServer(s, NewServer())
	
	// Enable reflection for debugging with tools like grpcurl
	reflection.Register(s)

	fmt.Printf("gRPC server listening on port %s\n", port)
	return s.Serve(lis)
}

// CreatePersona creates a new persona
func (s *Server) CreatePersona(ctx context.Context, req *pb.CreatePersonaRequest) (*pb.CreatePersonaResponse, error) {
	if req.Persona == nil {
		return nil, fmt.Errorf("persona is required")
	}

	p := persona.Persona{
		Name:    req.Persona.Name,
		Topic:   req.Persona.Topic,
		Prompt:  req.Persona.Prompt,
		Context: req.Persona.Context,
		RAG:     req.Persona.Rag,
	}

	if err := persona.CreatePersona(&p); err != nil {
		return nil, err
	}

	return &pb.CreatePersonaResponse{
		Persona: &pb.Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}, nil
}

// GetPersona retrieves a persona by ID
func (s *Server) GetPersona(ctx context.Context, req *pb.GetPersonaRequest) (*pb.GetPersonaResponse, error) {
	if req.Id == "" {
		return nil, fmt.Errorf("persona ID is required")
	}

	p, err := persona.GetPersona(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetPersonaResponse{
		Persona: &pb.Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}, nil
}

// ListPersonas returns all personas
func (s *Server) ListPersonas(ctx context.Context, req *pb.ListPersonasRequest) (*pb.ListPersonasResponse, error) {
	personas := persona.ListPersonas()
	
	pbPersonas := make([]*pb.Persona, len(personas))
	for i, p := range personas {
		pbPersonas[i] = &pb.Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		}
	}

	return &pb.ListPersonasResponse{
		Personas: pbPersonas,
	}, nil
}

// UpdatePersona updates an existing persona
func (s *Server) UpdatePersona(ctx context.Context, req *pb.UpdatePersonaRequest) (*pb.UpdatePersonaResponse, error) {
	if req.Id == "" {
		return nil, fmt.Errorf("persona ID is required")
	}
	if req.Persona == nil {
		return nil, fmt.Errorf("persona is required")
	}

	p := persona.Persona{
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

	return &pb.UpdatePersonaResponse{
		Persona: &pb.Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}, nil
}

// DeletePersona removes a persona by ID
func (s *Server) DeletePersona(ctx context.Context, req *pb.DeletePersonaRequest) (*pb.DeletePersonaResponse, error) {
	if req.Id == "" {
		return nil, fmt.Errorf("persona ID is required")
	}

	if err := persona.DeletePersona(req.Id); err != nil {
		return nil, err
	}

	return &pb.DeletePersonaResponse{}, nil
}
