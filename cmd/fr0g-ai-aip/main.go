package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fr0g-ai/fr0g-ai-aip/internal/api"
	"github.com/fr0g-ai/fr0g-ai-aip/internal/cli"
	grpcserver "github.com/fr0g-ai/fr0g-ai-aip/internal/grpc"
)

func main() {
	var (
		serverMode = flag.Bool("server", false, "Run HTTP REST API server")
		grpcMode   = flag.Bool("grpc", false, "Run gRPC server")
		httpPort   = flag.String("port", "8080", "HTTP server port")
		grpcPort   = flag.String("grpc-port", "9090", "gRPC server port")
		help       = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// If both server modes are requested, run them concurrently
	if *serverMode && *grpcMode {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			fmt.Printf("Starting fr0g-ai-aip HTTP server on port %s\n", *httpPort)
			if err := api.StartServer(*httpPort); err != nil {
				log.Fatalf("Failed to start HTTP server: %v", err)
			}
		}()

		go func() {
			defer wg.Done()
			fmt.Printf("Starting fr0g-ai-aip gRPC server on port %s\n", *grpcPort)
			if err := grpcserver.StartGRPCServer(*grpcPort); err != nil {
				log.Fatalf("Failed to start gRPC server: %v", err)
			}
		}()

		wg.Wait()
	} else if *serverMode {
		fmt.Printf("Starting fr0g-ai-aip HTTP server on port %s\n", *httpPort)
		if err := api.StartServer(*httpPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	} else if *grpcMode {
		fmt.Printf("Starting fr0g-ai-aip gRPC server on port %s\n", *grpcPort)
		if err := grpcserver.StartGRPCServer(*grpcPort); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	} else {
		// CLI mode
		if err := cli.Execute(); err != nil {
			log.Fatalf("CLI error: %v", err)
		}
	}
}
