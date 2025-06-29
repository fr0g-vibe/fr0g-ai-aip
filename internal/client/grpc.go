package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// GRPCClient implements gRPC client for persona service
type GRPCClient struct {
	conn   *grpc.ClientConn
	client PersonaServiceClient
}

// NewGRPCClient creates a new gRPC client
func NewGRPCClient(address string) (*GRPCClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	client := NewPersonaServiceClient(conn)

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

	req := &CreatePersonaRequest{
		Persona: &Persona{
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

	req := &GetPersonaRequest{Id: id}
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

	req := &ListPersonasRequest{}
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

	req := &UpdatePersonaRequest{
		Id: id,
		Persona: &Persona{
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

	req := &DeletePersonaRequest{Id: id}
	_, err := g.client.DeletePersona(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete persona: %v", err)
	}

	return nil
}

// Placeholder types - these will be replaced by generated protobuf code
type PersonaServiceClient interface {
	CreatePersona(ctx context.Context, in *CreatePersonaRequest, opts ...grpc.CallOption) (*CreatePersonaResponse, error)
	GetPersona(ctx context.Context, in *GetPersonaRequest, opts ...grpc.CallOption) (*GetPersonaResponse, error)
	ListPersonas(ctx context.Context, in *ListPersonasRequest, opts ...grpc.CallOption) (*ListPersonasResponse, error)
	UpdatePersona(ctx context.Context, in *UpdatePersonaRequest, opts ...grpc.CallOption) (*UpdatePersonaResponse, error)
	DeletePersona(ctx context.Context, in *DeletePersonaRequest, opts ...grpc.CallOption) (*DeletePersonaResponse, error)
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
func NewPersonaServiceClient(cc grpc.ClientConnInterface) PersonaServiceClient {
	return &personaServiceClient{cc}
}

type personaServiceClient struct {
	cc grpc.ClientConnInterface
}

func (c *personaServiceClient) CreatePersona(ctx context.Context, in *CreatePersonaRequest, opts ...grpc.CallOption) (*CreatePersonaResponse, error) {
	return nil, fmt.Errorf("gRPC client not fully implemented - run 'make build-with-grpc' to generate protobuf code")
}

func (c *personaServiceClient) GetPersona(ctx context.Context, in *GetPersonaRequest, opts ...grpc.CallOption) (*GetPersonaResponse, error) {
	return nil, fmt.Errorf("gRPC client not fully implemented - run 'make build-with-grpc' to generate protobuf code")
}

func (c *personaServiceClient) ListPersonas(ctx context.Context, in *ListPersonasRequest, opts ...grpc.CallOption) (*ListPersonasResponse, error) {
	return nil, fmt.Errorf("gRPC client not fully implemented - run 'make build-with-grpc' to generate protobuf code")
}

func (c *personaServiceClient) UpdatePersona(ctx context.Context, in *UpdatePersonaRequest, opts ...grpc.CallOption) (*UpdatePersonaResponse, error) {
	return nil, fmt.Errorf("gRPC client not fully implemented - run 'make build-with-grpc' to generate protobuf code")
}

func (c *personaServiceClient) DeletePersona(ctx context.Context, in *DeletePersonaRequest, opts ...grpc.CallOption) (*DeletePersonaResponse, error) {
	return nil, fmt.Errorf("gRPC client not fully implemented - run 'make build-with-grpc' to generate protobuf code")
}
