package openai

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

func ExplainCode(code string) (string, error) {
	// Configura o cliente para usar a API local do Ollama
	cfg := openai.DefaultConfig("")
	cfg.BaseURL = "http://localhost:11434/v1" // Ollama API local
	client := openai.NewClientWithConfig(cfg)

	// Monta a requisição de chat com o modelo codellama
	req := openai.ChatCompletionRequest{
		Model: "codellama",
		Messages: []openai.ChatCompletionMessage{
			{
				Role: "user",
				Content: `Explique de forma clara o que esse código faz:

` + code,
			},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
