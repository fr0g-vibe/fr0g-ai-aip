package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/fr0g-vibe/fr0g-ai-aip/internal/grpc/pb"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// GRPCClient implements a real gRPC client using protobuf
type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.PersonaServiceClient
}

// NewGRPCClient creates a new gRPC client
func NewGRPCClient(address string) (*GRPCClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	client := pb.NewPersonaServiceClient(conn)

	return &GRPCClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (g *GRPCClient) Close() error {
	return g.conn.Close()
}

func (g *GRPCClient) Create(p *types.Persona) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.CreatePersonaRequest{
		Persona: &pb.Persona{
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}

	resp, err := g.client.CreatePersona(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create persona: %v", err)
	}

	// Update the persona with the returned ID
	p.ID = resp.Persona.Id
	return nil
}

func (g *GRPCClient) Get(id string) (types.Persona, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetPersonaRequest{Id: id}

	resp, err := g.client.GetPersona(ctx, req)
	if err != nil {
		return types.Persona{}, fmt.Errorf("failed to get persona: %v", err)
	}

	return types.Persona{
		ID:      resp.Persona.Id,
		Name:    resp.Persona.Name,
		Topic:   resp.Persona.Topic,
		Prompt:  resp.Persona.Prompt,
		Context: resp.Persona.Context,
		RAG:     resp.Persona.Rag,
	}, nil
}

func (g *GRPCClient) List() ([]types.Persona, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.ListPersonasRequest{}

	resp, err := g.client.ListPersonas(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list personas: %v", err)
	}

	var personas []types.Persona
	for _, p := range resp.Personas {
		personas = append(personas, types.Persona{
			ID:      p.Id,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			RAG:     p.Rag,
		})
	}

	return personas, nil
}

func (g *GRPCClient) Update(id string, p types.Persona) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.UpdatePersonaRequest{
		Id: id,
		Persona: &pb.Persona{
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}

	_, err := g.client.UpdatePersona(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update persona: %v", err)
	}

	return nil
}

func (g *GRPCClient) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.DeletePersonaRequest{Id: id}

	_, err := g.client.DeletePersona(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete persona: %v", err)
	}

	return nil
}
