package genai

import "log/slog"

type NullGenAI struct{}

func (n *NullGenAI) Generate(systemPrompt string, prompt string) string {
	slog.Warn("GenAI invoked but not enabled - returning placeholder content")
	return "GenAI is not enabled!"
}

func NewNullGenAI() GenAI {
	return &NullGenAI{}
}
