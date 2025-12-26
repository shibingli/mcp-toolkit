// Package client 提供了 MCP Toolkit 的客户端实现
//
// 本包实现了 MCP 协议的客户端，用于连接和调用 MCP 服务器。
//
// # 主要功能
//
//   - 连接到 MCP 服务器（HTTP 传输）
//   - 初始化 MCP 会话
//   - 列出可用工具
//   - 调用服务器工具
//   - 自动处理 JSON-RPC 协议
//
// # 核心接口
//
// Client 接口定义了客户端的基本操作：
//
//	type Client interface {
//	    Initialize(ctx context.Context, protocolVersion string) (*InitializeResponse, error)
//	    ListTools(ctx context.Context) (*ListToolsResponse, error)
//	    CallTool(ctx context.Context, name string, arguments interface{}) (*CallToolResponse, error)
//	    Close() error
//	}
//
// # 实现
//
// HTTPClient：
//   - HTTP 传输的客户端实现
//   - 支持自定义主机、端口和路径
//   - 自动处理 JSON 序列化/反序列化
//   - 集成日志记录
//
// # 使用示例
//
// 创建客户端并连接：
//
//	logger := zap.NewNop()
//	client := client.NewHTTPClient("127.0.0.1", 8080, "/mcp", logger)
//	defer client.Close()
//
//	ctx := context.Background()
//
//	// 初始化连接
//	initResp, err := client.Initialize(ctx, types.ProtocolVersion)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Connected to: %s v%s\n",
//	    initResp.ServerInfo.Name,
//	    initResp.ServerInfo.Version)
//
// 列出工具：
//
//	toolsResp, err := client.ListTools(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, tool := range toolsResp.Tools {
//	    fmt.Printf("Tool: %s - %s\n", tool.Name, tool.Description)
//	}
//
// 调用工具：
//
//	resp, err := client.CallTool(ctx, "create_file", types.CreateFileRequest{
//	    Path:    "test.txt",
//	    Content: "Hello, World!",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(resp.Content[0].Text)
//
// # 数据类型
//
// 请求/响应类型（types.go）：
//
//   - InitializeResponse：初始化响应
//   - ServerInfo：服务器信息
//   - ListToolsResponse：工具列表响应
//   - Tool：工具定义
//   - CallToolResponse：工具调用响应
//   - Content：响应内容
//
// # 错误处理
//
// 客户端会处理以下错误：
//
//   - 网络连接错误
//   - JSON 序列化/反序列化错误
//   - MCP 协议错误
//   - 服务器返回的业务错误
//
// 所有错误都会包含详细的上下文信息。
//
// # 测试
//
// 本包包含完整的单元测试：
//
//   - client_test.go：客户端功能测试
//
// 运行测试：
//
//	go test -v ./pkg/client/...
//
// # 测试客户端
//
// cmd/client 目录提供了一个完整的测试客户端程序，
// 可以用于测试所有 MCP 工具的功能。
//
// 运行测试客户端：
//
//	./mcp-toolkit-client -host 127.0.0.1 -port 8080
//
// # 扩展
//
// 要添加新的传输方式（如 WebSocket），只需：
//
//  1. 实现 Client 接口
//  2. 处理相应的传输协议
//  3. 添加相应的测试
//
// # 性能
//
//   - 使用连接池复用 HTTP 连接
//   - 支持高性能 JSON 库（Sonic）
//   - 最小化内存分配
//
// # 日志
//
// 客户端使用 zap 进行结构化日志记录：
//
//   - 请求/响应日志
//   - 错误日志
//   - 性能指标
//
// 可以传入 nil logger 禁用日志。
//
// # 线程安全
//
// HTTPClient 是线程安全的，可以在多个 goroutine 中并发使用。
package client
