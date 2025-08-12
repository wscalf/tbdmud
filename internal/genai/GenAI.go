package genai

type GenAI interface {
	Generate(systemPrompt string, prompt string) string
	Chat(systemPrompt string, history []ChatMessage, actor, pose string) []ChatMessage
}

type ChatMessage struct {
	Actor string
	Body  string
}

func NewGenAI(address string, model string) (GenAI, error) {
	return NewOllamaGenAI(address, model)
}
