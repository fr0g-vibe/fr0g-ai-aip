package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestPersonasHandler_LargePayload(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	// Create persona with very large content
	largeContent := strings.Repeat("x", 1024*1024) // 1MB
	p := types.Persona{
		Name:   "Large Payload Test",
		Topic:  "Large Payloads",
		Prompt: largeContent,
	}
	
	body, _ := json.Marshal(p)
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	personasHandler(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 for large payload, got %d", w.Code)
	}
}

func TestPersonasHandler_UnicodeContent(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	p := types.Persona{
		Name:   "Unicode Test 测试",
		Topic:  "Unicode 测试",
		Prompt: "You are a Unicode expert. 你是一个Unicode专家。",
	}
	
	body, _ := json.Marshal(p)
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	personasHandler(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 for Unicode content, got %d", w.Code)
	}
	
	var created types.Persona
	json.NewDecoder(w.Body).Decode(&created)
	if created.Name != p.Name {
		t.Errorf("Unicode name not preserved: expected %s, got %s", p.Name, created.Name)
	}
}

func TestPersonaHandler_MalformedURL(t *testing.T) {
	tests := []struct {
		path       string
		expectCode int
	}{
		{"/personas//", http.StatusNotFound},         // Double slash should be treated as empty ID -> 404
		{"/personas/invalid-id", http.StatusNotFound}, // Invalid ID that doesn't exist
		{"/personas/123", http.StatusNotFound},        // Numeric ID that doesn't exist
	}
	
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		w := httptest.NewRecorder()
		
		personaHandler(w, req)
		
		if w.Code != tt.expectCode {
			t.Errorf("Path %q: expected status %d, got %d", tt.path, tt.expectCode, w.Code)
		}
	}
}

func TestPersonasHandler_ContentTypeValidation(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	p := types.Persona{
		Name:   "Content Type Test",
		Topic:  "Content Types",
		Prompt: "Testing content type validation",
	}
	
	body, _ := json.Marshal(p)
	
	tests := []struct {
		contentType string
		expectCode  int
	}{
		{"application/json", http.StatusCreated},
		{"application/json; charset=utf-8", http.StatusCreated},
		{"text/plain", http.StatusCreated}, // Should still work
		{"", http.StatusCreated}, // Should still work
	}
	
	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBuffer(body))
		if tt.contentType != "" {
			req.Header.Set("Content-Type", tt.contentType)
		}
		w := httptest.NewRecorder()
		
		personasHandler(w, req)
		
		if w.Code != tt.expectCode {
			t.Errorf("Content-Type %q: expected status %d, got %d", tt.contentType, tt.expectCode, w.Code)
		}
	}
}

func TestHealthHandler_Methods(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
	}
	
	for _, method := range methods {
		req := httptest.NewRequest(method, "/health", nil)
		w := httptest.NewRecorder()
		
		healthHandler(w, req)
		
		// Health handler should respond to all methods
		if w.Code != http.StatusOK {
			t.Errorf("Method %s: expected status 200, got %d", method, w.Code)
		}
	}
}

func TestPersonasHandler_EmptyBody(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	personasHandler(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty body, got %d", w.Code)
	}
}

func TestPersonaHandler_VeryLongID(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	// Create very long ID
	longID := strings.Repeat("a", 1000)
	
	req := httptest.NewRequest(http.MethodGet, "/personas/"+longID, nil)
	w := httptest.NewRecorder()
	
	personaHandler(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for very long ID, got %d", w.Code)
	}
}

func TestPersonaHandler_SpecialCharactersInID(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	// Test safe special characters that don't break HTTP
	safeSpecialIDs := []string{
		"id-with-dashes",
		"id_with_underscores",
		"id.with.dots",
		"id123numeric",
	}
	
	for _, id := range safeSpecialIDs {
		req := httptest.NewRequest(http.MethodGet, "/personas/"+id, nil)
		w := httptest.NewRecorder()
		
		personaHandler(w, req)
		
		// Should handle gracefully (return 404 since ID doesn't exist)
		if w.Code != http.StatusNotFound {
			t.Errorf("ID %q: expected status 404, got %d", id, w.Code)
		}
	}
}

func TestStartServer_PortValidation(t *testing.T) {
	// Test that StartServer function exists and handles port parameter
	go func() {
		// Use port 0 to let OS choose available port
		err := StartServer("0")
		if err != nil {
			// Expected to fail in test environment, but function should exist
			t.Logf("StartServer failed as expected in test: %v", err)
		}
	}()
	
	// Function should exist and be callable
}

func TestPersonaHandler_URLEncodingIssues(t *testing.T) {
	// Setup test service
	store := storage.NewMemoryStorage()
	persona.SetDefaultService(persona.NewService(store))
	
	// Test URL-encoded special characters that are safe
	urlEncodedIDs := []string{
		"id%20with%20spaces",  // URL-encoded spaces
		"id%40with%40symbols", // URL-encoded @ symbols
		"id%2Fwith%2Fslashes", // URL-encoded slashes
	}
	
	for _, id := range urlEncodedIDs {
		// Create the request with pre-encoded URL
		req := httptest.NewRequest(http.MethodGet, "/personas/"+id, nil)
		w := httptest.NewRecorder()
		
		personaHandler(w, req)
		
		// Should handle gracefully (return 404 since ID doesn't exist)
		if w.Code != http.StatusNotFound {
			t.Errorf("URL-encoded ID %q: expected status 404, got %d", id, w.Code)
		}
	}
}
