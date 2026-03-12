package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type FileReader struct{}

func NewFileReader() *FileReader {
	return &FileReader{}
}

func (f *FileReader) Name() string {
	return "read_file"
}

func (f *FileReader) Description() string {
	return "Reads the content of a file from the local filesystem."
}

func (f *FileReader) Parameters() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "The absolute or relative path to the file to read."
			}
		},
		"required": ["path"]
	}`)
}

func (f *FileReader) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}

	content, err := os.ReadFile(params.Path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err), nil
	}
	
	// Limit content size to avoid token limit issues
	// 50KB limit as a safety measure
	const maxBytes = 50 * 1024
	if len(content) > maxBytes {
		return fmt.Sprintf("File content truncated (showing first %d bytes):\n%s", maxBytes, string(content[:maxBytes])), nil
	}

	return string(content), nil
}

type FileWriter struct{}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

func (f *FileWriter) Name() string {
	return "write_file"
}

func (f *FileWriter) Description() string {
	return "Writes content to a file. Overwrites the file if it exists, or creates it if it doesn't."
}

func (f *FileWriter) Parameters() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "The absolute or relative path to the file to write."
			},
			"content": {
				"type": "string",
				"description": "The content to write to the file."
			}
		},
		"required": ["path", "content"]
	}`)
}

func (f *FileWriter) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}

	// Create directory if not exists
	dir := filepath.Dir(params.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Sprintf("Error creating directory: %v", err), nil
	}

	err := os.WriteFile(params.Path, []byte(params.Content), 0644)
	if err != nil {
		return fmt.Sprintf("Error writing file: %v", err), nil
	}

	return fmt.Sprintf("Successfully wrote to %s", params.Path), nil
}

type FileSearcher struct{}

func NewFileSearcher() *FileSearcher {
	return &FileSearcher{}
}

func (f *FileSearcher) Name() string {
	return "search_files"
}

func (f *FileSearcher) Description() string {
	return "Searches for a text pattern in files within a directory (recursive). Returns file paths and matching lines with line numbers."
}

func (f *FileSearcher) Parameters() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"pattern": {
				"type": "string",
				"description": "The text or regex pattern to search for."
			},
			"path": {
				"type": "string",
				"description": "The directory to search in (default is current directory)."
			}
		},
		"required": ["pattern"]
	}`)
}

func (f *FileSearcher) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		Pattern string `json:"pattern"`
		Path    string `json:"path"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %v", err)
	}

	if params.Path == "" {
		params.Path = "."
	}

	// Use grep for efficient searching
	// -r: recursive
	// -n: line numbers
	// -I: ignore binary files
	cmd := exec.CommandContext(ctx, "grep", "-rnI", params.Pattern, params.Path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// grep returns exit code 1 if no matches found, which is not an error for us
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return "No matches found.", nil
		}
		return fmt.Sprintf("Error searching files: %v\nOutput: %s", err, string(output)), nil
	}

	// Limit output size
	outStr := string(output)
	if len(outStr) > 5000 {
		outStr = outStr[:5000] + "\n... (output truncated)"
	}

	return outStr, nil
}
