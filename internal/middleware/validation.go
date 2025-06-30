package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, ", ")
}

// ValidatePersona validates a persona struct
func ValidatePersona(p *types.Persona) error {
	var errors []ValidationError
	
	if strings.TrimSpace(p.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "name is required and cannot be empty",
		})
	}
	
	if len(p.Name) > 100 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "name cannot exceed 100 characters",
		})
	}
	
	if strings.TrimSpace(p.Topic) == "" {
		errors = append(errors, ValidationError{
			Field:   "topic",
			Message: "topic is required and cannot be empty",
		})
	}
	
	if len(p.Topic) > 200 {
		errors = append(errors, ValidationError{
			Field:   "topic",
			Message: "topic cannot exceed 200 characters",
		})
	}
	
	if strings.TrimSpace(p.Prompt) == "" {
		errors = append(errors, ValidationError{
			Field:   "prompt",
			Message: "prompt is required and cannot be empty",
		})
	}
	
	if len(p.Prompt) > 10000 {
		errors = append(errors, ValidationError{
			Field:   "prompt",
			Message: "prompt cannot exceed 10000 characters",
		})
	}
	
	// Validate context keys and values
	for key, value := range p.Context {
		if strings.TrimSpace(key) == "" {
			errors = append(errors, ValidationError{
				Field:   "context",
				Message: "context keys cannot be empty",
			})
		}
		if len(key) > 50 {
			errors = append(errors, ValidationError{
				Field:   "context",
				Message: "context keys cannot exceed 50 characters",
			})
		}
		if len(value) > 500 {
			errors = append(errors, ValidationError{
				Field:   "context",
				Message: "context values cannot exceed 500 characters",
			})
		}
	}
	
	// Validate RAG entries
	if len(p.RAG) > 100 {
		errors = append(errors, ValidationError{
			Field:   "rag",
			Message: "cannot have more than 100 RAG entries",
		})
	}
	
	for i, rag := range p.RAG {
		if strings.TrimSpace(rag) == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("rag[%d]", i),
				Message: "RAG entries cannot be empty",
			})
		}
		if len(rag) > 1000 {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("rag[%d]", i),
				Message: "RAG entries cannot exceed 1000 characters",
			})
		}
	}
	
	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}
	
	return nil
}

// ValidationMiddleware returns an HTTP middleware that validates request bodies
func ValidationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only validate POST and PUT requests with JSON content
		if (r.Method == http.MethodPost || r.Method == http.MethodPut) && 
		   strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			
			var p types.Persona
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}
			
			if err := ValidatePersona(&p); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": "Validation failed",
					"details": err,
				})
				return
			}
			
			// Store validated persona in request context for handlers to use
			// For now, we'll let the handler re-decode, but this could be optimized
		}
		
		next.ServeHTTP(w, r)
	}
}

// SanitizePersona sanitizes persona input by trimming whitespace
func SanitizePersona(p *types.Persona) {
	p.Name = strings.TrimSpace(p.Name)
	p.Topic = strings.TrimSpace(p.Topic)
	p.Prompt = strings.TrimSpace(p.Prompt)
	
	// Sanitize context
	for key, value := range p.Context {
		delete(p.Context, key)
		cleanKey := strings.TrimSpace(key)
		cleanValue := strings.TrimSpace(value)
		if cleanKey != "" {
			p.Context[cleanKey] = cleanValue
		}
	}
	
	// Sanitize RAG entries
	var cleanRAG []string
	for _, rag := range p.RAG {
		if clean := strings.TrimSpace(rag); clean != "" {
			cleanRAG = append(cleanRAG, clean)
		}
	}
	p.RAG = cleanRAG
}
