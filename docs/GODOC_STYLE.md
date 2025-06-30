# Go Documentation Standards

This document outlines the documentation standards used throughout the fr0g-ai-aip codebase, following Go best practices for godoc generation.

## Package Documentation

Every package should have a package comment that describes its purpose:

```go
// Package persona provides AI persona management functionality.
// 
// This package implements the core persona service that handles creation,
// retrieval, updating, and deletion of AI subject matter experts. Each
// persona consists of a name, topic, system prompt, and optional context
// and RAG (Retrieval-Augmented Generation) documents.
//
// Example usage:
//
//	storage := storage.NewMemoryStorage()
//	service := persona.NewService(storage)
//	
//	p := &types.Persona{
//		Name:   "Go Expert",
//		Topic:  "Golang Programming",
//		Prompt: "You are an expert Go programmer...",
//	}
//	
//	err := service.CreatePersona(p)
//	if err != nil {
//		log.Fatal(err)
//	}
package persona
```

## Function Documentation

All exported functions must have documentation comments:

```go
// CreatePersona creates a new AI persona with validation.
//
// The persona must have a non-empty name, topic, and prompt. The function
// will generate a unique ID and store the persona using the configured
// storage backend.
//
// Returns an error if:
//   - persona is nil
//   - required fields are empty or contain only whitespace
//   - storage operation fails
//
// Example:
//
//	p := &types.Persona{
//		Name:   "Security Expert",
//		Topic:  "Cybersecurity", 
//		Prompt: "You are a cybersecurity expert...",
//	}
//	err := service.CreatePersona(p)
func (s *Service) CreatePersona(p *types.Persona) error
```

## Type Documentation

All exported types should be documented:

```go
// Persona represents an AI subject matter expert with specific knowledge
// and capabilities. Each persona consists of a system prompt, topic area,
// and optional context and RAG documents.
//
// The ID field is automatically generated when the persona is created and
// should not be set manually. The Name, Topic, and Prompt fields are
// required and will be validated.
//
// Context provides additional key-value pairs that can be used to customize
// the persona's behavior. RAG contains references to documents that should
// be used for retrieval-augmented generation.
type Persona struct {
	// ID is the unique identifier for this persona (auto-generated)
	ID string `json:"id"`
	
	// Name is the display name for this persona (required, 1-100 chars)
	Name string `json:"name"`
	
	// Topic describes the subject area or domain (required, 1-100 chars)
	Topic string `json:"topic"`
	
	// Prompt is the system prompt for the AI (required, 1-10000 chars)
	Prompt string `json:"prompt"`
	
	// Context provides additional key-value context (optional)
	Context map[string]string `json:"context,omitempty"`
	
	// RAG contains document references for retrieval (optional)
	RAG []string `json:"rag,omitempty"`
}
```

## Interface Documentation

Interfaces should clearly document their contract:

```go
// Storage defines the interface for persona storage backends.
//
// Implementations must be safe for concurrent use and should handle
// validation of input data. All methods should return appropriate
// errors for invalid input or storage failures.
//
// The interface supports both personas and identities, with identities
// being instances of personas with additional demographic attributes.
type Storage interface {
	// Create stores a new persona and generates a unique ID.
	// Returns an error if the persona is invalid or storage fails.
	Create(p *types.Persona) error
	
	// Get retrieves a persona by ID.
	// Returns an error if the persona is not found.
	Get(id string) (types.Persona, error)
	
	// List returns all stored personas.
	// Returns an empty slice if no personas exist.
	List() ([]types.Persona, error)
	
	// Update modifies an existing persona.
	// Returns an error if the persona is not found or invalid.
	Update(id string, p types.Persona) error
	
	// Delete removes a persona by ID.
	// Returns an error if the persona is not found.
	Delete(id string) error
}
```

## Error Documentation

Document error conditions clearly:

```go
// ErrPersonaNotFound is returned when a requested persona does not exist.
var ErrPersonaNotFound = errors.New("persona not found")

// ErrInvalidPersona is returned when persona validation fails.
var ErrInvalidPersona = errors.New("invalid persona")

// ValidationError represents a persona validation error with details
// about which field failed validation and why.
type ValidationError struct {
	Field   string // The field that failed validation
	Message string // Human-readable error message
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error in field %s: %s", e.Field, e.Message)
}
```

## Example Documentation

Include comprehensive examples:

```go
// Example_basicUsage demonstrates basic persona management operations.
func Example_basicUsage() {
	// Create a new service with memory storage
	storage := storage.NewMemoryStorage()
	service := persona.NewService(storage)
	
	// Create a persona
	p := &types.Persona{
		Name:   "Go Expert",
		Topic:  "Golang Programming",
		Prompt: "You are an expert Go programmer with deep knowledge of best practices.",
		Context: map[string]string{
			"experience": "10 years",
			"specialty":  "backend development",
		},
		RAG: []string{
			"go-best-practices.md",
			"effective-go.pdf",
		},
	}
	
	err := service.CreatePersona(p)
	if err != nil {
		log.Fatal(err)
	}
	
	// Retrieve the persona
	retrieved, err := service.GetPersona(p.ID)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Created persona: %s\n", retrieved.Name)
	// Output: Created persona: Go Expert
}

// Example_communityGeneration demonstrates community generation.
func Example_communityGeneration() {
	storage := storage.NewMemoryStorage()
	service := community.NewService(storage)
	
	// First create some personas
	personas := []*types.Persona{
		{Name: "Tech Lead", Topic: "Technology", Prompt: "You are a tech lead..."},
		{Name: "Designer", Topic: "Design", Prompt: "You are a UX designer..."},
		{Name: "PM", Topic: "Product", Prompt: "You are a product manager..."},
	}
	
	for _, p := range personas {
		storage.Create(p)
	}
	
	// Generate a tech community
	config := types.CommunityGenerationConfig{
		PersonaWeights: map[string]float64{
			personas[0].ID: 0.5, // 50% tech leads
			personas[1].ID: 0.3, // 30% designers  
			personas[2].ID: 0.2, // 20% PMs
		},
		AgeDistribution: types.AgeDistribution{
			Mean:   32,
			StdDev: 8,
			MinAge: 25,
			MaxAge: 50,
		},
		PoliticalSpread:    0.6,
		InterestSpread:     0.8,
		SocioeconomicRange: 0.7,
		ActivityLevel:      0.8,
	}
	
	community, err := service.GenerateCommunity(
		config,
		"Tech Startup Team",
		"A diverse tech startup team",
		"professional",
		20,
	)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Generated community with %d members\n", community.Size)
	// Output: Generated community with 20 members
}
```

## Benchmark Documentation

Document performance characteristics:

```go
// BenchmarkCreatePersona measures persona creation performance.
//
// Results on a typical development machine:
//   - Memory storage: ~100ns per operation
//   - File storage: ~1ms per operation
//
// The benchmark creates personas with minimal data to measure
// the overhead of the creation process itself.
func BenchmarkCreatePersona(b *testing.B) {
	service := persona.NewService(storage.NewMemoryStorage())
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := &types.Persona{
			Name:   fmt.Sprintf("Persona %d", i),
			Topic:  "Testing",
			Prompt: "Test prompt",
		}
		service.CreatePersona(p)
	}
}
```

## Documentation Generation

To generate and view documentation:

```bash
# Generate documentation for all packages
go doc -all ./...

# Generate HTML documentation
godoc -http=:6060

# View specific package documentation
go doc ./internal/persona

# View specific function documentation  
go doc ./internal/persona.Service.CreatePersona
```

## Documentation Checklist

For each exported symbol:

- [ ] Has a comment starting with the symbol name
- [ ] Describes what the symbol does
- [ ] Documents parameters and return values
- [ ] Lists possible error conditions
- [ ] Includes usage examples where helpful
- [ ] Uses proper Go comment formatting
- [ ] Follows the "what, not how" principle
- [ ] Is written for the user, not the implementer

## Best Practices

1. **Start with the symbol name**: `// CreatePersona creates...`
2. **Be concise but complete**: Cover the essential information
3. **Use examples**: Show how to use the API correctly
4. **Document errors**: List when and why errors occur
5. **Keep it current**: Update docs when code changes
6. **Use proper formatting**: Follow Go comment conventions
7. **Link related concepts**: Reference other relevant types/functions
8. **Avoid implementation details**: Focus on the public contract

## Tools Integration

The documentation integrates with:

- **godoc**: Standard Go documentation tool
- **pkg.go.dev**: Public package documentation
- **IDE support**: IntelliSense and hover documentation
- **API documentation**: Swagger/OpenAPI generation
- **MCP integration**: Tool and resource descriptions
