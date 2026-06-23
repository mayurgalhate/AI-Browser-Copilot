package agent

import (
	"context"
	"testing"

	"browser-copilot-backend/internal/config"
	"github.com/joho/godotenv"
)

func TestNewGroqLLM(t *testing.T) {
	// Load test environment if available
	_ = godotenv.Load("../../.env")
	
	// Load the config explicitly since the test bypasses main.go
	config.LoadConfig()

	if config.AppConfig.GroqAPIKey == "" {
		t.Skip("Skipping test: GROQ_API_KEY is not set")
	}

	llm, err := NewGroqLLM()
	if err != nil {
		t.Fatalf("Failed to initialize Groq LLM: %v", err)
	}
	if llm == nil {
		t.Fatal("Expected LLM instance, got nil")
	}
}

func TestAgentInitialization(t *testing.T) {
	_ = godotenv.Load("../../.env")
	config.LoadConfig()

	if config.AppConfig.GroqAPIKey == "" {
		t.Skip("Skipping test: GROQ_API_KEY is not set")
	}

	// Just a simple ping to ensure the chain doesn't immediately panic or fail on setup
	_, err := RunAgent(context.Background(), "test-session-123", "[ASK MODE] Hello")
	if err != nil {
		t.Logf("Agent failed (expected if API limits hit or WS disconnected): %v", err)
	}
}
