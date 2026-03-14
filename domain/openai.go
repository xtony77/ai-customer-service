package domain

type OpenAIInterface interface {
	GenerateEmbedding(input string) ([]float32, error)
	ChatCompletion(question string, references []QAEmbedding) (string, error)
}
