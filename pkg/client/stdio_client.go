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
	"context"
	"fmt"
	"io"
	"os/exec"

	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"go.uber.org/zap"
)

// StdioClient Stdio传输的MCP客户端 / MCP client with Stdio transport
type StdioClient struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	scanner   *bufio.Scanner
	logger    *zap.Logger
	requestID int64
}

// NewStdioClient 创建Stdio客户端 / Create Stdio client
func NewStdioClient(command string, args []string, logger *zap.Logger) (*StdioClient, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	// 创建命令 / Create command
	cmd := exec.Command(command, args...)

	// 获取stdin和stdout / Get stdin and stdout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// 启动进程 / Start process
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process: %w", err)
	}

	return &StdioClient{
		cmd:       cmd,
		stdin:     stdin,
		stdout:    stdout,
		scanner:   bufio.NewScanner(stdout),
		logger:    logger,
		requestID: 0,
	}, nil
}

// Initialize 初始化连接 / Initialize connection
func (c *StdioClient) Initialize(ctx context.Context, protocolVersion string) (*InitializeResponse, error) {
	c.requestID++
	req := types.MCPRequest{
		JSONRPC: "2.0",
		ID:      c.requestID,
		Method:  "initialize",
		Params: types.MCPInitializeRequest{
			ProtocolVersion: protocolVersion,
			Capabilities:    map[string]interface{}{},
			ClientInfo: map[string]interface{}{
				"name":    "mcp-stdio-client",
				"version": "1.0.0",
			},
		},
	}

	var resp InitializeResponse
	if err := c.sendRequest(&req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListTools 列出所有可用工具 / List all available tools
func (c *StdioClient) ListTools(ctx context.Context) (*ListToolsResponse, error) {
	c.requestID++
	req := types.MCPRequest{
		JSONRPC: "2.0",
		ID:      c.requestID,
		Method:  "tools/list",
		Params:  types.MCPToolsListRequest{},
	}

	var resp ListToolsResponse
	if err := c.sendRequest(&req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CallTool 调用工具 / Call tool
func (c *StdioClient) CallTool(ctx context.Context, name string, arguments interface{}) (*CallToolResponse, error) {
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
	if err := c.sendRequest(&req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Close 关闭客户端 / Close client
func (c *StdioClient) Close() error {
	if c.stdin != nil {
		_ = c.stdin.Close()
	}
	if c.stdout != nil {
		_ = c.stdout.Close()
	}
	if c.cmd != nil && c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
		_ = c.cmd.Wait()
	}
	return nil
}

// sendRequest 发送请求并接收响应 / Send request and receive response
func (c *StdioClient) sendRequest(mcpReq *types.MCPRequest, result interface{}) error {
	// 序列化请求 / Serialize request
	reqBody, err := json.Marshal(mcpReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	c.logger.Debug("sending request",
		zap.String("method", mcpReq.Method),
		zap.Any("id", mcpReq.ID),
		zap.String("body", string(reqBody)))

	// 发送请求 / Send request
	if _, err := c.stdin.Write(reqBody); err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}
	if _, err := c.stdin.Write([]byte("\n")); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	// 读取响应 / Read response
	if !c.scanner.Scan() {
		if err := c.scanner.Err(); err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
		return fmt.Errorf("no response received")
	}

	respBody := c.scanner.Bytes()
	c.logger.Debug("received response", zap.String("body", string(respBody)))

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
