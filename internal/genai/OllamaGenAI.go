package genai

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/ollama/ollama/api"
)

type OllamaGenAI struct {
	model  string
	client *api.Client
}

func (o *OllamaGenAI) Generate(systemPrompt string, prompt string) string {
	ctx := context.Background()
	result := ""
	err := o.client.Generate(ctx, &api.GenerateRequest{
		Model:  o.model,
		Prompt: prompt,
		System: systemPrompt,
		Stream: new(bool),
		KeepAlive: &api.Duration{
			Duration: time.Hour,
		},
	}, func(gr api.GenerateResponse) error {
		result = gr.Response
		return nil
	})

	if err != nil {
		slog.Error("Failed to generate result", "err", err)
	}

	return result
}

func NewOllamaGenAI(address, model string) (GenAI, error) {
	url, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	ollama := api.NewClient(url, client)

	return &OllamaGenAI{
		model:  model,
		client: ollama,
	}, nil
}
