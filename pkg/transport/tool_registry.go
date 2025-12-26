// Copyright 2024 MCP Toolkit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transport

import (
	"context"
	"fmt"
	"sync"

	"mcp-toolkit/pkg/utils/recovery"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// ToolHandler 工具处理器函数类型 / Tool handler function type
type ToolHandler func(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error)

// ToolRegistry 工具注册表 / Tool registry
type ToolRegistry struct {
	mu              sync.RWMutex
	handlers        map[string]ToolHandler
	tools           map[string]*mcp.Tool
	recoveryHandler *recovery.RecoveryHandler
}

// NewToolRegistry 创建工具注册表 / Create tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		handlers:        make(map[string]ToolHandler),
		tools:           make(map[string]*mcp.Tool),
		recoveryHandler: nil, // 将在SetLogger中设置 / Will be set in SetLogger
	}
}

// SetLogger 设置logger并初始化recovery handler / Set logger and initialize recovery handler
func (r *ToolRegistry) SetLogger(logger *zap.Logger) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recoveryHandler = recovery.NewRecoveryHandler(logger)
}

// RegisterTool 注册工具 / Register tool
func (r *ToolRegistry) RegisterTool(tool *mcp.Tool, handler ToolHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tools[tool.Name] = tool
	r.handlers[tool.Name] = handler
}

// CallTool 调用工具 / Call tool
// 自动捕获并恢复panic，防止单个工具异常导致整个服务崩溃
func (r *ToolRegistry) CallTool(ctx context.Context, name string, arguments interface{}) (*mcp.CallToolResult, error) {
	r.mu.RLock()
	handler, exists := r.handlers[name]
	recoveryHandler := r.recoveryHandler
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	// 如果没有设置recovery handler，直接调用 / If recovery handler not set, call directly
	if recoveryHandler == nil {
		return handler(ctx, arguments)
	}

	// 使用recovery handler包装调用 / Wrap call with recovery handler
	result, err := recoveryHandler.RecoverWithValue(func() (interface{}, error) {
		return handler(ctx, arguments)
	})

	if err != nil {
		return nil, err
	}

	// 类型断言结果 / Type assert result
	if result == nil {
		return nil, nil
	}

	callResult, ok := result.(*mcp.CallToolResult)
	if !ok {
		return nil, fmt.Errorf("invalid result type from tool %s", name)
	}

	return callResult, nil
}

// ListTools 列出所有工具 / List all tools
func (r *ToolRegistry) ListTools() []*mcp.Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]*mcp.Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}

	return tools
}

// GetTool 获取工具 / Get tool
func (r *ToolRegistry) GetTool(name string) (*mcp.Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	return tool, exists
}
