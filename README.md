# Go Quick Star CLI (Coding Agent)

A powerful, autonomous Coding Agent CLI built with Go. This tool acts as your intelligent pair programmer, capable of writing code, debugging issues, managing files, and handling git operations directly from your terminal.

[中文文档](README_CN.md)

## 🚀 Features

*   **Intelligent Coding Assistance**: Powered by Large Language Models (LLM), it understands your intent and helps you write code in various languages.
*   **Context Awareness**: Reads local files to understand your project structure and existing code.
*   **Codebase Search**: Efficiently searches your entire codebase for functions, variables, or patterns using `search_files`.
*   **Git Integration**: Checks status, views diffs, and commits changes with meaningful messages directly through the agent.
*   **Persistent Memory**: Remembers your conversation context across sessions (stored in `.agent_memory.json`).
*   **Rich Terminal UI**:
    *   Markdown rendering with syntax highlighting for code blocks.
    *   Animated spinner during "Thinking" phase.
    *   Command history and auto-completion (using `readline`).
    *   Multi-line input support.

## 🏗 Architecture

The project follows a clean, modular architecture:

```
.
├── cmd/
│   └── root.go           # Entry point, REPL loop, and CLI initialization
├── internal/
│   ├── agent/
│   │   ├── agent.go      # Core agent logic, conversation loop, and tool dispatching
│   │   └── memory.go     # Persistent memory management (JSON-based)
│   ├── config/
│   │   └── config.go     # Configuration loading (.env)
│   ├── llm/
│   │   └── client.go     # OpenAI-compatible LLM client wrapper
│   ├── tools/
│   │   ├── files.go      # File system tools (read, write, search)
│   │   ├── git.go        # Git operation tools (status, diff, commit)
│   │   ├── shell.go      # Shell command execution
│   │   └── toolset.go    # Tool definitions and registration
│   └── ui/
│       ├── markdown.go   # Markdown rendering for terminal
│       └── spinner.go    # Loading animation
├── main.go               # Main application entry
├── .env                  # Configuration file (API keys, etc.)
└── go.mod                # Go module definition
```

## 🛠 Supported Tools

The Agent has access to the following tools to perform tasks:

1.  **File Operations**:
    *   `read_file`: Read content of a specific file.
    *   `write_file`: Create or overwrite a file with content.
    *   `search_files`: Recursively search for text patterns in the codebase (grep-like).

2.  **Git Operations**:
    *   `git_status`: Check repository status.
    *   `git_diff`: View changes (staged or unstaged).
    *   `git_commit`: Commit changes with a message.

3.  **System**:
    *   `run_shell_command`: Execute arbitrary shell commands (ls, go run, etc.).

## 🏁 Getting Started

### Prerequisites

*   Go 1.21 or higher
*   An API Key for an OpenAI-compatible LLM provider (e.g., OpenAI, Minimax, DeepSeek).

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/go_quick_star_cli.git
    cd go_quick_star_cli
    ```

2.  Install dependencies:
    ```bash
    go mod tidy
    ```

### Configuration

1.  Create a `.env` file in the project root:
    ```bash
    cp .env.example .env # If example exists, otherwise create new
    ```

2.  Edit `.env` and configure your LLM provider:

    ```ini
    # Example for Minimax
    OPENAI_API_KEY=your_api_key_here
    OPENAI_BASE_URL=https://api.minimax.chat/v1
    OPENAI_MODEL=abab6.5s-chat
    
    # Optional: Enable debug logging
    # DEBUG=true
    ```

### Usage

Run the CLI:

```bash
go run main.go
```

Once inside the CLI, you can type natural language commands.

**Examples:**

*   **Writing Code**: "Write a Python script to calculate Fibonacci numbers."
*   **Debugging**: "Read main.go and check for error handling issues."
*   **Searching**: "Find where the `NewAgent` function is defined."
*   **Git**: "Check what I changed, then commit it with a message 'Update README'."

**Controls:**

*   `Up/Down Arrow`: Navigate command history.
*   `Ctrl+R`: Search history.
*   `\` (at end of line): Continue input on next line (Multi-line mode).
*   `exit` or `quit`: Exit the CLI.

## 📝 License

MIT
