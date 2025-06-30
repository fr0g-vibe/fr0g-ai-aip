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
	
	if unmarshaled.Context == nil {
		unmarshaled.Context = make(map[string]string)
	}
	if unmarshaled.RAG == nil {
		unmarshaled.RAG = make([]string, 0)
	}
	
	if len(unmarshaled.Context) != 0 {
		t.Errorf("Expected empty context, got %v", unmarshaled.Context)
	}
	if len(unmarshaled.RAG) != 0 {
		t.Errorf("Expected empty RAG, got %v", unmarshaled.RAG)
	}
}
