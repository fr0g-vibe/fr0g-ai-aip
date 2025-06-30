package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestRESTClient_Create(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/personas" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		var p types.Persona
		json.NewDecoder(r.Body).Decode(&p)
		p.ID = "test-id"
		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}))
	defer server.Close()
	
	client := NewRESTClient(server.URL)
	p := &types.Persona{
		Name:   "REST Test",
		Topic:  "REST Testing",
		Prompt: "You are a REST testing expert.",
	}
	
	err := client.Create(p)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	if p.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", p.ID)
	}
}

func TestRESTClient_Get(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/personas/test-id" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		p := types.Persona{
			ID:     "test-id",
			Name:   "Test Persona",
			Topic:  "Testing",
			Prompt: "Test prompt",
		}
		
		json.NewEncoder(w).Encode(p)
	}))
	defer server.Close()
	
	client := NewRESTClient(server.URL)
	p, err := client.Get("test-id")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	
	if p.Name != "Test Persona" {
		t.Errorf("Expected name 'Test Persona', got %s", p.Name)
	}
}

func TestRESTClient_List(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/personas" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		personas := []types.Persona{
			{ID: "1", Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"},
			{ID: "2", Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"},
		}
		
		json.NewEncoder(w).Encode(personas)
	}))
	defer server.Close()
	
	client := NewRESTClient(server.URL)
	personas, err := client.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	
	if len(personas) != 2 {
		t.Errorf("Expected 2 personas, got %d", len(personas))
	}
}

func TestRESTClient_Update(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/personas/test-id" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		var p types.Persona
		json.NewDecoder(r.Body).Decode(&p)
		json.NewEncoder(w).Encode(p)
	}))
	defer server.Close()
	
	client := NewRESTClient(server.URL)
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

func TestRESTClient_Delete(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/personas/test-id" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	
	client := NewRESTClient(server.URL)
	err := client.Delete("test-id")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestRESTClient_ErrorHandling(t *testing.T) {
	// Mock server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}))
	defer server.Close()
	
	client := NewRESTClient(server.URL)
	
	// Test Create error
	p := &types.Persona{Name: "Test", Topic: "Test", Prompt: "Test"}
	err := client.Create(p)
	if err == nil {
		t.Error("Expected error for Create")
	}
	
	// Test Get error
	_, err = client.Get("test-id")
	if err == nil {
		t.Error("Expected error for Get")
	}
	
	// Test List error
	_, err = client.List()
	if err == nil {
		t.Error("Expected error for List")
	}
	
	// Test Update error
	err = client.Update("test-id", types.Persona{})
	if err == nil {
		t.Error("Expected error for Update")
	}
	
	// Test Delete error
	err = client.Delete("test-id")
	if err == nil {
		t.Error("Expected error for Delete")
	}
}
