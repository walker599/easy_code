package llm

import (
	"context"
	"encoding/json"
	"log"

	"go_quick_star_cli/internal/config"
	"go_quick_star_cli/internal/tools"

	"github.com/sashabaranov/go-openai"
)

// Provider defines the interface for an LLM provider
type Provider interface {
	Chat(ctx context.Context, messages []openai.ChatCompletionMessage, toolSet tools.ToolSet) (openai.ChatCompletionMessage, error)
}

type Client struct {
	client *openai.Client
	model  string
	debug  bool
}

func NewClient(cfg *config.Config) *Client {
	openaiConfig := openai.DefaultConfig(cfg.OpenAIKey)
	if cfg.OpenAIBaseURL != "" {
		openaiConfig.BaseURL = cfg.OpenAIBaseURL
	}

	c := openai.NewClientWithConfig(openaiConfig)
	return &Client{
		client: c,
		model:  cfg.Model,
		debug:  cfg.Debug,
	}
}

// Chat sends messages to LLM and returns the response message.
// It handles the conversion of tools to OpenAI format.
func (c *Client) Chat(ctx context.Context, messages []openai.ChatCompletionMessage, toolSet tools.ToolSet) (openai.ChatCompletionMessage, error) {
	var openAITools []openai.Tool
	if len(toolSet) > 0 {
		for _, t := range toolSet {
			openAITools = append(openAITools, openai.Tool{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        t.Name(),
					Description: t.Description(),
					Parameters:  t.Parameters(),
				},
			})
		}
	}

	req := openai.ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
	}

	// Only set Tools if we have them, otherwise API might complain or behave differently
	if len(openAITools) > 0 {
		req.Tools = openAITools
	}

	if c.debug {
		reqJSON, _ := json.MarshalIndent(req, "", "  ")
		log.Printf("DEBUG: OpenAI Request:\n%s\n", string(reqJSON))
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}

	if c.debug {
		respJSON, _ := json.MarshalIndent(resp, "", "  ")
		log.Printf("DEBUG: OpenAI Response:\n%s\n", string(respJSON))
	}

	return resp.Choices[0].Message, nil
}
