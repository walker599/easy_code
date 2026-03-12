package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"go_quick_star_cli/internal/llm"
	"go_quick_star_cli/internal/tools"

	"github.com/sashabaranov/go-openai"
)

type Agent struct {
	client llm.Provider
	tools  tools.ToolSet
	memory *Memory
}

const systemPrompt = `You are an expert Coding Agent, designed to help users with software development tasks. 
Your capabilities include:
1. Writing clean, efficient, and idiomatic code in various languages.
2. Debugging and fixing errors in existing code.
3. Performing code reviews and suggesting improvements for performance, security, and maintainability.
4. Explaining complex code concepts in simple terms.
5. Searching through the codebase to find relevant files and code snippets.
6. Managing git repositories (status, diff, commit).

You have access to the following tools:
- run_shell_command: Execute shell commands to list files, run tests, or build projects.
- read_file: Read the content of code files to understand the context.
- write_file: Write or modify code files.
- search_files: Search for text patterns in the codebase to locate relevant code.
- git_status: Check the status of the git repository.
- git_diff: Show changes in the git repository.
- git_commit: Commit changes to the git repository.

Guidelines:
- When asked to write code, always try to save it to a file using 'write_file' if the user implies a file name, otherwise print it.
- Before editing a file, always read it first using 'read_file' to understand the existing code.
- Use 'search_files' when you need to find where a function or variable is defined but don't know the file path.
- When fixing bugs, first analyze the error, then propose a fix, and finally apply it.
- Always answer in the user's language (e.g., Chinese if the user asks in Chinese).
- Be concise but thorough in your explanations.
- Use Markdown for formatting your responses, including code blocks with language identifiers.
- When asked to commit code, first check the status and diff to ensure you know what you are committing, then generate a meaningful commit message.
`

func NewAgent(client llm.Provider, ts tools.ToolSet, memoryPath string) *Agent {
	mem := NewMemory(memoryPath)
	// Try to load existing memory
	if err := mem.Load(); err != nil || len(mem.History) == 0 {
		// Initialize with system prompt if empty
		mem.Add(openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	return &Agent{
		client: client,
		tools:  ts,
		memory: mem,
	}
}

// Run executes a single turn of conversation
func (a *Agent) Run(ctx context.Context, input string) (string, error) {
	// Add user message to memory
	a.memory.Add(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	// Max turns to prevent infinite loops
	maxTurns := 10

	for i := 0; i < maxTurns; i++ {
		// Get current history from memory
		history := a.memory.Get()

		// Call LLM
		respMsg, err := a.client.Chat(ctx, history, a.tools)
		if err != nil {
			return "", fmt.Errorf("LLM error: %w", err)
		}

		// Add assistant response to memory
		a.memory.Add(respMsg)

		// Check for tool calls
		if len(respMsg.ToolCalls) == 0 {
			return respMsg.Content, nil
		}

		// Handle tool calls
		for _, toolCall := range respMsg.ToolCalls {
			toolName := toolCall.Function.Name
			toolArgs := toolCall.Function.Arguments

			tool, ok := a.tools[toolName]
			if !ok {
				// Tool not found, report error to LLM
				errorMsg := openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    fmt.Sprintf("Tool %s not found", toolName),
					ToolCallID: toolCall.ID,
				}
				a.memory.Add(errorMsg)
				continue
			}

			// Execute tool
			result, err := tool.Execute(ctx, json.RawMessage(toolArgs))
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			// Add tool result to memory
			resultMsg := openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    result,
				ToolCallID: toolCall.ID,
			}
			a.memory.Add(resultMsg)
		}
	}

	return "", fmt.Errorf("max turns exceeded")
}
