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
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %v", err)
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

// Persona operations
func (g *GRPCClient) Create(p *types.Persona) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.CreatePersonaRequest{
		Persona: &pb.Persona{
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.Rag,
		},
	}

	resp, err := g.client.CreatePersona(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create persona: %v", err)
	}

	// Update the persona with the returned ID
	p.Id = resp.Persona.Id
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
		Id:      resp.Persona.Id,
		Name:    resp.Persona.Name,
		Topic:   resp.Persona.Topic,
		Prompt:  resp.Persona.Prompt,
		Context: resp.Persona.Context,
		Rag:     resp.Persona.Rag,
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
			Id:      p.Id,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.Rag,
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
			Rag:     p.Rag,
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

// Identity operations
func (g *GRPCClient) CreateIdentity(i *types.Identity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.CreateIdentityRequest{
		Identity: &pb.Identity{
			PersonaId:      i.PersonaId,
			Name:           i.Name,
			Description:    i.Description,
			RichAttributes: i.RichAttributes,
			IsActive:       i.IsActive,
			Tags:           i.Tags,
		},
	}

	resp, err := g.client.CreateIdentity(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create identity: %v", err)
	}

	// Update the identity with the returned ID
	i.Id = resp.Identity.Id
	return nil
}

func (g *GRPCClient) GetIdentity(id string) (types.Identity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetIdentityRequest{Id: id}

	resp, err := g.client.GetIdentity(ctx, req)
	if err != nil {
		return types.Identity{}, fmt.Errorf("failed to get identity: %v", err)
	}

	var createdAt, updatedAt time.Time
	if resp.Identity.CreatedAt != nil {
		createdAt = resp.Identity.CreatedAt.AsTime()
	}
	if resp.Identity.UpdatedAt != nil {
		updatedAt = resp.Identity.UpdatedAt.AsTime()
	}

	return types.Identity{
		Id:             resp.Identity.Id,
		PersonaId:      resp.Identity.PersonaId,
		Name:           resp.Identity.Name,
		Description:    resp.Identity.Description,
		RichAttributes: resp.Identity.RichAttributes,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		IsActive:       resp.Identity.IsActive,
		Tags:           resp.Identity.Tags,
	}, nil
}

func (g *GRPCClient) ListIdentities(filter *types.IdentityFilter) ([]types.Identity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var pbFilter *pb.IdentityFilter
	if filter != nil {
		pbFilter = &pb.IdentityFilter{
			PersonaId: filter.PersonaID,
			Tags:      filter.Tags,
			Search:    filter.Search,
		}
		if filter.IsActive != nil {
			pbFilter.IsActive = *filter.IsActive
		}
	}

	req := &pb.ListIdentitiesRequest{Filter: pbFilter}

	resp, err := g.client.ListIdentities(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list identities: %v", err)
	}

	var identities []types.Identity
	for _, i := range resp.Identities {
		var createdAt, updatedAt time.Time
		if i.CreatedAt != nil {
			createdAt = i.CreatedAt.AsTime()
		}
		if i.UpdatedAt != nil {
			updatedAt = i.UpdatedAt.AsTime()
		}

		identities = append(identities, types.Identity{
			Id:             i.Id,
			PersonaId:      i.PersonaId,
			Name:           i.Name,
			Description:    i.Description,
			RichAttributes: i.RichAttributes,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
			IsActive:       i.IsActive,
			Tags:           i.Tags,
		})
	}

	return identities, nil
}

func (g *GRPCClient) UpdateIdentity(id string, i types.Identity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.UpdateIdentityRequest{
		Id: id,
		Identity: &pb.Identity{
			PersonaId:      i.PersonaId,
			Name:           i.Name,
			Description:    i.Description,
			RichAttributes: i.RichAttributes,
			IsActive:       i.IsActive,
			Tags:           i.Tags,
		},
	}

	_, err := g.client.UpdateIdentity(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update identity: %v", err)
	}

	return nil
}

func (g *GRPCClient) DeleteIdentity(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.DeleteIdentityRequest{Id: id}

	_, err := g.client.DeleteIdentity(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete identity: %v", err)
	}

	return nil
}

func (g *GRPCClient) GetIdentityWithPersona(id string) (types.IdentityWithPersona, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetIdentityWithPersonaRequest{Id: id}

	resp, err := g.client.GetIdentityWithPersona(ctx, req)
	if err != nil {
		return types.IdentityWithPersona{}, fmt.Errorf("failed to get identity with persona: %v", err)
	}

	var createdAt, updatedAt time.Time
	if resp.IdentityWithPersona.Identity.CreatedAt != nil {
		createdAt = resp.IdentityWithPersona.Identity.CreatedAt.AsTime()
	}
	if resp.IdentityWithPersona.Identity.UpdatedAt != nil {
		updatedAt = resp.IdentityWithPersona.Identity.UpdatedAt.AsTime()
	}

	identity := types.Identity{
		Id:             resp.IdentityWithPersona.Identity.Id,
		PersonaId:      resp.IdentityWithPersona.Identity.PersonaId,
		Name:           resp.IdentityWithPersona.Identity.Name,
		Description:    resp.IdentityWithPersona.Identity.Description,
		RichAttributes: resp.IdentityWithPersona.Identity.RichAttributes,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		IsActive:       resp.IdentityWithPersona.Identity.IsActive,
		Tags:           resp.IdentityWithPersona.Identity.Tags,
	}

	persona := types.Persona{
		Id:      resp.IdentityWithPersona.Persona.Id,
		Name:    resp.IdentityWithPersona.Persona.Name,
		Topic:   resp.IdentityWithPersona.Persona.Topic,
		Prompt:  resp.IdentityWithPersona.Persona.Prompt,
		Context: resp.IdentityWithPersona.Persona.Context,
		Rag:     resp.IdentityWithPersona.Persona.Rag,
	}

	return types.IdentityWithPersona{
		Identity: identity,
		Persona:  persona,
	}, nil
}
