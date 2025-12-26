// Package transport 实现了 MCP Toolkit 的传输层
//
// 本包提供了多种传输方式的实现，用于在客户端和服务器之间传输 MCP 协议消息。
//
// # 支持的传输方式
//
// 1. Stdio 传输（默认）
//   - 使用标准输入/输出进行通信
//   - 适用于本地进程间通信
//   - 由 MCP SDK 直接支持
//
// 2. HTTP 传输
//   - 基于 HTTP POST 请求
//   - 支持 CORS 跨域访问
//   - 适用于远程访问和 Web 应用
//   - 实现：HTTPTransportServer
//
// 3. SSE 传输（Server-Sent Events）
//   - 支持服务器推送
//   - 实时事件流
//   - 心跳机制保持连接
//   - 实现：SSETransportServer
//
// # 核心组件
//
// HTTPTransportServer：
//   - HTTP 传输服务器实现
//   - 处理 MCP 协议的 HTTP 请求
//   - 支持 CORS 配置
//   - 集成工具注册表
//
// SSETransportServer：
//   - SSE 传输服务器实现
//   - 支持长连接和服务器推送
//   - 心跳机制维持连接
//   - 连接管理和限制
//
// ToolRegistry：
//   - 工具注册表，管理所有可用工具
//   - 线程安全的工具注册和调用
//   - 自动 panic recovery 保护
//   - 支持动态工具注册
//
// # 使用示例
//
// 创建 HTTP 传输服务器：
//
//	config := &types.HTTPConfig{
//	    Host:           "127.0.0.1",
//	    Port:           8080,
//	    Path:           "/mcp",
//	    EnableCORS:     true,
//	    AllowedOrigins: []string{"*"},
//	    Timeout:        30,
//	}
//	server, err := transport.NewHTTPTransportServer(config, logger)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 注册工具
//	fsService.RegisterToolsToRegistry(server.GetToolRegistry())
//
//	// 启动服务器
//	err = server.Start(ctx, mcpServer)
//
// 创建 SSE 传输服务器：
//
//	config := &types.SSEConfig{
//	    Host:              "127.0.0.1",
//	    Port:              8081,
//	    Path:              "/sse",
//	    EnableCORS:        true,
//	    AllowedOrigins:    []string{"*"},
//	    HeartbeatInterval: 30,
//	    MaxConnections:    100,
//	}
//	server, err := transport.NewSSETransportServer(config, logger)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 注册工具
//	fsService.RegisterToolsToRegistry(server.GetToolRegistry())
//
//	// 启动服务器
//	err = server.Start(ctx, mcpServer)
//
// 使用工具注册表：
//
//	registry := transport.NewToolRegistry()
//	registry.SetLogger(logger)
//
//	// 注册工具
//	registry.RegisterTool(&mcp.Tool{
//	    Name:        "my_tool",
//	    Description: "My custom tool",
//	}, func(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
//	    // 工具实现
//	    return &mcp.CallToolResult{
//	        Content: []mcp.Content{
//	            &mcp.TextContent{Text: "result"},
//	        },
//	    }, nil
//	})
//
//	// 调用工具
//	result, err := registry.CallTool(ctx, "my_tool", args)
//
// # MCP 协议处理
//
// 所有传输层都实现了以下 MCP 方法：
//
//   - initialize：初始化连接
//   - tools/list：列出所有可用工具
//   - tools/call：调用指定工具
//
// # 错误处理
//
// 传输层集成了 panic recovery 机制：
//
//   - 工具调用自动捕获 panic
//   - 记录完整的堆栈信息
//   - 转换为错误返回，不影响服务稳定性
//   - 详见 pkg/utils/recovery 包
//
// # 安全特性
//
//   - CORS 配置支持
//   - 来源白名单验证
//   - 请求超时控制
//   - 连接数限制（SSE）
//   - 工具调用隔离（panic recovery）
//
// # 性能优化
//
//   - 工具注册表缓存
//   - 并发安全的读写锁
//   - 连接复用（HTTP Keep-Alive）
//   - 心跳优化（SSE）
//
// # 测试
//
// 本包包含完整的单元测试：
//
//   - http_test.go：HTTP 传输测试
//   - sse_test.go：SSE 传输测试
//   - tool_registry_recovery_test.go：工具注册表 recovery 测试
//
// 运行测试：
//
//	go test -v ./pkg/transport/...
package transport
