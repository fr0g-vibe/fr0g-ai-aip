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
