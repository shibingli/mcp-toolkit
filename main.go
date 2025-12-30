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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"mcp-toolkit/internal/services/sandbox"
	"mcp-toolkit/pkg/transport"
	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// initLogger 初始化日志记录器 / Initialize logger
func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return logger, nil
}

func main() {
	// 解析命令行参数 / Parse command line arguments
	sandboxDir := flag.String("sandbox", types.GetDefaultSandboxDir(), "沙箱目录路径 / Sandbox directory path")
	transportType := flag.String("transport", "stdio", "传输类型: stdio, http, sse / Transport type: stdio, http, sse")
	httpHost := flag.String("http-host", "127.0.0.1", "HTTP监听地址 / HTTP listen address")
	httpPort := flag.Int("http-port", 8080, "HTTP监听端口 / HTTP listen port")
	sseHost := flag.String("sse-host", "127.0.0.1", "SSE监听地址 / SSE listen address")
	ssePort := flag.Int("sse-port", 8081, "SSE监听端口 / SSE listen port")
	flag.Parse()

	// 初始化日志记录器 / Initialize logger
	logger, err := initLogger()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()

	logger.Info("starting MCP server",
		zap.String("name", types.ServerName),
		zap.String("version", types.ServerVersion),
		zap.String("sandbox_dir", *sandboxDir),
		zap.String("transport", *transportType))

	// 预热JSON结构体(仅对sonic有效) / Pretouch JSON structures (only effective for sonic)
	logger.Info("preheating JSON structures",
		zap.String("json_library", json.Name()))
	failedCount, pretouchErr := json.PretouchAll()
	if pretouchErr != nil {
		logger.Warn("some types failed to preheat",
			zap.Int("failed_count", failedCount),
			zap.Error(pretouchErr))
	} else {
		logger.Info("JSON structures preheated successfully")
	}

	// 获取沙箱目录的绝对路径 / Get absolute path of sandbox directory
	absSandboxDir, err := filepath.Abs(*sandboxDir)
	if err != nil {
		logger.Fatal("failed to get absolute path of sandbox directory", zap.Error(err))
	}

	// 创建沙箱服务 / Create sandbox service
	sandboxService, err := sandbox.NewService(absSandboxDir, logger)
	if err != nil {
		logger.Fatal("failed to create sandbox service", zap.Error(err))
	}

	logger.Info("sandbox service initialized", zap.String("sandbox_dir", absSandboxDir))

	// 创建MCP服务器 / Create MCP server
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    types.ServerName,
			Version: types.ServerVersion,
		},
		&mcp.ServerOptions{
			Capabilities: &mcp.ServerCapabilities{
				Tools: &mcp.ToolCapabilities{},
			},
		},
	)

	// 注册沙箱工具 / Register sandbox tools
	sandboxService.RegisterTools(mcpServer)

	logger.Info("MCP tools registered successfully")

	// 创建上下文和信号处理 / Create context and signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号 / Listen for system signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 根据传输类型启动服务器 / Start server based on transport type
	errChan := make(chan error, 1)
	go func() {
		switch types.TransportType(*transportType) {
		case types.TransportHTTP:
			// 启动HTTP传输服务器 / Start HTTP transport server
			httpConfig := &types.HTTPConfig{
				Host:           *httpHost,
				Port:           *httpPort,
				Path:           "/mcp",
				EnableCORS:     true,
				AllowedOrigins: []string{"*"},
				ReadTimeout:    30,
				WriteTimeout:   30,
				MaxHeaderBytes: 1 << 20,
			}
			httpServer, err := transport.NewHTTPTransportServer(httpConfig, logger)
			if err != nil {
				errChan <- fmt.Errorf("failed to create HTTP transport server: %w", err)
				return
			}
			// 注册工具到工具注册表 / Register tools to tool registry
			sandboxService.RegisterToolsToRegistry(httpServer.GetToolRegistry())
			errChan <- httpServer.Start(ctx, mcpServer)

		case types.TransportSSE:
			// 启动SSE传输服务器 / Start SSE transport server
			sseConfig := &types.SSEConfig{
				Host:              *sseHost,
				Port:              *ssePort,
				Path:              "/sse",
				EnableCORS:        true,
				AllowedOrigins:    []string{"*"},
				HeartbeatInterval: 30,
				MaxConnections:    100,
			}
			sseServer, err := transport.NewSSETransportServer(sseConfig, logger)
			if err != nil {
				errChan <- fmt.Errorf("failed to create SSE transport server: %w", err)
				return
			}
			// 注册工具到工具注册表 / Register tools to tool registry
			sandboxService.RegisterToolsToRegistry(sseServer.GetToolRegistry())
			errChan <- sseServer.Start(ctx, mcpServer)

		case types.TransportStdio:
			fallthrough
		default:
			// 启动stdio传输服务器(默认) / Start stdio transport server (default)
			logger.Info("using stdio transport")
			errChan <- mcpServer.Run(ctx, &mcp.StdioTransport{})
		}
	}()

	// 等待信号或错误 / Wait for signal or error
	select {
	case sig := <-sigChan:
		logger.Info("received signal, shutting down", zap.String("signal", sig.String()))
		cancel()
	case err = <-errChan:
		if err != nil {
			logger.Fatal("server error", zap.Error(err))
		}
	}

	logger.Info("MCP server stopped")
}
