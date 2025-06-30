package grpc

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/fr0g-vibe/fr0g-ai-aip/internal/grpc/pb"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
)

const bufSize = 1024 * 1024

func setupTestServer(t *testing.T) (pb.PersonaServiceClient, func()) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	
	// Set up clean storage for each test
	memStorage := storage.NewMemoryStorage()
	service := persona.NewService(memStorage)
	persona.SetDefaultService(service)
	
	pb.RegisterPersonaServiceServer(s, &PersonaServer{})
	
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()
	
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	
	conn, err := grpc.DialContext(context.Background(), "bufnet", 
		grpc.WithContextDialer(bufDialer), 
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	
	cleanup := func() {
		conn.Close()
		s.Stop()
		lis.Close()
	}
	
	return pb.NewPersonaServiceClient(conn), cleanup
}

func TestPersonaServer_CreatePersona(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	req := &pb.CreatePersonaRequest{
		Persona: &pb.Persona{
			Name:   "Test Persona",
			Topic:  "Testing",
			Prompt: "You are a test assistant",
			Context: map[string]string{
				"key": "value",
			},
			Rag: []string{"doc1", "doc2"},
		},
	}
	
	resp, err := client.CreatePersona(context.Background(), req)
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}
	
	if resp.Persona.Name != "Test Persona" {
		t.Errorf("Expected name 'Test Persona', got %s", resp.Persona.Name)
	}
	if resp.Persona.Id == "" {
		t.Error("Expected non-empty ID")
	}
	if resp.Persona.Topic != "Testing" {
		t.Errorf("Expected topic 'Testing', got %s", resp.Persona.Topic)
	}
	if resp.Persona.Prompt != "You are a test assistant" {
		t.Errorf("Expected prompt 'You are a test assistant', got %s", resp.Persona.Prompt)
	}
}

func TestPersonaServer_CreatePersona_NilPersona(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	req := &pb.CreatePersonaRequest{
		Persona: nil,
	}
	
	_, err := client.CreatePersona(context.Background(), req)
	if err == nil {
		t.Error("Expected error for nil persona")
	}
	
	st, ok := status.FromError(err)
	if !ok {
		t.Error("Expected gRPC status error")
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument, got %v", st.Code())
	}
}

func TestPersonaServer_GetPersona(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	// First create a persona
	createReq := &pb.CreatePersonaRequest{
		Persona: &pb.Persona{
			Name:   "Get Test",
			Topic:  "Testing",
			Prompt: "Test prompt",
			Context: map[string]string{
				"env": "test",
			},
		},
	}
	
	createResp, err := client.CreatePersona(context.Background(), createReq)
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}
	
	// Now get it
	getReq := &pb.GetPersonaRequest{
		Id: createResp.Persona.Id,
	}
	
	getResp, err := client.GetPersona(context.Background(), getReq)
	if err != nil {
		t.Fatalf("GetPersona failed: %v", err)
	}
	
	if getResp.Persona.Name != "Get Test" {
		t.Errorf("Expected name 'Get Test', got %s", getResp.Persona.Name)
	}
	if getResp.Persona.Id != createResp.Persona.Id {
		t.Errorf("Expected ID %s, got %s", createResp.Persona.Id, getResp.Persona.Id)
	}
}

func TestPersonaServer_GetPersona_NotFound(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	getReq := &pb.GetPersonaRequest{
		Id: "nonexistent",
	}
	
	_, err := client.GetPersona(context.Background(), getReq)
	if err == nil {
		t.Error("Expected error for nonexistent persona")
	}
	
	st, ok := status.FromError(err)
	if !ok {
		t.Error("Expected gRPC status error")
	}
	if st.Code() != codes.NotFound {
		t.Errorf("Expected NotFound, got %v", st.Code())
	}
}

func TestPersonaServer_ListPersonas(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	// Create two personas
	personas := []*pb.Persona{
		{Name: "Persona 1", Topic: "Topic 1", Prompt: "Prompt 1"},
		{Name: "Persona 2", Topic: "Topic 2", Prompt: "Prompt 2"},
	}
	
	var createdIds []string
	for _, p := range personas {
		req := &pb.CreatePersonaRequest{Persona: p}
		resp, err := client.CreatePersona(context.Background(), req)
		if err != nil {
			t.Fatalf("CreatePersona failed: %v", err)
		}
		createdIds = append(createdIds, resp.Persona.Id)
	}
	
	// List all personas
	listReq := &pb.ListPersonasRequest{}
	listResp, err := client.ListPersonas(context.Background(), listReq)
	if err != nil {
		t.Fatalf("ListPersonas failed: %v", err)
	}
	
	if len(listResp.Personas) != 2 {
		t.Errorf("Expected 2 personas, got %d", len(listResp.Personas))
	}
	
	// Verify the personas are correct
	foundNames := make(map[string]bool)
	for _, p := range listResp.Personas {
		foundNames[p.Name] = true
	}
	
	if !foundNames["Persona 1"] || !foundNames["Persona 2"] {
		t.Error("Expected to find both created personas")
	}
}

func TestPersonaServer_UpdatePersona(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	// Create a persona
	createReq := &pb.CreatePersonaRequest{
		Persona: &pb.Persona{
			Name:   "Original Name",
			Topic:  "Original Topic",
			Prompt: "Original Prompt",
		},
	}
	
	createResp, err := client.CreatePersona(context.Background(), createReq)
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}
	
	// Update it
	updateReq := &pb.UpdatePersonaRequest{
		Id: createResp.Persona.Id,
		Persona: &pb.Persona{
			Name:   "Updated Name",
			Topic:  "Updated Topic",
			Prompt: "Updated Prompt",
			Context: map[string]string{
				"updated": "true",
			},
		},
	}
	
	updateResp, err := client.UpdatePersona(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("UpdatePersona failed: %v", err)
	}
	
	if updateResp.Persona.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", updateResp.Persona.Name)
	}
	if updateResp.Persona.Topic != "Updated Topic" {
		t.Errorf("Expected topic 'Updated Topic', got %s", updateResp.Persona.Topic)
	}
	if updateResp.Persona.Prompt != "Updated Prompt" {
		t.Errorf("Expected prompt 'Updated Prompt', got %s", updateResp.Persona.Prompt)
	}
}

func TestPersonaServer_UpdatePersona_NotFound(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	updateReq := &pb.UpdatePersonaRequest{
		Id: "nonexistent",
		Persona: &pb.Persona{
			Name:   "Updated Name",
			Topic:  "Updated Topic",
			Prompt: "Updated Prompt",
		},
	}
	
	_, err := client.UpdatePersona(context.Background(), updateReq)
	if err == nil {
		t.Error("Expected error for nonexistent persona")
	}
}

func TestPersonaServer_DeletePersona(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	// Create a persona
	createReq := &pb.CreatePersonaRequest{
		Persona: &pb.Persona{
			Name:   "To Delete",
			Topic:  "Testing",
			Prompt: "Will be deleted",
		},
	}
	
	createResp, err := client.CreatePersona(context.Background(), createReq)
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}
	
	// Delete it
	deleteReq := &pb.DeletePersonaRequest{
		Id: createResp.Persona.Id,
	}
	
	_, err = client.DeletePersona(context.Background(), deleteReq)
	if err != nil {
		t.Fatalf("DeletePersona failed: %v", err)
	}
	
	// Verify it's gone
	getReq := &pb.GetPersonaRequest{
		Id: createResp.Persona.Id,
	}
	
	_, err = client.GetPersona(context.Background(), getReq)
	if err == nil {
		t.Error("Expected error when getting deleted persona")
	}
	
	st, ok := status.FromError(err)
	if !ok {
		t.Error("Expected gRPC status error")
	}
	if st.Code() != codes.NotFound {
		t.Errorf("Expected NotFound, got %v", st.Code())
	}
}

func TestPersonaServer_DeletePersona_NotFound(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	deleteReq := &pb.DeletePersonaRequest{
		Id: "nonexistent",
	}
	
	_, err := client.DeletePersona(context.Background(), deleteReq)
	if err == nil {
		t.Error("Expected error for nonexistent persona")
	}
	
	st, ok := status.FromError(err)
	if !ok {
		t.Error("Expected gRPC status error")
	}
	if st.Code() != codes.NotFound {
		t.Errorf("Expected NotFound, got %v", st.Code())
	}
}
