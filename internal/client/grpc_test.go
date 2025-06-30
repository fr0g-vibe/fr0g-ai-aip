package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestGRPCClient_Create(t *testing.T) {
	// Mock gRPC-over-HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/PersonaService/CreatePersona" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		var req CreatePersonaRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		resp := CreatePersonaResponse{
			Persona: &Persona{
				Id:     "grpc-test-id",
				Name:   req.Persona.Name,
				Topic:  req.Persona.Topic,
				Prompt: req.Persona.Prompt,
			},
		}
		
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	// Extract address from server URL
	address := server.URL[7:] // Remove "http://" prefix
	client, err := NewGRPCClient(address)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	
	// Override baseURL to use test server
	client.baseURL = server.URL
	
	p := &types.Persona{
		Name:   "gRPC Test",
		Topic:  "gRPC Testing",
		Prompt: "You are a gRPC testing expert.",
	}
	
	err = client.Create(p)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	if p.ID != "grpc-test-id" {
		t.Errorf("Expected ID 'grpc-test-id', got %s", p.ID)
	}
}

func TestGRPCClient_Get(t *testing.T) {
	// Mock gRPC-over-HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/PersonaService/GetPersona" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		resp := GetPersonaResponse{
			Persona: &Persona{
				Id:     "test-id",
				Name:   "Test Persona",
				Topic:  "Testing",
				Prompt: "Test prompt",
			},
		}
		
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	address := server.URL[7:]
	client, _ := NewGRPCClient(address)
	client.baseURL = server.URL
	
	p, err := client.Get("test-id")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	
	if p.Name != "Test Persona" {
		t.Errorf("Expected name 'Test Persona', got %s", p.Name)
	}
}

func TestGRPCClient_List(t *testing.T) {
	// Mock gRPC-over-HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/PersonaService/ListPersonas" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		resp := ListPersonasResponse{
			Personas: []*Persona{
				{Id: "1", Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"},
				{Id: "2", Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"},
			},
		}
		
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	address := server.URL[7:]
	client, _ := NewGRPCClient(address)
	client.baseURL = server.URL
	
	personas, err := client.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	
	if len(personas) != 2 {
		t.Errorf("Expected 2 personas, got %d", len(personas))
	}
}

func TestGRPCClient_Update(t *testing.T) {
	// Mock gRPC-over-HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/PersonaService/UpdatePersona" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		var req UpdatePersonaRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		resp := UpdatePersonaResponse{
			Persona: req.Persona,
		}
		
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	address := server.URL[7:]
	client, _ := NewGRPCClient(address)
	client.baseURL = server.URL
	
	p := types.Persona{
		ID:     "test-id",
		Name:   "Updated Name",
		Topic:  "Updated Topic",
		Prompt: "Updated prompt",
	}
	
	err := client.Update("test-id", p)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
}

func TestGRPCClient_Delete(t *testing.T) {
	// Mock gRPC-over-HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/PersonaService/DeletePersona" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		resp := DeletePersonaResponse{}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	address := server.URL[7:]
	client, _ := NewGRPCClient(address)
	client.baseURL = server.URL
	
	err := client.Delete("test-id")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestNewGRPCClient(t *testing.T) {
	client, err := NewGRPCClient("localhost:9090")
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	
	if client == nil {
		t.Error("Expected client to be created")
	}
	
	expectedURL := "http://localhost:9090"
	if client.baseURL != expectedURL {
		t.Errorf("Expected baseURL %s, got %s", expectedURL, client.baseURL)
	}
}
