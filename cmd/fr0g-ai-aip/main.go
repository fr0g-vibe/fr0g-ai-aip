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
		serverMode = flag.Bool("server", false, "Run in server mode")
		port       = flag.String("port", "8080", "Server port")
		help       = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *serverMode {
		fmt.Printf("Starting fr0g-ai-aip server on port %s\n", *port)
		if err := api.StartServer(*port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	} else {
		// CLI mode
		if err := cli.Execute(); err != nil {
			log.Fatalf("CLI error: %v", err)
		}
	}
}
