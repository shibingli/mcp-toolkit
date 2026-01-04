# MCP 工具列表查询指南 / MCP Tools List Query Guide

本文档介绍如何使用 MCP Client 获取 MCP Server 的工具集，支持所有传输模式。

This document describes how to use MCP Client to retrieve the MCP Server's tool set, supporting all transport modes.

## 支持的传输模式 / Supported Transport Modes

- **Stdio** - 标准输入/输出，适用于本地进程通信
- **HTTP** - HTTP 传输，适用于远程访问
- **SSE** - Server-Sent Events，支持服务器推送
- **Streamable HTTP** - HTTP with session management and SSE streaming

## 快速开始 / Quick Start

### 方法 1：使用自动化脚本 / Method 1: Using Automated Scripts

#### Windows

```bash
# HTTP 模式（默认）/ HTTP mode (default)
scripts\list_tools.bat

# SSE 模式 / SSE mode
scripts\list_tools.bat sse

# Stdio 模式 / Stdio mode
scripts\list_tools.bat stdio

# 详细模式 / Detailed mode
scripts\list_tools.bat http detailed
scripts\list_tools.bat sse detailed
scripts\list_tools.bat stdio detailed
```

#### Linux/macOS

```bash
# HTTP 模式（默认）/ HTTP mode (default)
chmod +x scripts/list_tools.sh
./scripts/list_tools.sh

# SSE 模式 / SSE mode
./scripts/list_tools.sh sse

# Stdio 模式 / Stdio mode
./scripts/list_tools.sh stdio

# 详细模式 / Detailed mode
./scripts/list_tools.sh http detailed
./scripts/list_tools.sh sse detailed
./scripts/list_tools.sh stdio detailed
```

### 方法 2：手动运行 / Method 2: Manual Execution

#### HTTP 模式 / HTTP Mode

```bash
# 1. 启动 HTTP 服务器 / Start HTTP server
./mcp-toolkit.exe -transport http -http-port 8080

# 2. 运行客户端 / Run client
go run cmd/client/list_tools.go -transport http

# 详细模式 / Detailed mode
go run cmd/client/list_tools.go -transport http -detailed

# 自定义地址 / Custom address
go run cmd/client/list_tools.go -transport http -host 192.168.1.100 -port 8080
```

#### SSE 模式 / SSE Mode

```bash
# 1. 启动 SSE 服务器 / Start SSE server
./mcp-toolkit.exe -transport sse -sse-port 8081

# 2. 运行客户端 / Run client
go run cmd/client/list_tools.go -transport sse -port 8081 -path /sse

# 详细模式 / Detailed mode
go run cmd/client/list_tools.go -transport sse -port 8081 -path /sse -detailed
```

#### Stdio 模式 / Stdio Mode

```bash
# 直接运行（服务器作为子进程启动）/ Direct run (server as subprocess)
go run cmd/client/list_tools.go -transport stdio -command ./mcp-toolkit.exe -args "-transport,stdio"

# 详细模式 / Detailed mode
go run cmd/client/list_tools.go -transport stdio -command ./mcp-toolkit.exe -args "-transport,stdio" -detailed
```

## 命令行参数 / Command Line Parameters

| 参数 / Parameter | 说明 / Description                                        | 默认值 / Default |
|----------------|---------------------------------------------------------|---------------|
| `-transport`   | 传输模式：stdio, http, sse / Transport mode                  | `http`        |
| `-host`        | MCP 服务器主机地址（http/sse）/ MCP server host (http/sse)       | `127.0.0.1`   |
| `-port`        | MCP 服务器端口（http/sse）/ MCP server port (http/sse)         | `8080`        |
| `-path`        | MCP 服务器路径（http/sse）/ MCP server path (http/sse)         | `/mcp`        |
| `-command`     | 启动命令（stdio）/ Launch command (stdio)                     | -             |
| `-args`        | 命令参数（stdio，逗号分隔）/ Command args (stdio, comma-separated) | -             |
| `-detailed`    | 显示详细信息（包括参数）/ Show detailed info (including parameters) | `false`       |

## 输出示例 / Output Examples

### 简单模式输出 / Simple Mode Output

```
╔════════════════════════════════════════════════════════════════╗
║           MCP Server 工具集查询 / MCP Tools List              ║
╚════════════════════════════════════════════════════════════════╝

连接地址 / Server URL: http://127.0.0.1:8080/mcp

正在连接服务器... / Connecting to server...
✅ 连接成功 / Connected successfully
   服务器 / Server: mcp-toolkit v1.0.1
   协议版本 / Protocol: 2025-12-26

正在获取工具列表... / Fetching tools list...
✅ 成功获取 33 个工具 / Successfully retrieved 33 tools

════════════════════════════════════════════════════════════════
                        工具列表 / Tools List                   
════════════════════════════════════════════════════════════════

【文件操作 / File Operations】(9 个工具)
────────────────────────────────────────────────────────────────
1. create_file
   Create a new file with the specified content...
2. read_file
   Read and return the content of a file...
...
```

### 详细模式输出 / Detailed Mode Output

```
【文件操作 / File Operations】(9 个工具)
────────────────────────────────────────────────────────────────
1. create_file
   描述 / Description: Create a new file with the specified content...
   参数 / Parameters:
     - path (string) [必填 / Required]
     - content (string) [必填 / Required]

2. read_file
   描述 / Description: Read and return the content of a file...
   参数 / Parameters:
     - path (string) [必填 / Required]
...
```

## 工具分类 / Tool Categories

工具列表按以下类别组织：

Tools are organized into the following categories:

1. **文件操作 / File Operations** - 文件的创建、读取、写入、删除、复制、移动等
2. **目录操作 / Directory Operations** - 目录的创建、列出、删除、复制、移动等
3. **文件信息 / File Information** - 文件状态、存在性检查、搜索等
4. **系统信息 / System Information** - 获取系统信息（OS、CPU、内存、GPU、网卡）
5. **命令执行 / Command Execution** - 同步和异步命令执行
6. **命令管理 / Command Management** - 命令黑名单、历史记录、任务管理
7. **权限管理 / Permission** - 权限级别的获取和设置
8. **时间工具 / Time** - 获取当前时间
9. **下载工具 / Download** - HTTP/HTTPS 文件下载

## 技术细节 / Technical Details

### 客户端实现 / Client Implementation

- **文件位置 / File Location**: `cmd/client/list_tools.go`
- **传输协议 / Transport Protocol**: HTTP with JSON-RPC 2.0
- **会话管理 / Session Management**: 自动管理会话 ID / Automatic session ID management
- **响应格式 / Response Format**: JSON (通过 `Accept: application/json` 头指定)

### 服务器要求 / Server Requirements

- **传输模式 / Transport Mode**: HTTP
- **端口 / Port**: 默认 8080 / Default 8080
- **路径 / Path**: 默认 `/mcp` / Default `/mcp`

## 故障排除 / Troubleshooting

### 连接失败 / Connection Failed

```
❌ 连接失败 / Connection failed: dial tcp 127.0.0.1:8080: connectex: No connection could be made
```

**解决方案 / Solution**:

1. 确保 MCP Server 正在运行 / Ensure MCP Server is running
2. 检查端口是否正确 / Check if port is correct
3. 确认使用 HTTP 传输模式 / Confirm using HTTP transport mode

### 获取工具列表失败 / Failed to Get Tools List

```
❌ 获取工具列表失败 / Failed to get tools list
```

**解决方案 / Solution**:

1. 检查服务器日志 / Check server logs
2. 确认客户端和服务器版本兼容 / Confirm client and server version compatibility
3. 尝试重启服务器 / Try restarting the server

## 相关文档 / Related Documentation

- [MCP Server 使用指南](../README.md)
- [HTTP 传输配置](TRANSPORT.md)
- [Streamable HTTP 文档](STREAMABLE_HTTP.md)

