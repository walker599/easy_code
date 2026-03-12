package agent

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"go_quick_star_cli/internal/tools"

	"github.com/sashabaranov/go-openai"
)

// MockLLMProvider is a mock implementation of llm.Provider
type MockLLMProvider struct {
	Responses []openai.ChatCompletionMessage
	CallCount int
}

func (m *MockLLMProvider) Chat(ctx context.Context, messages []openai.ChatCompletionMessage, toolSet tools.ToolSet) (openai.ChatCompletionMessage, error) {
	if m.CallCount >= len(m.Responses) {
		// Default response if run out of mock responses
		return openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: "Mock response",
		}, nil
	}
	resp := m.Responses[m.CallCount]
	m.CallCount++
	return resp, nil
}

// MockTool is a mock tool
type MockTool struct {
	Executed bool
}

func (m *MockTool) Name() string                { return "mock_tool" }
func (m *MockTool) Description() string         { return "A mock tool" }
func (m *MockTool) Parameters() json.RawMessage { return json.RawMessage(`{}`) }
func (m *MockTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	m.Executed = true
	return "mock tool executed", nil
}

func TestAgent_Run_Simple(t *testing.T) {
	mockLLM := &MockLLMProvider{
		Responses: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "Hello, world!",
			},
		},
	}

	ts := tools.NewToolSet()
	memPath := filepath.Join(t.TempDir(), "memory.json")
	ag := NewAgent(mockLLM, ts, memPath)

	resp, err := ag.Run(context.Background(), "Hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "Hello, world!" {
		t.Errorf("expected 'Hello, world!', got '%s'", resp)
	}
}

func TestAgent_Run_WithTool(t *testing.T) {
	mockTool := &MockTool{}
	ts := tools.NewToolSet(mockTool)

	mockLLM := &MockLLMProvider{
		Responses: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleAssistant,
				ToolCalls: []openai.ToolCall{
					{
						ID:   "call_1",
						Type: openai.ToolTypeFunction,
						Function: openai.FunctionCall{
							Name:      "mock_tool",
							Arguments: "{}",
						},
					},
				},
			},
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "Tool executed successfully",
			},
		},
	}

	memPath := filepath.Join(t.TempDir(), "memory.json")
	ag := NewAgent(mockLLM, ts, memPath)

	resp, err := ag.Run(context.Background(), "Run tool")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mockTool.Executed {
		t.Error("expected tool to be executed")
	}

	if resp != "Tool executed successfully" {
		t.Errorf("expected 'Tool executed successfully', got '%s'", resp)
	}
}
