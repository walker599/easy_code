package tools

import (
	"context"
	"encoding/json"
)

// Tool 定义了一个 Agent 可以调用的工具
type Tool interface {
	// Name 返回工具的名称，如 "run_shell"
	Name() string
	// Description 返回工具的描述，用于告诉 LLM 什么时候使用它
	Description() string
	// Parameters 返回工具的参数 schema (JSON Schema)
	Parameters() json.RawMessage
	// Execute 执行工具逻辑
	Execute(ctx context.Context, args json.RawMessage) (string, error)
}

// ToolSet 是工具的集合，方便查找
type ToolSet map[string]Tool

func NewToolSet(tools ...Tool) ToolSet {
	ts := make(ToolSet)
	for _, t := range tools {
		ts[t.Name()] = t
	}
	return ts
}
