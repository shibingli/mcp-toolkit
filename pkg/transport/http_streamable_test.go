package transport

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupStreamableHTTPServer 创建支持Streamable HTTP的测试服务器 / Create test server with Streamable HTTP support
func setupStreamableHTTPServer(t *testing.T) (*HTTPTransportServer, *mcp.Server) {
	logger := zap.NewNop()
	config := types.DefaultHTTPConfig()
	config.Port = 18081 // 使用不同的测试端口 / Use different test port
	config.EnableSessionManagement = true
	config.EnableSSE = true

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

func TestStreamableHTTP_ProtocolVersion(t *testing.T) {
	server, _ := setupStreamableHTTPServer(t)

	tests := []struct {
		name           string
		version        string
		expectError    bool
		expectedStatus int
	}{
		{
			name:           "latest version 2025-12-26",
			version:        "2025-12-26",
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "version 2025-06-18",
			version:        "2025-06-18",
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "version 2025-03-26",
			version:        "2025-03-26",
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "version 2024-11-05",
			version:        "2024-11-05",
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "no version header (backward compatible)",
			version:        "",
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "unsupported version",
			version:        "2020-01-01",
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			httpReq.Header.Set("Accept", "application/json")
			if tt.version != "" {
				httpReq.Header.Set("MCP-Protocol-Version", tt.version)
			}

			// 测试版本验证方法 / Test version validation method
			err = server.validateProtocolVersion(httpReq)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStreamableHTTP_SessionManagement(t *testing.T) {
	server, _ := setupStreamableHTTPServer(t)

	t.Run("initialize creates session", func(t *testing.T) {
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
		httpReq.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()

		server.handleMCPRequest(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		sessionID := w.Header().Get("Mcp-Session-Id")
		assert.NotEmpty(t, sessionID)

		var resp types.MCPResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Nil(t, resp.Error)
	})

	t.Run("request without session fails", func(t *testing.T) {
		req := types.MCPRequest{
			JSONRPC: "2.0",
			ID:      2,
			Method:  "tools/list",
		}

		reqJSON, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(reqJSON))
		httpReq.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()

		server.handleMCPRequest(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("concurrent session access", func(t *testing.T) {
		// 首先创建一个会话 / First create a session
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
		httpReq.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()

		server.handleMCPRequest(w, httpReq)
		sessionID := w.Header().Get("Mcp-Session-Id")
		require.NotEmpty(t, sessionID)

		// 并发访问同一个会话 / Concurrent access to the same session
		const numGoroutines = 10
		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				req := types.MCPRequest{
					JSONRPC: "2.0",
					ID:      id + 100,
					Method:  "tools/list",
				}

				reqJSON, err := json.Marshal(req)
				if err != nil {
					return
				}

				httpReq := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(reqJSON))
				httpReq.Header.Set("Content-Type", "application/json")
				httpReq.Header.Set("Accept", "application/json")
				httpReq.Header.Set("Mcp-Session-Id", sessionID)
				w := httptest.NewRecorder()

				server.handleMCPRequest(w, httpReq)
			}(i)
		}

		// 等待所有goroutine完成 / Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
	})
}

func TestStreamableHTTP_CORS(t *testing.T) {
	server, _ := setupStreamableHTTPServer(t)

	t.Run("CORS headers are set correctly", func(t *testing.T) {
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
		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("Origin", "http://example.com")
		w := httptest.NewRecorder()

		// 确保CORS已启用 / Ensure CORS is enabled
		server.config.EnableCORS = true
		server.setCORSHeaders(w, httpReq)

		assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Mcp-Session-Id")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "MCP-Protocol-Version")
		assert.Contains(t, w.Header().Get("Access-Control-Expose-Headers"), "Mcp-Session-Id")
	})

	t.Run("OPTIONS request returns correct headers", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodOptions, "/mcp", nil)
		httpReq.Header.Set("Origin", "http://example.com")
		w := httptest.NewRecorder()

		// 直接调用处理器 / Call handler directly
		server.config.EnableCORS = true
		server.setCORSHeaders(w, httpReq)

		assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "DELETE")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	})
}
