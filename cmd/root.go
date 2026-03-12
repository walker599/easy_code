package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go_quick_star_cli/internal/agent"
	"go_quick_star_cli/internal/config"
	"go_quick_star_cli/internal/llm"
	"go_quick_star_cli/internal/tools"
	"go_quick_star_cli/internal/ui"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agent",
	Short: "A smart Coding Agent CLI",
	Long:  `A smart Coding Agent CLI that can help you write code, fix bugs, and review code by executing shell commands and file operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		startREPL()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startREPL() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		fmt.Println("Please make sure you have a .env file with OPENAI_API_KEY set.")
		return
	}

	client := llm.NewClient(cfg)
	toolSet := tools.NewToolSet(
		tools.NewShellRunner(),
		tools.NewFileReader(),
		tools.NewFileWriter(),
		tools.NewFileSearcher(),
		tools.NewGitStatus(),
		tools.NewGitDiff(),
		tools.NewGitCommit(),
	)

	cwd, _ := os.Getwd()
	memoryPath := filepath.Join(cwd, ".agent_memory.json")

	ag := agent.NewAgent(client, toolSet, memoryPath)

	// Setup readline
	historyFile := filepath.Join(cwd, ".agent_cli_history")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:            "> ",
		HistoryFile:       historyFile,
		HistorySearchFold: true,
		// Fix for cursor movement issues in some terminals
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	fmt.Println("Welcome to the Coding Agent CLI!")
	fmt.Println("You can ask me to write code, debug issues, or review files.")
	fmt.Println("Type 'exit' or 'quit' to leave. Use '\\' at the end of a line for multi-line input.")

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		input := strings.TrimSpace(line)

		// Handle multi-line input
		for strings.HasSuffix(input, "\\") {
			input = strings.TrimSuffix(input, "\\")
			rl.SetPrompt("... ")
			nextLine, err := rl.Readline()
			rl.SetPrompt("> ")
			if err != nil {
				break
			}
			input += "\n" + strings.TrimSpace(nextLine)
		}

		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			break
		}

		ctx := context.Background()

		// Start Spinner
		spinner := ui.NewSpinner()
		spinner.Start("Thinking...")

		response, err := ag.Run(ctx, input)

		// Stop Spinner
		spinner.Stop()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Render Markdown
		renderedResponse := ui.RenderMarkdown(response)
		fmt.Println(renderedResponse)
	}
}
