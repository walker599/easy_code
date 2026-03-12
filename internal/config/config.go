package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIKey     string
	OpenAIBaseURL string
	Model         string
	Debug         bool
}

func Load() (*Config, error) {
	// 尝试加载 .env 文件，如果不存在也不报错（可能通过环境变量直接传入）
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	// fmt.Println(apiKey) // Remove sensitive info printing
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	debug := os.Getenv("DEBUG") == "true"

	return &Config{
		OpenAIKey:     apiKey,
		OpenAIBaseURL: baseURL,
		Model:         model,
		Debug:         debug,
	}, nil
}
