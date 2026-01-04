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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"go.uber.org/zap"
)

// SSEClient SSE传输的MCP客户端 / MCP client with SSE transport
type SSEClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
	requestID  int64
}

// NewSSEClient 创建SSE客户端 / Create SSE client
func NewSSEClient(host string, port int, path string, logger *zap.Logger) *SSEClient {
	if logger == nil {
		logger = zap.NewNop()
	}

	// SSE 服务器使用 /message 端点接收 POST 请求 / SSE server uses /message endpoint for POST requests
	baseURL := fmt.Sprintf("http://%s:%d%s", host, port, path)
	if !strings.HasSuffix(baseURL, "/message") {
		baseURL = baseURL + "/message"
	}

	return &SSEClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:    logger,
		requestID: 0,
	}
}

// Initialize 初始化连接 / Initialize connection
func (c *SSEClient) Initialize(ctx context.Context, protocolVersion string) (*InitializeResponse, error) {
	c.requestID++
	req := types.MCPRequest{
		JSONRPC: "2.0",
		ID:      c.requestID,
		Method:  "initialize",
		Params: types.MCPInitializeRequest{
			ProtocolVersion: protocolVersion,
			Capabilities:    map[string]interface{}{},
			ClientInfo: map[string]interface{}{
				"name":    "mcp-sse-client",
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
func (c *SSEClient) ListTools(ctx context.Context) (*ListToolsResponse, error) {
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
func (c *SSEClient) CallTool(ctx context.Context, name string, arguments interface{}) (*CallToolResponse, error) {
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
func (c *SSEClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// sendRequest 发送HTTP请求并处理SSE响应 / Send HTTP request and handle SSE response
func (c *SSEClient) sendRequest(ctx context.Context, mcpReq *types.MCPRequest, result interface{}) error {
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
	httpReq.Header.Set("Accept", "application/json, text/event-stream")
	httpReq.Header.Set("MCP-Protocol-Version", "2025-12-26")

	// 发送请求 / Send request
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	// 检查响应类型 / Check response type
	contentType := httpResp.Header.Get("Content-Type")

	// 先尝试读取响应体 / Try to read response body first
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	c.logger.Debug("received response",
		zap.String("content_type", contentType),
		zap.String("body", string(respBody)))

	// 如果是JSON响应，直接解析 / If JSON response, parse directly
	if strings.Contains(contentType, "application/json") || !strings.Contains(contentType, "text/event-stream") {
		var mcpResp types.MCPResponse
		if err = json.Unmarshal(respBody, &mcpResp); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		if mcpResp.Error != nil {
			return fmt.Errorf("MCP error %d: %s (data: %v)", mcpResp.Error.Code, mcpResp.Error.Message, mcpResp.Error.Data)
		}

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

	// 处理SSE流 / Handle SSE stream
	scanner := bufio.NewScanner(bytes.NewReader(respBody))
	var dataLines []string

	for scanner.Scan() {
		line := scanner.Text()

		// 跳过空行和注释 / Skip empty lines and comments
		if line == "" {
			// 空行表示事件结束 / Empty line indicates end of event
			if len(dataLines) > 0 {
				// 解析数据 / Parse data
				data := strings.Join(dataLines, "\n")
				dataLines = nil

				// 解析MCP响应 / Parse MCP response
				var mcpResp types.MCPResponse
				if err = json.UnmarshalFromString(data, &mcpResp); err != nil {
					c.logger.Warn("failed to unmarshal SSE data", zap.Error(err), zap.String("data", data))
					continue
				}

				c.logger.Debug("received SSE response", zap.String("data", data))

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
			continue
		}

		if strings.HasPrefix(line, ":") {
			// 注释行，跳过 / Comment line, skip
			continue
		}

		if strings.HasPrefix(line, "data:") {
			// 数据行 / Data line
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			dataLines = append(dataLines, data)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read SSE stream: %w", err)
	}

	return fmt.Errorf("no response received from SSE stream")
}
