# MCP 客户端使用指南 / MCP Client Usage Guide

## 概述 / Overview

本项目提供了一个完整的 MCP (Model Context Protocol) 客户端实现，用于测试和验证 MCP 服务器的所有功能。

This project provides a complete MCP (Model Context Protocol) client implementation for testing and validating all MCP server functionalities.

## 功能特性 / Features

- ✅ 支持 HTTP 传输 / HTTP transport support
- ✅ 完整的 MCP 协议实现 / Complete MCP protocol implementation
- ✅ 所有文件系统工具测试 / All filesystem tools testing
- ✅ 详细的测试报告 / Detailed test reports
- ✅ 单元测试覆盖 / Unit test coverage

## 快速开始 / Quick Start

### 1. 启动 MCP 服务器 / Start MCP Server

首先启动 MCP 服务器（HTTP 模式）：

```bash
# Windows
mcp-toolkit.exe -transport http -http-host 127.0.0.1 -http-port 8080

# Linux/macOS
./mcp-toolkit -transport http -http-host 127.0.0.1 -http-port 8080
```

### 2. 运行客户端测试 / Run Client Tests

在另一个终端运行客户端测试程序：

```bash
# 编译客户端 / Build client
go build -tags="sonic" -o mcp_client cmd/client/main.go

# 运行测试（使用默认配置） / Run tests (with default config)
./mcp_client

# 运行测试（自定义配置） / Run tests (with custom config)
./mcp_client -host 127.0.0.1 -port 8080 -path /mcp

# 启用详细日志 / Enable verbose logging
./mcp_client -verbose
```

## 命令行参数 / Command Line Arguments

| 参数 / Parameter | 默认值 / Default | 说明 / Description |
|-----------------|-----------------|-------------------|
| `-host` | 127.0.0.1 | MCP 服务器地址 / MCP server host |
| `-port` | 8080 | MCP 服务器端口 / MCP server port |
| `-path` | /mcp | MCP 服务器路径 / MCP server path |
| `-verbose` | false | 启用详细日志 / Enable verbose logging |

## 测试覆盖 / Test Coverage

客户端测试程序会自动测试以下所有 MCP 工具：

The client test program automatically tests all the following MCP tools:

1. ✅ **create_file** - 创建文件 / Create file
2. ✅ **read_file** - 读取文件 / Read file
3. ✅ **write_file** - 写入文件 / Write file
4. ✅ **create_directory** - 创建目录 / Create directory
5. ✅ **list_directory** - 列出目录 / List directory
6. ✅ **file_stat** - 文件状态 / File stat
7. ✅ **file_exists** - 文件是否存在 / File exists
8. ✅ **search_files** - 搜索文件 / Search files
9. ✅ **copy** - 复制文件 / Copy file
10. ✅ **move** - 移动文件 / Move file
11. ✅ **batch_delete** - 批量删除 / Batch delete
12. ✅ **delete** - 删除 / Delete
13. ✅ **get_current_time** - 获取当前时间 / Get current time

## 编程使用 / Programmatic Usage

### 创建客户端 / Create Client

```go
package main

import (
    "context"
    "fmt"

    "mcp-toolkit/pkg/client"
    "mcp-toolkit/pkg/types"

    "go.uber.org/zap"
)

func main() {
    // 创建日志记录器 / Create logger
    logger := zap.NewNop()
    
    // 创建 HTTP 客户端 / Create HTTP client
    mcpClient := client.NewHTTPClient("127.0.0.1", 8080, "/mcp", logger)
    defer mcpClient.Close()
    
    ctx := context.Background()

    // 初始化连接 / Initialize connection
    initResp, err := mcpClient.Initialize(ctx, types.ProtocolVersion)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Connected to: %s v%s\n",
        initResp.ServerInfo.Name,
        initResp.ServerInfo.Version)
}
```

### 列出工具 / List Tools

```go
// 列出所有可用工具 / List all available tools
toolsResp, err := mcpClient.ListTools(ctx)
if err != nil {
    panic(err)
}

for _, tool := range toolsResp.Tools {
    fmt.Printf("Tool: %s - %s\n", tool.Name, tool.Description)
}
```

### 调用工具 / Call Tool

```go
// 创建文件 / Create file
resp, err := mcpClient.CallTool(ctx, "create_file", types.CreateFileRequest{
    Path:    "hello.txt",
    Content: "Hello, World!",
})
if err != nil {
    panic(err)
}

fmt.Println(resp.Content[0].Text)
```

### 读取文件 / Read File

```go
// 读取文件 / Read file
resp, err := mcpClient.CallTool(ctx, "read_file", types.ReadFileRequest{
    Path: "hello.txt",
})
if err != nil {
    panic(err)
}

// 解析响应 / Parse response
var readResp types.ReadFileResponse
json.Unmarshal([]byte(resp.Content[0].Text), &readResp)
fmt.Println(readResp.Content)
```

## 单元测试 / Unit Tests

运行客户端单元测试：

```bash
# 运行所有测试 / Run all tests
go test -v ./pkg/client/...

# 运行测试并查看覆盖率 / Run tests with coverage
go test -v -cover ./pkg/client/...

# 生成覆盖率报告 / Generate coverage report
go test -v -coverprofile=coverage.out ./pkg/client/...
go tool cover -html=coverage.out
```

## 示例输出 / Example Output

```
=== MCP Client Test Suite ===
Connecting to: http://127.0.0.1:8080/mcp

1. Testing Initialize...
✅ Initialize successful
   Server: mcp-toolkit v1.0.1
   Protocol: 2024-11-05

2. Testing List Tools...
✅ List tools successful
   Found 13 tools:
   1. create_file - 创建新文件并写入内容 / Create a new file and write content
   2. create_directory - 创建新目录 / Create a new directory
   ...

3. Testing create_file...
   Response: {"success":true,"message":"文件创建成功 / File created successfully"}

...

=== Test Results Summary ===
✅ create_file
✅ read_file
✅ write_file
✅ create_directory
✅ list_directory
✅ file_stat
✅ file_exists
✅ search_files
✅ copy
✅ move
✅ batch_delete
✅ delete
✅ get_current_time

Total: 13, Passed: 13, Failed: 0
```

## 故障排除 / Troubleshooting

### 连接失败 / Connection Failed

如果客户端无法连接到服务器：

1. 确认服务器正在运行 / Confirm server is running
2. 检查主机和端口配置 / Check host and port configuration
3. 确认防火墙设置 / Verify firewall settings

### 测试失败 / Test Failed

如果某个测试失败：

1. 查看详细错误信息 / Check detailed error message
2. 使用 `-verbose` 参数查看详细日志 / Use `-verbose` flag for detailed logs
3. 检查服务器日志 / Check server logs
4. 确认沙箱目录权限 / Verify sandbox directory permissions

## 开发指南 / Development Guide

### 添加新测试 / Adding New Tests

要添加新的工具测试，在 `cmd/client/main.go` 中：

1. 创建测试函数 / Create test function
2. 在 `runAllToolTests` 中调用 / Call in `runAllToolTests`
3. 返回 `TestResult` / Return `TestResult`

示例 / Example:

```go
func testNewTool(ctx context.Context, mcpClient client.Client) TestResult {
    fmt.Println("\nTesting new_tool...")
    
    resp, err := mcpClient.CallTool(ctx, "new_tool", map[string]interface{}{
        "param": "value",
    })
    
    if err != nil {
        return TestResult{Name: "new_tool", Success: false, Error: err.Error()}
    }
    
    fmt.Printf("   Response: %s\n", resp.Content[0].Text)
    return TestResult{Name: "new_tool", Success: true}
}
```

## 相关文档 / Related Documentation

- [README.md](README.md) - 项目主文档 / Main project documentation
- [EXAMPLE.md](EXAMPLE.md) - MCP 工具使用示例 / MCP tools usage examples
- [QUICKSTART.md](QUICKSTART.md) - 快速启动指南 / Quick start guide

## 许可证 / License

本项目采用与主项目相同的许可证。

This project uses the same license as the main project.

