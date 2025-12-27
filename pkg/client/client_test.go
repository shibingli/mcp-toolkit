package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mcp-toolkit/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestNewHTTPClient 测试创建HTTP客户端 / Test creating HTTP client
func TestNewHTTPClient(t *testing.T) {
	logger := zap.NewNop()
	client := NewHTTPClient("localhost", 8080, "/mcp", logger)

	assert.NotNil(t, client)
	assert.Equal(t, "http://localhost:8080/mcp", client.baseURL)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.logger)
}

// TestNewHTTPClient_NilLogger 测试使用nil logger创建客户端 / Test creating client with nil logger
func TestNewHTTPClient_NilLogger(t *testing.T) {
	client := NewHTTPClient("localhost", 8080, "/mcp", nil)

	assert.NotNil(t, client)
	assert.NotNil(t, client.logger)
}

// TestHTTPClient_Initialize 测试初始化 / Test initialize
func TestHTTPClient_Initialize(t *testing.T) {
	// 创建模拟服务器 / Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// 解析请求 / Parse request
		var req types.MCPRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "2.0", req.JSONRPC)
		assert.Equal(t, "initialize", req.Method)

		// 返回响应 / Return response
		resp := types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "test-server",
					"version": "1.0.0",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建客户端 / Create client
	logger := zap.NewNop()
	client := NewHTTPClient("", 0, "", logger)
	client.baseURL = server.URL

	// 测试初始化 / Test initialize
	ctx := context.Background()
	initResp, err := client.Initialize(ctx, types.ProtocolVersion)

	require.NoError(t, err)
	assert.NotNil(t, initResp)
	assert.Equal(t, types.ProtocolVersion, initResp.ProtocolVersion)
	assert.Equal(t, "test-server", initResp.ServerInfo.Name)
	assert.Equal(t, "1.0.0", initResp.ServerInfo.Version)
}

// TestHTTPClient_ListTools 测试列出工具 / Test list tools
func TestHTTPClient_ListTools(t *testing.T) {
	// 创建模拟服务器 / Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		// 解析请求 / Parse request
		var req types.MCPRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "tools/list", req.Method)

		// 返回响应 / Return response
		resp := types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": []map[string]interface{}{
					{
						"name":        "create_file",
						"description": "Create a new file",
					},
					{
						"name":        "read_file",
						"description": "Read file content",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建客户端 / Create client
	logger := zap.NewNop()
	client := NewHTTPClient("", 0, "", logger)
	client.baseURL = server.URL

	// 测试列出工具 / Test list tools
	ctx := context.Background()
	toolsResp, err := client.ListTools(ctx)

	require.NoError(t, err)
	assert.NotNil(t, toolsResp)
	assert.Len(t, toolsResp.Tools, 2)
	assert.Equal(t, "create_file", toolsResp.Tools[0].Name)
	assert.Equal(t, "read_file", toolsResp.Tools[1].Name)
}

// TestHTTPClient_CallTool 测试调用工具 / Test call tool
func TestHTTPClient_CallTool(t *testing.T) {
	// 创建模拟服务器 / Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		// 解析请求 / Parse request
		var req types.MCPRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "tools/call", req.Method)

		// 返回响应 / Return response
		resp := types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": `{"success":true,"message":"File created successfully"}`,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建客户端 / Create client
	logger := zap.NewNop()
	client := NewHTTPClient("", 0, "", logger)
	client.baseURL = server.URL

	// 测试调用工具 / Test call tool
	ctx := context.Background()
	callResp, err := client.CallTool(ctx, "create_file", types.CreateFileRequest{
		Path:    "test.txt",
		Content: "Hello",
	})

	require.NoError(t, err)
	assert.NotNil(t, callResp)
	assert.Len(t, callResp.Content, 1)
	assert.Equal(t, "text", callResp.Content[0].Type)
}

// TestHTTPClient_CallTool_Error 测试调用工具错误 / Test call tool error
func TestHTTPClient_CallTool_Error(t *testing.T) {
	// 创建模拟服务器 / Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 解析请求 / Parse request
		var req types.MCPRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		// 返回错误响应 / Return error response
		resp := types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &types.MCPError{
				Code:    types.MCPErrorCodeInvalidParams,
				Message: "Invalid parameters",
				Data:    "Path is required",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建客户端 / Create client
	logger := zap.NewNop()
	client := NewHTTPClient("", 0, "", logger)
	client.baseURL = server.URL

	// 测试调用工具 / Test call tool
	ctx := context.Background()
	_, err := client.CallTool(ctx, "create_file", map[string]interface{}{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "MCP error")
	assert.Contains(t, err.Error(), "Invalid parameters")
}

// TestHTTPClient_Close 测试关闭客户端 / Test close client
func TestHTTPClient_Close(t *testing.T) {
	logger := zap.NewNop()
	client := NewHTTPClient("localhost", 8080, "/mcp", logger)

	err := client.Close()
	assert.NoError(t, err)
}

// TestHTTPClient_SendRequest_InvalidJSON 测试发送无效JSON / Test send invalid JSON
func TestHTTPClient_SendRequest_InvalidJSON(t *testing.T) {
	// 创建模拟服务器 / Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	// 创建客户端 / Create client
	logger := zap.NewNop()
	client := NewHTTPClient("", 0, "", logger)
	client.baseURL = server.URL

	// 测试初始化 / Test initialize
	ctx := context.Background()
	_, err := client.Initialize(ctx, "2024-11-05")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal response")
}

// TestHTTPClient_SendRequest_NetworkError 测试网络错误 / Test network error
func TestHTTPClient_SendRequest_NetworkError(t *testing.T) {
	// 创建客户端连接到不存在的服务器 / Create client connecting to non-existent server
	logger := zap.NewNop()
	client := NewHTTPClient("localhost", 9999, "/mcp", logger)

	// 测试初始化 / Test initialize
	ctx := context.Background()
	_, err := client.Initialize(ctx, "2024-11-05")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send request")
}

// TestHTTPClient_GetSystemInfo 测试获取系统信息 / Test get system info
func TestHTTPClient_GetSystemInfo(t *testing.T) {
	// 创建模拟服务器 / Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.MCPRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "tools/call", req.Method)

		// 返回系统信息响应 / Return system info response
		resp := types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": `{"os":{"platform":"linux","architecture":"amd64","hostname":"test-host"},"cpu":{"logical_cores":8},"memory":{"total":16000000000}}`,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建客户端 / Create client
	logger := zap.NewNop()
	client := NewHTTPClient("", 0, "", logger)
	client.baseURL = server.URL

	// 测试获取系统信息 / Test get system info
	ctx := context.Background()
	result, err := client.CallTool(ctx, "get_system_info", map[string]interface{}{})

	require.NoError(t, err)
	assert.NotNil(t, result)
}
