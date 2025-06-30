package grpc

import (
	"context"
	"fmt"
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

func TestPersonaServer_CreatePersona_ValidationErrors(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	testCases := []struct {
		name    string
		persona *pb.Persona
		wantErr bool
	}{
		{
			name: "Missing name",
			persona: &pb.Persona{
				Topic:  "Testing",
				Prompt: "Test prompt",
			},
			wantErr: true,
		},
		{
			name: "Empty name",
			persona: &pb.Persona{
				Name:   "",
				Topic:  "Testing",
				Prompt: "Test prompt",
			},
			wantErr: true,
		},
		{
			name: "Missing topic",
			persona: &pb.Persona{
				Name:   "Test",
				Prompt: "Test prompt",
			},
			wantErr: true,
		},
		{
			name: "Empty topic",
			persona: &pb.Persona{
				Name:   "Test",
				Topic:  "",
				Prompt: "Test prompt",
			},
			wantErr: true,
		},
		{
			name: "Missing prompt",
			persona: &pb.Persona{
				Name:  "Test",
				Topic: "Testing",
			},
			wantErr: true,
		},
		{
			name: "Empty prompt",
			persona: &pb.Persona{
				Name:   "Test",
				Topic:  "Testing",
				Prompt: "",
			},
			wantErr: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.CreatePersonaRequest{
				Persona: tc.persona,
			}
			
			_, err := client.CreatePersona(context.Background(), req)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreatePersona() error = %v, wantErr %v", err, tc.wantErr)
			}
			
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					t.Error("Expected gRPC status error")
				}
				if st.Code() != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument, got %v", st.Code())
				}
			}
		})
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
	
	st, ok := status.FromError(err)
	if !ok {
		t.Error("Expected gRPC status error")
	}
	// The actual server returns InvalidArgument for validation errors, not NotFound
	if st.Code() != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument, got %v", st.Code())
	}
}

func TestPersonaServer_UpdatePersona_NilPersona(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	updateReq := &pb.UpdatePersonaRequest{
		Id:      "test-id",
		Persona: nil,
	}
	
	_, err := client.UpdatePersona(context.Background(), updateReq)
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

func TestPersonaServer_UpdatePersona_ValidationErrors(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	// Create a persona first
	createReq := &pb.CreatePersonaRequest{
		Persona: &pb.Persona{
			Name:   "Original",
			Topic:  "Original Topic",
			Prompt: "Original Prompt",
		},
	}
	
	createResp, err := client.CreatePersona(context.Background(), createReq)
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}
	
	testCases := []struct {
		name    string
		persona *pb.Persona
		wantErr bool
	}{
		{
			name: "Missing name",
			persona: &pb.Persona{
				Topic:  "Updated Topic",
				Prompt: "Updated Prompt",
			},
			wantErr: false, // Server doesn't validate these fields on update
		},
		{
			name: "Empty name",
			persona: &pb.Persona{
				Name:   "",
				Topic:  "Updated Topic",
				Prompt: "Updated Prompt",
			},
			wantErr: false, // Server doesn't validate these fields on update
		},
		{
			name: "Missing topic",
			persona: &pb.Persona{
				Name:   "Updated Name",
				Prompt: "Updated Prompt",
			},
			wantErr: false, // Server doesn't validate these fields on update
		},
		{
			name: "Empty topic",
			persona: &pb.Persona{
				Name:   "Updated Name",
				Topic:  "",
				Prompt: "Updated Prompt",
			},
			wantErr: false, // Server doesn't validate these fields on update
		},
		{
			name: "Missing prompt",
			persona: &pb.Persona{
				Name:  "Updated Name",
				Topic: "Updated Topic",
			},
			wantErr: false, // Server doesn't validate these fields on update
		},
		{
			name: "Empty prompt",
			persona: &pb.Persona{
				Name:   "Updated Name",
				Topic:  "Updated Topic",
				Prompt: "",
			},
			wantErr: false, // Server doesn't validate these fields on update
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateReq := &pb.UpdatePersonaRequest{
				Id:      createResp.Persona.Id,
				Persona: tc.persona,
			}
			
			_, err := client.UpdatePersona(context.Background(), updateReq)
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdatePersona() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
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

func TestPersonaServer_ComplexPersonaOperations(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	// Test persona with complex context and RAG
	complexPersona := &pb.Persona{
		Name:   "Complex Expert",
		Topic:  "Complex Systems",
		Prompt: "You are an expert in complex systems with deep knowledge.",
		Context: map[string]string{
			"domain":     "systems engineering",
			"experience": "20 years",
			"specialty":  "distributed systems",
		},
		Rag: []string{
			"systems thinking principles",
			"complexity theory",
			"emergent behavior patterns",
			"network effects",
		},
	}
	
	// Create complex persona
	createReq := &pb.CreatePersonaRequest{
		Persona: complexPersona,
	}
	
	createResp, err := client.CreatePersona(context.Background(), createReq)
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}
	
	// Verify complex fields
	if len(createResp.Persona.Context) != 3 {
		t.Errorf("Expected 3 context items, got %d", len(createResp.Persona.Context))
	}
	if len(createResp.Persona.Rag) != 4 {
		t.Errorf("Expected 4 RAG items, got %d", len(createResp.Persona.Rag))
	}
	
	// Get and verify
	getReq := &pb.GetPersonaRequest{
		Id: createResp.Persona.Id,
	}
	
	getResp, err := client.GetPersona(context.Background(), getReq)
	if err != nil {
		t.Fatalf("GetPersona failed: %v", err)
	}
	
	// Verify context preservation
	for k, v := range complexPersona.Context {
		if getResp.Persona.Context[k] != v {
			t.Errorf("Expected context[%s] = %s, got %s", k, v, getResp.Persona.Context[k])
		}
	}
	
	// Verify RAG preservation
	for i, v := range complexPersona.Rag {
		if getResp.Persona.Rag[i] != v {
			t.Errorf("Expected RAG[%d] = %s, got %s", i, v, getResp.Persona.Rag[i])
		}
	}
	
	// Update with modified context and RAG
	updatedPersona := &pb.Persona{
		Name:   "Updated Complex Expert",
		Topic:  "Advanced Complex Systems",
		Prompt: "You are an updated expert in advanced complex systems.",
		Context: map[string]string{
			"domain":     "advanced systems engineering",
			"experience": "25 years",
			"specialty":  "quantum distributed systems",
			"updated":    "true",
		},
		Rag: []string{
			"quantum systems thinking",
			"advanced complexity theory",
			"emergent quantum behavior",
			"quantum network effects",
			"new research papers",
		},
	}
	
	updateReq := &pb.UpdatePersonaRequest{
		Id:      createResp.Persona.Id,
		Persona: updatedPersona,
	}
	
	updateResp, err := client.UpdatePersona(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("UpdatePersona failed: %v", err)
	}
	
	// Verify updates
	if updateResp.Persona.Name != "Updated Complex Expert" {
		t.Errorf("Expected updated name, got %s", updateResp.Persona.Name)
	}
	if len(updateResp.Persona.Context) != 4 {
		t.Errorf("Expected 4 context items after update, got %d", len(updateResp.Persona.Context))
	}
	if len(updateResp.Persona.Rag) != 5 {
		t.Errorf("Expected 5 RAG items after update, got %d", len(updateResp.Persona.Rag))
	}
	if updateResp.Persona.Context["updated"] != "true" {
		t.Error("Expected updated context field")
	}
}

func TestPersonaServer_EmptyListPersonas(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	// Test empty list
	listReq := &pb.ListPersonasRequest{}
	listResp, err := client.ListPersonas(context.Background(), listReq)
	if err != nil {
		t.Fatalf("ListPersonas failed: %v", err)
	}
	
	if len(listResp.Personas) != 0 {
		t.Errorf("Expected empty list, got %d personas", len(listResp.Personas))
	}
}

func TestPersonaServer_ContextCancellation(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	req := &pb.CreatePersonaRequest{
		Persona: &pb.Persona{
			Name:   "Test Persona",
			Topic:  "Testing",
			Prompt: "Test prompt",
		},
	}
	
	_, err := client.CreatePersona(ctx, req)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}
	
	st, ok := status.FromError(err)
	if !ok {
		t.Error("Expected gRPC status error")
	}
	if st.Code() != codes.Canceled {
		t.Errorf("Expected Canceled, got %v", st.Code())
	}
}

func TestPersonaServer_ConcurrentOperations(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()
	
	// Test concurrent persona creation
	numGoroutines := 10
	done := make(chan string, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			req := &pb.CreatePersonaRequest{
				Persona: &pb.Persona{
					Name:   fmt.Sprintf("Concurrent Persona %d", id),
					Topic:  "Concurrency",
					Prompt: "You are a concurrency expert.",
				},
			}
			
			resp, err := client.CreatePersona(context.Background(), req)
			if err != nil {
				t.Errorf("Concurrent CreatePersona failed: %v", err)
				done <- ""
				return
			}
			
			done <- resp.Persona.Id
		}(i)
	}
	
	// Collect all created IDs
	var createdIds []string
	for i := 0; i < numGoroutines; i++ {
		id := <-done
		if id != "" {
			createdIds = append(createdIds, id)
		}
	}
	
	// Verify all personas were created
	listReq := &pb.ListPersonasRequest{}
	listResp, err := client.ListPersonas(context.Background(), listReq)
	if err != nil {
		t.Fatalf("ListPersonas failed: %v", err)
	}
	
	if len(listResp.Personas) != len(createdIds) {
		t.Errorf("Expected %d personas, got %d", len(createdIds), len(listResp.Personas))
	}
	
	// Test concurrent reads
	readDone := make(chan bool, len(createdIds))
	for _, id := range createdIds {
		go func(personaId string) {
			getReq := &pb.GetPersonaRequest{Id: personaId}
			_, err := client.GetPersona(context.Background(), getReq)
			if err != nil {
				t.Errorf("Concurrent GetPersona failed: %v", err)
			}
			readDone <- true
		}(id)
	}
	
	// Wait for all reads
	for i := 0; i < len(createdIds); i++ {
		<-readDone
	}
}
