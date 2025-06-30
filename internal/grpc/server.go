package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/grpc/pb/proto"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// PersonaServer implements the gRPC PersonaService
type PersonaServer struct {
	proto.UnimplementedPersonaServiceServer
}

// StartGRPCServer starts a real gRPC server using protobuf
func StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterPersonaServiceServer(s, &PersonaServer{})

	fmt.Printf("gRPC server listening on port %s\n", port)
	fmt.Println("Using real gRPC with protobuf")

	return s.Serve(lis)
}

// CreatePersona creates a new persona
func (s *PersonaServer) CreatePersona(ctx context.Context, req *proto.CreatePersonaRequest) (*proto.CreatePersonaResponse, error) {
	if req.Persona == nil {
		return nil, status.Errorf(codes.InvalidArgument, "persona is required")
	}

	p := &types.Persona{
		Name:    req.Persona.Name,
		Topic:   req.Persona.Topic,
		Prompt:  req.Persona.Prompt,
		Context: req.Persona.Context,
		RAG:     req.Persona.Rag,
	}

	if err := persona.CreatePersona(p); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &proto.CreatePersonaResponse{
		Persona: &proto.Persona{
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
func (s *PersonaServer) GetPersona(ctx context.Context, req *proto.GetPersonaRequest) (*proto.GetPersonaResponse, error) {
	p, err := persona.GetPersona(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &proto.GetPersonaResponse{
		Persona: &proto.Persona{
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
func (s *PersonaServer) ListPersonas(ctx context.Context, req *proto.ListPersonasRequest) (*proto.ListPersonasResponse, error) {
	personas := persona.ListPersonas()

	var protoPersonas []*proto.Persona
	for _, p := range personas {
		protoPersonas = append(protoPersonas, &proto.Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		})
	}

	return &proto.ListPersonasResponse{
		Personas: protoPersonas,
	}, nil
}

// UpdatePersona updates an existing persona
func (s *PersonaServer) UpdatePersona(ctx context.Context, req *proto.UpdatePersonaRequest) (*proto.UpdatePersonaResponse, error) {
	if req.Persona == nil {
		return nil, status.Errorf(codes.InvalidArgument, "persona is required")
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
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &proto.UpdatePersonaResponse{
		Persona: &proto.Persona{
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
func (s *PersonaServer) DeletePersona(ctx context.Context, req *proto.DeletePersonaRequest) (*proto.DeletePersonaResponse, error) {
	if err := persona.DeletePersona(req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &proto.DeletePersonaResponse{}, nil
}
