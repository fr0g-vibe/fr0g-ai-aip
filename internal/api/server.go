package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/config"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/middleware"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// Server holds the HTTP server configuration and dependencies
type Server struct {
	config  *config.Config
	service *persona.Service
	server  *http.Server
}

// NewServer creates a new HTTP server instance
func NewServer(cfg *config.Config, service *persona.Service) *Server {
	return &Server{
		config:  cfg,
		service: service,
	}
}

// Start starts the HTTP server with graceful shutdown support
func (s *Server) Start() error {
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", s.healthHandler)
	
	// Persona endpoints
	mux.HandleFunc("/personas", s.personasHandler)
	mux.HandleFunc("/personas/", s.personaHandler)
	
	// Identity endpoints
	mux.HandleFunc("/identities", s.identitiesHandler)
	mux.HandleFunc("/identities/", s.identityHandler)
	
	// Community endpoints
	mux.HandleFunc("/communities", s.communitiesHandler)
	mux.HandleFunc("/communities/", s.communityHandler)
	mux.HandleFunc("/communities/generate", s.generateCommunityHandler)
	
	// Apply middleware
	var handler http.Handler = mux
	
	// Add CORS middleware
	handler = middleware.CORSMiddleware(handler)
	
	// Add authentication middleware if enabled
	if s.config.Security.EnableAuth {
		handler = middleware.AuthMiddleware(s.config.Security.APIKey)(handler)
	}
	
	s.server = &http.Server{
		Addr:         ":" + s.config.HTTP.Port,
		Handler:      handler,
		ReadTimeout:  s.config.HTTP.ReadTimeout,
		WriteTimeout: s.config.HTTP.WriteTimeout,
	}
	
	if s.config.HTTP.EnableTLS {
		return s.server.ListenAndServeTLS(s.config.HTTP.CertFile, s.config.HTTP.KeyFile)
	}
	
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Enhanced health check
	health := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0", // This could be injected at build time
		"storage":   s.config.Storage.Type,
	}
	
	// Check storage health
	if personas, err := s.service.ListPersonas(); err != nil {
		health["status"] = "degraded"
		health["storage_error"] = err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		// Add storage stats
		health["persona_count"] = len(personas)
	}
	
	json.NewEncoder(w).Encode(health)
}

func (s *Server) personasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		personas, err := s.service.ListPersonas()
		if err != nil {
			http.Error(w, "Failed to list personas", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(personas)
		
	case http.MethodPost:
		var p types.Persona
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		if err := s.service.CreatePersona(&p); err != nil {
			// Check if it's a validation error
			if validationErr, ok := err.(middleware.ValidationErrors); ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":   "Validation failed",
					"details": validationErr.Errors,
				})
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) personaHandler(w http.ResponseWriter, r *http.Request) {
	// Extract persona ID from URL path
	id := r.URL.Path[len("/personas/"):]
	if id == "" {
		http.Error(w, "Persona ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		p, err := s.service.GetPersona(id)
		if err != nil {
			http.Error(w, "Persona not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
		
	case http.MethodPut:
		var p types.Persona
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		if err := s.service.UpdatePersona(id, p); err != nil {
			// Check if it's a validation error
			if validationErr, ok := err.(middleware.ValidationErrors); ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":   "Validation failed",
					"details": validationErr.Errors,
				})
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
		
	case http.MethodDelete:
		if err := s.service.DeletePersona(id); err != nil {
			http.Error(w, "Persona not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) identitiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Parse query parameters for filtering
		filter := &types.IdentityFilter{}
		if personaID := r.URL.Query().Get("persona_id"); personaID != "" {
			filter.PersonaID = personaID
		}
		if search := r.URL.Query().Get("search"); search != "" {
			filter.Search = search
		}
		if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
			if isActive := isActiveStr == "true"; isActiveStr == "true" || isActiveStr == "false" {
				filter.IsActive = &isActive
			}
		}
		
		// Get identities from storage
		identities, err := s.service.GetStorage().ListIdentities(filter)
		if err != nil {
			http.Error(w, "Failed to list identities", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(identities)
		
	case http.MethodPost:
		var req struct {
			PersonaID   string                 `json:"persona_id"`
			Name        string                 `json:"name"`
			Description string                 `json:"description"`
			Background  string                 `json:"background"`
			Tags        []string               `json:"tags"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		// TODO: Implement identity creation
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Identity creation not yet implemented"})
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) identityHandler(w http.ResponseWriter, r *http.Request) {
	// Extract identity ID from URL path
	path := r.URL.Path[len("/identities/"):]
	if path == "" {
		http.Error(w, "Identity ID required", http.StatusBadRequest)
		return
	}
	
	// Handle special endpoints
	if path == "with-persona" {
		// TODO: Implement get identity with persona
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Identity with persona not yet implemented"})
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		// TODO: Implement get identity
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Get identity not yet implemented", "id": path})
		
	case http.MethodPut:
		// TODO: Implement update identity
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Update identity not yet implemented", "id": path})
		
	case http.MethodDelete:
		// TODO: Implement delete identity
		w.WriteHeader(http.StatusNoContent)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) communitiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Parse query parameters for filtering
		filter := &types.CommunityFilter{}
		if communityType := r.URL.Query().Get("type"); communityType != "" {
			filter.Type = communityType
		}
		if search := r.URL.Query().Get("search"); search != "" {
			filter.Search = search
		}
		
		// TODO: Implement community service methods
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]types.Community{})
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) communityHandler(w http.ResponseWriter, r *http.Request) {
	// Extract community ID from URL path
	path := r.URL.Path[len("/communities/"):]
	if path == "" {
		http.Error(w, "Community ID required", http.StatusBadRequest)
		return
	}
	
	// Handle special endpoints
	if path == "stats" {
		// TODO: Implement get community stats
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Community stats not yet implemented"})
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		// TODO: Implement get community
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Get community not yet implemented", "id": path})
		
	case http.MethodPut:
		// TODO: Implement update community
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Update community not yet implemented", "id": path})
		
	case http.MethodDelete:
		// TODO: Implement delete community
		w.WriteHeader(http.StatusNoContent)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) generateCommunityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		Name             string                              `json:"name"`
		Description      string                              `json:"description"`
		Type             string                              `json:"type"`
		TargetSize       int                                 `json:"target_size"`
		GenerationConfig types.CommunityGenerationConfig    `json:"generation_config"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// TODO: Implement community generation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Community generation not yet implemented"})
}

// StartServer starts the HTTP API server (legacy function for backward compatibility)
func StartServer(port string) error {
	cfg := &config.Config{
		HTTP: config.HTTPConfig{
			Port:         port,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Security: config.SecurityConfig{
			EnableAuth: false,
		},
	}
	
	// Use legacy global service
	server := NewServer(cfg, nil)
	return server.Start()
}

// StartServerWithConfig starts the HTTP server with full configuration
func StartServerWithConfig(cfg *config.Config, service *persona.Service) error {
	server := NewServer(cfg, service)
	return server.Start()
}
