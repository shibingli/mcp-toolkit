# Streamable HTTP 使用指南 / Streamable HTTP Usage Guide

## 概述 / Overview

MCP Toolkit 完全实现了 MCP 规范 (2025-06-18) 的 Streamable HTTP 功能，提供会话管理、SSE 流支持和协议版本验证。

MCP Toolkit fully implements MCP specification (2025-06-18) Streamable HTTP features, providing session management, SSE
streaming support, and protocol version validation.

## 核心特性 / Core Features

### 1. 会话管理 / Session Management

- 自动生成加密安全的会话 ID
- 通过 `Mcp-Session-Id` 头管理会话
- 可配置的会话超时时间 (默认 30 分钟)
- 支持会话终止

### 2. 协议版本支持 / Protocol Version Support

- 支持 `MCP-Protocol-Version` 头
- 兼容多个协议版本: 2025-06-18, 2025-03-26, 2024-11-05
- 向后兼容旧版本

### 3. SSE 流支持 / SSE Streaming Support

- POST 请求支持 JSON 和 SSE 响应
- GET 请求打开 SSE 流用于服务器推送
- 可配置的心跳间隔 (默认 30 秒)
- 自动连接保持

## 快速开始 / Quick Start

### 启动服务器 / Start Server

```bash
# 使用默认配置 (启用会话管理和SSE)
./mcp-toolkit -transport http

# 自定义配置
./mcp-toolkit -transport http \
  -http-host 0.0.0.0 \
  -http-port 8080 \
  -http-session-timeout 3600 \
  -http-sse-heartbeat 30
```

### 基本工作流程 / Basic Workflow

#### 1. 初始化会话 / Initialize Session

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "MCP-Protocol-Version: 2025-06-18" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2025-06-18",
      "capabilities": {},
      "clientInfo": {
        "name": "example-client",
        "version": "1.0.0"
      }
    }
  }'
```

**响应示例 / Response Example:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "2025-06-18",
    "capabilities": {
      "tools": {}
    },
    "serverInfo": {
      "name": "mcp-toolkit",
      "version": "1.0.0"
    }
  }
}
```

响应头包含 / Response headers include:

```
Mcp-Session-Id: <session-id>
```

#### 2. 调用工具 / Call Tools

```bash
# 列出所有工具
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "Mcp-Session-Id: <session-id>" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }'

# 调用特定工具
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "Mcp-Session-Id: <session-id>" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "get_current_time"
    }
  }'
```

#### 3. 使用 SSE 流 / Use SSE Streaming

```bash
# POST 请求返回 SSE 流
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -H "Mcp-Session-Id: <session-id>" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/list"
  }'

# GET 请求打开 SSE 流监听服务器消息
curl -X GET http://localhost:8080/mcp \
  -H "Accept: text/event-stream" \
  -H "Mcp-Session-Id: <session-id>"
```

**SSE 响应格式 / SSE Response Format:**

```
data: {"jsonrpc":"2.0","id":4,"result":{...}}

```

#### 4. 终止会话 / Terminate Session

```bash
curl -X DELETE http://localhost:8080/mcp \
  -H "Mcp-Session-Id: <session-id>"
```

## 配置选项 / Configuration Options

### 启用/禁用功能 / Enable/Disable Features

```bash
# 禁用会话管理 (不推荐)
./mcp-toolkit -transport http -http-disable-session

# 禁用 SSE 流
./mcp-toolkit -transport http -http-disable-sse

# 同时禁用两者 (回退到基本 HTTP)
./mcp-toolkit -transport http -http-disable-session -http-disable-sse
```

### 调整超时和心跳 / Adjust Timeouts and Heartbeat

```bash
# 设置会话超时为 1 小时
./mcp-toolkit -transport http -http-session-timeout 3600

# 设置 SSE 心跳为 60 秒
./mcp-toolkit -transport http -http-sse-heartbeat 60
```

## HTTP 方法说明 / HTTP Methods

| 方法 / Method | 用途 / Purpose   | 必需头 / Required Headers                |
|-------------|----------------|---------------------------------------|
| POST        | 发送 JSON-RPC 消息 | Content-Type, Accept, Mcp-Session-Id* |
| GET         | 打开 SSE 流       | Accept, Mcp-Session-Id                |
| DELETE      | 终止会话           | Mcp-Session-Id                        |
| OPTIONS     | CORS 预检        | -                                     |

\* 初始化请求不需要 Mcp-Session-Id

## 错误处理 / Error Handling

### 常见错误 / Common Errors

**1. 缺少会话 ID / Missing Session ID**

```bash
HTTP 400 Bad Request
{"error": "Missing Mcp-Session-Id header"}
```

**2. 无效的会话 ID / Invalid Session ID**

```bash
HTTP 400 Bad Request
{"error": "Invalid session"}
```

**3. 不支持的协议版本 / Unsupported Protocol Version**

```bash
HTTP 400 Bad Request
{"error": "Unsupported protocol version: xxx"}
```

**4. 会话超时 / Session Timeout**

```bash
HTTP 400 Bad Request
{"error": "Session expired"}
```

## 最佳实践 / Best Practices

1. **始终使用协议版本头 / Always Use Protocol Version Header**
   ```bash
   -H "MCP-Protocol-Version: 2025-06-18"
   ```

2. **保存会话 ID / Save Session ID**
    - 从初始化响应中提取会话 ID
    - 在所有后续请求中使用相同的会话 ID

3. **处理会话过期 / Handle Session Expiration**
    - 实现会话过期检测
    - 自动重新初始化会话

4. **选择合适的响应格式 / Choose Appropriate Response Format**
    - 使用 `Accept: application/json` 用于一次性请求
    - 使用 `Accept: text/event-stream` 用于流式响应

5. **正确终止会话 / Properly Terminate Sessions**
    - 在完成工作后发送 DELETE 请求
    - 释放服务器资源

## 示例代码 / Example Code

完整的示例代码请参考 [EXAMPLE.md](EXAMPLE.md)。

For complete example code, see [EXAMPLE.md](EXAMPLE.md).

