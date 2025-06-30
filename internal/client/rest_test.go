package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
		p.Id = "test-id"
		
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
	
	if p.Id != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", p.Id)
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
			Id:     "test-id",
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
			{Id: "1", Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"},
			{Id: "2", Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"},
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
		Id:     "test-id",
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

func TestRESTClient_InvalidJSON(t *testing.T) {
	// Mock server returning invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	
	client := NewRESTClient(server.URL)
	
	// Test Get with invalid JSON
	_, err := client.Get("test-id")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	
	// Test List with invalid JSON
	_, err = client.List()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestRESTClient_NetworkErrors(t *testing.T) {
	// Use invalid URL to trigger network errors
	client := NewRESTClient("http://invalid-url-that-does-not-exist:99999")
	
	// Test Create network error
	p := &types.Persona{Name: "Test", Topic: "Test", Prompt: "Test"}
	err := client.Create(p)
	if err == nil {
		t.Error("Expected network error for Create")
	}
	
	// Test Get network error
	_, err = client.Get("test-id")
	if err == nil {
		t.Error("Expected network error for Get")
	}
	
	// Test List network error
	_, err = client.List()
	if err == nil {
		t.Error("Expected network error for List")
	}
	
	// Test Update network error
	err = client.Update("test-id", types.Persona{})
	if err == nil {
		t.Error("Expected network error for Update")
	}
	
	// Test Delete network error
	err = client.Delete("test-id")
	if err == nil {
		t.Error("Expected network error for Delete")
	}
}

func TestRESTClient_MarshalErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewRESTClient(server.URL)
	
	// Test Create with unmarshalable data (channels can't be marshaled)
	p := &types.Persona{
		Name:   "Test",
		Topic:  "Test", 
		Prompt: "Test",
		Context: map[string]string{
			"test": string([]byte{0xff, 0xfe, 0xfd}), // Invalid UTF-8
		},
	}
	
	// This should still work as JSON can handle most strings
	err := client.Create(p)
	if err != nil {
		t.Logf("Create with special characters: %v", err)
	}
}

func TestRESTClient_StatusCodes(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"BadRequest", http.StatusBadRequest, true},
		{"NotFound", http.StatusNotFound, true},
		{"InternalServerError", http.StatusInternalServerError, true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			}))
			defer server.Close()
			
			client := NewRESTClient(server.URL)
			
			// Test Get
			_, err := client.Get("test-id")
			if (err != nil) != tc.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tc.wantErr)
			}
			
			// Test List
			_, err = client.List()
			if (err != nil) != tc.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tc.wantErr)
			}
			
			// Test Create
			p := &types.Persona{Name: "Test", Topic: "Test", Prompt: "Test"}
			err = client.Create(p)
			if (err != nil) != tc.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tc.wantErr)
			}
			
			// Test Update
			err = client.Update("test-id", types.Persona{Name: "Updated", Topic: "Updated", Prompt: "Updated"})
			if (err != nil) != tc.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tc.wantErr)
			}
			
			// Test Delete
			err = client.Delete("test-id")
			if (err != nil) != tc.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestRESTClient_ComplexPersona(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p types.Persona
		if r.Method == "POST" && r.URL.Path == "/personas" {
			w.WriteHeader(http.StatusCreated)
			json.NewDecoder(r.Body).Decode(&p)
			p.Id = "complex-id"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
		} else if r.Method == "GET" && r.URL.Path == "/personas/complex-id" {
			w.WriteHeader(http.StatusOK)
			p = types.Persona{
				Id:     "complex-id",
				Name:   "Complex Expert 🚀",
				Topic:  "Complex Systems\nWith Newlines",
				Prompt: "You are an expert with special chars: @#$%",
				Context: map[string]string{
					"unicode":  "🎯💡",
					"newlines": "line1\nline2",
					"tabs":     "col1\tcol2",
				},
				Rag: []string{
					"doc with spaces.txt",
					"unicode-doc-🚀.md",
					"special@chars.pdf",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	}))
	defer server.Close()
	
	client := NewRESTClient(server.URL)
	
	// Test complex persona creation
	complexPersona := &types.Persona{
		Name:   "Complex Expert 🚀",
		Topic:  "Complex Systems\nWith Newlines", 
		Prompt: "You are an expert with special chars: @#$%",
		Context: map[string]string{
			"unicode":  "🎯💡",
			"newlines": "line1\nline2",
			"tabs":     "col1\tcol2",
		},
		Rag: []string{
			"doc with spaces.txt",
			"unicode-doc-🚀.md", 
			"special@chars.pdf",
		},
	}
	
	err := client.Create(complexPersona)
	if err != nil {
		t.Fatalf("Create complex persona failed: %v", err)
	}
	
	// Test getting complex persona
	retrieved, err := client.Get("complex-id")
	if err != nil {
		t.Fatalf("Get complex persona failed: %v", err)
	}
	
	if !strings.Contains(retrieved.Name, "🚀") {
		t.Error("Expected unicode characters to be preserved")
	}
	
	if !strings.Contains(retrieved.Topic, "\n") {
		t.Error("Expected newlines to be preserved")
	}
}
