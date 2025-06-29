package persona

import (
	"testing"
)

func TestCreatePersona(t *testing.T) {
	p := Persona{
		Name:   "Test Expert",
		Topic:  "Testing",
		Prompt: "You are a testing expert.",
	}
	
	err := CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}
	
	if p.ID == "" {
		t.Error("Expected persona ID to be generated")
	}
}

func TestGetPersona(t *testing.T) {
	// Create a persona first
	p := Persona{
		Name:   "Get Test Expert",
		Topic:  "Getting",
		Prompt: "You are a getting expert.",
	}
	
	err := CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}
	
	// Get the persona
	retrieved, err := GetPersona(p.ID)
	if err != nil {
		t.Fatalf("Failed to get persona: %v", err)
	}
	
	if retrieved.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, retrieved.Name)
	}
}

func TestListPersonas(t *testing.T) {
	// Clear existing personas for clean test
	personas = make(map[string]Persona)
	
	// Create test personas
	p1 := Persona{Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"}
	p2 := Persona{Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"}
	
	CreatePersona(&p1)
	CreatePersona(&p2)
	
	list := ListPersonas()
	if len(list) != 2 {
		t.Errorf("Expected 2 personas, got %d", len(list))
	}
}

func TestDeletePersona(t *testing.T) {
	p := Persona{
		Name:   "Delete Test Expert",
		Topic:  "Deleting",
		Prompt: "You are a deleting expert.",
	}
	
	err := CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}
	
	err = DeletePersona(p.ID)
	if err != nil {
		t.Fatalf("Failed to delete persona: %v", err)
	}
	
	// Try to get deleted persona
	_, err = GetPersona(p.ID)
	if err == nil {
		t.Error("Expected error when getting deleted persona")
	}
}
