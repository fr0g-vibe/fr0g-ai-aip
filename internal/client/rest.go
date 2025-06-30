package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// RESTClient implements REST API client for persona and identity service
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

// Persona operations
func (r *RESTClient) Create(p *types.Persona) error {
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

func (r *RESTClient) Get(id string) (types.Persona, error) {
	resp, err := r.client.Get(r.baseURL + "/personas/" + id)
	if err != nil {
		return types.Persona{}, fmt.Errorf("failed to get persona: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Persona{}, fmt.Errorf("persona not found: %s", id)
	}

	var p types.Persona
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return types.Persona{}, fmt.Errorf("failed to decode persona: %v", err)
	}

	return p, nil
}

func (r *RESTClient) List() ([]types.Persona, error) {
	resp, err := r.client.Get(r.baseURL + "/personas")
	if err != nil {
		return nil, fmt.Errorf("failed to list personas: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list personas")
	}

	var personas []types.Persona
	if err := json.NewDecoder(resp.Body).Decode(&personas); err != nil {
		return nil, fmt.Errorf("failed to decode personas: %v", err)
	}

	return personas, nil
}

func (r *RESTClient) Update(id string, p types.Persona) error {
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

// Identity operations
func (r *RESTClient) CreateIdentity(i *types.Identity) error {
	data, err := json.Marshal(i)
	if err != nil {
		return fmt.Errorf("failed to marshal identity: %v", err)
	}

	resp, err := r.client.Post(r.baseURL+"/identities", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create identity: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create identity: %s", string(body))
	}

	return json.NewDecoder(resp.Body).Decode(i)
}

func (r *RESTClient) GetIdentity(id string) (types.Identity, error) {
	resp, err := r.client.Get(r.baseURL + "/identities/" + id)
	if err != nil {
		return types.Identity{}, fmt.Errorf("failed to get identity: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Identity{}, fmt.Errorf("identity not found: %s", id)
	}

	var i types.Identity
	if err := json.NewDecoder(resp.Body).Decode(&i); err != nil {
		return types.Identity{}, fmt.Errorf("failed to decode identity: %v", err)
	}

	return i, nil
}

func (r *RESTClient) ListIdentities(filter *types.IdentityFilter) ([]types.Identity, error) {
	u, err := url.Parse(r.baseURL + "/identities")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Add query parameters for filtering
	if filter != nil {
		q := u.Query()
		if filter.PersonaID != "" {
			q.Set("persona_id", filter.PersonaID)
		}
		if filter.Search != "" {
			q.Set("search", filter.Search)
		}
		if filter.IsActive != nil {
			q.Set("is_active", fmt.Sprintf("%t", *filter.IsActive))
		}
		if len(filter.Tags) > 0 {
			for _, tag := range filter.Tags {
				q.Add("tags", tag)
			}
		}
		u.RawQuery = q.Encode()
	}

	resp, err := r.client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to list identities: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list identities")
	}

	var identities []types.Identity
	if err := json.NewDecoder(resp.Body).Decode(&identities); err != nil {
		return nil, fmt.Errorf("failed to decode identities: %v", err)
	}

	return identities, nil
}

func (r *RESTClient) UpdateIdentity(id string, i types.Identity) error {
	data, err := json.Marshal(i)
	if err != nil {
		return fmt.Errorf("failed to marshal identity: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, r.baseURL+"/identities/"+id, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update identity: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update identity: %s", string(body))
	}

	return nil
}

func (r *RESTClient) DeleteIdentity(id string) error {
	req, err := http.NewRequest(http.MethodDelete, r.baseURL+"/identities/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete identity: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete identity: %s", id)
	}

	return nil
}

func (r *RESTClient) GetIdentityWithPersona(id string) (types.IdentityWithPersona, error) {
	resp, err := r.client.Get(r.baseURL + "/identities/" + id + "/with-persona")
	if err != nil {
		return types.IdentityWithPersona{}, fmt.Errorf("failed to get identity with persona: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.IdentityWithPersona{}, fmt.Errorf("identity not found: %s", id)
	}

	var iwp types.IdentityWithPersona
	if err := json.NewDecoder(resp.Body).Decode(&iwp); err != nil {
		return types.IdentityWithPersona{}, fmt.Errorf("failed to decode identity with persona: %v", err)
	}

	return iwp, nil
}

func (r *RESTClient) Close() error {
	// REST client doesn't need cleanup
	return nil
}
