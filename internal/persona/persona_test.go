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

func TestServiceWithFileStorage(t *testing.T) {
	tmpDir := t.TempDir()
	fileStorage, err := storage.NewFileStorage(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	service := NewService(fileStorage)

	// Test service with file storage backend
	p := types.Persona{
		Name:   "File Storage Expert",
		Topic:  "File Storage",
		Prompt: "You are a file storage expert.",
		Context: map[string]string{
			"backend": "file",
		},
		RAG: []string{"file-doc1", "file-doc2"},
	}

	err = service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona with file storage: %v", err)
	}

	// Verify persistence
	retrieved, err := service.GetPersona(p.ID)
	if err != nil {
		t.Fatalf("Failed to get persona from file storage: %v", err)
	}

	if retrieved.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, retrieved.Name)
	}
	if retrieved.Context["backend"] != "file" {
		t.Error("Expected context to be preserved in file storage")
	}
}

func TestServiceEdgeCases(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())

	// Test with moderately long strings that fit within validation limits
	longString := fmt.Sprintf("%0*d", 400, 0)    // 400 character string of zeros (within 500 char limit for context)
	promptString := fmt.Sprintf("%0*d", 1000, 0) // 1000 character string for prompt (within 10000 char limit)
	ragString := fmt.Sprintf("%0*d", 500, 0)     // 500 character string for RAG (within 1000 char limit)

	p := types.Persona{
		Name:   "Long String Test",
		Topic:  "Long Strings",
		Prompt: promptString,
		Context: map[string]string{
			"long_key": longString,
		},
		RAG: []string{ragString},
	}

	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona with long strings: %v", err)
	}

	// Verify long strings are preserved
	retrieved, err := service.GetPersona(p.ID)
	if err != nil {
		t.Fatalf("Failed to get persona with long strings: %v", err)
	}

	if len(retrieved.Prompt) != 1000 {
		t.Errorf("Expected prompt length 1000, got %d", len(retrieved.Prompt))
	}
}

func TestServiceStorageErrors(t *testing.T) {
	// Create a mock storage that returns errors
	mockStorage := &errorStorage{}
	service := NewService(mockStorage)

	p := types.Persona{Name: "Test", Topic: "Test", Prompt: "Test"}

	// Test create error
	err := service.CreatePersona(&p)
	if err == nil {
		t.Error("Expected error from mock storage")
	}

	// Test get error
	_, err = service.GetPersona("test-id")
	if err == nil {
		t.Error("Expected error from mock storage")
	}

	// Test list error
	_, err = service.ListPersonas()
	if err == nil {
		t.Error("Expected error from mock storage")
	}

	// Test update error
	err = service.UpdatePersona("test-id", p)
	if err == nil {
		t.Error("Expected error from mock storage")
	}

	// Test delete error
	err = service.DeletePersona("test-id")
	if err == nil {
		t.Error("Expected error from mock storage")
	}
}

// Mock storage that always returns errors
type errorStorage struct{}

func (e *errorStorage) Create(p *types.Persona) error {
	return fmt.Errorf("mock create error")
}

func (e *errorStorage) Get(id string) (types.Persona, error) {
	return types.Persona{}, fmt.Errorf("mock get error")
}

func (e *errorStorage) List() ([]types.Persona, error) {
	return nil, fmt.Errorf("mock list error")
}

func (e *errorStorage) Update(id string, p types.Persona) error {
	return fmt.Errorf("mock update error")
}

func (e *errorStorage) Delete(id string) error {
	return fmt.Errorf("mock delete error")
}

// Identity methods for errorStorage mock
func (e *errorStorage) CreateIdentity(i *types.Identity) error {
	return fmt.Errorf("mock create identity error")
}

func (e *errorStorage) GetIdentity(id string) (types.Identity, error) {
	return types.Identity{}, fmt.Errorf("mock get identity error")
}

func (e *errorStorage) ListIdentities(filter *types.IdentityFilter) ([]types.Identity, error) {
	return nil, fmt.Errorf("mock list identities error")
}

func (e *errorStorage) UpdateIdentity(id string, i types.Identity) error {
	return fmt.Errorf("mock update identity error")
}

func (e *errorStorage) DeleteIdentity(id string) error {
	return fmt.Errorf("mock delete identity error")
}

func (e *errorStorage) GetIdentityWithPersona(id string) (types.IdentityWithPersona, error) {
	return types.IdentityWithPersona{}, fmt.Errorf("mock get identity with persona error")
}

func TestServiceNewService(t *testing.T) {
	memStorage := storage.NewMemoryStorage()
	service := NewService(memStorage)

	if service == nil {
		t.Error("Expected non-nil service")
	}

	// Test that service uses the provided storage
	p := types.Persona{Name: "Storage Test", Topic: "Testing", Prompt: "Test"}
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}

	// Verify it was stored in the provided storage
	retrieved, err := memStorage.Get(p.ID)
	if err != nil {
		t.Fatalf("Direct storage get failed: %v", err)
	}

	if retrieved.Name != p.Name {
		t.Errorf("Expected name %s, got %s", p.Name, retrieved.Name)
	}
}

// Test identity functionality
func TestServiceCreateIdentity(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())

	// First create a persona
	p := types.Persona{
		Name:   "Test Persona",
		Topic:  "Test Topic",
		Prompt: "Test prompt",
	}
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}

	// Create an identity based on the persona
	i := types.Identity{
		PersonaID:   p.ID,
		Name:        "Test Identity",
		Description: "Test identity description",
		Tags:        []string{"test", "identity"},
	}

	err = service.CreateIdentity(&i)
	if err != nil {
		t.Fatalf("Failed to create identity: %v", err)
	}

	if i.ID == "" {
		t.Error("Identity ID should be set after creation")
	}

	// Verify the identity was created
	retrieved, err := service.GetIdentity(i.ID)
	if err != nil {
		t.Fatalf("Failed to get identity: %v", err)
	}

	if retrieved.Name != i.Name {
		t.Errorf("Expected name %s, got %s", i.Name, retrieved.Name)
	}
	if retrieved.PersonaID != p.ID {
		t.Errorf("Expected persona ID %s, got %s", p.ID, retrieved.PersonaID)
	}
}

func TestServiceCreateIdentityWithInvalidPersona(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())

	// Try to create an identity with a non-existent persona
	i := types.Identity{
		PersonaID: "non-existent-id",
		Name:      "Test Identity",
	}

	err := service.CreateIdentity(&i)
	if err == nil {
		t.Error("Expected error when creating identity with invalid persona ID")
	}
}

func TestServiceListIdentities(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())

	// Create a persona
	p := types.Persona{
		Name:   "Test Persona",
		Topic:  "Test Topic",
		Prompt: "Test prompt",
	}
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}

	// Create multiple identities
	identities := []types.Identity{
		{PersonaID: p.ID, Name: "Identity 1", Tags: []string{"tag1"}},
		{PersonaID: p.ID, Name: "Identity 2", Tags: []string{"tag2"}},
		{PersonaID: p.ID, Name: "Identity 3", Tags: []string{"tag1", "tag2"}},
	}

	for i := range identities {
		err = service.CreateIdentity(&identities[i])
		if err != nil {
			t.Fatalf("Failed to create identity %d: %v", i, err)
		}
	}

	// Test listing all identities
	allIdentities, err := service.ListIdentities(nil)
	if err != nil {
		t.Fatalf("Failed to list identities: %v", err)
	}

	if len(allIdentities) != 3 {
		t.Errorf("Expected 3 identities, got %d", len(allIdentities))
	}

	// Test filtering by persona ID
	filter := &types.IdentityFilter{PersonaID: p.ID}
	filteredIdentities, err := service.ListIdentities(filter)
	if err != nil {
		t.Fatalf("Failed to list identities with filter: %v", err)
	}

	if len(filteredIdentities) != 3 {
		t.Errorf("Expected 3 identities for persona %s, got %d", p.ID, len(filteredIdentities))
	}

	// Test filtering by tags
	tagFilter := &types.IdentityFilter{Tags: []string{"tag1"}}
	tagFilteredIdentities, err := service.ListIdentities(tagFilter)
	if err != nil {
		t.Fatalf("Failed to list identities with tag filter: %v", err)
	}

	if len(tagFilteredIdentities) != 2 {
		t.Errorf("Expected 2 identities with tag1, got %d", len(tagFilteredIdentities))
	}
}

func TestServiceUpdateIdentity(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())

	// Create a persona
	p := types.Persona{
		Name:   "Test Persona",
		Topic:  "Test Topic",
		Prompt: "Test prompt",
	}
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}

	// Create an identity
	i := types.Identity{
		PersonaID:   p.ID,
		Name:        "Original Name",
		Description: "Original description",
	}
	err = service.CreateIdentity(&i)
	if err != nil {
		t.Fatalf("Failed to create identity: %v", err)
	}

	// Update the identity
	updatedIdentity := i
	updatedIdentity.Name = "Updated Name"
	updatedIdentity.Description = "Updated description"

	err = service.UpdateIdentity(i.ID, updatedIdentity)
	if err != nil {
		t.Fatalf("Failed to update identity: %v", err)
	}

	// Verify the update
	retrieved, err := service.GetIdentity(i.ID)
	if err != nil {
		t.Fatalf("Failed to get updated identity: %v", err)
	}

	if retrieved.Name != "Updated Name" {
		t.Errorf("Expected updated name 'Updated Name', got %s", retrieved.Name)
	}
	if retrieved.Description != "Updated description" {
		t.Errorf("Expected updated description 'Updated description', got %s", retrieved.Description)
	}
}

func TestServiceDeleteIdentity(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())

	// Create a persona
	p := types.Persona{
		Name:   "Test Persona",
		Topic:  "Test Topic",
		Prompt: "Test prompt",
	}
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}

	// Create an identity
	i := types.Identity{
		PersonaID: p.ID,
		Name:      "Test Identity",
	}
	err = service.CreateIdentity(&i)
	if err != nil {
		t.Fatalf("Failed to create identity: %v", err)
	}

	// Delete the identity
	err = service.DeleteIdentity(i.ID)
	if err != nil {
		t.Fatalf("Failed to delete identity: %v", err)
	}

	// Verify it's deleted
	_, err = service.GetIdentity(i.ID)
	if err == nil {
		t.Error("Expected error when getting deleted identity")
	}
}

func TestServiceGetIdentityWithPersona(t *testing.T) {
	service := NewService(storage.NewMemoryStorage())

	// Create a persona
	p := types.Persona{
		Name:    "Test Persona",
		Topic:   "Test Topic",
		Prompt:  "Test prompt",
		Context: map[string]string{"key": "value"},
		RAG:     []string{"doc1", "doc2"},
	}
	err := service.CreatePersona(&p)
	if err != nil {
		t.Fatalf("Failed to create persona: %v", err)
	}

	// Create an identity
	i := types.Identity{
		PersonaID:   p.ID,
		Name:        "Test Identity",
		Description: "Test description",
		Background:  "Test background",
		Tags:        []string{"test"},
	}
	err = service.CreateIdentity(&i)
	if err != nil {
		t.Fatalf("Failed to create identity: %v", err)
	}

	// Get identity with persona
	iwp, err := service.GetIdentityWithPersona(i.ID)
	if err != nil {
		t.Fatalf("Failed to get identity with persona: %v", err)
	}

	// Verify both identity and persona data
	if iwp.Identity.ID != i.ID {
		t.Errorf("Expected identity ID %s, got %s", i.ID, iwp.Identity.ID)
	}
	if iwp.Identity.Name != i.Name {
		t.Errorf("Expected identity name %s, got %s", i.Name, iwp.Identity.Name)
	}
	if iwp.Persona.ID != p.ID {
		t.Errorf("Expected persona ID %s, got %s", p.ID, iwp.Persona.ID)
	}
	if iwp.Persona.Name != p.Name {
		t.Errorf("Expected persona name %s, got %s", p.Name, iwp.Persona.Name)
	}
	if iwp.Persona.Topic != p.Topic {
		t.Errorf("Expected persona topic %s, got %s", p.Topic, iwp.Persona.Topic)
	}
}
