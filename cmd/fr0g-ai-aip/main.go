package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/api"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/cli"
	grpcserver "github.com/fr0g-vibe/fr0g-ai-aip/internal/grpc"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
)

func main() {
	var (
		serverMode  = flag.Bool("server", false, "Run HTTP REST API server")
		grpcMode    = flag.Bool("grpc", false, "Run gRPC server")
		httpPort    = flag.String("port", "8080", "HTTP server port")
		grpcPort    = flag.String("grpc-port", "9090", "gRPC server port")
		storageType = flag.String("storage", "memory", "Storage type: memory, file")
		dataDir     = flag.String("data-dir", "./data", "Data directory for file storage")
		help        = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Initialize storage for server mode
	if *serverMode || *grpcMode {
		var store storage.Storage
		var err error
		
		switch *storageType {
		case "memory":
			store = storage.NewMemoryStorage()
		case "file":
			store, err = storage.NewFileStorage(*dataDir)
			if err != nil {
				log.Fatalf("Failed to initialize file storage: %v", err)
			}
		default:
			log.Fatalf("Unknown storage type: %s", *storageType)
		}
		
		// Set default service for API handlers
		persona.SetDefaultService(persona.NewService(store))
	}

	// If both server modes are requested, run them concurrently
	if *serverMode && *grpcMode {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			fmt.Printf("Starting fr0g-ai-aip HTTP server on port %s (storage: %s)\n", *httpPort, *storageType)
			if err := api.StartServer(*httpPort); err != nil {
				log.Fatalf("Failed to start HTTP server: %v", err)
			}
		}()

		go func() {
			defer wg.Done()
			if err := grpcserver.StartGRPCServer(*grpcPort); err != nil {
				log.Fatalf("Failed to start gRPC server: %v", err)
			}
		}()

		wg.Wait()
	} else if *serverMode {
		fmt.Printf("Starting fr0g-ai-aip HTTP server on port %s (storage: %s)\n", *httpPort, *storageType)
		if err := api.StartServer(*httpPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	} else if *grpcMode {
		fmt.Printf("Starting fr0g-ai-aip gRPC server on port %s (storage: %s)\n", *grpcPort, *storageType)
		if err := grpcserver.StartGRPCServer(*grpcPort); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	} else {
		// CLI mode - use environment configuration
		config := cli.GetConfigFromEnv()
		if err := cli.ExecuteWithConfig(config); err != nil {
			log.Fatalf("CLI error: %v", err)
		}
	}
}
