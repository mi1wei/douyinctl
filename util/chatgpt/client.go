package chatgpt

import (
	"github.com/sashabaranov/go-openai"
)

type GptClient struct {
	index        int
	contents     []openai.ChatCompletionMessage
	openaiClient []*openai.Client
}

func getOpenaiClients() []*openai.Client {
	clients := make([]*openai.Client, 0)
	for _, token := range gptTokens {
		config := openai.DefaultConfig(token)
		config.BaseURL = "https://api.openai-proxy.com/v1"
		clients = append(clients, openai.NewClientWithConfig(config))
	}
	return clients
}

func NewClient() *GptClient {
	return &GptClient{
		index:        0,
		contents:     make([]openai.ChatCompletionMessage, 0),
		openaiClient: getOpenaiClients(),
	}
}
