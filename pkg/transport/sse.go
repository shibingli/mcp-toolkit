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
	"io"
	"mcp-toolkit/pkg/utils/recovery"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// SSEConnection SSE连接 / SSE connection
type SSEConnection struct {
	id       string
	writer   http.ResponseWriter
	flusher  http.Flusher
	messages chan *types.SSEMessage
	done     chan struct{}
}

// SSETransportServer SSE传输服务器 / SSE transport server
type SSETransportServer struct {
	config       *types.SSEConfig
	logger       *zap.Logger
	server       *http.Server
	mcpServer    *mcp.Server
	toolRegistry *ToolRegistry
	connections  sync.Map // 存储活跃连接 / Store active connections
	connCount    int32    // 连接计数 / Connection count
}

// NewSSETransportServer 创建SSE传输服务器 / Create SSE transport server
func NewSSETransportServer(config *types.SSEConfig, logger *zap.Logger) (*SSETransportServer, error) {
	if config == nil {
		config = types.DefaultSSEConfig()
	}

	// 验证配置 / Validate configuration
	if config.Port <= 0 || config.Port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", config.Port)
	}

	if config.HeartbeatInterval <= 0 {
		config.HeartbeatInterval = 30
	}

	if config.MaxConnections <= 0 {
		config.MaxConnections = 100
	}

	// 创建工具注册表 / Create tool registry
	toolRegistry := NewToolRegistry()
	toolRegistry.SetLogger(logger)

	return &SSETransportServer{
		config:       config,
		logger:       logger,
		mcpServer:    nil, // 将在Start方法中设置 / Will be set in Start method
		toolRegistry: toolRegistry,
	}, nil
}

// GetToolRegistry 获取工具注册表 / Get tool registry
func (s *SSETransportServer) GetToolRegistry() *ToolRegistry {
	return s.toolRegistry
}

// Start 启动SSE服务器 / Start SSE server
func (s *SSETransportServer) Start(ctx context.Context, mcpServer *mcp.Server) error {
	s.mcpServer = mcpServer
	mux := http.NewServeMux()

	// 注册MCP消息处理端点(POST) / Register MCP message handler endpoint (POST)
	mux.HandleFunc(s.config.Path+"/message", func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头 / Set CORS headers
		if s.config.EnableCORS {
			s.setCORSHeaders(w, r)
		}

		// 处理OPTIONS预检请求 / Handle OPTIONS preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 仅允许POST请求 / Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 处理MCP请求 / Handle MCP request
		s.handleMCPMessage(w, r)
	})

	// 注册SSE处理器 / Register SSE handler
	mux.HandleFunc(s.config.Path, func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头 / Set CORS headers
		if s.config.EnableCORS {
			s.setCORSHeaders(w, r)
		}

		// 处理OPTIONS预检请求 / Handle OPTIONS preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 仅允许GET请求 / Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 设置SSE响应头 / Set SSE response headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 创建flusher / Create flusher
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		s.logger.Info("new SSE connection",
			zap.String("remote_addr", r.RemoteAddr))

		// 发送心跳 / Send heartbeat
		ticker := time.NewTicker(time.Duration(s.config.HeartbeatInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-r.Context().Done():
				s.logger.Info("SSE connection closed",
					zap.String("remote_addr", r.RemoteAddr))
				return
			case <-ticker.C:
				// 发送心跳事件 / Send heartbeat event
				_, err := fmt.Fprintf(w, ": heartbeat\n\n")
				if err != nil {
					s.logger.Error("failed to send heartbeat",
						zap.Error(err))
					return
				}
				flusher.Flush()
			}
		}
	})

	// 健康检查端点 / Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	s.logger.Info("starting SSE transport server",
		zap.String("address", addr),
		zap.String("path", s.config.Path))

	// 在goroutine中启动服务器 / Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// 等待上下文取消或错误 / Wait for context cancellation or error
	select {
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

// handleMCPMessage 处理MCP消息 / Handle MCP message
func (s *SSETransportServer) handleMCPMessage(w http.ResponseWriter, r *http.Request) {
	// 添加panic recovery保护 / Add panic recovery protection
	recoveryHandler := recovery.NewRecoveryHandler(s.logger)
	err := recoveryHandler.Recover(func() error {
		return s.handleMCPMessageInternal(w, r)
	})

	if err != nil {
		// 如果是panic导致的错误，发送内部错误响应 / If error from panic, send internal error response
		s.logger.Error("panic in MCP message handler", zap.Error(err))
		s.sendJSONError(w, nil, types.MCPErrorCodeInternalError, "Internal server error", err.Error())
	}
}

// handleMCPMessageInternal 内部MCP消息处理 / Internal MCP message handling
func (s *SSETransportServer) handleMCPMessageInternal(w http.ResponseWriter, r *http.Request) error {
	// 读取请求体 / Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("failed to read request body", zap.Error(err))
		s.sendJSONError(w, nil, types.MCPErrorCodeParseError, types.MCPErrorMsgParseError, err.Error())
		return nil
	}
	defer func() { _ = r.Body.Close() }()

	// 解析MCP请求 / Parse MCP request
	var mcpReq types.MCPRequest
	if err = json.Unmarshal(body, &mcpReq); err != nil {
		s.logger.Error("failed to parse MCP request", zap.Error(err))
		s.sendJSONError(w, nil, types.MCPErrorCodeParseError, types.MCPErrorMsgParseError, err.Error())
		return nil
	}

	s.logger.Info("received MCP message",
		zap.String("method", mcpReq.Method),
		zap.Any("id", mcpReq.ID),
		zap.String("remote_addr", r.RemoteAddr))

	// 根据方法路由请求 / Route request based on method
	var result interface{}
	switch mcpReq.Method {
	case "initialize":
		result = s.handleInitialize(mcpReq)
	case "tools/list":
		result = s.handleToolsList(mcpReq)
	case "tools/call":
		result, err = s.handleToolsCall(mcpReq)
		if err != nil {
			s.sendJSONError(w, mcpReq.ID, types.MCPErrorCodeInternalError, err.Error(), nil)
			return nil
		}
	default:
		s.sendJSONError(w, mcpReq.ID, types.MCPErrorCodeMethodNotFound, types.MCPErrorMsgMethodNotFound, mcpReq.Method)
		return nil
	}

	// 发送成功响应 / Send success response
	resp := types.NewMCPResponse(mcpReq.ID, result)
	respJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("failed to marshal response", zap.Error(err))
		s.sendJSONError(w, mcpReq.ID, types.MCPErrorCodeInternalError, types.MCPErrorMsgInternalError, err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJSON)
	return nil
}

// handleInitialize 处理初始化请求 / Handle initialize request
func (s *SSETransportServer) handleInitialize(req types.MCPRequest) interface{} {
	return map[string]interface{}{
		"protocolVersion": types.ProtocolVersion,
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    types.ServerName,
			"version": types.ServerVersion,
		},
	}
}

// handleToolsList 处理工具列表请求 / Handle tools list request
func (s *SSETransportServer) handleToolsList(req types.MCPRequest) interface{} {
	// 从工具注册表获取所有工具 / Get all tools from registry
	registeredTools := s.toolRegistry.ListTools()

	tools := make([]map[string]interface{}, 0, len(registeredTools))
	for _, tool := range registeredTools {
		toolMap := map[string]interface{}{
			"name":        tool.Name,
			"description": tool.Description,
		}
		if tool.InputSchema != nil {
			toolMap["inputSchema"] = tool.InputSchema
		}
		tools = append(tools, toolMap)
	}

	return map[string]interface{}{
		"tools": tools,
	}
}

// handleToolsCall 处理工具调用请求 / Handle tools call request
func (s *SSETransportServer) handleToolsCall(req types.MCPRequest) (interface{}, error) {
	// 解析工具调用参数 / Parse tool call params
	paramsJSON, err := json.Marshal(req.Params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	var callReq types.MCPToolsCallRequest
	if err = json.Unmarshal(paramsJSON, &callReq); err != nil {
		return nil, fmt.Errorf("failed to parse tool call request: %w", err)
	}

	s.logger.Info("calling tool",
		zap.String("tool", callReq.Name),
		zap.Any("arguments", callReq.Arguments))

	// 调用工具注册表中的工具 / Call tool from registry
	ctx := context.Background()
	result, err := s.toolRegistry.CallTool(ctx, callReq.Name, callReq.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool: %w", err)
	}

	// 转换结果为响应格式 / Convert result to response format
	content := make([]map[string]interface{}, len(result.Content))
	for i, c := range result.Content {
		if textContent, ok := c.(*mcp.TextContent); ok {
			content[i] = map[string]interface{}{
				"type": "text",
				"text": textContent.Text,
			}
		}
	}

	return map[string]interface{}{
		"content": content,
	}, nil
}

// sendJSONError 发送JSON错误响应 / Send JSON error response
func (s *SSETransportServer) sendJSONError(w http.ResponseWriter, id interface{}, code int, message string, data interface{}) {
	mcpErr := types.NewMCPError(code, message, data)
	resp := types.NewMCPErrorResponse(id, mcpErr)

	respJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("failed to marshal error response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // JSON-RPC总是返回200 / JSON-RPC always returns 200
	_, _ = w.Write(respJSON)
}

// Shutdown 关闭SSE服务器 / Shutdown SSE server
func (s *SSETransportServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.logger.Info("shutting down SSE transport server")
	return s.server.Shutdown(ctx)
}

// setCORSHeaders 设置CORS响应头 / Set CORS response headers
func (s *SSETransportServer) setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}

	// 检查是否允许该来源 / Check if origin is allowed
	allowed := false
	for _, allowedOrigin := range s.config.AllowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			allowed = true
			break
		}
	}

	if allowed {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")
	}
}

// getNextConnID 获取下一个连接ID / Get next connection ID
func (s *SSETransportServer) getNextConnID() string {
	count := atomic.AddInt32(&s.connCount, 1)
	return fmt.Sprintf("conn-%d", count)
}
