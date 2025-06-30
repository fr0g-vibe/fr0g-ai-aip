package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/client"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// Config holds CLI configuration
type Config struct {
	ClientType  string // "local", "rest", "grpc"
	StorageType string // "memory", "file"
	DataDir     string
	ServerURL   string
}

var defaultConfig = Config{
	ClientType:  "grpc",
	StorageType: "file", // Changed to file for persistence
	DataDir:     "./data",
	ServerURL:   "localhost:9090", // Default to gRPC
}

// Execute runs the CLI interface
func Execute() error {
	return ExecuteWithConfig(defaultConfig)
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
	default:
		printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}
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
		return client.NewRESTClient(config.ServerURL), nil
	case "grpc":
		// Use gRPC-specific default or extract from config
		address := "localhost:9090"
		if config.ServerURL != "" && config.ServerURL != "http://localhost:8080" {
			// If a custom server URL is provided and it's not the REST default, use it
			address = config.ServerURL
		}
		grpcClient, err := client.NewGRPCClient(address)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC client: %v", err)
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
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  fr0g-ai-aip [command] [options]")
	fmt.Println("  fr0g-ai-aip [flags]")
	fmt.Println()
	fmt.Println("COMMANDS:")
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
		fmt.Printf("  ID: %s, Name: %s, Topic: %s\n", p.ID, p.Name, p.Topic)
	}
	return nil
}

func createPersona(c client.Client) error {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	name := fs.String("name", "", "Persona name")
	topic := fs.String("topic", "", "Persona topic/expertise")
	prompt := fs.String("prompt", "", "System prompt")

	fs.Parse(os.Args[2:])

	if *name == "" || *topic == "" || *prompt == "" {
		fmt.Println("Usage: fr0g-ai-aip create -name <name> -topic <topic> -prompt <prompt>")
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

	fmt.Printf("Created persona: %s (ID: %s)\n", p.Name, p.ID)
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

	fmt.Printf("ID: %s\n", p.ID)
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Topic: %s\n", p.Topic)
	fmt.Printf("Prompt: %s\n", p.Prompt)
	if len(p.Context) > 0 {
		fmt.Println("Context:")
		for k, v := range p.Context {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
	if len(p.RAG) > 0 {
		fmt.Println("RAG:")
		for _, r := range p.RAG {
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

	fs := flag.NewFlagSet("update", flag.ExitOnError)
	name := fs.String("name", "", "Persona name")
	topic := fs.String("topic", "", "Persona topic/expertise")
	prompt := fs.String("prompt", "", "System prompt")

	fs.Parse(os.Args[3:])

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
