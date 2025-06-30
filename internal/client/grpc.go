package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// Local gRPC-compatible types (using only standard library)
type Persona struct {
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	Topic   string            `json:"topic"`
	Prompt  string            `json:"prompt"`
	Context map[string]string `json:"context,omitempty"`
	Rag     []string          `json:"rag,omitempty"`
}

type CreatePersonaRequest struct {
	Persona *Persona `json:"persona"`
}

type CreatePersonaResponse struct {
	Persona *Persona `json:"persona"`
}

type GetPersonaRequest struct {
	Id string `json:"id"`
}

type GetPersonaResponse struct {
	Persona *Persona `json:"persona"`
}

type ListPersonasRequest struct{}

type ListPersonasResponse struct {
	Personas []*Persona `json:"personas"`
}

type UpdatePersonaRequest struct {
	Id      string   `json:"id"`
	Persona *Persona `json:"persona"`
}

type UpdatePersonaResponse struct {
	Persona *Persona `json:"persona"`
}

type DeletePersonaRequest struct {
	Id string `json:"id"`
}

type DeletePersonaResponse struct{}

// GRPCClient implements a gRPC-like client using standard library HTTP
type GRPCClient struct {
	baseURL string
	client  *http.Client
}

// NewGRPCClient creates a new local gRPC-compatible client
func NewGRPCClient(address string) (*GRPCClient, error) {
	// Convert gRPC address to HTTP URL
	baseURL := fmt.Sprintf("http://%s", address)
	
	return &GRPCClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}, nil
}

func (g *GRPCClient) Create(p *types.Persona) error {
	req := CreatePersonaRequest{
		Persona: &Persona{
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}
	
	var resp CreatePersonaResponse
	if err := g.makeRequest("/PersonaService/CreatePersona", req, &resp); err != nil {
		return fmt.Errorf("failed to create persona: %v", err)
	}
	
	// Update the persona with the returned ID
	p.ID = resp.Persona.Id
	return nil
}

func (g *GRPCClient) Get(id string) (types.Persona, error) {
	req := GetPersonaRequest{Id: id}
	
	var resp GetPersonaResponse
	if err := g.makeRequest("/PersonaService/GetPersona", req, &resp); err != nil {
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
	req := ListPersonasRequest{}
	
	var resp ListPersonasResponse
	if err := g.makeRequest("/PersonaService/ListPersonas", req, &resp); err != nil {
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
	req := UpdatePersonaRequest{
		Id: id,
		Persona: &Persona{
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}
	
	var resp UpdatePersonaResponse
	if err := g.makeRequest("/PersonaService/UpdatePersona", req, &resp); err != nil {
		return fmt.Errorf("failed to update persona: %v", err)
	}
	
	return nil
}

func (g *GRPCClient) Delete(id string) error {
	req := DeletePersonaRequest{Id: id}
	
	var resp DeletePersonaResponse
	if err := g.makeRequest("/PersonaService/DeletePersona", req, &resp); err != nil {
		return fmt.Errorf("failed to delete persona: %v", err)
	}
	
	return nil
}

func (g *GRPCClient) makeRequest(endpoint string, req interface{}, resp interface{}) error {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	
	httpResp, err := g.client.Post(g.baseURL+endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer httpResp.Body.Close()
	
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %d", httpResp.StatusCode)
	}
	
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	
	return nil
}
