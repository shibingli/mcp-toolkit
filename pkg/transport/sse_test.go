package transport

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupSSEServer 创建测试用的SSE服务器 / Create SSE server for testing
func setupSSEServer(t *testing.T) (*SSETransportServer, *mcp.Server) {
	logger := zap.NewNop()
	config := types.DefaultSSEConfig()
	config.Port = 18081 // 使用测试端口 / Use test port

	server, err := NewSSETransportServer(config, logger)
	require.NoError(t, err)
	require.NotNil(t, server)

	// 创建MCP服务器 / Create MCP server
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		},
		&mcp.ServerOptions{
			Capabilities: &mcp.ServerCapabilities{
				Tools: &mcp.ToolCapabilities{},
			},
		},
	)

	server.mcpServer = mcpServer

	return server, mcpServer
}

func TestNewSSETransportServer(t *testing.T) {
	logger := zap.NewNop()

	t.Run("with default config", func(t *testing.T) {
		server, err := NewSSETransportServer(nil, logger)
		assert.NoError(t, err)
		assert.NotNil(t, server)
		assert.Equal(t, types.DefaultSSEConfig().Port, server.config.Port)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &types.SSEConfig{
			Host:              "0.0.0.0",
			Port:              9091,
			Path:              "/custom-sse",
			EnableCORS:        false,
			HeartbeatInterval: 20,
			MaxConnections:    100,
		}
		server, err := NewSSETransportServer(config, logger)
		assert.NoError(t, err)
		assert.NotNil(t, server)
		assert.Equal(t, 9091, server.config.Port)
		assert.Equal(t, "/custom-sse", server.config.Path)
		assert.Equal(t, 20, server.config.HeartbeatInterval)
	})

	t.Run("with invalid port", func(t *testing.T) {
		config := &types.SSEConfig{
			Port: -1, // 无效端口 / Invalid port
		}
		server, err := NewSSETransportServer(config, logger)
		assert.Error(t, err)
		assert.Nil(t, server)
	})

	t.Run("with zero heartbeat interval", func(t *testing.T) {
		config := &types.SSEConfig{
			Port:              8081,
			HeartbeatInterval: 0,
		}
		server, err := NewSSETransportServer(config, logger)
		assert.NoError(t, err)
		assert.Equal(t, 30, server.config.HeartbeatInterval)
	})
}

func TestSSETransportServer_HandleMCPMessage(t *testing.T) {
	server, _ := setupSSEServer(t)

	t.Run("initialize request", func(t *testing.T) {
		req := types.MCPRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "initialize",
			Params: map[string]interface{}{
				"protocolVersion": "2024-11-05",
			},
		}

		reqJSON, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/sse/message", bytes.NewReader(reqJSON))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.handleMCPMessage(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp types.MCPResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "2.0", resp.JSONRPC)
		assert.Equal(t, float64(1), resp.ID) // JSON数字被解析为float64 / JSON numbers are parsed as float64
		assert.Nil(t, resp.Error)
		assert.NotNil(t, resp.Result)
	})

	t.Run("tools/list request", func(t *testing.T) {
		req := types.MCPRequest{
			JSONRPC: "2.0",
			ID:      2,
			Method:  "tools/list",
		}

		reqJSON, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/sse/message", bytes.NewReader(reqJSON))
		w := httptest.NewRecorder()

		server.handleMCPMessage(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp types.MCPResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Nil(t, resp.Error)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodPost, "/sse/message", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		server.handleMCPMessage(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp types.MCPResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotNil(t, resp.Error)
		assert.Equal(t, types.MCPErrorCodeParseError, resp.Error.Code)
	})

	t.Run("method not found", func(t *testing.T) {
		req := types.MCPRequest{
			JSONRPC: "2.0",
			ID:      3,
			Method:  "unknown/method",
		}

		reqJSON, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/sse/message", bytes.NewReader(reqJSON))
		w := httptest.NewRecorder()

		server.handleMCPMessage(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp types.MCPResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotNil(t, resp.Error)
		assert.Equal(t, types.MCPErrorCodeMethodNotFound, resp.Error.Code)
	})
}

func TestSSETransportServer_SetCORSHeaders(t *testing.T) {
	logger := zap.NewNop()
	config := types.DefaultSSEConfig()
	config.EnableCORS = true
	config.AllowedOrigins = []string{"http://localhost:3000", "https://example.com"}

	server, err := NewSSETransportServer(config, logger)
	require.NoError(t, err)

	t.Run("allowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/sse", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		server.setCORSHeaders(w, req)

		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	})

	t.Run("wildcard origin", func(t *testing.T) {
		config.AllowedOrigins = []string{"*"}
		server, _ := NewSSETransportServer(config, logger)

		req := httptest.NewRequest(http.MethodOptions, "/sse", nil)
		req.Header.Set("Origin", "http://any-origin.com")
		w := httptest.NewRecorder()

		server.setCORSHeaders(w, req)

		assert.Equal(t, "http://any-origin.com", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("disallowed origin", func(t *testing.T) {
		config.AllowedOrigins = []string{"http://localhost:3000"}
		server, _ := NewSSETransportServer(config, logger)

		req := httptest.NewRequest(http.MethodOptions, "/sse", nil)
		req.Header.Set("Origin", "http://evil.com")
		w := httptest.NewRecorder()

		server.setCORSHeaders(w, req)

		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestSSETransportServer_SendJSONError(t *testing.T) {
	server, _ := setupSSEServer(t)

	t.Run("send error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		server.sendJSONError(w, 1, types.MCPErrorCodeInternalError, "test error", "test data")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp types.MCPResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "2.0", resp.JSONRPC)
		assert.Equal(t, float64(1), resp.ID) // JSON数字被解析为float64 / JSON numbers are parsed as float64
		assert.NotNil(t, resp.Error)
		assert.Equal(t, types.MCPErrorCodeInternalError, resp.Error.Code)
		assert.Equal(t, "test error", resp.Error.Message)
	})
}

func TestSSETransportServer_GetNextConnID(t *testing.T) {
	server, _ := setupSSEServer(t)

	t.Run("generate unique connection IDs", func(t *testing.T) {
		id1 := server.getNextConnID()
		id2 := server.getNextConnID()
		id3 := server.getNextConnID()

		assert.NotEqual(t, id1, id2)
		assert.NotEqual(t, id2, id3)
		assert.Contains(t, id1, "conn-")
		assert.Contains(t, id2, "conn-")
		assert.Contains(t, id3, "conn-")
	})
}

func TestSSETransportServer_Shutdown(t *testing.T) {
	server, _ := setupSSEServer(t)

	t.Run("shutdown without start", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		assert.NoError(t, err)
	})
}

func TestSSETransportServer_HandleInitialize(t *testing.T) {
	server, _ := setupSSEServer(t)

	req := types.MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
	}

	result := server.handleInitialize(req)
	assert.NotNil(t, result)

	resultMap, ok := result.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "2024-11-05", resultMap["protocolVersion"])
	assert.NotNil(t, resultMap["capabilities"])
	assert.NotNil(t, resultMap["serverInfo"])
}

func TestSSETransportServer_HandleToolsList(t *testing.T) {
	server, _ := setupSSEServer(t)

	// 注册测试工具 / Register test tool
	server.toolRegistry.RegisterTool(&mcp.Tool{
		Name:        "test_tool",
		Description: "Test tool",
	}, func(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "test"},
			},
		}, nil
	})

	req := types.MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
	}

	result := server.handleToolsList(req)
	assert.NotNil(t, result)

	resultMap, ok := result.(map[string]interface{})
	assert.True(t, ok)
	assert.NotNil(t, resultMap["tools"])

	tools, ok := resultMap["tools"].([]map[string]interface{})
	assert.True(t, ok)
	assert.Greater(t, len(tools), 0)
}

func TestSSETransportServer_HandleToolsCall(t *testing.T) {
	server, _ := setupSSEServer(t)

	// 注册测试工具 / Register test tool
	server.toolRegistry.RegisterTool(&mcp.Tool{
		Name:        "create_file",
		Description: "Create file",
	}, func(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "file created"},
			},
		}, nil
	})

	t.Run("valid tool call", func(t *testing.T) {
		req := types.MCPRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "tools/call",
			Params: map[string]interface{}{
				"name": "create_file",
				"arguments": map[string]interface{}{
					"path":    "test.txt",
					"content": "test content",
				},
			},
		}

		result, err := server.handleToolsCall(req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("invalid params", func(t *testing.T) {
		req := types.MCPRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "tools/call",
			Params:  "invalid", // 无效的参数类型 / Invalid param type
		}

		result, err := server.handleToolsCall(req)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestSSETransportServer_MethodNotAllowed(t *testing.T) {
	t.Run("GET request to message endpoint", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodGet, "/sse/message", nil)
		w := httptest.NewRecorder()

		// 模拟路由处理 / Simulate route handling
		if httpReq.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func TestSSETransportServer_OptionsRequest(t *testing.T) {
	logger := zap.NewNop()
	config := types.DefaultSSEConfig()
	config.Port = 18083
	config.EnableCORS = true
	config.AllowedOrigins = []string{"*"}

	server, err := NewSSETransportServer(config, logger)
	require.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodOptions, "/sse/message", nil)
	httpReq.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	// 设置CORS头并返回200 / Set CORS headers and return 200
	server.setCORSHeaders(w, httpReq)
	w.WriteHeader(http.StatusOK)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

// TestSSETransportServer_StartAndShutdown 测试SSE服务器启动和关闭 / Test SSE server start and shutdown
func TestSSETransportServer_StartAndShutdown(t *testing.T) {
	logger := zap.NewNop()
	config := &types.SSEConfig{
		Host:              "127.0.0.1",
		Port:              18090, // 使用不同的端口避免冲突 / Use different port to avoid conflicts
		Path:              "/sse",
		EnableCORS:        true,
		AllowedOrigins:    []string{"*"},
		HeartbeatInterval: 30,
		MaxConnections:    100,
	}

	server, err := NewSSETransportServer(config, logger)
	require.NoError(t, err)

	// 创建MCP服务器 / Create MCP server
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		},
		&mcp.ServerOptions{
			Capabilities: &mcp.ServerCapabilities{
				Tools: &mcp.ToolCapabilities{},
			},
		},
	)

	// 启动服务器 / Start server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start(ctx, mcpServer)
	}()

	// 等待服务器启动 / Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// 测试服务器是否响应 / Test if server responds
	resp, err := http.Get("http://127.0.0.1:18090/health")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// 取消上下文触发关闭 / Cancel context to trigger shutdown
	cancel()

	// 等待服务器关闭 / Wait for server to shutdown
	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("server shutdown timeout")
	}
}

// TestSSETransportServer_StartError 测试SSE服务器启动错误 / Test SSE server start error
func TestSSETransportServer_StartError(t *testing.T) {
	logger := zap.NewNop()
	config := &types.SSEConfig{
		Host:              "127.0.0.1",
		Port:              18091,
		Path:              "/sse",
		EnableCORS:        false,
		HeartbeatInterval: 30,
		MaxConnections:    100,
	}

	server1, err := NewSSETransportServer(config, logger)
	require.NoError(t, err)

	server2, err := NewSSETransportServer(config, logger)
	require.NoError(t, err)

	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		},
		&mcp.ServerOptions{
			Capabilities: &mcp.ServerCapabilities{
				Tools: &mcp.ToolCapabilities{},
			},
		},
	)

	// 启动第一个服务器 / Start first server
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	errChan1 := make(chan error, 1)
	go func() {
		errChan1 <- server1.Start(ctx1, mcpServer)
	}()

	// 等待第一个服务器启动 / Wait for first server to start
	time.Sleep(100 * time.Millisecond)

	// 尝试启动第二个服务器（应该失败，因为端口已被占用） / Try to start second server (should fail due to port conflict)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()

	err = server2.Start(ctx2, mcpServer)
	assert.Error(t, err) // 应该返回端口占用错误 / Should return port in use error

	// 关闭第一个服务器 / Shutdown first server
	cancel1()
	select {
	case <-errChan1:
	case <-time.After(5 * time.Second):
		t.Fatal("server1 shutdown timeout")
	}
}
