# Go Quick Star CLI (Coding Agent)

这是一个基于 Go 语言构建的强大、自主的 Coding Agent CLI。它就像你的智能结对程序员，能够直接在终端中帮你编写代码、调试问题、管理文件以及处理 Git 操作。

[English Documentation](README.md)

## 🚀 功能特性

*   **智能编程辅助**：基于大语言模型 (LLM)，理解你的意图并帮你编写多种语言的代码。
*   **上下文感知**：能够读取本地文件，理解你的项目结构和现有代码。
*   **代码库搜索**：使用 `search_files` 高效搜索整个代码库中的函数、变量或特定模式。
*   **Git 集成**：直接通过 Agent 查看状态、对比差异 (diff) 并提交代码（自动生成 Commit Message）。
*   **持久化记忆**：跨会话记住你的对话上下文（存储在 `.agent_memory.json` 中）。
*   **丰富的终端 UI**：
    *   Markdown 渲染：代码块支持语法高亮，排版精美。
    *   动态 Spinner：在“思考”阶段显示旋转动画。
    *   命令历史与自动补全：基于 `readline` 实现。
    *   支持多行输入。

## 🏗 架构设计

项目采用清晰的模块化架构：

```
.
├── cmd/
│   └── root.go           # 入口点，REPL 循环和 CLI 初始化
├── internal/
│   ├── agent/
│   │   ├── agent.go      # 核心 Agent 逻辑，对话循环和工具分发
│   │   └── memory.go     # 持久化记忆管理 (基于 JSON)
│   ├── config/
│   │   └── config.go     # 配置加载 (.env)
│   ├── llm/
│   │   └── client.go     # 兼容 OpenAI 接口的 LLM 客户端封装
│   ├── tools/
│   │   ├── files.go      # 文件系统工具 (读、写、搜)
│   │   ├── git.go        # Git 操作工具 (status, diff, commit)
│   │   ├── shell.go      # Shell 命令执行
│   │   └── toolset.go    # 工具定义与注册
│   └── ui/
│       ├── markdown.go   # 终端 Markdown 渲染
│       └── spinner.go    # 加载动画
├── main.go               # 应用程序主入口
├── .env                  # 配置文件 (API Key 等)
└── go.mod                # Go 模块定义
```

## 🛠 支持的工具 (Tools)

Agent 可以调用以下工具来完成任务：

1.  **文件操作**：
    *   `read_file`: 读取指定文件的内容。
    *   `write_file`: 创建文件或覆盖写入内容。
    *   `search_files`: 递归搜索代码库中的文本模式 (类似 grep)。

2.  **Git 操作**：
    *   `git_status`: 查看仓库状态。
    *   `git_diff`: 查看变更内容 (支持 staged/unstaged)。
    *   `git_commit`: 提交变更并附带消息。

3.  **系统操作**：
    *   `run_shell_command`: 执行任意 Shell 命令 (如 ls, go run 等)。

## 🏁 快速开始

### 前置要求

*   Go 1.21 或更高版本
*   一个兼容 OpenAI 接口的大模型 API Key (例如 OpenAI, Minimax, DeepSeek 等)。

### 安装步骤

1.  克隆仓库：
    ```bash
    git clone https://github.com/yourusername/go_quick_star_cli.git
    cd go_quick_star_cli
    ```

2.  安装依赖：
    ```bash
    go mod tidy
    ```

### 配置

1.  在项目根目录创建 `.env` 文件：
    ```bash
    cp .env.example .env # 如果有示例文件，否则新建一个
    ```

2.  编辑 `.env` 并配置你的 LLM 提供商：

    ```ini
    # Minimax 配置示例
    OPENAI_API_KEY=your_api_key_here
    OPENAI_BASE_URL=https://api.minimax.chat/v1
    OPENAI_MODEL=abab6.5s-chat
    
    # 可选：开启调试日志
    # DEBUG=true
    ```

### 使用方法

运行 CLI：

```bash
go run main.go
```

进入 CLI 后，你可以直接输入自然语言指令。

**交互示例：**

*   **写代码**: "帮我写一个 Python 脚本来计算斐波那契数列。"
*   **调试**: "读取 main.go 文件，帮我检查一下有没有错误处理缺失的问题。"
*   **搜索**: "帮我找一下 `NewAgent` 函数是在哪个文件里定义的。"
*   **Git**: "看一下我改了什么，然后帮我提交代码，Commit Message 写 '更新 README'。"

**操作控制：**

*   `上/下箭头`: 翻阅历史命令。
*   `Ctrl+R`: 搜索历史命令。
*   `\` (在行尾): 开启多行输入模式（换行继续输入）。
*   `exit` 或 `quit`: 退出 CLI。

## 📝 许可证

MIT
