package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
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
		fmt.Printf("\nReceived signal %v, shutting down gracefully...\n", sig)
		cancel()
	case err := <-errChan:
		return err
	}
	
	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		fmt.Println("Shutdown complete")
	case <-sigChan:
		fmt.Println("Force shutdown")
	}
	
	return nil
}

func createStorage(cfg config.StorageConfig) (storage.Storage, error) {
	switch cfg.Type {
	case "memory":
		return storage.NewMemoryStorage(), nil
	case "file":
		return storage.NewFileStorage(cfg.DataDir)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
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
