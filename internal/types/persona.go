package types

// Persona represents an AI persona with specific expertise
type Persona struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Topic   string            `json:"topic"`
	Prompt  string            `json:"prompt"`
	Context map[string]string `json:"context,omitempty"`
	RAG     []string          `json:"rag,omitempty"`
}
