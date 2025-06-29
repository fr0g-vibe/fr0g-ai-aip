package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/fr0g-ai/fr0g-ai-aip/internal/persona"
)

// Execute runs the CLI interface
func Execute() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	command := os.Args[1]
	
	switch command {
	case "list":
		return listPersonas()
	case "create":
		return createPersona()
	case "get":
		return getPersona()
	case "delete":
		return deletePersona()
	default:
		printUsage()
		return fmt.Errorf("unknown command: %s", command)
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
	fmt.Println("  delete <id>         Delete persona by ID")
	fmt.Println()
	fmt.Println("Server mode:")
	fmt.Println("  -server             Run in server mode")
	fmt.Println("  -port <port>        Server port (default: 8080)")
}

func listPersonas() error {
	personas := persona.ListPersonas()
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

func createPersona() error {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	name := fs.String("name", "", "Persona name")
	topic := fs.String("topic", "", "Persona topic/expertise")
	prompt := fs.String("prompt", "", "System prompt")
	
	fs.Parse(os.Args[2:])
	
	if *name == "" || *topic == "" || *prompt == "" {
		fmt.Println("Usage: fr0g-ai-aip create -name <name> -topic <topic> -prompt <prompt>")
		return fmt.Errorf("missing required parameters")
	}
	
	p := persona.Persona{
		Name:   *name,
		Topic:  *topic,
		Prompt: *prompt,
	}
	
	if err := persona.CreatePersona(&p); err != nil {
		return err
	}
	
	fmt.Printf("Created persona: %s (ID: %s)\n", p.Name, p.ID)
	return nil
}

func getPersona() error {
	if len(os.Args) < 3 {
		fmt.Println("Usage: fr0g-ai-aip get <id>")
		return fmt.Errorf("persona ID required")
	}
	
	id := os.Args[2]
	p, err := persona.GetPersona(id)
	if err != nil {
		return err
	}
	
	fmt.Printf("ID: %s\n", p.ID)
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Topic: %s\n", p.Topic)
	fmt.Printf("Prompt: %s\n", p.Prompt)
	return nil
}

func deletePersona() error {
	if len(os.Args) < 3 {
		fmt.Println("Usage: fr0g-ai-aip delete <id>")
		return fmt.Errorf("persona ID required")
	}
	
	id := os.Args[2]
	if err := persona.DeletePersona(id); err != nil {
		return err
	}
	
	fmt.Printf("Deleted persona: %s\n", id)
	return nil
}
