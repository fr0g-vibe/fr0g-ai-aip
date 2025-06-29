package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
)

// RESTClient implements REST API client for persona service
type RESTClient struct {
	baseURL string
	client  *http.Client
}

// NewRESTClient creates a new REST client
func NewRESTClient(baseURL string) *RESTClient {
	return &RESTClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (r *RESTClient) Create(p *persona.Persona) error {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal persona: %v", err)
	}
	
	resp, err := r.client.Post(r.baseURL+"/personas", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create persona: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create persona: %s", string(body))
	}
	
	return json.NewDecoder(resp.Body).Decode(p)
}

func (r *RESTClient) Get(id string) (persona.Persona, error) {
	resp, err := r.client.Get(r.baseURL + "/personas/" + id)
	if err != nil {
		return persona.Persona{}, fmt.Errorf("failed to get persona: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return persona.Persona{}, fmt.Errorf("persona not found: %s", id)
	}
	
	var p persona.Persona
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return persona.Persona{}, fmt.Errorf("failed to decode persona: %v", err)
	}
	
	return p, nil
}

func (r *RESTClient) List() ([]persona.Persona, error) {
	resp, err := r.client.Get(r.baseURL + "/personas")
	if err != nil {
		return nil, fmt.Errorf("failed to list personas: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list personas")
	}
	
	var personas []persona.Persona
	if err := json.NewDecoder(resp.Body).Decode(&personas); err != nil {
		return nil, fmt.Errorf("failed to decode personas: %v", err)
	}
	
	return personas, nil
}

func (r *RESTClient) Update(id string, p persona.Persona) error {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal persona: %v", err)
	}
	
	req, err := http.NewRequest(http.MethodPut, r.baseURL+"/personas/"+id, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update persona: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update persona: %s", string(body))
	}
	
	return nil
}

func (r *RESTClient) Delete(id string) error {
	req, err := http.NewRequest(http.MethodDelete, r.baseURL+"/personas/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete persona: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete persona: %s", id)
	}
	
	return nil
}
