# MCP 传输方式说明 / MCP Transport Documentation

## 概述 / Overview

MCP服务器支持多种传输方式，可以根据不同的使用场景选择合适的传输协议。

The MCP server supports multiple transport methods, allowing you to choose the appropriate transport protocol based on different usage scenarios.

## 支持的传输方式 / Supported Transport Methods

### 1. Stdio (标准输入输出) / Standard Input/Output

**默认传输方式 / Default Transport Method**

通过标准输入输出进行通信，适用于命令行工具和进程间通信。

Communicates through standard input/output, suitable for command-line tools and inter-process communication.

**启动方式 / How to Start:**

```bash
# 使用默认stdio传输 / Use default stdio transport
./mcp-toolkit

# 或明确指定stdio / Or explicitly specify stdio
./mcp-toolkit -transport stdio
```

**特点 / Features:**
- ✅ 简单易用 / Simple and easy to use
- ✅ 适合本地进程通信 / Suitable for local process communication
- ✅ 低延迟 / Low latency
- ❌ 不支持远程访问 / Does not support remote access

---

### 2. HTTP (支持 Streamable HTTP / Supports Streamable HTTP)

通过HTTP协议提供服务，支持远程访问和跨域请求。完全实现 MCP 规范的 Streamable HTTP 功能。

Provides services through HTTP protocol, supports remote access and cross-origin requests. Fully implements MCP
specification's Streamable HTTP features.

**启动方式 / How to Start:**

```bash
# 使用默认配置启动HTTP服务器 / Start HTTP server with default configuration
./mcp-toolkit -transport http

# 自定义主机和端口 / Customize host and port
./mcp-toolkit -transport http -http-host 0.0.0.0 -http-port 8080

# 禁用会话管理 / Disable session management
./mcp-toolkit -transport http -http-disable-session

# 禁用SSE流 / Disable SSE streaming
./mcp-toolkit -transport http -http-disable-sse

# 自定义会话超时和SSE心跳 / Customize session timeout and SSE heartbeat
./mcp-toolkit -transport http -http-session-timeout 3600 -http-sse-heartbeat 30
```

**默认配置 / Default Configuration:**
- Host: `127.0.0.1`
- Port: `8080`
- Path: `/mcp`
- CORS: 启用 / Enabled
- 读取超时 / Read Timeout: 30秒 / 30 seconds
- 写入超时 / Write Timeout: 30秒 / 30 seconds
- 会话管理 / Session Management: 启用 / Enabled
- 会话超时 / Session Timeout: 1800秒 / 1800 seconds (30分钟 / 30 minutes)
- SSE流 / SSE Streaming: 启用 / Enabled
- SSE心跳 / SSE Heartbeat: 30秒 / 30 seconds

**特点 / Features:**
- ✅ 支持远程访问 / Supports remote access
- ✅ 支持CORS跨域 / Supports CORS
- ✅ 标准HTTP协议 / Standard HTTP protocol
- ✅ 易于集成 / Easy to integrate
- ✅ **会话管理** / Session management
- ✅ **SSE流支持** / SSE streaming support
- ✅ **协议版本验证** / Protocol version validation
- ✅ **多种响应格式** / Multiple response formats (JSON/SSE)
- ❌ 相对stdio有额外开销 / Additional overhead compared to stdio

**健康检查 / Health Check:**
```bash
curl http://127.0.0.1:8080/health
```

#### Streamable HTTP 详细说明 / Streamable HTTP Details

HTTP 传输完全实现了 MCP 规范 (2025-06-18) 的 Streamable HTTP 功能。

HTTP transport fully implements MCP specification (2025-06-18) Streamable HTTP features.

**1. 会话管理 / Session Management**

会话通过 `Mcp-Session-Id` 头进行管理，在初始化时创建，后续请求必须携带。

Sessions are managed through the `Mcp-Session-Id` header, created during initialization, and must be included in
subsequent requests.

```bash
# 初始化会话 / Initialize session
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

# 响应头包含: Mcp-Session-Id: <session-id>
# Response header includes: Mcp-Session-Id: <session-id>

# 使用会话调用工具 / Call tool with session
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "Mcp-Session-Id: <session-id>" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }'

# 终止会话 / Terminate session
curl -X DELETE http://localhost:8080/mcp \
  -H "Mcp-Session-Id: <session-id>"
```

**2. 协议版本支持 / Protocol Version Support**

通过 `MCP-Protocol-Version` 头指定协议版本。

Specify protocol version through `MCP-Protocol-Version` header.

支持的版本 / Supported versions:

- `2025-06-18` (最新 / Latest)
- `2025-03-26`
- `2024-11-05`

```bash
# 使用特定协议版本 / Use specific protocol version
curl -X POST http://localhost:8080/mcp \
  -H "MCP-Protocol-Version: 2025-06-18" \
  -H "Content-Type: application/json" \
  -d '...'
```

**3. SSE 流支持 / SSE Streaming Support**

支持通过 `Accept` 头选择响应格式。

Support response format selection through `Accept` header.

```bash
# JSON 响应 / JSON response
curl -X POST http://localhost:8080/mcp \
  -H "Accept: application/json" \
  -H "Mcp-Session-Id: <session-id>" \
  -d '...'

# SSE 流响应 / SSE stream response
curl -X POST http://localhost:8080/mcp \
  -H "Accept: text/event-stream" \
  -H "Mcp-Session-Id: <session-id>" \
  -d '...'

# 打开 SSE 流监听服务器消息 / Open SSE stream to listen for server messages
curl -X GET http://localhost:8080/mcp \
  -H "Accept: text/event-stream" \
  -H "Mcp-Session-Id: <session-id>"
```

**4. HTTP 方法 / HTTP Methods**

| 方法 / Method | 用途 / Purpose             | 说明 / Description                    |
|-------------|--------------------------|-------------------------------------|
| POST        | 发送消息 / Send message      | 发送 JSON-RPC 消息到服务器，支持 JSON 和 SSE 响应 |
| GET         | 打开流 / Open stream        | 打开 SSE 流监听服务器推送的消息                  |
| DELETE      | 终止会话 / Terminate session | 终止当前会话并清理资源                         |
| OPTIONS     | CORS 预检 / CORS preflight | 处理跨域请求预检                            |

**5. 配置选项 / Configuration Options**

| 参数 / Parameter          | 说明 / Description                              | 默认值 / Default |
|-------------------------|-----------------------------------------------|---------------|
| `-http-enable-session`  | 启用会话管理 / Enable session management            | `true`        |
| `-http-disable-session` | 禁用会话管理 / Disable session management           | -             |
| `-http-session-timeout` | 会话超时(秒) / Session timeout (seconds)           | `1800`        |
| `-http-enable-sse`      | 启用SSE流 / Enable SSE streaming                 | `true`        |
| `-http-disable-sse`     | 禁用SSE流 / Disable SSE streaming                | -             |
| `-http-sse-heartbeat`   | SSE心跳间隔(秒) / SSE heartbeat interval (seconds) | `30`          |

**6. 错误处理 / Error Handling**

```bash
# 缺少会话ID / Missing session ID
HTTP 400 Bad Request
{"error": "Missing Mcp-Session-Id header"}

# 无效的会话ID / Invalid session ID
HTTP 400 Bad Request
{"error": "Invalid session"}

# 不支持的协议版本 / Unsupported protocol version
HTTP 400 Bad Request
{"error": "Unsupported protocol version"}
```

---

### 3. SSE (Server-Sent Events)

通过SSE协议提供实时事件流，适合需要服务器主动推送的场景。

Provides real-time event streams through SSE protocol, suitable for scenarios requiring server-initiated push.

**启动方式 / How to Start:**

```bash
# 使用默认配置启动SSE服务器 / Start SSE server with default configuration
./mcp-toolkit -transport sse

# 自定义主机和端口 / Customize host and port
./mcp-toolkit -transport sse -sse-host 0.0.0.0 -sse-port 8081
```

**默认配置 / Default Configuration:**
- Host: `127.0.0.1`
- Port: `8081`
- Path: `/sse`
- CORS: 启用 / Enabled
- 心跳间隔 / Heartbeat Interval: 30秒 / 30 seconds
- 最大连接数 / Max Connections: 100

**特点 / Features:**
- ✅ 支持服务器推送 / Supports server push
- ✅ 长连接 / Long-lived connections
- ✅ 自动重连 / Automatic reconnection
- ✅ 支持CORS跨域 / Supports CORS
- ❌ 单向通信 / One-way communication

**健康检查 / Health Check:**
```bash
curl http://127.0.0.1:8081/health
```

---

## 命令行参数 / Command Line Arguments

### 基础参数 / Basic Parameters

| 参数 / Parameter | 说明 / Description                        | 默认值 / Default |
|----------------|-----------------------------------------|---------------|
| `-transport`   | 传输类型: stdio, http, sse / Transport type | `stdio`       |
| `-sandbox`     | 沙箱目录路径 / Sandbox directory path         | `./sandbox`   |

### HTTP 参数 / HTTP Parameters

| 参数 / Parameter          | 说明 / Description                              | 默认值 / Default |
|-------------------------|-----------------------------------------------|---------------|
| `-http-host`            | HTTP监听地址 / HTTP listen address                | `127.0.0.1`   |
| `-http-port`            | HTTP监听端口 / HTTP listen port                   | `8080`        |
| `-http-enable-session`  | 启用会话管理 / Enable session management            | `true`        |
| `-http-disable-session` | 禁用会话管理 / Disable session management           | -             |
| `-http-session-timeout` | 会话超时(秒) / Session timeout (seconds)           | `1800`        |
| `-http-enable-sse`      | 启用SSE流 / Enable SSE streaming                 | `true`        |
| `-http-disable-sse`     | 禁用SSE流 / Disable SSE streaming                | -             |
| `-http-sse-heartbeat`   | SSE心跳间隔(秒) / SSE heartbeat interval (seconds) | `30`          |

### SSE 参数 / SSE Parameters

| 参数 / Parameter | 说明 / Description             | 默认值 / Default |
|----------------|------------------------------|---------------|
| `-sse-host`    | SSE监听地址 / SSE listen address | `127.0.0.1`   |
| `-sse-port`    | SSE监听端口 / SSE listen port    | `8081`        |

## 使用示例 / Usage Examples

### 本地开发 / Local Development
```bash
# 使用stdio进行本地开发和测试
./mcp-toolkit -transport stdio -sandbox ./dev-sandbox
```

### 远程服务 / Remote Service
```bash
# 启动HTTP服务器供远程客户端访问
./mcp-toolkit -transport http -http-host 0.0.0.0 -http-port 8080
```

### 实时推送 / Real-time Push
```bash
# 启动SSE服务器用于实时事件推送
./mcp-toolkit -transport sse -sse-host 0.0.0.0 -sse-port 8081
```

## 安全建议 / Security Recommendations

1. **生产环境 / Production Environment:**
   - 使用反向代理（如Nginx）处理HTTPS
   - 配置适当的CORS策略
   - 限制访问IP范围
   - 启用身份验证和授权

2. **开发环境 / Development Environment:**
   - 使用localhost绑定
   - 启用详细日志
   - 使用沙箱目录隔离

## 性能对比 / Performance Comparison

| 传输方式 / Transport | 延迟 / Latency | 吞吐量 / Throughput | 并发 / Concurrency |
|---------------------|---------------|--------------------|--------------------|
| Stdio | 最低 / Lowest | 最高 / Highest | 单进程 / Single process |
| HTTP | 中等 / Medium | 高 / High | 高 / High |
| SSE | 中等 / Medium | 中等 / Medium | 中等 / Medium |

