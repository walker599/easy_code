package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

type ShellRunner struct{}

func NewShellRunner() *ShellRunner {
	return &ShellRunner{}
}

func (s *ShellRunner) Name() string {
	return "run_shell_command"
}

func (s *ShellRunner) Description() string {
	return "Executes a shell command on the local machine. Use this to list files, read file contents, or run system commands. The command will be executed in the current working directory."
}

func (s *ShellRunner) Parameters() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"command": {
				"type": "string",
				"description": "The command to execute, e.g., 'ls -la' or 'cat main.go'"
			}
		},
		"required": ["command"]
	}`)
}

func (s *ShellRunner) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", params.Command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", params.Command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Return error as result so LLM can see what happened
		return fmt.Sprintf("Command failed with error: %v\nOutput:\n%s", err, string(output)), nil
	}

	// Limit output size to avoid blowing up context window
	outStr := string(output)
	if len(outStr) > 2000 {
		outStr = outStr[:2000] + "\n... (output truncated)"
	}

	if outStr == "" {
		return "(command executed successfully with no output)", nil
	}

	return outStr, nil
}
