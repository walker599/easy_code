package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type GitStatus struct{}

func NewGitStatus() *GitStatus {
	return &GitStatus{}
}

func (g *GitStatus) Name() string {
	return "git_status"
}

func (g *GitStatus) Description() string {
	return "Shows the working tree status (git status)."
}

func (g *GitStatus) Parameters() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {},
		"required": []
	}`)
}

func (g *GitStatus) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running git status: %v\nOutput: %s", err, string(output)), nil
	}
	return string(output), nil
}

type GitDiff struct{}

func NewGitDiff() *GitDiff {
	return &GitDiff{}
}

func (g *GitDiff) Name() string {
	return "git_diff"
}

func (g *GitDiff) Description() string {
	return "Shows changes between commits, commit and working tree, etc (git diff)."
}

func (g *GitDiff) Parameters() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"staged": {
				"type": "boolean",
				"description": "Show staged changes (--staged)."
			}
		}
	}`)
}

func (g *GitDiff) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		Staged bool `json:"staged"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}

	cmdArgs := []string{"diff"}
	if params.Staged {
		cmdArgs = append(cmdArgs, "--staged")
	}

	cmd := exec.CommandContext(ctx, "git", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running git diff: %v\nOutput: %s", err, string(output)), nil
	}
	
	// Limit output size
	outStr := string(output)
	if len(outStr) > 5000 {
		outStr = outStr[:5000] + "\n... (output truncated)"
	}
	
	return outStr, nil
}

type GitCommit struct{}

func NewGitCommit() *GitCommit {
	return &GitCommit{}
}

func (g *GitCommit) Name() string {
	return "git_commit"
}

func (g *GitCommit) Description() string {
	return "Record changes to the repository (git commit -m)."
}

func (g *GitCommit) Parameters() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"message": {
				"type": "string",
				"description": "The commit message."
			}
		},
		"required": ["message"]
	}`)
}

func (g *GitCommit) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}

	cmd := exec.CommandContext(ctx, "git", "commit", "-m", params.Message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running git commit: %v\nOutput: %s", err, string(output)), nil
	}
	return string(output), nil
}
