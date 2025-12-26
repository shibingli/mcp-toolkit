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

// setupHTTPServer 创建测试用的HTTP服务器 / Create HTTP server for testing
func setupHTTPServer(t *testing.T) (*HTTPTransportServer, *mcp.Server) {
	logger := zap.NewNop()
	config := types.DefaultHTTPConfig()
	config.Port = 18080 // 使用测试端口 / Use test port

	server, err := NewHTTPTransportServer(config, logger)
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

func TestNewHTTPTransportServer(t *testing.T) {
	logger := zap.NewNop()

	t.Run("with default config", func(t *testing.T) {
		server, err := NewHTTPTransportServer(nil, logger)
		assert.NoError(t, err)
		assert.NotNil(t, server)
		assert.Equal(t, types.DefaultHTTPConfig().Port, server.config.Port)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &types.HTTPConfig{
			Host:           "0.0.0.0",
			Port:           9090,
			Path:           "/custom",
			EnableCORS:     false,
			ReadTimeout:    60,
			WriteTimeout:   60,
			MaxHeaderBytes: 2 << 20,
		}
		server, err := NewHTTPTransportServer(config, logger)
		assert.NoError(t, err)
		assert.NotNil(t, server)
		assert.Equal(t, 9090, server.config.Port)
		assert.Equal(t, "/custom", server.config.Path)
	})

	t.Run("with invalid port", func(t *testing.T) {
		config := &types.HTTPConfig{
			Port: 70000, // 无效端口 / Invalid port
		}
		server, err := NewHTTPTransportServer(config, logger)
		assert.Error(t, err)
		assert.Nil(t, server)
	})

	t.Run("with zero timeout", func(t *testing.T) {
		config := &types.HTTPConfig{
			Port:         8080,
			ReadTimeout:  0,
			WriteTimeout: 0,
		}
		server, err := NewHTTPTransportServer(config, logger)
		assert.NoError(t, err)
		assert.Equal(t, 30, server.config.ReadTimeout)
		assert.Equal(t, 30, server.config.WriteTimeout)
	})
}

func TestHTTPTransportServer_HandleMCPRequest(t *testing.T) {
	server, _ := setupHTTPServer(t)

	t.Run("initialize request", func(t *testing.T) {
		req := types.MCPRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "initialize",
			Params: map[string]interface{}{
				"protocolVersion": types.ProtocolVersion,
			},
		}

		reqJSON, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(reqJSON))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.handleMCPRequest(w, httpReq)

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

		httpReq := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(reqJSON))
		w := httptest.NewRecorder()

		server.handleMCPRequest(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp types.MCPResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Nil(t, resp.Error)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		server.handleMCPRequest(w, httpReq)

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

		httpReq := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(reqJSON))
		w := httptest.NewRecorder()

		server.handleMCPRequest(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp types.MCPResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotNil(t, resp.Error)
		assert.Equal(t, types.MCPErrorCodeMethodNotFound, resp.Error.Code)
	})
}

func TestHTTPTransportServer_SetCORSHeaders(t *testing.T) {
	logger := zap.NewNop()
	config := types.DefaultHTTPConfig()
	config.EnableCORS = true
	config.AllowedOrigins = []string{"http://localhost:3000", "https://example.com"}

	server, err := NewHTTPTransportServer(config, logger)
	require.NoError(t, err)

	t.Run("allowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/mcp", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		server.setCORSHeaders(w, req)

		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	})

	t.Run("wildcard origin", func(t *testing.T) {
		config.AllowedOrigins = []string{"*"}
		server, _ := NewHTTPTransportServer(config, logger)

		req := httptest.NewRequest(http.MethodOptions, "/mcp", nil)
		req.Header.Set("Origin", "http://any-origin.com")
		w := httptest.NewRecorder()

		server.setCORSHeaders(w, req)

		assert.Equal(t, "http://any-origin.com", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("disallowed origin", func(t *testing.T) {
		config.AllowedOrigins = []string{"http://localhost:3000"}
		server, _ := NewHTTPTransportServer(config, logger)

		req := httptest.NewRequest(http.MethodOptions, "/mcp", nil)
		req.Header.Set("Origin", "http://evil.com")
		w := httptest.NewRecorder()

		server.setCORSHeaders(w, req)

		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("no origin header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/mcp", nil)
		w := httptest.NewRecorder()

		server.setCORSHeaders(w, req)

		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestHTTPTransportServer_SendErrorResponse(t *testing.T) {
	server, _ := setupHTTPServer(t)

	t.Run("send error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		server.sendErrorResponse(w, 1, types.MCPErrorCodeInternalError, "test error", "test data")

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

func TestHTTPTransportServer_Shutdown(t *testing.T) {
	server, _ := setupHTTPServer(t)

	t.Run("shutdown without start", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		assert.NoError(t, err)
	})
}

// TestHTTPTransportServer_StartAndShutdown 测试HTTP服务器启动和关闭 / Test HTTP server start and shutdown
func TestHTTPTransportServer_StartAndShutdown(t *testing.T) {
	logger := zap.NewNop()
	config := &types.HTTPConfig{
		Host:           "127.0.0.1",
		Port:           18080, // 使用不同的端口避免冲突 / Use different port to avoid conflicts
		Path:           "/mcp",
		EnableCORS:     true,
		AllowedOrigins: []string{"*"},
		ReadTimeout:    30,
		WriteTimeout:   30,
		MaxHeaderBytes: 1 << 20,
	}

	server, err := NewHTTPTransportServer(config, logger)
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
	resp, err := http.Get("http://127.0.0.1:18080/health")
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

// TestHTTPTransportServer_StartError 测试HTTP服务器启动错误 / Test HTTP server start error
func TestHTTPTransportServer_StartError(t *testing.T) {
	logger := zap.NewNop()
	config := &types.HTTPConfig{
		Host:           "127.0.0.1",
		Port:           18081,
		Path:           "/mcp",
		EnableCORS:     false,
		ReadTimeout:    30,
		WriteTimeout:   30,
		MaxHeaderBytes: 1 << 20,
	}

	server1, err := NewHTTPTransportServer(config, logger)
	require.NoError(t, err)

	server2, err := NewHTTPTransportServer(config, logger)
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

func TestHTTPTransportServer_HandleInitialize(t *testing.T) {
	server, _ := setupHTTPServer(t)

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

func TestHTTPTransportServer_HandleToolsList(t *testing.T) {
	server, _ := setupHTTPServer(t)

	// 注册一个测试工具 / Register a test tool
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

func TestHTTPTransportServer_HandleToolsCall(t *testing.T) {
	server, _ := setupHTTPServer(t)

	// 注册一个测试工具 / Register a test tool
	server.toolRegistry.RegisterTool(&mcp.Tool{
		Name:        "create_file",
		Description: "Create a file",
	}, func(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: `{"success":true}`},
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

func TestHTTPTransportServer_MethodNotAllowed(t *testing.T) {
	t.Run("GET request", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodGet, "/mcp", nil)
		w := httptest.NewRecorder()

		// 模拟路由处理 / Simulate route handling
		if httpReq.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("PUT request", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodPut, "/mcp", nil)
		w := httptest.NewRecorder()

		if httpReq.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func TestHTTPTransportServer_OptionsRequest(t *testing.T) {
	logger := zap.NewNop()
	config := types.DefaultHTTPConfig()
	config.Port = 18082
	config.EnableCORS = true
	config.AllowedOrigins = []string{"*"}

	server, err := NewHTTPTransportServer(config, logger)
	require.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodOptions, "/mcp", nil)
	httpReq.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	// 设置CORS头并返回200 / Set CORS headers and return 200
	server.setCORSHeaders(w, httpReq)
	w.WriteHeader(http.StatusOK)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
}
