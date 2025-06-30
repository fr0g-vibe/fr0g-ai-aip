package persona

import (
	"fmt"
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

func TestUpdatePersona(t *testing.T) {
	// Reset default service for clean test
	defaultService = NewService(storage.NewMemoryStorage())
	
	p := types.Persona{
		Name:   "Update Test Expert",
		Topic:  "Updating",
		Prompt: "You are an updating expert.",
	}
	
	err := CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}
	
	// Update the persona
	p.Name = "Updated Expert"
	err = UpdatePersona(p.ID, p)
	if err != nil {
		t.Fatalf("Failed to update persona: %v", err)
	}
	
	// Get the updated persona
	retrieved, err := GetPersona(p.ID)
	if err != nil {
		t.Fatalf("Failed to get updated persona: %v", err)
	}
	
	if retrieved.Name != "Updated Expert" {
		t.Errorf("Expected name 'Updated Expert', got %s", retrieved.Name)
	}
}

func TestServiceValidation(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())
	
	// Test validation errors
	tests := []struct {
		name    string
		persona *types.Persona
		wantErr bool
	}{
		{"nil persona", nil, true},
		{"missing name", &types.Persona{Topic: "Test", Prompt: "Test"}, true},
		{"empty name", &types.Persona{Name: "", Topic: "Test", Prompt: "Test"}, true},
		{"missing topic", &types.Persona{Name: "Test", Prompt: "Test"}, true},
		{"empty topic", &types.Persona{Name: "Test", Topic: "", Prompt: "Test"}, true},
		{"missing prompt", &types.Persona{Name: "Test", Topic: "Test"}, true},
		{"empty prompt", &types.Persona{Name: "Test", Topic: "Test", Prompt: ""}, true},
		{"valid persona", &types.Persona{Name: "Test", Topic: "Test", Prompt: "Test"}, false},
		{"valid with context", &types.Persona{Name: "Test", Topic: "Test", Prompt: "Test", Context: map[string]string{"key": "value"}}, false},
		{"valid with RAG", &types.Persona{Name: "Test", Topic: "Test", Prompt: "Test", RAG: []string{"doc1"}}, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreatePersona(tt.persona)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePersona() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceErrorHandling(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())
	
	// Test get non-existent persona
	_, err := service.GetPersona("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent persona")
	}
	
	// Test update non-existent persona
	p := types.Persona{Name: "Test", Topic: "Test", Prompt: "Test"}
	err = service.UpdatePersona("nonexistent", p)
	if err == nil {
		t.Error("Expected error for updating non-existent persona")
	}
	
	// Test delete non-existent persona
	err = service.DeletePersona("nonexistent")
	if err == nil {
		t.Error("Expected error for deleting non-existent persona")
	}
}

func TestSetDefaultService(t *testing.T) {
	// Test SetDefaultService function
	store := storage.NewMemoryStorage()
	service := NewService(store)
	
	SetDefaultService(service)
	
	// Verify the service was set by using a legacy function
	p := types.Persona{Name: "Default Test", Topic: "Testing", Prompt: "Test"}
	err := CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona with default service: %v", err)
	}
	
	if p.ID == "" {
		t.Error("Expected ID to be generated")
	}
}

func TestLegacyFunctionsWithNilService(t *testing.T) {
	// Reset default service to nil
	defaultService = nil
	
	// Test that legacy functions create a default service when nil
	p := types.Persona{Name: "Nil Test", Topic: "Testing", Prompt: "Test"}
	err := CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona with nil default service: %v", err)
	}
	
	// Verify default service was created
	if defaultService == nil {
		t.Error("Expected default service to be created")
	}
}

func TestServiceWithComplexPersona(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())
	
	// Test persona with context and RAG
	p := types.Persona{
		Name:   "Complex Expert",
		Topic:  "Complex Systems",
		Prompt: "You are an expert in complex systems.",
		Context: map[string]string{
			"domain":     "systems engineering",
			"experience": "20 years",
		},
		RAG: []string{
			"systems thinking principles",
			"complexity theory",
			"emergent behavior patterns",
		},
	}
	
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create complex persona: %v", err)
	}
	
	// Retrieve and verify
	retrieved, err := service.GetPersona(p.ID)
	if err != nil {
		t.Fatalf("Failed to get complex persona: %v", err)
	}
	
	if len(retrieved.Context) != 2 {
		t.Errorf("Expected 2 context items, got %d", len(retrieved.Context))
	}
	if len(retrieved.RAG) != 3 {
		t.Errorf("Expected 3 RAG items, got %d", len(retrieved.RAG))
	}
}


func TestLegacyFunctionsCoverage(t *testing.T) {
	// Reset default service
	defaultService = NewService(storage.NewMemoryStorage())
	
	// Test all legacy functions for complete coverage
	p1 := types.Persona{Name: "Legacy 1", Topic: "Testing", Prompt: "Test"}
	p2 := types.Persona{Name: "Legacy 2", Topic: "Testing", Prompt: "Test"}
	
	// Create multiple personas
	err := CreatePersona(&p1)
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}
	err = CreatePersona(&p2)
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}
	
	// List all
	list := ListPersonas()
	if len(list) != 2 {
		t.Errorf("Expected 2 personas, got %d", len(list))
	}
	
	// Update one
	p1.Name = "Updated Legacy 1"
	err = UpdatePersona(p1.ID, p1)
	if err != nil {
		t.Fatalf("UpdatePersona failed: %v", err)
	}
	
	// Get updated
	retrieved, err := GetPersona(p1.ID)
	if err != nil {
		t.Fatalf("GetPersona failed: %v", err)
	}
	if retrieved.Name != "Updated Legacy 1" {
		t.Errorf("Expected updated name, got %s", retrieved.Name)
	}
	
	// Delete one
	err = DeletePersona(p2.ID)
	if err != nil {
		t.Fatalf("DeletePersona failed: %v", err)
	}
	
	// Verify final count
	finalList := ListPersonas()
	if len(finalList) != 1 {
		t.Errorf("Expected 1 persona after delete, got %d", len(finalList))
	}
}

func TestServiceConcurrentOperations(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())
	
	// Test concurrent persona creation
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(id int) {
			p := types.Persona{
				Name:   fmt.Sprintf("Concurrent Expert %d", id),
				Topic:  "Concurrency",
				Prompt: "You are a concurrency expert.",
			}
			service.CreatePersona(&p)
			done <- true
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
	
	// Verify all personas were created
	personas, err := service.ListPersonas()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(personas) != 5 {
		t.Errorf("Expected 5 personas, got %d", len(personas))
	}
}
