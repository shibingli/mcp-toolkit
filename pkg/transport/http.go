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
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mcp-toolkit/pkg/utils/recovery"
	"net/http"
	"sync"
	"time"

	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// HTTPSession HTTP会话 / HTTP session
// 使用互斥锁保护并发访问 / Use mutex to protect concurrent access
type HTTPSession struct {
	ID        string
	CreatedAt time.Time
	lastUsed  time.Time
	mu        sync.RWMutex // 保护 lastUsed 字段 / Protect lastUsed field
}

// GetLastUsed 获取最后使用时间 / Get last used time
func (s *HTTPSession) GetLastUsed() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastUsed
}

// UpdateLastUsed 更新最后使用时间 / Update last used time
func (s *HTTPSession) UpdateLastUsed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastUsed = time.Now()
}

// HTTPTransportServer HTTP传输服务器 / HTTP transport server
type HTTPTransportServer struct {
	config       *types.HTTPConfig
	logger       *zap.Logger
	server       *http.Server
	mcpServer    *mcp.Server
	toolRegistry *ToolRegistry
	sessions     sync.Map // map[string]*HTTPSession
	sessionMu    sync.RWMutex
	rateLimiter  *RateLimiter // 请求频率限制器 / Request rate limiter
}

// GetToolRegistry 获取工具注册表 / Get tool registry
func (s *HTTPTransportServer) GetToolRegistry() *ToolRegistry {
	return s.toolRegistry
}

// generateSessionID 生成会话ID / Generate session ID
func (s *HTTPTransportServer) generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// createSession 创建新会话 / Create new session
func (s *HTTPTransportServer) createSession() (*HTTPSession, error) {
	sessionID, err := s.generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	now := time.Now()
	session := &HTTPSession{
		ID:        sessionID,
		CreatedAt: now,
		lastUsed:  now,
	}

	s.sessions.Store(sessionID, session)
	s.logger.Info("created new session", zap.String("session_id", sessionID))
	return session, nil
}

// getSession 获取会话 / Get session
func (s *HTTPTransportServer) getSession(sessionID string) (*HTTPSession, error) {
	if sessionID == "" {
		return nil, errors.New("session ID is empty")
	}

	value, ok := s.sessions.Load(sessionID)
	if !ok {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	session := value.(*HTTPSession)
	session.UpdateLastUsed() // 使用线程安全的方法更新 / Use thread-safe method to update
	return session, nil
}

// deleteSession 删除会话 / Delete session
func (s *HTTPTransportServer) deleteSession(sessionID string) {
	s.sessions.Delete(sessionID)
	s.logger.Info("deleted session", zap.String("session_id", sessionID))
}

// cleanupExpiredSessions 清理过期会话 / Cleanup expired sessions
func (s *HTTPTransportServer) cleanupExpiredSessions() {
	if s.config.SessionTimeout <= 0 {
		return
	}

	timeout := time.Duration(s.config.SessionTimeout) * time.Second
	now := time.Now()

	// 收集需要删除的会话ID / Collect session IDs to delete
	var toDelete []string
	s.sessions.Range(func(key, value interface{}) bool {
		session := value.(*HTTPSession)
		if now.Sub(session.GetLastUsed()) > timeout {
			toDelete = append(toDelete, session.ID)
		}
		return true
	})

	// 批量删除过期会话 / Batch delete expired sessions
	for _, sessionID := range toDelete {
		s.deleteSession(sessionID)
	}

	if len(toDelete) > 0 {
		s.logger.Info("cleaned up expired sessions",
			zap.Int("count", len(toDelete)))
	}
}

// NewHTTPTransportServer 创建HTTP传输服务器 / Create HTTP transport server
func NewHTTPTransportServer(config *types.HTTPConfig, logger *zap.Logger) (*HTTPTransportServer, error) {
	if config == nil {
		config = types.DefaultHTTPConfig()
	}

	// 验证配置 / Validate configuration
	if config.Port <= 0 || config.Port > 65535 {
		return nil, fmt.Errorf("invalid port: %d (must be 1-65535)", config.Port)
	}

	if config.ReadTimeout < 0 {
		return nil, fmt.Errorf("invalid read timeout: %d (must be >= 0)", config.ReadTimeout)
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 30
	}

	if config.WriteTimeout < 0 {
		return nil, fmt.Errorf("invalid write timeout: %d (must be >= 0)", config.WriteTimeout)
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 30
	}

	if config.MaxHeaderBytes < 0 {
		return nil, fmt.Errorf("invalid max header bytes: %d (must be >= 0)", config.MaxHeaderBytes)
	}
	if config.MaxHeaderBytes == 0 {
		config.MaxHeaderBytes = 1 << 20 // 1MB
	}

	// 验证会话配置 / Validate session configuration
	if config.SessionTimeout < 0 {
		return nil, fmt.Errorf("invalid session timeout: %d (must be >= 0)", config.SessionTimeout)
	}
	if config.SessionTimeout == 0 && config.EnableSessionManagement {
		config.SessionTimeout = 1800 // 默认30分钟 / Default 30 minutes
	}

	// 验证SSE配置 / Validate SSE configuration
	if config.SSEHeartbeatInterval < 0 {
		return nil, fmt.Errorf("invalid SSE heartbeat interval: %d (must be >= 0)", config.SSEHeartbeatInterval)
	}
	if config.SSEHeartbeatInterval == 0 && config.EnableSSE {
		config.SSEHeartbeatInterval = 30
	}

	// 验证频率限制配置 / Validate rate limit configuration
	if config.RateLimitRequests < 0 {
		return nil, fmt.Errorf("invalid rate limit requests: %d (must be >= 0)", config.RateLimitRequests)
	}
	if config.RateLimitWindow < 0 {
		return nil, fmt.Errorf("invalid rate limit window: %d (must be >= 0)", config.RateLimitWindow)
	}
	if config.EnableRateLimit && config.RateLimitRequests == 0 {
		config.RateLimitRequests = 100 // 默认每分钟100个请求 / Default 100 requests per minute
	}
	if config.EnableRateLimit && config.RateLimitWindow == 0 {
		config.RateLimitWindow = 60 // 默认60秒窗口 / Default 60 second window
	}

	registry := NewToolRegistry()
	registry.SetLogger(logger)

	server := &HTTPTransportServer{
		config:       config,
		logger:       logger,
		mcpServer:    nil, // 将在Start方法中设置 / Will be set in Start method
		toolRegistry: registry,
	}

	// 如果启用频率限制，创建限制器 / If rate limiting enabled, create limiter
	if config.EnableRateLimit {
		server.rateLimiter = NewRateLimiter(
			config.RateLimitRequests,
			time.Duration(config.RateLimitWindow)*time.Second,
			logger,
		)
		logger.Info("rate limiter enabled",
			zap.Int("max_requests", config.RateLimitRequests),
			zap.Int("window_seconds", config.RateLimitWindow))
	}

	return server, nil
}

// Start 启动HTTP服务器 / Start HTTP server
func (s *HTTPTransportServer) Start(ctx context.Context, mcpServer *mcp.Server) error {
	s.mcpServer = mcpServer
	mux := http.NewServeMux()

	// 注册MCP处理器 (支持POST和GET) / Register MCP handler (supports POST and GET)
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

		// 检查频率限制 / Check rate limit
		if s.rateLimiter != nil {
			clientID := s.getClientID(r)
			if !s.rateLimiter.Allow(clientID) {
				s.logger.Warn("rate limit exceeded",
					zap.String("client_id", clientID),
					zap.String("remote_addr", r.RemoteAddr))
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
		}

		// 验证协议版本头 / Validate protocol version header
		if err := s.validateProtocolVersion(r); err != nil {
			s.logger.Error("invalid protocol version", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 根据请求方法路由 / Route based on request method
		switch r.Method {
		case http.MethodPost:
			// POST: 发送消息到服务器 / POST: Send message to server
			s.handleMCPRequest(w, r)
		case http.MethodGet:
			// GET: 打开SSE流监听服务器消息 / GET: Open SSE stream to listen for server messages
			if s.config.EnableSSE {
				s.handleSSEStream(w, r)
			} else {
				http.Error(w, "SSE not enabled", http.StatusMethodNotAllowed)
			}
		case http.MethodDelete:
			// DELETE: 终止会话 / DELETE: Terminate session
			if s.config.EnableSessionManagement {
				s.handleSessionDelete(w, r)
			} else {
				http.Error(w, "Session management not enabled", http.StatusMethodNotAllowed)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
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
		zap.String("path", s.config.Path),
		zap.Bool("sse_enabled", s.config.EnableSSE),
		zap.Bool("session_management_enabled", s.config.EnableSessionManagement))

	// 启动会话清理goroutine / Start session cleanup goroutine
	if s.config.EnableSessionManagement && s.config.SessionTimeout > 0 {
		go func() {
			ticker := time.NewTicker(time.Duration(s.config.SessionTimeout/2) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					s.cleanupExpiredSessions()
				}
			}
		}()
	}

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
// 优雅关闭：等待正在处理的请求完成，清理会话和资源
// Graceful shutdown: wait for in-flight requests to complete, cleanup sessions and resources
func (s *HTTPTransportServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.logger.Info("shutting down HTTP transport server")

	// 清理所有会话 / Cleanup all sessions
	sessionCount := 0
	s.sessions.Range(func(key, value interface{}) bool {
		sessionCount++
		return true
	})

	if sessionCount > 0 {
		s.logger.Info("cleaning up sessions",
			zap.Int("count", sessionCount))
		s.sessions = sync.Map{}
	}

	// 优雅关闭HTTP服务器 / Gracefully shutdown HTTP server
	// server.Shutdown 会等待所有活跃连接完成
	// server.Shutdown waits for all active connections to complete
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("error during server shutdown", zap.Error(err))
		return err
	}

	s.logger.Info("HTTP transport server shutdown complete")
	return nil
}

// validateProtocolVersion 验证协议版本头 / Validate protocol version header
func (s *HTTPTransportServer) validateProtocolVersion(r *http.Request) error {
	version := r.Header.Get("MCP-Protocol-Version")
	if version == "" {
		// 向后兼容：如果没有版本头，假设客户端支持所有版本
		// Backward compatibility: if no version header, assume client supports all versions
		return nil
	}

	// 验证版本是否在支持列表中 / Validate version is in supported list
	for _, v := range types.SupportedProtocolVersions {
		if version == v {
			return nil
		}
	}

	return fmt.Errorf("unsupported protocol version: %s (supported: %v)", version, types.SupportedProtocolVersions)
}

// handleSessionDelete 处理会话删除请求 / Handle session delete request
func (s *HTTPTransportServer) handleSessionDelete(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("Mcp-Session-Id")
	if sessionID == "" {
		http.Error(w, "Missing Mcp-Session-Id header", http.StatusBadRequest)
		return
	}

	s.deleteSession(sessionID)
	w.WriteHeader(http.StatusOK)
}

// handleSSEStream 处理SSE流请求 / Handle SSE stream request
func (s *HTTPTransportServer) handleSSEStream(w http.ResponseWriter, r *http.Request) {
	// 验证Accept头 / Validate Accept header
	accept := r.Header.Get("Accept")
	if accept != "text/event-stream" && accept != "*/*" {
		http.Error(w, "Accept header must include text/event-stream", http.StatusBadRequest)
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

	s.logger.Info("new SSE stream connection",
		zap.String("remote_addr", r.RemoteAddr))

	// 发送心跳 / Send heartbeat
	ticker := time.NewTicker(time.Duration(s.config.SSEHeartbeatInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			s.logger.Info("SSE stream connection closed",
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
}

// handleMCPRequest 处理MCP JSON-RPC请求 / Handle MCP JSON-RPC request
func (s *HTTPTransportServer) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	// 验证Accept头 / Validate Accept header
	accept := r.Header.Get("Accept")
	supportsSSE := false
	supportsJSON := false

	if accept == "*/*" || accept == "" {
		supportsSSE = s.config.EnableSSE
		supportsJSON = true
	} else {
		// 检查Accept头中是否包含支持的类型
		if contains(accept, "text/event-stream") {
			supportsSSE = s.config.EnableSSE
		}
		if contains(accept, "application/json") {
			supportsJSON = true
		}
	}

	if !supportsSSE && !supportsJSON {
		http.Error(w, "Accept header must include application/json or text/event-stream", http.StatusBadRequest)
		return
	}

	// 添加panic recovery保护 / Add panic recovery protection
	recoveryHandler := recovery.NewRecoveryHandler(s.logger)
	err := recoveryHandler.Recover(func() error {
		return s.handleMCPRequestInternal(w, r, supportsSSE, supportsJSON)
	})

	if err != nil {
		// 如果是panic导致的错误，发送内部错误响应 / If error from panic, send internal error response
		s.logger.Error("panic in MCP request handler", zap.Error(err))
		s.sendErrorResponse(w, nil, types.MCPErrorCodeInternalError, "Internal server error", err.Error())
	}
}

// contains 检查字符串是否包含子串 / Check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// handleMCPRequestInternal 内部MCP请求处理 / Internal MCP request handling
func (s *HTTPTransportServer) handleMCPRequestInternal(w http.ResponseWriter, r *http.Request, supportsSSE, supportsJSON bool) error {
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

	// 检查是否是通知（无ID但有方法）或响应（无方法）
	// Check if it's a notification (no ID but has method) or response (no method)
	isNotification := mcpReq.ID == nil && mcpReq.Method != ""
	isResponse := mcpReq.Method == ""
	isNotificationOrResponse := isNotification || isResponse

	// 验证会话 / Validate session
	sessionID := r.Header.Get("Mcp-Session-Id")
	var session *HTTPSession

	if s.config.EnableSessionManagement {
		if mcpReq.Method == "initialize" {
			// 初始化请求：创建新会话 / Initialize request: create new session
			session, err = s.createSession()
			if err != nil {
				s.logger.Error("failed to create session", zap.Error(err))
				s.sendErrorResponse(w, mcpReq.ID, types.MCPErrorCodeInternalError, "Failed to create session", err.Error())
				return nil
			}
		} else {
			// 其他请求：验证会话 / Other requests: validate session
			if sessionID == "" {
				http.Error(w, "Missing Mcp-Session-Id header", http.StatusBadRequest)
				return nil
			}
			session, err = s.getSession(sessionID)
			if err != nil {
				http.Error(w, "Session not found", http.StatusNotFound)
				return nil
			}
		}
	}

	s.logger.Info("received MCP request",
		zap.String("method", mcpReq.Method),
		zap.Any("id", mcpReq.ID),
		zap.String("session_id", sessionID),
		zap.String("remote_addr", r.RemoteAddr))

	// 如果是通知或响应，返回202 Accepted / If it's a notification or response, return 202 Accepted
	if isNotificationOrResponse {
		w.WriteHeader(http.StatusAccepted)
		return nil
	}

	// 根据方法路由请求 / Route request based on method
	var result interface{}
	switch mcpReq.Method {
	case "initialize":
		result = s.handleInitialize(mcpReq, session)
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

	// 根据客户端支持的类型发送响应 / Send response based on client support
	if supportsSSE && s.config.EnableSSE {
		// 发送SSE流响应 / Send SSE stream response
		s.sendSSEResponse(w, r, mcpReq.ID, result, session)
	} else if supportsJSON {
		// 发送JSON响应 / Send JSON response
		s.sendJSONResponse(w, mcpReq.ID, result, session)
	} else {
		http.Error(w, "No supported response format", http.StatusNotAcceptable)
	}

	return nil
}

// handleInitialize 处理初始化请求 / Handle initialize request
func (s *HTTPTransportServer) handleInitialize(req types.MCPRequest, session *HTTPSession) interface{} {
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

	// 调用工具注册表中的工具，使用带超时的context / Call tool from registry with timeout context
	// 注意：这里使用 Background 是因为工具调用不应该被HTTP请求取消中断
	// Note: Using Background here because tool calls should not be interrupted by HTTP request cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
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

// sendJSONResponse 发送JSON响应 / Send JSON response
func (s *HTTPTransportServer) sendJSONResponse(w http.ResponseWriter, id interface{}, result interface{}, session *HTTPSession) {
	resp := types.NewMCPResponse(id, result)
	respJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("failed to marshal response", zap.Error(err))
		s.sendErrorResponse(w, id, types.MCPErrorCodeInternalError, types.MCPErrorMsgInternalError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if session != nil {
		w.Header().Set("Mcp-Session-Id", session.ID)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJSON)
}

// sendSSEResponse 发送SSE流响应 / Send SSE stream response
func (s *HTTPTransportServer) sendSSEResponse(w http.ResponseWriter, r *http.Request, id interface{}, result interface{}, session *HTTPSession) {
	// 设置SSE响应头 / Set SSE response headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	if session != nil {
		w.Header().Set("Mcp-Session-Id", session.ID)
	}

	// 创建flusher / Create flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 发送响应事件 / Send response event
	resp := types.NewMCPResponse(id, result)
	respJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("failed to marshal response", zap.Error(err))
		return
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", string(respJSON))
	if err != nil {
		s.logger.Error("failed to write SSE response", zap.Error(err))
		return
	}
	flusher.Flush()

	s.logger.Info("sent SSE response", zap.Any("id", id))
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

// getClientID 获取客户端标识 / Get client identifier
// 优先使用 X-Forwarded-For，否则使用 RemoteAddr
// Prefer X-Forwarded-For, otherwise use RemoteAddr
func (s *HTTPTransportServer) getClientID(r *http.Request) string {
	// 检查 X-Forwarded-For 头 / Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// 检查 X-Real-IP 头 / Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// 使用 RemoteAddr / Use RemoteAddr
	return r.RemoteAddr
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
		// 支持所有必要的HTTP方法 / Support all necessary HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		// 包含所有MCP相关的头 / Include all MCP-related headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Mcp-Session-Id, MCP-Protocol-Version")
		// 允许暴露会话ID头 / Allow exposing session ID header
		w.Header().Set("Access-Control-Expose-Headers", "Mcp-Session-Id")
		w.Header().Set("Access-Control-Max-Age", "3600")
	}
}
