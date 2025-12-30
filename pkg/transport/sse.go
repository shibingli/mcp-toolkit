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
	id        string
	writer    http.ResponseWriter
	flusher   http.Flusher
	messages  chan *types.SSEMessage
	done      chan struct{}
	createdAt time.Time
	lastPing  time.Time
	mu        sync.RWMutex // 保护lastPing / Protect lastPing
}

// SSETransportServer SSE传输服务器 / SSE transport server
type SSETransportServer struct {
	config       *types.SSEConfig
	logger       *zap.Logger
	server       *http.Server
	mcpServer    *mcp.Server
	toolRegistry *ToolRegistry
	connections  sync.Map     // 存储活跃连接 / Store active connections
	connCount    int32        // 连接计数 / Connection count
	rateLimiter  *RateLimiter // 请求频率限制器 / Request rate limiter
}

// NewSSETransportServer 创建SSE传输服务器 / Create SSE transport server
func NewSSETransportServer(config *types.SSEConfig, logger *zap.Logger) (*SSETransportServer, error) {
	if config == nil {
		config = types.DefaultSSEConfig()
	}

	// 验证配置 / Validate configuration
	if config.Port <= 0 || config.Port > 65535 {
		return nil, fmt.Errorf("invalid port: %d (must be 1-65535)", config.Port)
	}

	if config.HeartbeatInterval < 0 {
		return nil, fmt.Errorf("invalid heartbeat interval: %d (must be >= 0)", config.HeartbeatInterval)
	}
	if config.HeartbeatInterval == 0 {
		config.HeartbeatInterval = 30
	}

	if config.MaxConnections < 0 {
		return nil, fmt.Errorf("invalid max connections: %d (must be >= 0)", config.MaxConnections)
	}
	if config.MaxConnections == 0 {
		config.MaxConnections = 100
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

	// 创建工具注册表 / Create tool registry
	toolRegistry := NewToolRegistry()
	toolRegistry.SetLogger(logger)

	server := &SSETransportServer{
		config:       config,
		logger:       logger,
		mcpServer:    nil, // 将在Start方法中设置 / Will be set in Start method
		toolRegistry: toolRegistry,
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

		// 检查频率限制 / Check rate limit
		if s.rateLimiter != nil {
			clientID := s.getClientID(r)
			if !s.rateLimiter.Allow(clientID) {
				s.logger.Warn("rate limit exceeded for SSE connection",
					zap.String("client_id", clientID),
					zap.String("remote_addr", r.RemoteAddr))
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
		}

		s.handleSSEStream(w, r, ctx)
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

	// 启动连接清理goroutine / Start connection cleanup goroutine
	go s.startConnectionCleanup(ctx)

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

// handleSSEStream 处理SSE流连接 / Handle SSE stream connection
func (s *SSETransportServer) handleSSEStream(w http.ResponseWriter, r *http.Request, ctx context.Context) {
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

	// 创建新连接 / Create new connection
	now := time.Now()
	conn := &SSEConnection{
		id:        s.getNextConnID(),
		writer:    w,
		flusher:   flusher,
		messages:  make(chan *types.SSEMessage, 100), // 缓冲100条消息 / Buffer 100 messages
		done:      make(chan struct{}),
		createdAt: now,
		lastPing:  now,
	}

	// 添加连接到连接池 / Add connection to pool
	if err := s.addConnection(conn); err != nil {
		s.logger.Error("failed to add connection", zap.Error(err))
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer s.removeConnection(conn.id)

	s.logger.Info("new SSE connection established",
		zap.String("conn_id", conn.id),
		zap.String("remote_addr", r.RemoteAddr))

	// 发送连接ID给客户端 / Send connection ID to client
	_, _ = fmt.Fprintf(w, "event: connected\ndata: {\"connectionId\":\"%s\"}\n\n", conn.id)
	flusher.Flush()

	// 发送心跳 / Send heartbeat
	ticker := time.NewTicker(time.Duration(s.config.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("SSE connection closed by server shutdown",
				zap.String("conn_id", conn.id))
			return
		case <-r.Context().Done():
			s.logger.Info("SSE connection closed by client",
				zap.String("conn_id", conn.id))
			return
		case <-conn.done:
			s.logger.Info("SSE connection closed",
				zap.String("conn_id", conn.id))
			return
		case msg := <-conn.messages:
			// 发送消息 / Send message
			_, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", msg.Event, msg.Data)
			if err != nil {
				s.logger.Error("failed to send message",
					zap.String("conn_id", conn.id),
					zap.Error(err))
				return
			}
			flusher.Flush()
		case <-ticker.C:
			// 发送心跳事件 / Send heartbeat event
			_, err := fmt.Fprintf(w, ": heartbeat\n\n")
			if err != nil {
				s.logger.Error("failed to send heartbeat",
					zap.String("conn_id", conn.id),
					zap.Error(err))
				return
			}
			flusher.Flush()
			conn.UpdateLastPing()
		}
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
// 优雅关闭：关闭所有SSE连接，等待正在处理的请求完成
// Graceful shutdown: close all SSE connections, wait for in-flight requests to complete
func (s *SSETransportServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	connCount := atomic.LoadInt32(&s.connCount)
	s.logger.Info("shutting down SSE transport server",
		zap.Int32("active_connections", connCount))

	// 关闭所有活跃连接 / Close all active connections
	closedCount := 0
	s.connections.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*SSEConnection); ok {
			s.logger.Debug("closing SSE connection",
				zap.String("conn_id", conn.id))
			// 安全关闭channel / Safely close channels
			select {
			case <-conn.done:
				// 已经关闭 / Already closed
			default:
				close(conn.done)
			}
			closedCount++
		}
		return true
	})

	if closedCount > 0 {
		s.logger.Info("closed SSE connections",
			zap.Int("count", closedCount))
	}

	// 清空连接映射 / Clear connections map
	s.connections = sync.Map{}
	atomic.StoreInt32(&s.connCount, 0)

	// 优雅关闭HTTP服务器 / Gracefully shutdown HTTP server
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("error during server shutdown", zap.Error(err))
		return err
	}

	s.logger.Info("SSE transport server shutdown complete")
	return nil
}

// getClientID 获取客户端标识 / Get client identifier
// 优先使用 X-Forwarded-For，否则使用 RemoteAddr
// Prefer X-Forwarded-For, otherwise use RemoteAddr
func (s *SSETransportServer) getClientID(r *http.Request) string {
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
		// 包含所有MCP相关的头 / Include all MCP-related headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, MCP-Protocol-Version")
		w.Header().Set("Access-Control-Max-Age", "3600")
	}
}

// validateProtocolVersion 验证协议版本头 / Validate protocol version header
func (s *SSETransportServer) validateProtocolVersion(r *http.Request) error {
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

// getNextConnID 获取下一个连接ID / Get next connection ID
func (s *SSETransportServer) getNextConnID() string {
	count := atomic.AddInt32(&s.connCount, 1)
	return fmt.Sprintf("conn-%d", count)
}

// addConnection 添加新连接 / Add new connection
func (s *SSETransportServer) addConnection(conn *SSEConnection) error {
	// 原子性地检查并增加连接数 / Atomically check and increment connection count
	for {
		currentCount := atomic.LoadInt32(&s.connCount)
		if currentCount >= int32(s.config.MaxConnections) {
			return fmt.Errorf("maximum connections reached: %d", s.config.MaxConnections)
		}

		// 尝试增加计数器 / Try to increment counter
		if atomic.CompareAndSwapInt32(&s.connCount, currentCount, currentCount+1) {
			break
		}
		// CAS失败，重试 / CAS failed, retry
	}

	s.connections.Store(conn.id, conn)
	newCount := atomic.LoadInt32(&s.connCount)
	s.logger.Info("SSE connection added",
		zap.String("conn_id", conn.id),
		zap.Int32("total_connections", newCount))

	return nil
}

// removeConnection 移除连接 / Remove connection
func (s *SSETransportServer) removeConnection(connID string) {
	if conn, ok := s.connections.LoadAndDelete(connID); ok {
		if sseConn, ok := conn.(*SSEConnection); ok {
			close(sseConn.done)
			close(sseConn.messages)
		}
		atomic.AddInt32(&s.connCount, -1)
		s.logger.Info("SSE connection removed",
			zap.String("conn_id", connID),
			zap.Int32("remaining_connections", atomic.LoadInt32(&s.connCount)))
	}
}

// getConnection 获取连接 / Get connection
func (s *SSETransportServer) getConnection(connID string) (*SSEConnection, bool) {
	if conn, ok := s.connections.Load(connID); ok {
		if sseConn, ok := conn.(*SSEConnection); ok {
			return sseConn, true
		}
	}
	return nil, false
}

// broadcastMessage 向所有连接广播消息 / Broadcast message to all connections
func (s *SSETransportServer) broadcastMessage(msg *types.SSEMessage) {
	s.connections.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*SSEConnection); ok {
			select {
			case conn.messages <- msg:
			case <-conn.done:
				// 连接已关闭，移除它 / Connection closed, remove it
				s.removeConnection(conn.id)
			default:
				// 消息队列已满，记录警告 / Message queue full, log warning
				s.logger.Warn("message queue full for connection",
					zap.String("conn_id", conn.id))
			}
		}
		return true
	})
}

// sendToConnection 向指定连接发送消息 / Send message to specific connection
func (s *SSETransportServer) sendToConnection(connID string, msg *types.SSEMessage) error {
	conn, ok := s.getConnection(connID)
	if !ok {
		return fmt.Errorf("connection not found: %s", connID)
	}

	select {
	case conn.messages <- msg:
		return nil
	case <-conn.done:
		s.removeConnection(connID)
		return fmt.Errorf("connection closed: %s", connID)
	default:
		return fmt.Errorf("message queue full for connection: %s", connID)
	}
}

// cleanupStaleConnections 清理过期连接 / Cleanup stale connections
func (s *SSETransportServer) cleanupStaleConnections() {
	timeout := time.Duration(s.config.HeartbeatInterval*3) * time.Second
	now := time.Now()

	var staleConns []string
	s.connections.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*SSEConnection); ok {
			conn.mu.RLock()
			lastPing := conn.lastPing
			conn.mu.RUnlock()

			if now.Sub(lastPing) > timeout {
				staleConns = append(staleConns, conn.id)
			}
		}
		return true
	})

	// 移除过期连接 / Remove stale connections
	for _, connID := range staleConns {
		s.logger.Warn("removing stale connection",
			zap.String("conn_id", connID),
			zap.Duration("timeout", timeout))
		s.removeConnection(connID)
	}
}

// startConnectionCleanup 启动连接清理goroutine / Start connection cleanup goroutine
func (s *SSETransportServer) startConnectionCleanup(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.config.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.cleanupStaleConnections()
		}
	}
}

// UpdateLastPing 更新连接的最后ping时间 / Update connection's last ping time
func (c *SSEConnection) UpdateLastPing() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastPing = time.Now()
}

// GetLastPing 获取连接的最后ping时间 / Get connection's last ping time
func (c *SSEConnection) GetLastPing() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastPing
}
