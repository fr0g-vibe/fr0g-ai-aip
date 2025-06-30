package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/client"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/community"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// Config holds CLI configuration
type Config struct {
	ClientType  string // "local", "rest", "grpc"
	StorageType string // "memory", "file"
	DataDir     string
	ServerURL   string
	Service     interface{} // persona.Service interface
}

var defaultConfig = Config{
	ClientType:  "grpc",
	StorageType: "file", // Changed to file for persistence
	DataDir:     "./data",
	ServerURL:   "localhost:9090", // Default to gRPC
}

// Execute runs the CLI interface
func Execute() error {
	config := GetConfigFromEnv() // Use environment variables
	return ExecuteWithConfig(config)
}

// ExecuteWithConfig runs the CLI interface with the given configuration
func ExecuteWithConfig(config Config) error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	command := os.Args[1]

	// Check for help flags
	if command == "help" || command == "-h" || command == "--help" {
		printUsage()
		return nil
	}

	// Handle generate-identities command (requires direct service access)
	if command == "generate-identities" {
		return handleGenerateIdentities(config)
	}

	// Handle generate-random-community command (requires direct service access)
	if command == "generate-random-community" {
		return handleGenerateRandomCommunity(config)
	}

	// Create client based on configuration
	client, err := createClient(config)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	switch command {
	case "list":
		return listPersonas(client)
	case "create":
		return createPersona(client)
	case "get":
		return getPersona(client)
	case "delete":
		return deletePersona(client)
	case "update":
		return updatePersona(client)
	case "serve":
		return serveCommand()
	// Identity commands
	case "identity-list":
		return listIdentities(client)
	case "identity-create":
		return createIdentity(client)
	case "identity-get":
		return getIdentity(client)
	case "identity-delete":
		return deleteIdentity(client)
	case "identity-update":
		return updateIdentity(client)
	case "identity-get-with-persona":
		return getIdentityWithPersona(client)
	// Generation commands
	case "generate-identity":
		return generateIdentity(client)
	case "generate-community":
		return generateCommunity(client)
	default:
		printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}
}

func handleGenerateIdentities(config Config) error {
	if config.Service == nil {
		return fmt.Errorf("service not available for identity generation")
	}

	// Type assert to get the persona service
	service, ok := config.Service.(*persona.Service)
	if !ok {
		return fmt.Errorf("invalid service type for identity generation")
	}

	fmt.Println("Generating diverse set of sample identities...")

	// First, ensure we have some personas to work with
	personas, err := service.ListPersonas()
	if err != nil {
		return fmt.Errorf("failed to list personas: %v", err)
	}

	if len(personas) == 0 {
		fmt.Println("No personas found. Creating sample personas first...")
		if err := createSamplePersonas(service); err != nil {
			return fmt.Errorf("failed to create sample personas: %v", err)
		}
		personas, err = service.ListPersonas()
		if err != nil {
			return fmt.Errorf("failed to list personas after creation: %v", err)
		}
	}

	// Generate diverse identities
	identities := generateSampleIdentities(personas)
	
	created := 0
	for _, identity := range identities {
		if err := service.CreateIdentity(&identity); err != nil {
			fmt.Printf("Warning: Failed to create identity %s: %v\n", identity.Name, err)
			continue
		}
		created++
		fmt.Printf("Created identity: %s (based on %s persona)\n", identity.Name, getPersonaName(personas, identity.PersonaId))
	}

	fmt.Printf("\nSuccessfully created %d identities out of %d attempted.\n", created, len(identities))
	return nil
}

func handleGenerateRandomCommunity(config Config) error {
	if config.Service == nil {
		return fmt.Errorf("service not available for community generation")
	}

	// Parse command line flags
	fs := flag.NewFlagSet("generate-random-community", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: fr0g-ai-aip generate-random-community -size <number> [-name <name>] [-type <type>] [-location <city>] [-age-range <min>-<max>]")
		fmt.Println("  -size <number>        Number of identities to generate (required)")
		fmt.Println("  -name <name>          Community name (optional)")
		fmt.Println("  -type <type>          Community type (optional: geographic, demographic, interest, political, professional)")
		fmt.Println("  -location <city>      Location constraint (optional)")
		fmt.Println("  -age-range <min>-<max> Age range for members (optional)")
	}
	
	size := fs.Int("size", 0, "Number of identities to generate (required)")
	name := fs.String("name", "", "Community name (optional)")
	communityType := fs.String("type", "demographic", "Community type")
	location := fs.String("location", "", "Location constraint")
	ageRange := fs.String("age-range", "", "Age range (min-max)")

	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	if *size <= 0 {
		fs.Usage()
		return fmt.Errorf("size must be a positive number")
	}

	// Type assert to get the persona service
	service, ok := config.Service.(*persona.Service)
	if !ok {
		return fmt.Errorf("invalid service type for community generation")
	}

	// Get available personas
	personas, err := service.ListPersonas()
	if err != nil {
		return fmt.Errorf("failed to list personas: %v", err)
	}

	if len(personas) == 0 {
		fmt.Println("No personas found. Creating sample personas first...")
		if err := createSamplePersonas(service); err != nil {
			return fmt.Errorf("failed to create sample personas: %v", err)
		}
		personas, err = service.ListPersonas()
		if err != nil {
			return fmt.Errorf("failed to list personas after creation: %v", err)
		}
	}

	// Set default community name if not provided
	if *name == "" {
		*name = fmt.Sprintf("Random Community %d", time.Now().Unix())
	}

	// Parse age range if provided
	var ageDistribution types.AgeDistribution
	if *ageRange != "" {
		parts := strings.Split(*ageRange, "-")
		if len(parts) == 2 {
			min, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			max, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err1 == nil && err2 == nil && min <= max {
				ageDistribution = types.AgeDistribution{
					Mean:    float64(min+max) / 2,
					StdDev:  float64(max-min) / 4, // Reasonable spread
					MinAge:  min,
					MaxAge:  max,
					Skewness: 0,
				}
			}
		}
	} else {
		// Default age distribution
		ageDistribution = types.AgeDistribution{
			Mean:    35,
			StdDev:  12,
			MinAge:  18,
			MaxAge:  75,
			Skewness: 0,
		}
	}

	// Set up location constraint
	var locationConstraint types.LocationConstraint
	if *location != "" {
		locationConstraint = types.LocationConstraint{
			Type:      "city",
			Locations: []string{*location},
			Urban:     func() *bool { b := true; return &b }(), // Default to urban
		}
	} else {
		locationConstraint = types.LocationConstraint{
			Type: "global",
		}
	}

	// Create persona weights (equal distribution)
	personaWeights := make(map[string]float64)
	for _, persona := range personas {
		personaWeights[persona.Id] = 1.0
	}

	// Create generation config
	generationConfig := types.CommunityGenerationConfig{
		PersonaWeights:     personaWeights,
		AgeDistribution:    ageDistribution,
		LocationConstraint: locationConstraint,
		PoliticalSpread:    0.8,  // High political diversity
		InterestSpread:     0.9,  // High interest diversity
		SocioeconomicRange: 0.7,  // Moderate socioeconomic diversity
		ActivityLevel:      0.6,  // Moderate activity level
	}

	fmt.Printf("Generating random community '%s' with %d members...\n", *name, *size)
	fmt.Printf("Community type: %s\n", *communityType)
	if *location != "" {
		fmt.Printf("Location: %s\n", *location)
	}
	if *ageRange != "" {
		fmt.Printf("Age range: %s\n", *ageRange)
	}
	fmt.Println()

	// Create community service
	communityService := community.NewService(service.GetStorage())

	// Generate the community
	generatedCommunity, err := communityService.GenerateCommunity(
		generationConfig,
		*name,
		fmt.Sprintf("Randomly generated community with %d diverse members", *size),
		*communityType,
		*size,
	)
	if err != nil {
		return fmt.Errorf("failed to generate community: %v", err)
	}

	fmt.Printf("✅ Successfully generated community: %s (ID: %s)\n", generatedCommunity.Name, generatedCommunity.Id)
	fmt.Printf("   Members: %d\n", generatedCommunity.Size)
	fmt.Printf("   Diversity: %.2f\n", generatedCommunity.Diversity)
	fmt.Printf("   Cohesion: %.2f\n", generatedCommunity.Cohesion)
	fmt.Printf("   Type: %s\n", generatedCommunity.Type)

	// Show some member details
	fmt.Println("\n📊 Community Members:")
	memberCount := 0
	for _, memberId := range generatedCommunity.MemberIds {
		if memberCount >= 5 { // Show first 5 members
			fmt.Printf("   ... and %d more members\n", len(generatedCommunity.MemberIds)-5)
			break
		}
		
		identity, err := service.GetIdentity(memberId)
		if err != nil {
			continue
		}
		
		persona, err := service.GetPersona(identity.PersonaId)
		if err != nil {
			continue
		}
		
		fmt.Printf("   • %s (based on %s persona)\n", identity.Name, persona.Name)
		memberCount++
	}

	fmt.Printf("\n🎯 Use 'fr0g-ai-aip' with REST/gRPC API to explore community details\n")
	fmt.Printf("   Community ID: %s\n", generatedCommunity.Id)

	return nil
}

func createSamplePersonas(service *persona.Service) error {
	samplePersonas := []types.Persona{
		{
			Name:   "Tech Expert",
			Topic:  "Technology",
			Prompt: "You are a technology expert with deep knowledge of software development, AI, and emerging technologies.",
			Context: map[string]string{
				"experience": "15 years",
				"specialty":  "software architecture",
			},
		},
		{
			Name:   "Healthcare Professional",
			Topic:  "Healthcare",
			Prompt: "You are a healthcare professional with expertise in medical practices, patient care, and health policy.",
			Context: map[string]string{
				"experience": "12 years",
				"specialty":  "primary care",
			},
		},
		{
			Name:   "Education Specialist",
			Topic:  "Education",
			Prompt: "You are an education specialist with knowledge of teaching methods, curriculum development, and student engagement.",
			Context: map[string]string{
				"experience": "10 years",
				"specialty":  "K-12 education",
			},
		},
		{
			Name:   "Business Analyst",
			Topic:  "Business",
			Prompt: "You are a business analyst with expertise in market research, strategy development, and organizational management.",
			Context: map[string]string{
				"experience": "8 years",
				"specialty":  "strategic planning",
			},
		},
	}

	for _, persona := range samplePersonas {
		if err := service.CreatePersona(&persona); err != nil {
			return fmt.Errorf("failed to create persona %s: %v", persona.Name, err)
		}
		fmt.Printf("Created persona: %s\n", persona.Name)
	}

	return nil
}

func generateSampleIdentities(personas []types.Persona) []types.Identity {
	identities := []types.Identity{}

	// Sample names and attributes for diverse identities
	sampleData := []struct {
		name        string
		description string
		age         int32
		gender      string
		political   string
		education   string
		location    string
		interests   []string
	}{
		{
			name:        "Alex Chen",
			description: "Software engineer passionate about AI and machine learning",
			age:         28,
			gender:      "male",
			political:   "liberal",
			education:   "bachelor",
			location:    "San Francisco",
			interests:   []string{"technology", "gaming", "music"},
		},
		{
			name:        "Maria Rodriguez",
			description: "Nurse practitioner focused on community health",
			age:         35,
			gender:      "female",
			political:   "moderate",
			education:   "master",
			location:    "Austin",
			interests:   []string{"healthcare", "fitness", "cooking"},
		},
		{
			name:        "Jordan Smith",
			description: "High school teacher specializing in STEM education",
			age:         42,
			gender:      "non-binary",
			political:   "liberal",
			education:   "master",
			location:    "Portland",
			interests:   []string{"education", "science", "reading"},
		},
		{
			name:        "Robert Johnson",
			description: "Business consultant with expertise in small business development",
			age:         51,
			gender:      "male",
			political:   "conservative",
			education:   "graduate",
			location:    "Dallas",
			interests:   []string{"business", "golf", "travel"},
		},
		{
			name:        "Sarah Kim",
			description: "UX designer focused on accessible technology",
			age:         29,
			gender:      "female",
			political:   "very_liberal",
			education:   "bachelor",
			location:    "Seattle",
			interests:   []string{"design", "technology", "art"},
		},
		{
			name:        "Michael Brown",
			description: "Emergency room physician",
			age:         38,
			gender:      "male",
			political:   "moderate",
			education:   "graduate",
			location:    "Chicago",
			interests:   []string{"medicine", "sports", "photography"},
		},
		{
			name:        "Emily Davis",
			description: "Elementary school principal",
			age:         45,
			gender:      "female",
			political:   "liberal",
			education:   "master",
			location:    "Denver",
			interests:   []string{"education", "gardening", "community"},
		},
		{
			name:        "David Wilson",
			description: "Marketing director for tech startups",
			age:         33,
			gender:      "male",
			political:   "conservative",
			education:   "bachelor",
			location:    "Miami",
			interests:   []string{"marketing", "business", "fitness"},
		},
	}

	// Create identities by cycling through personas
	for i, data := range sampleData {
		persona := personas[i%len(personas)]
		
		identity := types.Identity{
			PersonaId:   persona.Id,
			Name:        data.name,
			Description: data.description,
			Background:  fmt.Sprintf("Generated identity based on %s persona", persona.Name),
			RichAttributes: &types.RichAttributes{
				Demographics: &types.Demographics{
					Age:    data.age,
					Gender: data.gender,
					Location: &types.Location{
						City:       data.location,
						UrbanRural: "urban",
					},
					Education: data.education,
				},
				Psychographics: &types.Psychographics{
					Values: data.interests, // Map interests to values for now
				},
			},
			Tags:     []string{"generated", "sample", persona.Topic},
			IsActive: true,
		}

		identities = append(identities, identity)
	}

	return identities
}

func getPersonaName(personas []types.Persona, personaId string) string {
	for _, persona := range personas {
		if persona.Id == personaId {
			return persona.Name
		}
	}
	return "Unknown"
}

func createClient(config Config) (client.Client, error) {
	switch config.ClientType {
	case "local":
		var store storage.Storage
		var err error

		switch config.StorageType {
		case "memory":
			store = storage.NewMemoryStorage()
		case "file":
			store, err = storage.NewFileStorage(config.DataDir)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown storage type: %s", config.StorageType)
		}

		return client.NewLocalClient(store), nil
	case "rest":
		// Ensure REST URL format
		serverURL := config.ServerURL
		if !strings.HasPrefix(serverURL, "http://") && !strings.HasPrefix(serverURL, "https://") {
			serverURL = "http://" + serverURL
		}
		return client.NewRESTClient(serverURL), nil
	case "grpc":
		// Use gRPC-specific default or extract from config
		address := "localhost:9090"
		if config.ServerURL != "" {
			// Remove http/https prefix if present for gRPC
			address = strings.TrimPrefix(config.ServerURL, "http://")
			address = strings.TrimPrefix(address, "https://")
		}
		grpcClient, err := client.NewGRPCClient(address)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC client for %s: %v\nTip: Make sure the gRPC server is running", address, err)
		}
		return grpcClient, nil
	default:
		return nil, fmt.Errorf("unknown client type: %s", config.ClientType)
	}
}

func printUsage() {
	fmt.Println("fr0g-ai-aip - AI Personas Management")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  A customizable AI subject matter expert system that provides specialized")
	fmt.Println("  AI personas for on-demand expertise in specific topics or domains.")
	fmt.Println("  Now with identity management for creating personalized personas.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  fr0g-ai-aip [command] [options]")
	fmt.Println("  fr0g-ai-aip [flags]")
	fmt.Println()
	fmt.Println("PERSONA COMMANDS:")
	fmt.Println("  list                List all personas")
	fmt.Println("  create              Create a new persona")
	fmt.Println("    -name <name>        Persona name (required)")
	fmt.Println("    -topic <topic>      Persona topic/expertise (required)")
	fmt.Println("    -prompt <prompt>    System prompt (required)")
	fmt.Println("  get <id>            Get persona by ID")
	fmt.Println("  update <id>         Update persona by ID")
	fmt.Println("    -name <name>        Update persona name")
	fmt.Println("    -topic <topic>      Update persona topic")
	fmt.Println("    -prompt <prompt>    Update system prompt")
	fmt.Println("  delete <id>         Delete persona by ID")
	fmt.Println()
	fmt.Println("IDENTITY COMMANDS:")
	fmt.Println("  identity-list       List all identities")
	fmt.Println("  identity-create     Create a new identity")
	fmt.Println("    -persona-id <id>    Persona ID (required)")
	fmt.Println("    -name <name>        Identity name (required)")
	fmt.Println("    -description <desc> Identity description")
	fmt.Println("    -tags <tag1,tag2>   Comma-separated tags")
	fmt.Println("    -background <story> Personal background story")
	fmt.Println("  identity-get <id>   Get identity by ID")
	fmt.Println("  identity-update <id> Update identity by ID")
	fmt.Println("    -name <name>        Update identity name")
	fmt.Println("    -description <desc> Update identity description")
	fmt.Println("    -tags <tag1,tag2>   Update comma-separated tags")
	fmt.Println("    -background <story> Update personal background story")
	fmt.Println("    -active <true|false> Update active status")
	fmt.Println("  identity-delete <id> Delete identity by ID")
	fmt.Println("  identity-get-with-persona <id> Get identity with associated persona")
	fmt.Println()
	fmt.Println("GENERATION COMMANDS:")
	fmt.Println("  generate-identities   Generate a diverse set of sample identities")
	fmt.Println("  generate-random-community Generate a random community with specified size")
	fmt.Println("    -size <number>        Number of identities to generate (required)")
	fmt.Println("    -name <name>          Community name (optional)")
	fmt.Println("    -type <type>          Community type (optional: geographic, demographic, interest, political, professional)")
	fmt.Println("    -location <city>      Location constraint (optional)")
	fmt.Println("    -age-range <min>-<max> Age range for members (optional)")
	fmt.Println("  generate-identity     Generate a random identity based on a persona")
	fmt.Println("    -persona-id <id>      Persona ID (required)")
	fmt.Println("    -name <name>          Identity name (optional, auto-generated if not provided)")
	fmt.Println("  generate-community     Generate a community of identities")
	fmt.Println("    -persona-id <id>      Persona ID (required)")
	fmt.Println("    -size <number>        Number of identities to generate (required)")
	fmt.Println("    -location <city,country> Location for the community (optional)")
	fmt.Println("    -age-range <min>-<max>   Age range for the community (optional)")
	fmt.Println()
	fmt.Println("SERVER COMMANDS:")
	fmt.Println("  serve               Start gRPC server")
	fmt.Println("  help                Show this help message")
	fmt.Println()
	fmt.Println("SERVER FLAGS:")
	fmt.Println("  -server             Run HTTP REST API server")
	fmt.Println("  -grpc               Run gRPC server")
	fmt.Println("  -port <port>        HTTP server port (default: 8080)")
	fmt.Println("  -grpc-port <port>   gRPC server port (default: 9090)")
	fmt.Println("  -storage <type>     Storage type: memory, file (default: file)")
	fmt.Println("  -data-dir <dir>     Data directory for file storage (default: ./data)")
	fmt.Println("  -help               Show help")
	fmt.Println()
	fmt.Println("ENVIRONMENT VARIABLES:")
	fmt.Println("  FR0G_CLIENT_TYPE    Client type: local, rest, grpc (default: grpc)")
	fmt.Println("                      - local: Use local file/memory storage directly")
	fmt.Println("                      - rest: Connect to HTTP REST API server")
	fmt.Println("                      - grpc: Connect to gRPC server")
	fmt.Println()
	fmt.Println("  FR0G_STORAGE_TYPE   Storage type: memory, file (default: file)")
	fmt.Println("                      - memory: Store personas in memory (not persistent)")
	fmt.Println("                      - file: Store personas in JSON files (persistent)")
	fmt.Println()
	fmt.Println("  FR0G_DATA_DIR       Data directory for file storage (default: ./data)")
	fmt.Println("                      Directory where persona JSON files are stored")
	fmt.Println()
	fmt.Println("  FR0G_SERVER_URL     Server URL for remote clients")
	fmt.Println("                      - REST: http://localhost:8080 (default)")
	fmt.Println("                      - gRPC: localhost:9090 (default)")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Show help")
	fmt.Println("  fr0g-ai-aip")
	fmt.Println("  fr0g-ai-aip help")
	fmt.Println()
	fmt.Println("  # Create a persona")
	fmt.Println("  fr0g-ai-aip create -name \"Go Expert\" -topic \"Golang Programming\" \\")
	fmt.Println("    -prompt \"You are an expert Go programmer with deep knowledge...\"")
	fmt.Println()
	fmt.Println("  # List all personas")
	fmt.Println("  fr0g-ai-aip list")
	fmt.Println()
	fmt.Println("  # Create an identity based on a persona")
	fmt.Println("  fr0g-ai-aip identity-create -persona-id <persona_id> -name \"John Doe\" \\")
	fmt.Println("    -description \"Senior Go developer with 10 years experience\" \\")
	fmt.Println("    -tags \"senior,backend,golang\" -background \"Started programming at age 12...\"")
	fmt.Println()
	fmt.Println("  # List all identities")
	fmt.Println("  fr0g-ai-aip identity-list")
	fmt.Println()
	fmt.Println("  # Get identity with its associated persona")
	fmt.Println("  fr0g-ai-aip identity-get-with-persona <identity_id>")
	fmt.Println()
	fmt.Println("  # Generate a diverse set of sample identities")
	fmt.Println("  fr0g-ai-aip generate-identities")
	fmt.Println()
	fmt.Println("  # Generate a random community with 15 members")
	fmt.Println("  fr0g-ai-aip generate-random-community -size 15 -name \"Tech Community\" \\")
	fmt.Println("    -type \"professional\" -location \"San Francisco\" -age-range 25-45")
	fmt.Println()
	fmt.Println("  # Generate a random identity")
	fmt.Println("  fr0g-ai-aip generate-identity -persona-id <persona_id> -name \"Alex Chen\"")
	fmt.Println()
	fmt.Println("  # Generate a community of 10 identities in San Francisco")
	fmt.Println("  fr0g-ai-aip generate-community -persona-id <persona_id> -size 10 \\")
	fmt.Println("    -location \"San Francisco,United States\" -age-range 25-65")
	fmt.Println()
	fmt.Println("  # Use local file storage")
	fmt.Println("  FR0G_CLIENT_TYPE=local FR0G_STORAGE_TYPE=file fr0g-ai-aip list")
	fmt.Println()
	fmt.Println("  # Connect to gRPC server")
	fmt.Println("  FR0G_CLIENT_TYPE=grpc FR0G_SERVER_URL=localhost:9090 fr0g-ai-aip list")
	fmt.Println()
	fmt.Println("  # Start gRPC server")
	fmt.Println("  fr0g-ai-aip -grpc")
	fmt.Println()
	fmt.Println("  # Start HTTP REST server with file storage")
	fmt.Println("  fr0g-ai-aip -server -storage file -data-dir ./personas")
}

func listPersonas(c client.Client) error {
	personas, err := c.List()
	if err != nil {
		return err
	}

	if len(personas) == 0 {
		fmt.Println("No personas found")
		return nil
	}

	fmt.Println("Personas:")
	for _, p := range personas {
		fmt.Printf("  ID: %s, Name: %s, Topic: %s\n", p.Id, p.Name, p.Topic)
	}
	return nil
}

func createPersona(c client.Client) error {
	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: fr0g-ai-aip create -name <name> -topic <topic> -prompt <prompt>")
	}
	name := fs.String("name", "", "Persona name")
	topic := fs.String("topic", "", "Persona topic/expertise")
	prompt := fs.String("prompt", "", "System prompt")

	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	if *name == "" || *topic == "" || *prompt == "" {
		fs.Usage()
		return fmt.Errorf("missing required parameters")
	}

	p := types.Persona{
		Name:   *name,
		Topic:  *topic,
		Prompt: *prompt,
	}

	if err := c.Create(&p); err != nil {
		return err
	}

	fmt.Printf("Created persona: %s (ID: %s)\n", p.Name, p.Id)
	return nil
}

func getPersona(c client.Client) error {
	if len(os.Args) < 3 {
		fmt.Println("Usage: fr0g-ai-aip get <id>")
		return fmt.Errorf("persona ID required")
	}

	id := os.Args[2]
	p, err := c.Get(id)
	if err != nil {
		return err
	}

	fmt.Printf("ID: %s\n", p.Id)
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Topic: %s\n", p.Topic)
	fmt.Printf("Prompt: %s\n", p.Prompt)
	if len(p.Context) > 0 {
		fmt.Println("Context:")
		for k, v := range p.Context {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
	if len(p.Rag) > 0 {
		fmt.Println("RAG:")
		for _, r := range p.Rag {
			fmt.Printf("  %s\n", r)
		}
	}
	return nil
}

func updatePersona(c client.Client) error {
	if len(os.Args) < 3 {
		fmt.Println("Usage: fr0g-ai-aip update <id> -name <name> -topic <topic> -prompt <prompt>")
		return fmt.Errorf("persona ID required")
	}

	id := os.Args[2]

	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: fr0g-ai-aip update <id> -name <name> -topic <topic> -prompt <prompt>")
	}
	name := fs.String("name", "", "Persona name")
	topic := fs.String("topic", "", "Persona topic/expertise")
	prompt := fs.String("prompt", "", "System prompt")

	if err := fs.Parse(os.Args[3:]); err != nil {
		return err
	}

	// Get existing persona first
	existing, err := c.Get(id)
	if err != nil {
		return err
	}

	// Update only provided fields
	if *name != "" {
		existing.Name = *name
	}
	if *topic != "" {
		existing.Topic = *topic
	}
	if *prompt != "" {
		existing.Prompt = *prompt
	}

	if err := c.Update(id, existing); err != nil {
		return err
	}

	fmt.Printf("Updated persona: %s\n", id)
	return nil
}

func deletePersona(c client.Client) error {
	if len(os.Args) < 3 {
		fmt.Println("Usage: fr0g-ai-aip delete <id>")
		return fmt.Errorf("persona ID required")
	}

	id := os.Args[2]
	if err := c.Delete(id); err != nil {
		return err
	}

	fmt.Printf("Deleted persona: %s\n", id)
	return nil
}

// GetConfigFromEnv reads configuration from environment variables
func GetConfigFromEnv() Config {
	config := defaultConfig

	if clientType := os.Getenv("FR0G_CLIENT_TYPE"); clientType != "" {
		config.ClientType = clientType
	}
	if storageType := os.Getenv("FR0G_STORAGE_TYPE"); storageType != "" {
		config.StorageType = storageType
	}
	if dataDir := os.Getenv("FR0G_DATA_DIR"); dataDir != "" {
		config.DataDir = dataDir
	}
	if serverURL := os.Getenv("FR0G_SERVER_URL"); serverURL != "" {
		config.ServerURL = serverURL
	}

	// Expand relative paths
	if !filepath.IsAbs(config.DataDir) {
		if abs, err := filepath.Abs(config.DataDir); err == nil {
			config.DataDir = abs
		}
	}

	return config
}

func serveCommand() error {
	fmt.Println("Starting gRPC server on port 9090...")
	fmt.Println("Use Ctrl+C to stop the server")

	// Import the grpc package here to avoid circular imports
	// We'll need to refactor this properly
	return fmt.Errorf("serve command not yet implemented - use './bin/fr0g-ai-aip -grpc' instead")
}

// Identity management functions
func listIdentities(c client.Client) error {
	identities, err := c.ListIdentities(nil)
	if err != nil {
		return err
	}

	if len(identities) == 0 {
		fmt.Println("No identities found")
		return nil
	}

	fmt.Println("Identities:")
	for _, i := range identities {
		status := "Active"
		if !i.IsActive {
			status = "Inactive"
		}
		fmt.Printf("  ID: %s, Name: %s, Persona: %s, Status: %s\n",
			i.Id, i.Name, i.PersonaId, status)
		if i.Description != "" {
			fmt.Printf("    Description: %s\n", i.Description)
		}
		if len(i.Tags) > 0 {
			fmt.Printf("    Tags: %s\n", strings.Join(i.Tags, ", "))
		}
	}
	return nil
}

func createIdentity(c client.Client) error {
	fs := flag.NewFlagSet("identity-create", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: fr0g-ai-aip identity-create -persona-id <id> -name <name> [-description <desc>] [-tags <tag1,tag2>]")
	}
	personaID := fs.String("persona-id", "", "Persona ID (required)")
	name := fs.String("name", "", "Identity name (required)")
	description := fs.String("description", "", "Identity description")
	tags := fs.String("tags", "", "Comma-separated tags")

	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	if *personaID == "" || *name == "" {
		fs.Usage()
		return fmt.Errorf("persona-id and name are required")
	}

	// Parse tags
	var tagList []string
	if *tags != "" {
		tagList = strings.Split(*tags, ",")
		for i, tag := range tagList {
			tagList[i] = strings.TrimSpace(tag)
		}
	}

	i := types.Identity{
		PersonaId:   *personaID,
		Name:        *name,
		Description: *description,
		Tags:        tagList,
		IsActive:    true,
	}

	if err := c.CreateIdentity(&i); err != nil {
		return fmt.Errorf("failed to create identity: %v", err)
	}

	fmt.Printf("Identity created successfully: %s\n", i.Id)
	return nil
}

func getIdentity(c client.Client) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("identity ID is required")
	}

	id := os.Args[2]
	i, err := c.GetIdentity(id)
	if err != nil {
		return err
	}

	fmt.Printf("Identity: %s\n", i.Id)
	fmt.Printf("  Name: %s\n", i.Name)
	fmt.Printf("  Persona ID: %s\n", i.PersonaId)
	fmt.Printf("  Description: %s\n", i.Description)
	fmt.Printf("  Status: %s\n", map[bool]string{true: "Active", false: "Inactive"}[i.IsActive])
	fmt.Printf("  Created: %s\n", i.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("  Updated: %s\n", i.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"))

	if len(i.Tags) > 0 {
		fmt.Printf("  Tags: %s\n", strings.Join(i.Tags, ", "))
	}

	if i.RichAttributes != nil {
		fmt.Printf("  Rich Attributes: Available\n")
		if i.RichAttributes.Demographics != nil {
			fmt.Printf("    Age: %d\n", i.RichAttributes.Demographics.Age)
			fmt.Printf("    Gender: %s\n", i.RichAttributes.Demographics.Gender)
		}
	}

	return nil
}

func updateIdentity(c client.Client) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("identity ID is required")
	}

	id := os.Args[2]

	// Get existing identity
	existing, err := c.GetIdentity(id)
	if err != nil {
		return err
	}

	fs := flag.NewFlagSet("identity-update", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: fr0g-ai-aip identity-update <id> [-name <name>] [-description <desc>] [-tags <tag1,tag2>] [-background <story>] [-active <true|false>]")
	}
	name := fs.String("name", "", "Identity name")
	description := fs.String("description", "", "Identity description")
	tags := fs.String("tags", "", "Comma-separated tags")
	background := fs.String("background", "", "Personal background story")
	active := fs.String("active", "", "Active status (true/false)")

	if err := fs.Parse(os.Args[3:]); err != nil {
		return err
	}

	// Update fields if provided
	if *name != "" {
		existing.Name = *name
	}
	if *description != "" {
		existing.Description = *description
	}
	if *background != "" {
		existing.Background = *background
	}
	if *tags != "" {
		tagList := strings.Split(*tags, ",")
		for i, tag := range tagList {
			tagList[i] = strings.TrimSpace(tag)
		}
		existing.Tags = tagList
	}
	if *active != "" {
		existing.IsActive = *active == "true"
	}

	if err := c.UpdateIdentity(id, existing); err != nil {
		return fmt.Errorf("failed to update identity: %v", err)
	}

	fmt.Printf("Identity updated successfully: %s\n", id)
	return nil
}

func deleteIdentity(c client.Client) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("identity ID is required")
	}

	id := os.Args[2]
	if err := c.DeleteIdentity(id); err != nil {
		return err
	}

	fmt.Printf("Identity deleted successfully: %s\n", id)
	return nil
}

func getIdentityWithPersona(c client.Client) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("identity ID is required")
	}

	id := os.Args[2]
	iwp, err := c.GetIdentityWithPersona(id)
	if err != nil {
		return err
	}

	fmt.Printf("Identity with Persona: %s\n", iwp.Identity.Id)
	fmt.Printf("  Identity Name: %s\n", iwp.Identity.Name)
	fmt.Printf("  Persona Name: %s\n", iwp.Persona.Name)
	fmt.Printf("  Persona Topic: %s\n", iwp.Persona.Topic)
	fmt.Printf("  Identity Description: %s\n", iwp.Identity.Description)
	fmt.Printf("  Persona Prompt: %s\n", iwp.Persona.Prompt)
	fmt.Printf("  Status: %s\n", map[bool]string{true: "Active", false: "Inactive"}[iwp.Identity.IsActive])

	if len(iwp.Identity.Tags) > 0 {
		fmt.Printf("  Tags: %s\n", strings.Join(iwp.Identity.Tags, ", "))
	}

	return nil
}

// Generation functions
func generateIdentity(c client.Client) error {
	fs := flag.NewFlagSet("generate-identity", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: fr0g-ai-aip generate-identity -persona-id <id> [-name <name>] [-random]")
		fmt.Println("  -persona-id <id>    Persona ID (required)")
		fmt.Println("  -name <name>        Identity name (optional, will generate if not provided)")
		fmt.Println("  -random             Generate random attributes (default)")
	}
	personaID := fs.String("persona-id", "", "Persona ID (required)")
	name := fs.String("name", "", "Identity name (optional)")
	// random := fs.Bool("random", true, "Generate random attributes") // TODO: implement random generation

	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	if *personaID == "" {
		fs.Usage()
		return fmt.Errorf("persona-id is required")
	}

	// Verify persona exists
	_, err := c.Get(*personaID)
	if err != nil {
		return fmt.Errorf("persona not found: %v", err)
	}

	// Generate identity name if not provided
	if *name == "" {
		*name = fmt.Sprintf("Generated Identity %d", time.Now().Unix())
	}

	// For now, we'll create a simple identity
	// In a full implementation, you'd use the generator package
	i := types.Identity{
		PersonaId:   *personaID,
		Name:        *name,
		Description: "A generated identity with rich attributes",
		Tags:        []string{"generated", "random"},
		IsActive:    true,
		RichAttributes: &types.RichAttributes{
			Demographics: &types.Demographics{
				Age:       30,
				Gender:    "non-binary",
				Education: "bachelors",
				Location: &types.Location{
					Country:    "United States",
					City:       "San Francisco",
					UrbanRural: "urban",
				},
			},
			Psychographics: &types.Psychographics{
				Personality: &types.Personality{
					Openness:          0.7,
					Conscientiousness: 0.6,
					Extraversion:      0.5,
					Agreeableness:     0.8,
					Neuroticism:       0.3,
				},
				Values:        []string{"curiosity", "compassion", "growth"},
				RiskTolerance: "medium",
			},
		},
	}

	if err := c.CreateIdentity(&i); err != nil {
		return fmt.Errorf("failed to create generated identity: %v", err)
	}

	fmt.Printf("Generated identity created successfully: %s\n", i.Id)
	fmt.Printf("  Name: %s\n", i.Name)
	fmt.Printf("  Persona: %s\n", i.PersonaId)
	if i.RichAttributes != nil && i.RichAttributes.Demographics != nil {
		fmt.Printf("  Age: %d\n", i.RichAttributes.Demographics.Age)
		if i.RichAttributes.Demographics.Location != nil {
			fmt.Printf("  Location: %s, %s\n", i.RichAttributes.Demographics.Location.City, i.RichAttributes.Demographics.Location.Country)
		}
	}

	return nil
}

func generateCommunity(c client.Client) error {
	fs := flag.NewFlagSet("generate-community", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: fr0g-ai-aip generate-community -persona-id <id> -size <number> [-location <city,country>] [-age-range <min>-<max>]")
		fmt.Println("  -persona-id <id>    Persona ID (required)")
		fmt.Println("  -size <number>      Number of identities to generate (required)")
		fmt.Println("  -location <city,country>  Location for the community (optional)")
		fmt.Println("  -age-range <min>-<max>    Age range for the community (optional)")
	}
	personaID := fs.String("persona-id", "", "Persona ID (required)")
	size := fs.Int("size", 0, "Number of identities to generate (required)")
	location := fs.String("location", "", "Location (city,country)")
	ageRange := fs.String("age-range", "", "Age range (min-max)")

	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	if *personaID == "" || *size <= 0 {
		fs.Usage()
		return fmt.Errorf("persona-id and size are required")
	}

	// Verify persona exists
	_, err := c.Get(*personaID)
	if err != nil {
		return fmt.Errorf("persona not found: %v", err)
	}

	// Parse location if provided
	var loc *types.Location
	if *location != "" {
		parts := strings.Split(*location, ",")
		if len(parts) >= 2 {
			loc = &types.Location{
				City:       strings.TrimSpace(parts[0]),
				Country:    strings.TrimSpace(parts[1]),
				UrbanRural: "urban",
			}
		}
	}

	// Parse age range if provided
	var ageRangeStruct *types.AgeRange
	if *ageRange != "" {
		parts := strings.Split(*ageRange, "-")
		if len(parts) == 2 {
			min, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			max, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err1 == nil && err2 == nil && min <= max {
				ageRangeStruct = &types.AgeRange{Min: int32(min), Max: int32(max)}
			}
		}
	}

	// Generate community identities
	createdCount := 0
	for i := 0; i < *size; i++ {
		name := fmt.Sprintf("Community Member %d", i+1)

		// Create identity with location and age range if specified
		identity := types.Identity{
			PersonaId:   *personaID,
			Name:        name,
			Description: fmt.Sprintf("Community member %d with generated attributes", i+1),
			Tags:        []string{"community", "generated"},
			IsActive:    true,
			RichAttributes: &types.RichAttributes{
				Demographics: &types.Demographics{
					Age:       int32(25 + i%45), // Spread ages
					Gender:    []string{"male", "female", "non-binary"}[i%3],
					Education: []string{"high_school", "bachelors", "masters"}[i%3],
				},
				Psychographics: &types.Psychographics{
					Personality: &types.Personality{
						Openness:          0.5 + float64(i%5)*0.1,
						Conscientiousness: 0.5 + float64(i%5)*0.1,
						Extraversion:      0.5 + float64(i%5)*0.1,
						Agreeableness:     0.5 + float64(i%5)*0.1,
						Neuroticism:       0.3 + float64(i%5)*0.1,
					},
					Values: []string{"community", "growth", "connection"},
				},
			},
		}

		// Apply location if specified
		if loc != nil {
			identity.RichAttributes.Demographics.Location = loc
		}

		// Apply age range if specified
		if ageRangeStruct != nil {
			ageSpread := int(ageRangeStruct.Max - ageRangeStruct.Min + 1)
			identity.RichAttributes.Demographics.Age = ageRangeStruct.Min + int32(i%ageSpread)
		}

		if err := c.CreateIdentity(&identity); err != nil {
			fmt.Printf("Warning: Failed to create identity %d: %v\n", i+1, err)
			continue
		}
		createdCount++
	}

	fmt.Printf("Successfully generated %d community identities for persona %s\n", createdCount, *personaID)
	if loc != nil {
		fmt.Printf("  Location: %s, %s\n", loc.City, loc.Country)
	}
	if ageRangeStruct != nil {
		fmt.Printf("  Age Range: %d-%d\n", ageRangeStruct.Min, ageRangeStruct.Max)
	}

	return nil
}
