package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/taskvanguard/taskvanguard/pkg/types"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Client struct {
	llm llms.Model
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewClient(cfg *types.LLMConfig) (*Client, error) {
	var model llms.Model
	var err error

	switch cfg.Provider {
	case "openai":
		opts := []openai.Option{
			openai.WithModel(cfg.Model),
		}
		if cfg.APIKey != "" {
			opts = append(opts, openai.WithToken(cfg.APIKey))
		}
		if cfg.BaseURL != "" {
			opts = append(opts, openai.WithBaseURL(cfg.BaseURL))
		}
		model, err = openai.New(opts...)
	case "deepseek":
		opts := []openai.Option{
			openai.WithModel(cfg.Model),
			openai.WithBaseURL("https://api.deepseek.com/v1"),
		}
		if cfg.APIKey != "" {
			opts = append(opts, openai.WithToken(cfg.APIKey))
		}
		if cfg.BaseURL != "" {
			opts = append(opts, openai.WithBaseURL(cfg.BaseURL))
		}
		model, err = openai.New(opts...)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	if err != nil {
		return nil, err
	}

	return &Client{
		llm: model,
	}, nil
}

func (c *Client) Chat(messages []Message) (string, error) {
	llmMessages := make([]llms.MessageContent, len(messages))
	for i, msg := range messages {
		switch msg.Role {
		case "user":
			llmMessages[i] = llms.TextParts(llms.ChatMessageTypeHuman, msg.Content)
		case "assistant":
			llmMessages[i] = llms.TextParts(llms.ChatMessageTypeAI, msg.Content)
		case "system":
			llmMessages[i] = llms.TextParts(llms.ChatMessageTypeSystem, msg.Content)
		default:
			llmMessages[i] = llms.TextParts(llms.ChatMessageTypeHuman, msg.Content)
		}
	}

	ctx := context.Background()
	completion, err := c.llm.GenerateContent(ctx, llmMessages)
	if err != nil {
		return "", err
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return completion.Choices[0].Content, nil
}

// CleanResponse removes markdown code blocks from LLM responses
func CleanResponse(response string) string {
	cleanResponse := strings.TrimSpace(response)
	
	// Remove opening code block markers
	if strings.HasPrefix(cleanResponse, "```json") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
	} else if strings.HasPrefix(cleanResponse, "```Json") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```Json")
	} else if strings.HasPrefix(cleanResponse, "```JSON") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```JSON")
	} else if strings.HasPrefix(cleanResponse, "```") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```")
	}
	
	// Remove closing code block marker
	if strings.HasSuffix(cleanResponse, "```") {
		cleanResponse = strings.TrimSuffix(cleanResponse, "```")
	}
	
	return strings.TrimSpace(cleanResponse)
}
