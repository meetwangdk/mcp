package mcp

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"hcs-agent/pkg/config"
	"sync"
)

var LLMClientO *LLMClient

// LLMClient LLM 客户端
type LLMClient struct {
	client      *openai.Client
	model       string
	stream      bool
	temperature float32
}

var NewLLMClientOnce = sync.OnceFunc(func() {
	llmClient := NewLLMClient(
		config.Conf.Model.Token,
		config.Conf.Model.URL,
		config.Conf.Model.Answer,
		false,
		0.7,
	)
	LLMClientO = llmClient
})

// NewLLMClient 创建 LLM 客户端
func NewLLMClient(apiKey, baseURL, model string, stream bool, temperature float32) *LLMClient {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	return &LLMClient{
		client:      openai.NewClientWithConfig(config),
		model:       model,
		stream:      stream,
		temperature: temperature,
	}
}

// CreateCompletion 调用 LLM 生成响应
func (l *LLMClient) CreateCompletion(ctx context.Context, messages []openai.ChatCompletionMessage, tools []openai.Tool) (interface{}, error) {
	req := openai.ChatCompletionRequest{
		Model:       l.model,
		Messages:    messages,
		Stream:      l.stream,
		Temperature: l.temperature,
		ChatTemplateKwargs: map[string]any{
			"enable_thinking": false,
		},
		ParallelToolCalls: true,
	}

	// 若有工具，启用工具调用
	if len(tools) > 0 {
		req.Tools = tools
		req.ToolChoice = "auto"
	}

	if l.stream {
		return l.client.CreateChatCompletionStream(ctx, req)
	}
	return l.client.CreateChatCompletion(ctx, req)
}
