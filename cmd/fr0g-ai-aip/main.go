package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/api"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/cli"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/config"
	grpcserver "github.com/fr0g-vibe/fr0g-ai-aip/internal/grpc"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
)

// App holds the application state
type App struct {
	config  *config.Config
	service *persona.Service
}

// NewApp creates a new application instance
func NewApp() (*App, error) {
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}
	
	// Add port conflict validation
	if cfg.HTTP.Port == cfg.GRPC.Port {
		return nil, fmt.Errorf("HTTP and gRPC ports cannot be the same: %s", cfg.HTTP.Port)
	}
	
	// Initialize storage
	store, err := createStorage(cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %v", err)
	}
	
	return &App{
		config:  cfg,
		service: persona.NewService(store),
	}, nil
}

// RunCLI runs the CLI interface
func (app *App) RunCLI() error {
	cliConfig := cli.Config{
		ClientType:  app.config.Client.Type,
		StorageType: app.config.Storage.Type,
		DataDir:     app.config.Storage.DataDir,
		ServerURL:   app.config.Client.ServerURL,
	}
	return cli.ExecuteWithConfig(cliConfig)
}

// RunServers runs the HTTP and/or gRPC servers
func (app *App) RunServers(httpMode, grpcMode bool) error {
	// Print startup banner
	app.printStartupBanner()
	
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	
	// Start HTTP server
	if httpMode {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("Starting fr0g-ai-aip HTTP server on port %s (storage: %s)\n", 
				app.config.HTTP.Port, app.config.Storage.Type)
			if err := api.StartServerWithConfig(app.config, app.service); err != nil {
				errChan <- fmt.Errorf("HTTP server error: %v", err)
			}
		}()
	}
	
	// Start gRPC server
	if grpcMode {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("Starting fr0g-ai-aip gRPC server on port %s (storage: %s)\n", 
				app.config.GRPC.Port, app.config.Storage.Type)
			if err := grpcserver.StartGRPCServerWithConfig(app.config, app.service); err != nil {
				errChan <- fmt.Errorf("gRPC server error: %v", err)
			}
		}()
	}
	
	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		fmt.Printf("\nReceived signal %v, shutting down immediately...\n", sig)
		os.Exit(0)
	case err := <-errChan:
		return err
	}
	
	return nil
}

// printStartupBanner displays a startup banner with configuration info
func (app *App) printStartupBanner() {
	fmt.Println("🐸 fr0g-ai-aip - AI Personas Management System")
	fmt.Printf("   Version: 1.0.0\n")
	fmt.Printf("   Storage: %s", app.config.Storage.Type)
	if app.config.Storage.Type == "file" {
		fmt.Printf(" (%s)", app.config.Storage.DataDir)
	}
	fmt.Println()
	fmt.Println("   Ready to manage your AI personas!")
	fmt.Println()
}

func createStorage(cfg config.StorageConfig) (storage.Storage, error) {
	switch cfg.Type {
	case "memory":
		return storage.NewMemoryStorage(), nil
	case "file":
		if cfg.DataDir == "" {
			return nil, fmt.Errorf("data directory is required for file storage")
		}
		store, err := storage.NewFileStorage(cfg.DataDir)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize file storage at %s: %v", cfg.DataDir, err)
		}
		return store, nil
	default:
		return nil, fmt.Errorf("unsupported storage type '%s' (supported: memory, file)", cfg.Type)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	var (
		serverMode = flag.Bool("server", false, "Run HTTP REST API server")
		grpcMode   = flag.Bool("grpc", false, "Run gRPC server")
		httpPort   = flag.String("port", "", "HTTP server port (overrides config)")
		grpcPort   = flag.String("grpc-port", "", "gRPC server port (overrides config)")
		help       = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		return nil
	}

	app, err := NewApp()
	if err != nil {
		return err
	}
	
	// Override config with command line flags
	if *httpPort != "" {
		app.config.HTTP.Port = *httpPort
	}
	if *grpcPort != "" {
		app.config.GRPC.Port = *grpcPort
	}

	// Determine mode
	if *serverMode || *grpcMode {
		return app.RunServers(*serverMode, *grpcMode)
	} else {
		return app.RunCLI()
	}
}
