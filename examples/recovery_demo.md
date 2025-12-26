# Recovery 功能演示 / Recovery Feature Demo

本文档演示 MCP Server 的 Panic Recovery 功能如何保护服务稳定性。

This document demonstrates how the MCP Server's Panic Recovery feature protects service stability.

## 场景 1: 工具层 Panic 恢复 / Scenario 1: Tool-Level Panic Recovery

### 问题场景 / Problem Scenario

假设某个工具在处理请求时发生了 panic：

Suppose a tool panics while processing a request:

```go
func (s *Service) handleReadFile(ctx context.Context, req *mcp.CallToolRequest, args types.ReadFileRequest) (*mcp.CallToolResult, *types.ReadFileResponse, error) {
    // 假设这里发生了 panic
    panic("unexpected nil pointer")
}
```

### 没有 Recovery 的情况 / Without Recovery

```
❌ 整个服务崩溃
❌ 所有正在处理的请求失败
❌ 需要重启服务
❌ 用户体验极差
```

### 有 Recovery 的情况 / With Recovery

```
✅ 只有当前请求失败
✅ 其他请求正常处理
✅ 服务继续运行
✅ 错误被记录到日志
✅ 返回友好的错误响应
```

**日志输出 / Log Output:**
```
ERROR   panic recovered
  panic: unexpected nil pointer
  stack: goroutine 123 [running]:
    mcp_server/internal/services/filesystem.(*Service).handleReadFile(...)
    ...
```

**客户端收到的响应 / Client Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32603,
    "message": "panic: unexpected nil pointer"
  }
}
```

## 场景 2: 传输层 Panic 恢复 / Scenario 2: Transport-Level Panic Recovery

### 问题场景 / Problem Scenario

HTTP 请求处理过程中发生 panic：

Panic occurs during HTTP request processing:

```go
func (s *HTTPTransportServer) handleMCPRequestInternal(w http.ResponseWriter, r *http.Request) error {
    // JSON 解析时发生 panic
    panic("invalid JSON structure")
}
```

### Recovery 保护 / Recovery Protection

```
✅ HTTP 连接不会断开
✅ 返回标准的 MCP 错误响应
✅ 其他并发请求不受影响
✅ 完整的错误堆栈被记录
```

## 场景 3: 多层防护演示 / Scenario 3: Multi-Layer Protection Demo

### 防护层级 / Protection Levels

```
┌─────────────────────────────────────┐
│   HTTP/SSE Transport Layer          │  ← Level 3: 传输层保护
│   (handleMCPRequest)                │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Tool Registry Layer                │  ← Level 2: 注册表层保护
│   (CallTool)                         │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Individual Tool Handler            │  ← Level 1: 工具层保护
│   (handleReadFile, etc.)             │
└─────────────────────────────────────┘
```

### 测试代码示例 / Test Code Example

```go
package main

import (
    "context"
    "fmt"
    "mcp_server/pkg/utils/recovery"
    "go.uber.org/zap"
)

func main() {
    // 创建 logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // 创建 recovery handler
    recoveryHandler := recovery.NewRecoveryHandler(logger)

    // 测试 1: 正常执行
    fmt.Println("Test 1: Normal execution")
    err := recoveryHandler.Recover(func() error {
        fmt.Println("  Executing normal operation...")
        return nil
    })
    fmt.Printf("  Result: %v\n\n", err)

    // 测试 2: Panic recovery
    fmt.Println("Test 2: Panic recovery")
    err = recoveryHandler.Recover(func() error {
        fmt.Println("  About to panic...")
        panic("something went wrong!")
    })
    fmt.Printf("  Result: %v\n\n", err)

    // 测试 3: 带返回值的 recovery
    fmt.Println("Test 3: Recovery with value")
    result, err := recoveryHandler.RecoverWithValue(func() (interface{}, error) {
        return "success", nil
    })
    fmt.Printf("  Result: %v, Error: %v\n\n", result, err)

    // 测试 4: 安全的 goroutine
    fmt.Println("Test 4: Safe goroutine")
    recoveryHandler.SafeGo(func() {
        fmt.Println("  Running in goroutine...")
        panic("goroutine panic - but service continues!")
    })

    // 等待一下让 goroutine 执行
    time.Sleep(100 * time.Millisecond)
    fmt.Println("  Main program continues running!")
}
```

### 预期输出 / Expected Output

```
Test 1: Normal execution
  Executing normal operation...
  Result: <nil>

Test 2: Panic recovery
  About to panic...
  Result: panic: something went wrong!

Test 3: Recovery with value
  Result: success, Error: <nil>

Test 4: Safe goroutine
  Running in goroutine...
  Main program continues running!
```

## 实际应用场景 / Real-World Use Cases

### 1. 文件读取异常 / File Read Exception

```go
// 即使文件系统出现意外错误，服务也能继续运行
result, err := registry.CallTool(ctx, "read_file", map[string]interface{}{
    "path": "/some/file.txt",
})
// 如果内部 panic，会被转换为 error 返回
```

### 2. JSON 解析异常 / JSON Parse Exception

```go
// 即使 JSON 解析出现问题，HTTP 连接也不会断开
// 客户端会收到标准的错误响应
```

### 3. 并发请求保护 / Concurrent Request Protection

```go
// 100 个并发请求中，即使某个请求 panic
// 其他 99 个请求仍然正常处理
for i := 0; i < 100; i++ {
    go func(id int) {
        result, err := client.CallTool(ctx, "some_tool", args)
        // ...
    }(i)
}
```

## 监控和告警 / Monitoring and Alerting

### 日志监控 / Log Monitoring

搜索关键字 "panic recovered" 来监控 panic 发生频率：

Search for "panic recovered" to monitor panic frequency:

```bash
# 查看最近的 panic 日志
grep "panic recovered" /var/log/mcp_server.log

# 统计 panic 次数
grep -c "panic recovered" /var/log/mcp_server.log
```

### 告警规则建议 / Alerting Rules Suggestion

```yaml
# Prometheus 告警规则示例
- alert: HighPanicRate
  expr: rate(mcp_panic_total[5m]) > 0.1
  for: 5m
  annotations:
    summary: "MCP Server panic rate is high"
    description: "Panic rate is {{ $value }} per second"
```

## 总结 / Summary

Recovery 功能提供了：
1. ✅ **多层防护**: 工具层、注册表层、传输层
2. ✅ **优雅降级**: Panic 转换为错误响应
3. ✅ **完整日志**: 堆栈跟踪帮助调试
4. ✅ **服务稳定**: 单点故障不影响整体
5. ✅ **零配置**: 自动启用，无需额外配置

The Recovery feature provides:
1. ✅ **Multi-layer protection**: Tool, registry, and transport layers
2. ✅ **Graceful degradation**: Panic converted to error response
3. ✅ **Complete logging**: Stack traces help debugging
4. ✅ **Service stability**: Single point failures don't affect overall service
5. ✅ **Zero configuration**: Automatically enabled, no extra configuration needed

