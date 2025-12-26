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
	"time"

	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// HTTPTransportServer HTTP传输服务器 / HTTP transport server
type HTTPTransportServer struct {
	config       *types.HTTPConfig
	logger       *zap.Logger
	server       *http.Server
	mcpServer    *mcp.Server
	toolRegistry *ToolRegistry
}

// GetToolRegistry 获取工具注册表 / Get tool registry
func (s *HTTPTransportServer) GetToolRegistry() *ToolRegistry {
	return s.toolRegistry
}

// NewHTTPTransportServer 创建HTTP传输服务器 / Create HTTP transport server
func NewHTTPTransportServer(config *types.HTTPConfig, logger *zap.Logger) (*HTTPTransportServer, error) {
	if config == nil {
		config = types.DefaultHTTPConfig()
	}

	// 验证配置 / Validate configuration
	if config.Port <= 0 || config.Port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", config.Port)
	}

	if config.ReadTimeout <= 0 {
		config.ReadTimeout = 30
	}

	if config.WriteTimeout <= 0 {
		config.WriteTimeout = 30
	}

	if config.MaxHeaderBytes <= 0 {
		config.MaxHeaderBytes = 1 << 20 // 1MB
	}

	registry := NewToolRegistry()
	registry.SetLogger(logger)

	return &HTTPTransportServer{
		config:       config,
		logger:       logger,
		mcpServer:    nil, // 将在Start方法中设置 / Will be set in Start method
		toolRegistry: registry,
	}, nil
}

// Start 启动HTTP服务器 / Start HTTP server
func (s *HTTPTransportServer) Start(ctx context.Context, mcpServer *mcp.Server) error {
	s.mcpServer = mcpServer
	mux := http.NewServeMux()

	// 注册MCP处理器 / Register MCP handler
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

		// 仅允许POST请求 / Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 处理MCP JSON-RPC请求 / Handle MCP JSON-RPC request
		s.handleMCPRequest(w, r)
	})

	// 健康检查端点 / Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.server = &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    time.Duration(s.config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.config.WriteTimeout) * time.Second,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}

	s.logger.Info("starting HTTP transport server",
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

// Shutdown 关闭HTTP服务器 / Shutdown HTTP server
func (s *HTTPTransportServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.logger.Info("shutting down HTTP transport server")
	return s.server.Shutdown(ctx)
}

// handleMCPRequest 处理MCP JSON-RPC请求 / Handle MCP JSON-RPC request
func (s *HTTPTransportServer) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	// 添加panic recovery保护 / Add panic recovery protection
	recoveryHandler := recovery.NewRecoveryHandler(s.logger)
	err := recoveryHandler.Recover(func() error {
		return s.handleMCPRequestInternal(w, r)
	})

	if err != nil {
		// 如果是panic导致的错误，发送内部错误响应 / If error from panic, send internal error response
		s.logger.Error("panic in MCP request handler", zap.Error(err))
		s.sendErrorResponse(w, nil, types.MCPErrorCodeInternalError, "Internal server error", err.Error())
	}
}

// handleMCPRequestInternal 内部MCP请求处理 / Internal MCP request handling
func (s *HTTPTransportServer) handleMCPRequestInternal(w http.ResponseWriter, r *http.Request) error {
	// 读取请求体 / Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("failed to read request body", zap.Error(err))
		s.sendErrorResponse(w, nil, types.MCPErrorCodeParseError, types.MCPErrorMsgParseError, err.Error())
		return nil
	}
	defer func() { _ = r.Body.Close() }()

	// 解析MCP请求 / Parse MCP request
	var mcpReq types.MCPRequest
	if err = json.Unmarshal(body, &mcpReq); err != nil {
		s.logger.Error("failed to parse MCP request", zap.Error(err))
		s.sendErrorResponse(w, nil, types.MCPErrorCodeParseError, types.MCPErrorMsgParseError, err.Error())
		return nil
	}

	s.logger.Info("received MCP request",
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
			s.sendErrorResponse(w, mcpReq.ID, types.MCPErrorCodeInternalError, err.Error(), nil)
			return nil
		}
	default:
		s.sendErrorResponse(w, mcpReq.ID, types.MCPErrorCodeMethodNotFound, types.MCPErrorMsgMethodNotFound, mcpReq.Method)
		return nil
	}

	// 发送成功响应 / Send success response
	resp := types.NewMCPResponse(mcpReq.ID, result)
	respJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("failed to marshal response", zap.Error(err))
		s.sendErrorResponse(w, mcpReq.ID, types.MCPErrorCodeInternalError, types.MCPErrorMsgInternalError, err.Error())
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJSON)
	return nil
}

// handleInitialize 处理初始化请求 / Handle initialize request
func (s *HTTPTransportServer) handleInitialize(req types.MCPRequest) interface{} {
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
func (s *HTTPTransportServer) handleToolsList(req types.MCPRequest) interface{} {
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
func (s *HTTPTransportServer) handleToolsCall(req types.MCPRequest) (interface{}, error) {
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

// sendErrorResponse 发送错误响应 / Send error response
func (s *HTTPTransportServer) sendErrorResponse(w http.ResponseWriter, id interface{}, code int, message string, data interface{}) {
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

// setCORSHeaders 设置CORS响应头 / Set CORS response headers
func (s *HTTPTransportServer) setCORSHeaders(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")
	}
}
