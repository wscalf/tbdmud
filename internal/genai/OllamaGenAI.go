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
		slog.Error("failed to generate result", "err", err)
	}

	return result
}

func (o *OllamaGenAI) Chat(systemPrompt string, history []ChatMessage, actor, pose string) []ChatMessage {
	ctx := context.Background()
	var result ChatMessage

	messages := make([]api.Message, 0, len(history)+2) //Existing chat history + system message + current pose
	messages = append(messages, api.Message{
		Role:    "system",
		Content: systemPrompt,
	})

	for _, msg := range history {
		messages = append(messages, api.Message{
			Role:    msg.Actor,
			Content: msg.Body,
		})
	}

	messages = append(messages, api.Message{
		Role:    "user",
		Content: pose,
	})

	err := o.client.Chat(ctx, &api.ChatRequest{
		Model:    o.model,
		Messages: messages,
		Stream:   new(bool),
		KeepAlive: &api.Duration{
			Duration: time.Hour,
		},
	}, func(cr api.ChatResponse) error {
		result = ChatMessage{
			Actor: cr.Message.Role,
			Body:  cr.Message.Content,
		}
		return nil
	})

	if err != nil {
		slog.Error("failed to generate next chat message", "err", err)
		return []ChatMessage{}
	}

	return []ChatMessage{{Actor: "user", Body: pose}, result}
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
