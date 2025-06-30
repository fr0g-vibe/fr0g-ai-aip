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

func TestStartServer(t *testing.T) {
	// Test that StartServer function exists and can be called
	// We can't actually start a server in tests, but we can test the setup
	go func() {
		// This would normally block, so run in goroutine
		StartServer("0") // Use port 0 for testing (will fail but that's expected)
	}()
	
	// Give it a moment to attempt startup
	// The function should exist and be callable
}

func TestPersonaHandler_URLParsing(t *testing.T) {
	// Test URL parsing edge cases
	tests := []struct {
		path       string
		expectCode int
	}{
		{"/personas/", http.StatusBadRequest},
		{"/personas/valid-id", http.StatusNotFound}, // ID doesn't exist
		{"/personas/123", http.StatusNotFound},      // ID doesn't exist
	}
	
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		w := httptest.NewRecorder()
		
		personaHandler(w, req)
		
		if w.Code != tt.expectCode {
			t.Errorf("Path %s: expected status %d, got %d", tt.path, tt.expectCode, w.Code)
		}
	}
}

func TestPersonasHandler_EmptyList(t *testing.T) {
	// Setup test service with empty storage
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	req := httptest.NewRequest(http.MethodGet, "/personas", nil)
	w := httptest.NewRecorder()
	
	personasHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var personas []types.Persona
	json.NewDecoder(w.Body).Decode(&personas)
	if len(personas) != 0 {
		t.Errorf("Expected empty list, got %d personas", len(personas))
	}
}


func TestPersonaHandler_ComplexPersona(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	// Create a persona with context and RAG
	p := &types.Persona{
		Name:   "Complex Test",
		Topic:  "Complex Testing",
		Prompt: "Test prompt",
		Context: map[string]string{
			"domain": "testing",
			"level":  "expert",
		},
		RAG: []string{
			"testing best practices",
			"test automation",
		},
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
	if len(retrieved.Context) != 2 {
		t.Errorf("Expected 2 context items, got %d", len(retrieved.Context))
	}
	if len(retrieved.RAG) != 2 {
		t.Errorf("Expected 2 RAG items, got %d", len(retrieved.RAG))
	}
}

func TestPersonaHandler_DELETE_NotFound(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	req := httptest.NewRequest(http.MethodDelete, "/personas/nonexistent", nil)
	w := httptest.NewRecorder()
	
	personaHandler(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
