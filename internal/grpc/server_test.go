package grpc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func setupTestService() {
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
}

func TestHandleCreatePersona(t *testing.T) {
	setupTestService()
	
	req := CreatePersonaRequest{
		Persona: &Persona{
			Name:   "gRPC Test",
			Topic:  "gRPC Testing",
			Prompt: "You are a gRPC testing expert.",
		},
	}
	
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/CreatePersona", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handleCreatePersona(w, httpReq)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var resp CreatePersonaResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Persona.Name != req.Persona.Name {
		t.Errorf("Expected name %s, got %s", req.Persona.Name, resp.Persona.Name)
	}
	if resp.Persona.Id == "" {
		t.Error("Expected ID to be generated")
	}
}

func TestHandleCreatePersona_InvalidMethod(t *testing.T) {
	httpReq := httptest.NewRequest(http.MethodGet, "/PersonaService/CreatePersona", nil)
	w := httptest.NewRecorder()
	
	handleCreatePersona(w, httpReq)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleCreatePersona_InvalidJSON(t *testing.T) {
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/CreatePersona", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()
	
	handleCreatePersona(w, httpReq)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleGetPersona(t *testing.T) {
	setupTestService()
	
	// Create a persona first
	p := &types.Persona{
		Name:   "Get Test",
		Topic:  "Getting",
		Prompt: "Test prompt",
	}
	persona.CreatePersona(p)
	
	req := GetPersonaRequest{Id: p.ID}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/GetPersona", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handleGetPersona(w, httpReq)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var resp GetPersonaResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Persona.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, resp.Persona.Name)
	}
}

func TestHandleListPersonas(t *testing.T) {
	setupTestService()
	
	// Create test personas
	p1 := &types.Persona{Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"}
	p2 := &types.Persona{Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"}
	persona.CreatePersona(p1)
	persona.CreatePersona(p2)
	
	req := ListPersonasRequest{}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/ListPersonas", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handleListPersonas(w, httpReq)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var resp ListPersonasResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp.Personas) != 2 {
		t.Errorf("Expected 2 personas, got %d", len(resp.Personas))
	}
}

func TestHandleUpdatePersona(t *testing.T) {
	setupTestService()
	
	// Create a persona first
	p := &types.Persona{
		Name:   "Update Test",
		Topic:  "Updating",
		Prompt: "Test prompt",
	}
	persona.CreatePersona(p)
	
	req := UpdatePersonaRequest{
		Id: p.ID,
		Persona: &Persona{
			Name:   "Updated Name",
			Topic:  p.Topic,
			Prompt: p.Prompt,
		},
	}
	
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/UpdatePersona", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handleUpdatePersona(w, httpReq)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var resp UpdatePersonaResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Persona.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", resp.Persona.Name)
	}
}

func TestHandleDeletePersona(t *testing.T) {
	setupTestService()
	
	// Create a persona first
	p := &types.Persona{
		Name:   "Delete Test",
		Topic:  "Deleting",
		Prompt: "Test prompt",
	}
	persona.CreatePersona(p)
	
	req := DeletePersonaRequest{Id: p.ID}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/DeletePersona", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handleDeletePersona(w, httpReq)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	// Verify persona was deleted
	_, err := persona.GetPersona(p.ID)
	if err == nil {
		t.Error("Expected error when getting deleted persona")
	}
}

func TestHandleGetPersona_NotFound(t *testing.T) {
	setupTestService()
	
	req := GetPersonaRequest{Id: "nonexistent"}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/GetPersona", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handleGetPersona(w, httpReq)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandleUpdatePersona_NotFound(t *testing.T) {
	setupTestService()
	
	req := UpdatePersonaRequest{
		Id: "nonexistent",
		Persona: &Persona{
			Name:   "Updated Name",
			Topic:  "Updated Topic",
			Prompt: "Updated prompt",
		},
	}
	
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/UpdatePersona", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handleUpdatePersona(w, httpReq)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleDeletePersona_NotFound(t *testing.T) {
	setupTestService()
	
	req := DeletePersonaRequest{Id: "nonexistent"}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/DeletePersona", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handleDeletePersona(w, httpReq)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandleCreatePersona_MissingPersona(t *testing.T) {
	req := CreatePersonaRequest{Persona: nil}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/CreatePersona", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handleCreatePersona(w, httpReq)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleUpdatePersona_MissingPersona(t *testing.T) {
	req := UpdatePersonaRequest{
		Id:      "test-id",
		Persona: nil,
	}
	
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/UpdatePersona", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	handleUpdatePersona(w, httpReq)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestStartGRPCServer(t *testing.T) {
	// Test that StartGRPCServer function exists and can be called
	// We can't actually start a server in tests, but we can test the setup
	go func() {
		// This would normally block, so run in goroutine
		StartGRPCServer("0") // Use port 0 for testing (will fail but that's expected)
	}()
	
	// Give it a moment to attempt startup
	// The function should exist and be callable
}

func TestGRPCHandlers_AllMethods(t *testing.T) {
	setupTestService()
	
	// Test all handler methods exist and respond to invalid methods correctly
	endpoints := []string{
		"/PersonaService/CreatePersona",
		"/PersonaService/GetPersona", 
		"/PersonaService/ListPersonas",
		"/PersonaService/UpdatePersona",
		"/PersonaService/DeletePersona",
	}
	
	for _, endpoint := range endpoints {
		req := httptest.NewRequest(http.MethodGet, endpoint, nil) // Wrong method
		w := httptest.NewRecorder()
		
		switch endpoint {
		case "/PersonaService/CreatePersona":
			handleCreatePersona(w, req)
		case "/PersonaService/GetPersona":
			handleGetPersona(w, req)
		case "/PersonaService/ListPersonas":
			handleListPersonas(w, req)
		case "/PersonaService/UpdatePersona":
			handleUpdatePersona(w, req)
		case "/PersonaService/DeletePersona":
			handleDeletePersona(w, req)
		}
		
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Endpoint %s: expected status 405, got %d", endpoint, w.Code)
		}
	}
}

func TestHandleCreatePersona_ValidationErrors(t *testing.T) {
	setupTestService()
	
	tests := []struct {
		name    string
		persona *Persona
		wantErr bool
	}{
		{"missing name", &Persona{Topic: "Test", Prompt: "Test"}, true},
		{"missing topic", &Persona{Name: "Test", Prompt: "Test"}, true},
		{"missing prompt", &Persona{Name: "Test", Topic: "Test"}, true},
		{"valid persona", &Persona{Name: "Test", Topic: "Test", Prompt: "Test"}, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreatePersonaRequest{Persona: tt.persona}
			body, _ := json.Marshal(req)
			httpReq := httptest.NewRequest(http.MethodPost, "/PersonaService/CreatePersona", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			
			handleCreatePersona(w, httpReq)
			
			if tt.wantErr && w.Code == http.StatusOK {
				t.Error("Expected error but got success")
			}
			if !tt.wantErr && w.Code != http.StatusOK {
				t.Errorf("Expected success but got status %d", w.Code)
			}
		})
	}
}
