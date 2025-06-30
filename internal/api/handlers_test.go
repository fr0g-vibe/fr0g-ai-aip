package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/config"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func createTestServer() *Server {
	cfg := &config.Config{
		HTTP: config.HTTPConfig{
			Port:         "8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Security: config.SecurityConfig{
			EnableAuth: false,
		},
	}
	
	store := storage.NewMemoryStorage()
	service := persona.NewService(store)
	return NewServer(cfg, service)
}

func TestHealthHandler(t *testing.T) {
	server := createTestServer()
	
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.healthHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	
	if response["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", response["status"])
	}
}

func TestCreatePersona(t *testing.T) {
	server := createTestServer()
	
	persona := types.Persona{
		Name:   "Test Expert",
		Topic:  "Testing",
		Prompt: "You are a testing expert",
	}
	
	jsonData, err := json.Marshal(persona)
	if err != nil {
		t.Fatal(err)
	}
	
	req, err := http.NewRequest("POST", "/personas", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.personasHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	
	var response types.Persona
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	
	if response.Name != persona.Name {
		t.Errorf("expected name %s, got %s", persona.Name, response.Name)
	}
	
	if response.Id == "" {
		t.Error("expected ID to be set")
	}
}

func TestListPersonas(t *testing.T) {
	server := createTestServer()
	
	// Create test personas
	personas := []types.Persona{
		{Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"},
		{Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"},
	}
	
	for _, p := range personas {
		if err := server.service.CreatePersona(&p); err != nil {
			t.Fatal(err)
		}
	}
	
	req, err := http.NewRequest("GET", "/personas", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.personasHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	var response []types.Persona
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	
	if len(response) != 2 {
		t.Errorf("expected 2 personas, got %d", len(response))
	}
}

func TestCreateIdentity(t *testing.T) {
	server := createTestServer()
	
	// Create a persona first
	persona := types.Persona{
		Name:   "Test Expert",
		Topic:  "Testing",
		Prompt: "You are a testing expert",
	}
	if err := server.service.CreatePersona(&persona); err != nil {
		t.Fatal(err)
	}
	
	identity := map[string]interface{}{
		"persona_id":  persona.Id,
		"name":        "Test Identity",
		"description": "A test identity",
		"tags":        []string{"test", "sample"},
	}
	
	jsonData, err := json.Marshal(identity)
	if err != nil {
		t.Fatal(err)
	}
	
	req, err := http.NewRequest("POST", "/identities", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.identitiesHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	
	var response types.Identity
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	
	if response.Name != "Test Identity" {
		t.Errorf("expected name 'Test Identity', got %s", response.Name)
	}
	
	if response.PersonaId != persona.Id {
		t.Errorf("expected persona_id %s, got %s", persona.Id, response.PersonaId)
	}
}

func TestInvalidPersonaCreation(t *testing.T) {
	server := createTestServer()
	
	// Test with missing required fields
	invalidPersona := map[string]interface{}{
		"name": "Test Expert",
		// Missing topic and prompt
	}
	
	jsonData, err := json.Marshal(invalidPersona)
	if err != nil {
		t.Fatal(err)
	}
	
	req, err := http.NewRequest("POST", "/personas", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.personasHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetNonexistentPersona(t *testing.T) {
	server := createTestServer()
	
	req, err := http.NewRequest("GET", "/personas/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.personaHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	server := createTestServer()
	
	req, err := http.NewRequest("PATCH", "/personas", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.personasHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestCommunityGeneration(t *testing.T) {
	server := createTestServer()
	
	// Create a persona first
	persona := types.Persona{
		Name:   "Community Expert",
		Topic:  "Community Building",
		Prompt: "You are a community building expert",
	}
	if err := server.service.CreatePersona(&persona); err != nil {
		t.Fatal(err)
	}
	
	request := map[string]interface{}{
		"name":        "Test Community",
		"description": "A test community",
		"type":        "demographic",
		"target_size": 5,
		"generation_config": map[string]interface{}{
			"persona_weights": map[string]float64{
				persona.Id: 1.0,
			},
			"age_distribution": map[string]interface{}{
				"mean":     35.0,
				"std_dev":  10.0,
				"min_age":  18,
				"max_age":  65,
				"skewness": 0.0,
			},
			"political_spread":    0.8,
			"interest_spread":     0.9,
			"socioeconomic_range": 0.7,
			"activity_level":      0.6,
		},
	}
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	
	req, err := http.NewRequest("POST", "/communities/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.generateCommunityHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	
	var response types.Community
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}
	
	if response.Name != "Test Community" {
		t.Errorf("expected name 'Test Community', got %s", response.Name)
	}
	
	if response.Size != 5 {
		t.Errorf("expected size 5, got %d", response.Size)
	}
	
	if len(response.MemberIds) != 5 {
		t.Errorf("expected 5 member IDs, got %d", len(response.MemberIds))
	}
}
