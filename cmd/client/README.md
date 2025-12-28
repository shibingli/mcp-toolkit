# MCP Client 测试工具 / MCP Client Test Tool

## 简介 / Introduction

这是一个用于测试MCP Server完整功能的客户端程序。它会自动测试所有25个MCP工具,验证服务器的各项功能是否正常工作。

This is a client program for testing the complete functionality of the MCP Server. It automatically tests all 25 MCP tools to verify that all server functions work properly.

## 功能特性 / Features

- ✅ 覆盖所有25个MCP工具 / Covers all 25 MCP tools
- ✅ 自动化测试流程 / Automated testing process
- ✅ 详细的测试报告 / Detailed test reports
- ✅ 支持自定义服务器配置 / Supports custom server configuration
- ✅ 可选的详细日志输出 / Optional verbose logging

## 编译 / Build

```bash
# 使用sonic标签编译 / Build with sonic tag
go build -tags="sonic" -o mcp_client.exe ./cmd/client

# Linux/macOS
go build -tags="sonic" -o mcp_client ./cmd/client
```

## 使用方法 / Usage

### 基本使用 / Basic Usage

```bash
# 连接到默认服务器 (127.0.0.1:8080)
# Connect to default server (127.0.0.1:8080)
./mcp_client
```

### 自定义服务器 / Custom Server

```bash
# 指定服务器地址和端口 / Specify server address and port
./mcp_client -host 192.168.1.100 -port 9090

# 指定完整的服务器配置 / Specify complete server configuration
./mcp_client -host example.com -port 8080 -path /api/mcp
```

### 详细日志 / Verbose Logging

```bash
# 启用详细日志输出 / Enable verbose logging
./mcp_client -verbose
```

## 命令行参数 / Command Line Arguments

| 参数 / Argument | 类型 / Type | 默认值 / Default | 说明 / Description |
|----------------|------------|-----------------|-------------------|
| `-host` | string | 127.0.0.1 | MCP服务器主机地址 / MCP server host address |
| `-port` | int | 8080 | MCP服务器端口 / MCP server port |
| `-path` | string | /mcp | MCP服务器路径 / MCP server path |
| `-verbose` | bool | false | 启用详细日志 / Enable verbose logging |

## 测试流程 / Test Flow

1. **初始化连接** / Initialize Connection
   - 连接到MCP服务器 / Connect to MCP server
   - 验证协议版本 / Verify protocol version

2. **列出工具** / List Tools
   - 获取所有可用工具 / Get all available tools
   - 验证工具数量 / Verify tool count

3. **执行测试** / Execute Tests
   - 按顺序测试每个工具 / Test each tool in sequence
   - 验证响应数据 / Verify response data
   - 记录测试结果 / Record test results

4. **生成报告** / Generate Report
   - 统计通过/失败数量 / Count passed/failed tests
   - 显示详细结果 / Display detailed results

## 测试用例 / Test Cases

### 文件操作测试 / File Operation Tests
- 创建文件并写入内容 / Create file and write content
- 读取文件内容 / Read file content
- 修改文件内容 / Modify file content
- 复制和移动文件 / Copy and move files
- 删除文件 / Delete files

### 目录操作测试 / Directory Operation Tests
- 创建目录 / Create directory
- 列出目录内容 / List directory contents
- 切换工作目录 / Change working directory
- 搜索文件 / Search files

### 命令执行测试 / Command Execution Tests
- 同步执行命令 / Synchronous command execution
- 异步执行命令 / Asynchronous command execution
- 获取命令任务状态 / Get command task status
- 取消命令任务 / Cancel command task

### 高级功能测试 / Advanced Feature Tests
- 命令黑名单管理 / Command blacklist management
- 命令历史记录 / Command history
- 权限级别管理 / Permission level management

## 输出示例 / Output Example

```
=== MCP Client Test Suite ===
Connecting to: http://127.0.0.1:8080/mcp

1. Testing Initialize...
✅ Initialize successful
   Server: filesystem-mcp-server v1.0.1
   Protocol: 2024-11-05

2. Testing List Tools...
✅ List tools successful
   Found 25 tools

...

=== Test Results Summary ===
✅ create_file
✅ read_file
✅ write_file
...
✅ get_current_time

Total: 25, Passed: 25, Failed: 0
```

## 退出码 / Exit Codes

- `0` - 所有测试通过 / All tests passed
- `1` - 有测试失败 / Some tests failed

## 故障排查 / Troubleshooting

### 连接失败 / Connection Failed

```
❌ Initialize failed: dial tcp 127.0.0.1:8080: connect: connection refused
```

**解决方案** / Solution:
1. 确认MCP Server正在运行 / Ensure MCP Server is running
2. 检查端口号是否正确 / Check if port number is correct
3. 检查防火墙设置 / Check firewall settings

### 工具测试失败 / Tool Test Failed

```
❌ create_file: failed to create file: permission denied
```

**解决方案** / Solution:
1. 检查沙箱目录权限 / Check sandbox directory permissions
2. 查看服务器日志 / Check server logs
3. 使用`-verbose`参数获取更多信息 / Use `-verbose` flag for more information

## 开发 / Development

### 添加新测试 / Adding New Tests

1. 在`runAllToolTests`函数中添加测试调用 / Add test call in `runAllToolTests` function
2. 实现测试函数 / Implement test function
3. 返回`TestResult`结构 / Return `TestResult` structure

示例 / Example:

```go
func testNewTool(ctx context.Context, mcpClient client.Client) TestResult {
    fmt.Println("\nTesting new_tool...")
    
    resp, err := mcpClient.CallTool(ctx, "new_tool", types.NewToolRequest{
        // 参数 / Parameters
    })
    
    if err != nil {
        return TestResult{Name: "new_tool", Success: false, Error: err.Error()}
    }
    
    fmt.Printf("   Response: %s\n", resp.Content[0].Text)
    return TestResult{Name: "new_tool", Success: true}
}
```

## 相关文档 / Related Documentation

- [完整测试文档](../../docs/CLIENT_TEST.md) / [Complete Test Documentation](../../docs/CLIENT_TEST.md)
- [MCP客户端文档](../../docs/CLIENT.md) / [MCP Client Documentation](../../docs/CLIENT.md)
- [快速开始指南](../../docs/GETTING_STARTED.md) / [Getting Started Guide](../../docs/GETTING_STARTED.md)

## 许可证 / License

与主项目相同 / Same as main project

