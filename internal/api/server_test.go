package api

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

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	
	healthHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %s", response["status"])
	}
}

func TestPersonasHandler_GET(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	req := httptest.NewRequest(http.MethodGet, "/personas", nil)
	w := httptest.NewRecorder()
	
	personasHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestPersonasHandler_POST(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	p := types.Persona{
		Name:   "API Test",
		Topic:  "API Testing",
		Prompt: "You are an API testing expert.",
	}
	
	body, _ := json.Marshal(p)
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	personasHandler(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestPersonasHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/personas", nil)
	w := httptest.NewRecorder()
	
	personasHandler(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestPersonasHandler_POST_InvalidJSON(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	personasHandler(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPersonaHandler_GET(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	// Create a persona first
	p := &types.Persona{
		Name:   "Handler Test",
		Topic:  "Testing",
		Prompt: "Test prompt",
	}
	persona.CreatePersona(p)
	
	req := httptest.NewRequest(http.MethodGet, "/personas/"+p.ID, nil)
	w := httptest.NewRecorder()
	
	personaHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var retrieved types.Persona
	json.NewDecoder(w.Body).Decode(&retrieved)
	if retrieved.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, retrieved.Name)
	}
}

func TestPersonaHandler_GET_NotFound(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	req := httptest.NewRequest(http.MethodGet, "/personas/nonexistent", nil)
	w := httptest.NewRecorder()
	
	personaHandler(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestPersonaHandler_PUT(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	// Create a persona first
	p := &types.Persona{
		Name:   "Update Test",
		Topic:  "Testing",
		Prompt: "Test prompt",
	}
	persona.CreatePersona(p)
	
	// Update the persona
	p.Name = "Updated Name"
	body, _ := json.Marshal(p)
	req := httptest.NewRequest(http.MethodPut, "/personas/"+p.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	personaHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestPersonaHandler_DELETE(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	// Create a persona first
	p := &types.Persona{
		Name:   "Delete Test",
		Topic:  "Testing",
		Prompt: "Test prompt",
	}
	persona.CreatePersona(p)
	
	req := httptest.NewRequest(http.MethodDelete, "/personas/"+p.ID, nil)
	w := httptest.NewRecorder()
	
	personaHandler(w, req)
	
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestPersonaHandler_EmptyID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/personas/", nil)
	w := httptest.NewRecorder()
	
	personaHandler(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPersonaHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/personas/test-id", nil)
	w := httptest.NewRecorder()
	
	personaHandler(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestPersonaHandler_PUT_InvalidJSON(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	req := httptest.NewRequest(http.MethodPut, "/personas/test-id", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	personaHandler(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPersonasHandler_POST_ValidationError(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	// Create persona with missing required fields
	p := types.Persona{
		Name: "Test", // Missing topic and prompt
	}
	
	body, _ := json.Marshal(p)
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	personasHandler(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPersonaHandler_PUT_NotFound(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	p := types.Persona{
		Name:   "Updated Name",
		Topic:  "Updated Topic",
		Prompt: "Updated prompt",
	}
	
	body, _ := json.Marshal(p)
	req := httptest.NewRequest(http.MethodPut, "/personas/nonexistent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	personaHandler(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
