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

### 2. HTTP

通过HTTP协议提供服务，支持远程访问和跨域请求。

Provides services through HTTP protocol, supports remote access and cross-origin requests.

**启动方式 / How to Start:**

```bash
# 使用默认配置启动HTTP服务器 / Start HTTP server with default configuration
./mcp-toolkit -transport http

# 自定义主机和端口 / Customize host and port
./mcp-toolkit -transport http -http-host 0.0.0.0 -http-port 8080
```

**默认配置 / Default Configuration:**
- Host: `127.0.0.1`
- Port: `8080`
- Path: `/mcp`
- CORS: 启用 / Enabled
- 读取超时 / Read Timeout: 30秒 / 30 seconds
- 写入超时 / Write Timeout: 30秒 / 30 seconds

**特点 / Features:**
- ✅ 支持远程访问 / Supports remote access
- ✅ 支持CORS跨域 / Supports CORS
- ✅ 标准HTTP协议 / Standard HTTP protocol
- ✅ 易于集成 / Easy to integrate
- ❌ 相对stdio有额外开销 / Additional overhead compared to stdio

**健康检查 / Health Check:**
```bash
curl http://127.0.0.1:8080/health
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

| 参数 / Parameter | 说明 / Description | 默认值 / Default |
|-----------------|-------------------|-----------------|
| `-transport` | 传输类型: stdio, http, sse / Transport type | `stdio` |
| `-sandbox` | 沙箱目录路径 / Sandbox directory path | `./sandbox` |
| `-http-host` | HTTP监听地址 / HTTP listen address | `127.0.0.1` |
| `-http-port` | HTTP监听端口 / HTTP listen port | `8080` |
| `-sse-host` | SSE监听地址 / SSE listen address | `127.0.0.1` |
| `-sse-port` | SSE监听端口 / SSE listen port | `8081` |

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

