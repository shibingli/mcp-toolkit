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

package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"go.uber.org/zap"
)

// Client MCP客户端接口 / MCP client interface
type Client interface {
	// Initialize 初始化连接 / Initialize connection
	Initialize(ctx context.Context, protocolVersion string) (*InitializeResponse, error)

	// ListTools 列出所有可用工具 / List all available tools
	ListTools(ctx context.Context) (*ListToolsResponse, error)

	// CallTool 调用工具 / Call tool
	CallTool(ctx context.Context, name string, arguments interface{}) (*CallToolResponse, error)

	// Close 关闭客户端 / Close client
	Close() error
}

// HTTPClient HTTP传输的MCP客户端 / MCP client with HTTP transport
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
	requestID  int64
}

// NewHTTPClient 创建HTTP客户端 / Create HTTP client
func NewHTTPClient(host string, port int, path string, logger *zap.Logger) *HTTPClient {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &HTTPClient{
		baseURL: fmt.Sprintf("http://%s:%d%s", host, port, path),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:    logger,
		requestID: 0,
	}
}

// Initialize 初始化连接 / Initialize connection
func (c *HTTPClient) Initialize(ctx context.Context, protocolVersion string) (*InitializeResponse, error) {
	c.requestID++
	req := types.MCPRequest{
		JSONRPC: "2.0",
		ID:      c.requestID,
		Method:  "initialize",
		Params: types.MCPInitializeRequest{
			ProtocolVersion: protocolVersion,
			Capabilities:    map[string]interface{}{},
			ClientInfo: map[string]interface{}{
				"name":    "mcp-test-client",
				"version": "1.0.0",
			},
		},
	}

	var resp InitializeResponse
	if err := c.sendRequest(ctx, &req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListTools 列出所有可用工具 / List all available tools
func (c *HTTPClient) ListTools(ctx context.Context) (*ListToolsResponse, error) {
	c.requestID++
	req := types.MCPRequest{
		JSONRPC: "2.0",
		ID:      c.requestID,
		Method:  "tools/list",
		Params:  types.MCPToolsListRequest{},
	}

	var resp ListToolsResponse
	if err := c.sendRequest(ctx, &req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CallTool 调用工具 / Call tool
func (c *HTTPClient) CallTool(ctx context.Context, name string, arguments interface{}) (*CallToolResponse, error) {
	c.requestID++
	req := types.MCPRequest{
		JSONRPC: "2.0",
		ID:      c.requestID,
		Method:  "tools/call",
		Params: types.MCPToolsCallRequest{
			Name:      name,
			Arguments: arguments,
		},
	}

	var resp CallToolResponse
	if err := c.sendRequest(ctx, &req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Close 关闭客户端 / Close client
func (c *HTTPClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// sendRequest 发送HTTP请求 / Send HTTP request
func (c *HTTPClient) sendRequest(ctx context.Context, mcpReq *types.MCPRequest, result interface{}) error {
	// 序列化请求 / Serialize request
	reqBody, err := json.Marshal(mcpReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	c.logger.Debug("sending request",
		zap.String("method", mcpReq.Method),
		zap.Any("id", mcpReq.ID),
		zap.String("body", string(reqBody)))

	// 创建HTTP请求 / Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求 / Send request
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	// 读取响应 / Read response
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	c.logger.Debug("received response",
		zap.Int("status", httpResp.StatusCode),
		zap.String("body", string(respBody)))

	// 解析MCP响应 / Parse MCP response
	var mcpResp types.MCPResponse
	if err = json.Unmarshal(respBody, &mcpResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 检查错误 / Check error
	if mcpResp.Error != nil {
		return fmt.Errorf("MCP error %d: %s (data: %v)", mcpResp.Error.Code, mcpResp.Error.Message, mcpResp.Error.Data)
	}

	// 解析结果 / Parse result
	if result != nil && mcpResp.Result != nil {
		resultJSON, err := json.Marshal(mcpResp.Result)
		if err != nil {
			return fmt.Errorf("failed to marshal result: %w", err)
		}

		if err = json.Unmarshal(resultJSON, result); err != nil {
			return fmt.Errorf("failed to unmarshal result: %w", err)
		}
	}

	return nil
}
