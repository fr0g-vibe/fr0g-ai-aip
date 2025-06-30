package grpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// Local gRPC-compatible types (using only standard library)
type Persona struct {
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	Topic   string            `json:"topic"`
	Prompt  string            `json:"prompt"`
	Context map[string]string `json:"context,omitempty"`
	Rag     []string          `json:"rag,omitempty"`
}

type CreatePersonaRequest struct {
	Persona *Persona `json:"persona"`
}

type CreatePersonaResponse struct {
	Persona *Persona `json:"persona"`
}

type GetPersonaRequest struct {
	Id string `json:"id"`
}

type GetPersonaResponse struct {
	Persona *Persona `json:"persona"`
}

type ListPersonasRequest struct{}

type ListPersonasResponse struct {
	Personas []*Persona `json:"personas"`
}

type UpdatePersonaRequest struct {
	Id      string   `json:"id"`
	Persona *Persona `json:"persona"`
}

type UpdatePersonaResponse struct {
	Persona *Persona `json:"persona"`
}

type DeletePersonaRequest struct {
	Id string `json:"id"`
}

type DeletePersonaResponse struct{}

// StartGRPCServer starts a JSON-over-HTTP server that mimics gRPC functionality
// This uses only Go standard library to avoid third-party dependency issues
func StartGRPCServer(port string) error {
	mux := http.NewServeMux()
	
	// gRPC-style endpoints using JSON over HTTP
	mux.HandleFunc("/PersonaService/CreatePersona", handleCreatePersona)
	mux.HandleFunc("/PersonaService/GetPersona", handleGetPersona)
	mux.HandleFunc("/PersonaService/ListPersonas", handleListPersonas)
	mux.HandleFunc("/PersonaService/UpdatePersona", handleUpdatePersona)
	mux.HandleFunc("/PersonaService/DeletePersona", handleDeletePersona)
	
	fmt.Printf("Local gRPC-compatible server listening on port %s\n", port)
	fmt.Println("Using JSON-over-HTTP protocol (standard library implementation)")
	
	return http.ListenAndServe(":"+port, mux)
}

func handleCreatePersona(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req CreatePersonaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if req.Persona == nil {
		http.Error(w, "persona is required", http.StatusBadRequest)
		return
	}
	
	p := &types.Persona{
		Name:    req.Persona.Name,
		Topic:   req.Persona.Topic,
		Prompt:  req.Persona.Prompt,
		Context: req.Persona.Context,
		RAG:     req.Persona.Rag,
	}
	
	if err := persona.CreatePersona(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	resp := CreatePersonaResponse{
		Persona: &Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleGetPersona(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req GetPersonaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	p, err := persona.GetPersona(req.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	resp := GetPersonaResponse{
		Persona: &Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleListPersonas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	personas := persona.ListPersonas()
	
	var grpcPersonas []*Persona
	for _, p := range personas {
		grpcPersonas = append(grpcPersonas, &Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		})
	}
	
	resp := ListPersonasResponse{
		Personas: grpcPersonas,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleUpdatePersona(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req UpdatePersonaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if req.Persona == nil {
		http.Error(w, "persona is required", http.StatusBadRequest)
		return
	}
	
	p := types.Persona{
		ID:      req.Id,
		Name:    req.Persona.Name,
		Topic:   req.Persona.Topic,
		Prompt:  req.Persona.Prompt,
		Context: req.Persona.Context,
		RAG:     req.Persona.Rag,
	}
	
	if err := persona.UpdatePersona(req.Id, p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	resp := UpdatePersonaResponse{
		Persona: &Persona{
			Id:      p.ID,
			Name:    p.Name,
			Topic:   p.Topic,
			Prompt:  p.Prompt,
			Context: p.Context,
			Rag:     p.RAG,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleDeletePersona(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req DeletePersonaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if err := persona.DeletePersona(req.Id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	resp := DeletePersonaResponse{}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
