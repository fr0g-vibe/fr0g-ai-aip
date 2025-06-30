package types

import (
	"encoding/json"
	"testing"
)

func TestPersonaJSONSerialization(t *testing.T) {
	p := Persona{
		ID:     "test-id",
		Name:   "Test Expert",
		Topic:  "Testing",
		Prompt: "You are a testing expert.",
		Context: map[string]string{
			"language": "Go",
			"domain":   "software testing",
		},
		RAG: []string{
			"testing best practices",
			"Go testing framework",
		},
	}
	
	// Test JSON marshaling
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal persona: %v", err)
	}
	
	// Test JSON unmarshaling
	var unmarshaled Persona
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal persona: %v", err)
	}
	
	// Verify fields
	if unmarshaled.ID != p.ID {
		t.Errorf("Expected ID %s, got %s", p.ID, unmarshaled.ID)
	}
	if unmarshaled.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, unmarshaled.Name)
	}
	if unmarshaled.Topic != p.Topic {
		t.Errorf("Expected topic %s, got %s", p.Topic, unmarshaled.Topic)
	}
	if unmarshaled.Prompt != p.Prompt {
		t.Errorf("Expected prompt %s, got %s", p.Prompt, unmarshaled.Prompt)
	}
	
	// Verify context
	if len(unmarshaled.Context) != len(p.Context) {
		t.Errorf("Expected context length %d, got %d", len(p.Context), len(unmarshaled.Context))
	}
	for k, v := range p.Context {
		if unmarshaled.Context[k] != v {
			t.Errorf("Expected context[%s] = %s, got %s", k, v, unmarshaled.Context[k])
		}
	}
	
	// Verify RAG
	if len(unmarshaled.RAG) != len(p.RAG) {
		t.Errorf("Expected RAG length %d, got %d", len(p.RAG), len(unmarshaled.RAG))
	}
	for i, v := range p.RAG {
		if unmarshaled.RAG[i] != v {
			t.Errorf("Expected RAG[%d] = %s, got %s", i, v, unmarshaled.RAG[i])
		}
	}
}

func TestPersonaEmptyFields(t *testing.T) {
	p := Persona{
		ID:     "test-id",
		Name:   "Test Expert",
		Topic:  "Testing",
		Prompt: "You are a testing expert.",
	}
	
	// Test with nil context and RAG
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal persona with empty fields: %v", err)
	}
	
	var unmarshaled Persona
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal persona with empty fields: %v", err)
	}
	
	// Verify required fields
	if unmarshaled.ID != p.ID {
		t.Errorf("Expected ID %s, got %s", p.ID, unmarshaled.ID)
	}
	if unmarshaled.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, unmarshaled.Name)
	}
	if unmarshaled.Topic != p.Topic {
		t.Errorf("Expected topic %s, got %s", p.Topic, unmarshaled.Topic)
	}
	if unmarshaled.Prompt != p.Prompt {
		t.Errorf("Expected prompt %s, got %s", p.Prompt, unmarshaled.Prompt)
	}
	
	// Verify optional fields are properly handled
	if unmarshaled.Context != nil && len(unmarshaled.Context) != 0 {
		t.Errorf("Expected empty context, got %v", unmarshaled.Context)
	}
	if unmarshaled.RAG != nil && len(unmarshaled.RAG) != 0 {
		t.Errorf("Expected empty RAG, got %v", unmarshaled.RAG)
	}
}

func TestPersonaJSONEdgeCases(t *testing.T) {
	// Test with special characters in strings
	p := Persona{
		ID:     "special-chars-id",
		Name:   "Test \"Expert\" with 'quotes'",
		Topic:  "Testing\nwith\nnewlines",
		Prompt: "You are a testing expert with unicode: 🚀 and symbols: @#$%",
		Context: map[string]string{
			"key with spaces": "value with\nnewlines",
			"unicode_key":     "unicode value: 🎯",
		},
		RAG: []string{
			"doc with spaces.txt",
			"unicode-doc-🚀.md",
		},
	}
	
	// Test marshaling
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal persona with special chars: %v", err)
	}
	
	// Test unmarshaling
	var unmarshaled Persona
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal persona with special chars: %v", err)
	}
	
	// Verify all fields are preserved
	if unmarshaled.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, unmarshaled.Name)
	}
	if unmarshaled.Topic != p.Topic {
		t.Errorf("Expected topic %s, got %s", p.Topic, unmarshaled.Topic)
	}
	if unmarshaled.Prompt != p.Prompt {
		t.Errorf("Expected prompt %s, got %s", p.Prompt, unmarshaled.Prompt)
	}
	
	// Verify context with special characters
	for k, v := range p.Context {
		if unmarshaled.Context[k] != v {
			t.Errorf("Expected context[%s] = %s, got %s", k, v, unmarshaled.Context[k])
		}
	}
	
	// Verify RAG with special characters
	for i, v := range p.RAG {
		if unmarshaled.RAG[i] != v {
			t.Errorf("Expected RAG[%d] = %s, got %s", i, v, unmarshaled.RAG[i])
		}
	}
}

func TestPersonaJSONInvalidCases(t *testing.T) {
	testCases := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name:     "Invalid JSON syntax",
			jsonData: `{"id": "test", "name": "Test"`,
			wantErr:  true,
		},
		{
			name:     "Missing quotes",
			jsonData: `{id: "test", name: "Test"}`,
			wantErr:  true,
		},
		{
			name:     "Invalid field type",
			jsonData: `{"id": 123, "name": "Test"}`,
			wantErr:  true, // JSON cannot unmarshal number into string field
		},
		{
			name:     "Empty JSON object",
			jsonData: `{}`,
			wantErr:  false,
		},
		{
			name:     "Null values",
			jsonData: `{"id": null, "name": null}`,
			wantErr:  false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var p Persona
			err := json.Unmarshal([]byte(tc.jsonData), &p)
			if (err != nil) != tc.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
