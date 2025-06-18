package genai

type GenAI interface {
	Generate(systemPrompt string, prompt string) string
}

func NewGenAI(address string, model string) (GenAI, error) {
	return NewOllamaGenAI(address, model)
}
