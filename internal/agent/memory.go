package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/sashabaranov/go-openai"
)

type Memory struct {
	History []openai.ChatCompletionMessage `json:"history"`
	mu      sync.RWMutex
	path    string
}

func NewMemory(path string) *Memory {
	return &Memory{
		History: []openai.ChatCompletionMessage{},
		path:    path,
	}
}

func (m *Memory) Add(msg openai.ChatCompletionMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.History = append(m.History, msg)
	m.save()
}

func (m *Memory) Get() []openai.ChatCompletionMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Return a copy to avoid data races
	history := make([]openai.ChatCompletionMessage, len(m.History))
	copy(history, m.History)
	return history
}

func (m *Memory) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &m.History)
}

func (m *Memory) save() error {
	data, err := json.MarshalIndent(m.History, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(m.path, data, 0644)
}

func (m *Memory) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.History = []openai.ChatCompletionMessage{}
	return os.Remove(m.path)
}
