package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fr0g-ai/fr0g-ai-aip/internal/persona"
)

// StartServer starts the HTTP API server
func StartServer(port string) error {
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", healthHandler)
	
	// Persona endpoints
	mux.HandleFunc("/personas", personasHandler)
	mux.HandleFunc("/personas/", personaHandler)
	
	return http.ListenAndServe(":"+port, mux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func personasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		personas := persona.ListPersonas()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(personas)
	case http.MethodPost:
		var p persona.Persona
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if err := persona.CreatePersona(p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func personaHandler(w http.ResponseWriter, r *http.Request) {
	// Extract persona ID from URL path
	id := r.URL.Path[len("/personas/"):]
	if id == "" {
		http.Error(w, "Persona ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		p, err := persona.GetPersona(id)
		if err != nil {
			http.Error(w, "Persona not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	case http.MethodDelete:
		if err := persona.DeletePersona(id); err != nil {
			http.Error(w, "Persona not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
