//go:build no_grpc

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fr0g-ai/fr0g-ai-aip/internal/api"
	"github.com/fr0g-ai/fr0g-ai-aip/internal/cli"
)

func main() {
	var (
		serverMode = flag.Bool("server", false, "Run HTTP REST API server")
		grpcMode   = flag.Bool("grpc", false, "Run gRPC server (not available in this build)")
		httpPort   = flag.String("port", "8080", "HTTP server port")
		help       = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *grpcMode {
		log.Fatal("gRPC support not available in this build. Run 'make build' to build with gRPC support.")
	}

	if *serverMode {
		fmt.Printf("Starting fr0g-ai-aip HTTP server on port %s\n", *httpPort)
		if err := api.StartServer(*httpPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	} else {
		// CLI mode
		if err := cli.Execute(); err != nil {
			log.Fatalf("CLI error: %v", err)
		}
	}
}
