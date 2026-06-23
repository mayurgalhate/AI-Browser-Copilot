package agent

import (
	"context"

	"browser-copilot-backend/internal/config"
	brtools "browser-copilot-backend/internal/tools"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
)

// NewGroqLLM creates an LLM client pointing to Groq via OpenAI provider
func NewGroqLLM() (*openai.LLM, error) {
	apiKey := config.AppConfig.GroqAPIKey
	baseURL := config.AppConfig.GroqBaseURL

	return openai.New(
		openai.WithBaseURL(baseURL),
		openai.WithToken(apiKey),
		openai.WithModel("llama-3.3-70b-versatile"),
	)
}

// RunAgent executes a query using the agent with tools
func RunAgent(ctx context.Context, sessionID, query string) (string, error) {
	llm, err := NewGroqLLM()
	if err != nil {
		return "", err
	}

	browserTools := brtools.CreateBrowserTools(sessionID)

	// Create conversational memory
	// In a real app, you would load history from Postgres/Redis here
	mem := memory.NewConversationBuffer()

	executor, err := agents.Initialize(
		llm,
		browserTools,
		agents.ConversationalReactDescription,
		agents.WithMemory(mem),
		agents.WithMaxIterations(5),
	)
	if err != nil {
		return "", err
	}

	res, err := chains.Run(ctx, executor, query)
	return res, err
}
