package agent

import (
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

// NewOllamaEmbedder initializes an embedder using a local Ollama instance
func NewOllamaEmbedder(model string) (*embeddings.EmbedderImpl, error) {
	if model == "" {
		model = "nomic-embed-text"
	}
	llm, err := ollama.New(ollama.WithModel(model))
	if err != nil {
		return nil, err
	}
	return embeddings.NewEmbedder(llm)
}
