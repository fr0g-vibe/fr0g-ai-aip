package persona

import (
	"testing"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

func TestServiceCreatePersona(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())
	
	p := types.Persona{
		Name:   "Test Expert",
		Topic:  "Testing",
		Prompt: "You are a testing expert.",
	}
	
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}
	
	if p.ID == "" {
		t.Error("Expected persona ID to be generated")
	}
}

func TestServiceGetPersona(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())
	
	// Create a persona first
	p := types.Persona{
		Name:   "Get Test Expert",
		Topic:  "Getting",
		Prompt: "You are a getting expert.",
	}
	
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}
	
	// Get the persona
	retrieved, err := service.GetPersona(p.ID)
	if err != nil {
		t.Fatalf("Failed to get persona: %v", err)
	}
	
	if retrieved.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, retrieved.Name)
	}
}

func TestServiceListPersonas(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())
	
	// Create test personas
	p1 := types.Persona{Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"}
	p2 := types.Persona{Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"}
	
	service.CreatePersona(&p1)
	service.CreatePersona(&p2)
	
	list, err := service.ListPersonas()
	if err != nil {
		t.Fatalf("Failed to list personas: %v", err)
	}
	
	if len(list) != 2 {
		t.Errorf("Expected 2 personas, got %d", len(list))
	}
}

func TestServiceDeletePersona(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())
	
	p := types.Persona{
		Name:   "Delete Test Expert",
		Topic:  "Deleting",
		Prompt: "You are a deleting expert.",
	}
	
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}
	
	err = service.DeletePersona(p.ID)
	if err != nil {
		t.Fatalf("Failed to delete persona: %v", err)
	}
	
	// Try to get deleted persona
	_, err = service.GetPersona(p.ID)
	if err == nil {
		t.Error("Expected error when getting deleted persona")
	}
}

func TestServiceUpdatePersona(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())
	
	p := types.Persona{
		Name:   "Update Test Expert",
		Topic:  "Updating",
		Prompt: "You are an updating expert.",
	}
	
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}
	
	// Update the persona
	p.Name = "Updated Expert"
	err = service.UpdatePersona(p.ID, p)
	if err != nil {
		t.Fatalf("Failed to update persona: %v", err)
	}
	
	// Get the updated persona
	retrieved, err := service.GetPersona(p.ID)
	if err != nil {
		t.Fatalf("Failed to get updated persona: %v", err)
	}
	
	if retrieved.Name != "Updated Expert" {
		t.Errorf("Expected name 'Updated Expert', got %s", retrieved.Name)
	}
}

// Legacy function tests for backward compatibility
func TestCreatePersona(t *testing.T) {
	p := types.Persona{
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
	p := types.Persona{
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
	// Reset default service for clean test
	defaultService = NewService(storage.NewMemoryStorage())
	
	// Create test personas
	p1 := types.Persona{Name: "Expert 1", Topic: "Topic 1", Prompt: "Prompt 1"}
	p2 := types.Persona{Name: "Expert 2", Topic: "Topic 2", Prompt: "Prompt 2"}
	
	CreatePersona(&p1)
	CreatePersona(&p2)
	
	list := ListPersonas()
	if len(list) != 2 {
		t.Errorf("Expected 2 personas, got %d", len(list))
	}
}

func TestDeletePersona(t *testing.T) {
	// Reset default service for clean test
	defaultService = NewService(storage.NewMemoryStorage())
	
	p := types.Persona{
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
