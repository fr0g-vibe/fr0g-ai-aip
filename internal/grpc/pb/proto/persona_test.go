package proto

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestPersonaProtoSerialization(t *testing.T) {
	p := &Persona{
		Id:     "test-id",
		Name:   "Test Expert",
		Topic:  "Testing",
		Prompt: "You are a testing expert.",
		Context: map[string]string{
			"language": "Go",
			"domain":   "software testing",
		},
		Rag: []string{
			"testing best practices",
			"Go testing framework",
		},
	}

	// Test protobuf marshaling
	data, err := proto.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal persona: %v", err)
	}

	// Test protobuf unmarshaling
	var unmarshaled Persona
	err = proto.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal persona: %v", err)
	}

	// Verify fields
	if unmarshaled.Id != p.Id {
		t.Errorf("Expected ID %s, got %s", p.Id, unmarshaled.Id)
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
	if len(unmarshaled.Rag) != len(p.Rag) {
		t.Errorf("Expected RAG length %d, got %d", len(p.Rag), len(unmarshaled.Rag))
	}
	for i, v := range p.Rag {
		if unmarshaled.Rag[i] != v {
			t.Errorf("Expected RAG[%d] = %s, got %s", i, v, unmarshaled.Rag[i])
		}
	}
}

func TestCreatePersonaRequest(t *testing.T) {
	req := &CreatePersonaRequest{
		Persona: &Persona{
			Name:   "Test",
			Topic:  "Testing",
			Prompt: "Test prompt",
		},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	var unmarshaled CreatePersonaRequest
	err = proto.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if unmarshaled.Persona.Name != req.Persona.Name {
		t.Errorf("Expected name %s, got %s", req.Persona.Name, unmarshaled.Persona.Name)
	}
}

func TestListPersonasResponse(t *testing.T) {
	resp := &ListPersonasResponse{
		Personas: []*Persona{
			{Id: "1", Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"},
			{Id: "2", Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"},
		},
	}

	data, err := proto.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var unmarshaled ListPersonasResponse
	err = proto.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(unmarshaled.Personas) != 2 {
		t.Errorf("Expected 2 personas, got %d", len(unmarshaled.Personas))
	}

	if unmarshaled.Personas[0].Name != "Expert 1" {
		t.Errorf("Expected name 'Expert 1', got %s", unmarshaled.Personas[0].Name)
	}
}

func TestEmptyPersona(t *testing.T) {
	p := &Persona{}

	data, err := proto.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal empty persona: %v", err)
	}

	var unmarshaled Persona
	err = proto.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty persona: %v", err)
	}

	if unmarshaled.Id != "" {
		t.Errorf("Expected empty ID, got %s", unmarshaled.Id)
	}
	if unmarshaled.Name != "" {
		t.Errorf("Expected empty name, got %s", unmarshaled.Name)
	}
}
