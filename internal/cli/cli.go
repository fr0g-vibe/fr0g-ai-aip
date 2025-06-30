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
	ClientType:  "local",
	StorageType: "memory",
	DataDir:     "./data",
	ServerURL:   "localhost:8080", // For REST, or "localhost:9090" for gRPC
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
		// Extract address from server URL or use default
		address := "localhost:9090"
		if config.ServerURL != "" {
			// Convert HTTP URL to gRPC address if needed
			// For now, assume ServerURL contains the gRPC address
			address = config.ServerURL
		}
		return client.NewGRPCClient(address)
	default:
		return nil, fmt.Errorf("unknown client type: %s", config.ClientType)
	}
}

func printUsage() {
	fmt.Println("fr0g-ai-aip - AI Personas Management")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  fr0g-ai-aip [command] [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list                List all personas")
	fmt.Println("  create              Create a new persona")
	fmt.Println("  get <id>            Get persona by ID")
	fmt.Println("  update <id>         Update persona by ID")
	fmt.Println("  delete <id>         Delete persona by ID")
	fmt.Println()
	fmt.Println("Server mode:")
	fmt.Println("  -server             Run in server mode")
	fmt.Println("  -port <port>        Server port (default: 8080)")
	fmt.Println("  -grpc               Run gRPC server")
	fmt.Println("  -grpc-port <port>   gRPC server port (default: 9090)")
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  FR0G_CLIENT_TYPE    Client type: local, rest, grpc (default: local)")
	fmt.Println("  FR0G_STORAGE_TYPE   Storage type: memory, file (default: memory)")
	fmt.Println("  FR0G_DATA_DIR       Data directory for file storage (default: ./data)")
	fmt.Println("  FR0G_SERVER_URL     Server URL for REST (http://localhost:8080) or gRPC (localhost:9090)")
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
