package transport

import (
	"context"
	"errors"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// setupTestRegistry 创建测试用的工具注册表 / Create test tool registry
func setupTestRegistry() (*ToolRegistry, *observer.ObservedLogs) {
	core, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)

	registry := NewToolRegistry()
	registry.SetLogger(logger)

	return registry, logs
}

func TestToolRegistry_CallTool_NoPanic(t *testing.T) {
	registry, logs := setupTestRegistry()

	// 注册一个正常的工具 / Register a normal tool
	tool := &mcp.Tool{
		Name:        "test_tool",
		Description: "Test tool",
	}

	handler := func(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "success"},
			},
		}, nil
	}

	registry.RegisterTool(tool, handler)

	// 调用工具 / Call tool
	result, err := registry.CallTool(context.Background(), "test_tool", nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, logs.Len(), "should not log when no panic")
}

func TestToolRegistry_CallTool_WithError(t *testing.T) {
	registry, logs := setupTestRegistry()

	// 注册一个返回错误的工具 / Register a tool that returns error
	tool := &mcp.Tool{
		Name:        "error_tool",
		Description: "Error tool",
	}

	expectedErr := errors.New("tool error")
	handler := func(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
		return nil, expectedErr
	}

	registry.RegisterTool(tool, handler)

	// 调用工具 / Call tool
	result, err := registry.CallTool(context.Background(), "error_tool", nil)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, logs.Len(), "should not log when returning error normally")
}

func TestToolRegistry_CallTool_PanicRecovery(t *testing.T) {
	registry, logs := setupTestRegistry()

	// 注册一个会panic的工具 / Register a tool that panics
	tool := &mcp.Tool{
		Name:        "panic_tool",
		Description: "Panic tool",
	}

	handler := func(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
		panic("tool panic")
	}

	registry.RegisterTool(tool, handler)

	// 调用工具 / Call tool
	result, err := registry.CallTool(context.Background(), "panic_tool", nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "panic: tool panic")
	assert.Nil(t, result)
	assert.Equal(t, 1, logs.Len(), "should log panic")

	logEntry := logs.All()[0]
	assert.Equal(t, "panic recovered", logEntry.Message)
}

func TestToolRegistry_CallTool_PanicWithError(t *testing.T) {
	registry, logs := setupTestRegistry()

	// 注册一个panic error的工具 / Register a tool that panics with error
	tool := &mcp.Tool{
		Name:        "panic_error_tool",
		Description: "Panic error tool",
	}

	panicErr := errors.New("panic error")
	handler := func(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
		panic(panicErr)
	}

	registry.RegisterTool(tool, handler)

	// 调用工具 / Call tool
	result, err := registry.CallTool(context.Background(), "panic_error_tool", nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "panic:")
	assert.ErrorIs(t, err, panicErr)
	assert.Nil(t, result)
	assert.Equal(t, 1, logs.Len(), "should log panic")
}

func TestToolRegistry_CallTool_NotFound(t *testing.T) {
	registry, logs := setupTestRegistry()

	// 调用不存在的工具 / Call non-existent tool
	result, err := registry.CallTool(context.Background(), "non_existent", nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "tool not found")
	assert.Nil(t, result)
	assert.Equal(t, 0, logs.Len(), "should not log for not found error")
}

func TestToolRegistry_CallTool_WithoutRecoveryHandler(t *testing.T) {
	// 创建没有设置logger的注册表 / Create registry without logger
	registry := NewToolRegistry()

	// 注册一个正常的工具 / Register a normal tool
	tool := &mcp.Tool{
		Name:        "test_tool",
		Description: "Test tool",
	}

	handler := func(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "success"},
			},
		}, nil
	}

	registry.RegisterTool(tool, handler)

	// 调用工具应该正常工作 / Call tool should work normally
	result, err := registry.CallTool(context.Background(), "test_tool", nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
