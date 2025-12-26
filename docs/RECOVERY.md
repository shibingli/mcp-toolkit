# Recovery 功能文档 / Recovery Feature Documentation

## 概述 / Overview

本项目实现了完整的 panic recovery 机制，防止 MCP 工具异常导致整个服务崩溃。Recovery 功能在多个层面提供保护，确保服务的稳定性和可靠性。

This project implements a complete panic recovery mechanism to prevent MCP tool exceptions from causing the entire service to crash. The recovery feature provides protection at multiple levels to ensure service stability and reliability.

## 架构设计 / Architecture Design

### 1. Recovery 中间件 / Recovery Middleware

位置：`pkg/utils/recovery/recovery.go`

核心组件：
- `RecoveryHandler`: 通用的 panic 恢复处理器
- 自动捕获 panic 并转换为 error
- 记录完整的堆栈跟踪信息
- 支持多种使用场景

Location: `pkg/utils/recovery/recovery.go`

Core components:
- `RecoveryHandler`: Generic panic recovery handler
- Automatically catches panic and converts to error
- Logs complete stack trace information
- Supports multiple use cases

### 2. 保护层级 / Protection Levels

#### Level 1: 工具注册表层 / Tool Registry Level

**位置 / Location**: `pkg/transport/tool_registry.go`

**功能 / Features**:
- 在 `CallTool` 方法中添加 panic recovery
- 确保单个工具的异常不会影响其他工具
- 自动记录异常信息和堆栈

**实现 / Implementation**:
```go
func (r *ToolRegistry) CallTool(ctx context.Context, name string, arguments interface{}) (*mcp.CallToolResult, error) {
    // 使用 recovery handler 包装调用
    result, err := recoveryHandler.RecoverWithValue(func() (interface{}, error) {
        return handler(ctx, arguments)
    })
    // ...
}
```

#### Level 2: HTTP 传输层 / HTTP Transport Level

**位置 / Location**: `pkg/transport/http.go`

**功能 / Features**:
- 在 `handleMCPRequest` 方法中添加 panic recovery
- 捕获请求处理过程中的所有异常
- 返回标准的 MCP 错误响应

**实现 / Implementation**:
```go
func (s *HTTPTransportServer) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
    recoveryHandler := recovery.NewRecoveryHandler(s.logger)
    err := recoveryHandler.Recover(func() error {
        return s.handleMCPRequestInternal(w, r)
    })
    // ...
}
```

#### Level 3: SSE 传输层 / SSE Transport Level

**位置 / Location**: `pkg/transport/sse.go`

**功能 / Features**:
- 在 `handleMCPMessage` 方法中添加 panic recovery
- 保护 SSE 连接不因异常而断开
- 确保实时事件流的稳定性

**实现 / Implementation**:
```go
func (s *SSETransportServer) handleMCPMessage(w http.ResponseWriter, r *http.Request) {
    recoveryHandler := recovery.NewRecoveryHandler(s.logger)
    err := recoveryHandler.Recover(func() error {
        return s.handleMCPMessageInternal(w, r)
    })
    // ...
}
```

## API 使用指南 / API Usage Guide

### 1. 基本用法 / Basic Usage

```go
import "mcp_server/pkg/utils/recovery"

// 创建 recovery handler
recoveryHandler := recovery.NewRecoveryHandler(logger)

// 包装可能 panic 的函数
err := recoveryHandler.Recover(func() error {
    // 可能会 panic 的代码
    return doSomething()
})

if err != nil {
    // 处理错误（包括从 panic 转换的错误）
    log.Error("operation failed", zap.Error(err))
}
```

### 2. 带返回值的用法 / Usage with Return Value

```go
result, err := recoveryHandler.RecoverWithValue(func() (interface{}, error) {
    // 可能会 panic 的代码
    return computeSomething()
})

if err != nil {
    // 处理错误
    return nil, err
}

// 使用结果
return result, nil
```

### 3. 安全启动 Goroutine / Safe Goroutine Launch

```go
// 启动一个受保护的 goroutine
recoveryHandler.SafeGo(func() {
    // 后台任务代码
    // 即使 panic 也不会导致程序崩溃
    processBackgroundTask()
})
```

### 4. 包装处理器函数 / Wrap Handler Functions

```go
// 包装普通处理器
wrappedHandler := recoveryHandler.WrapHandler(originalHandler)

// 包装带返回值的处理器
wrappedHandlerWithValue := recoveryHandler.WrapHandlerWithValue(originalHandlerWithValue)
```

## 错误处理 / Error Handling

### Panic 类型转换 / Panic Type Conversion

Recovery 机制会自动将不同类型的 panic 转换为 error：

The recovery mechanism automatically converts different types of panic to error:

| Panic 类型 / Type | 转换结果 / Conversion Result |
|------------------|----------------------------|
| `error` | `fmt.Errorf("panic: %w", err)` |
| `string` | `fmt.Errorf("panic: %s", str)` |
| 其他类型 / Other | `fmt.Errorf("panic: %v", value)` |

### 日志记录 / Logging

所有捕获的 panic 都会记录以下信息：
- Panic 值
- 完整的堆栈跟踪
- 错误级别日志

All caught panics log the following information:
- Panic value
- Complete stack trace
- Error level log

## 测试覆盖 / Test Coverage

### 单元测试 / Unit Tests

**位置 / Location**: `pkg/utils/recovery/recovery_test.go`

测试场景包括 / Test scenarios include:
- ✅ 正常执行（无 panic）
- ✅ 返回错误（非 panic）
- ✅ String 类型 panic
- ✅ Error 类型 panic
- ✅ 整数类型 panic
- ✅ Nil panic
- ✅ 结构体 panic
- ✅ Goroutine panic recovery
- ✅ 堆栈跟踪记录

### 集成测试 / Integration Tests

**位置 / Location**: `pkg/transport/tool_registry_recovery_test.go`

测试场景包括 / Test scenarios include:
- ✅ 工具正常调用
- ✅ 工具返回错误
- ✅ 工具 panic recovery
- ✅ 工具不存在错误
- ✅ 无 recovery handler 的降级处理

## 性能影响 / Performance Impact

Recovery 机制的性能开销极小：
- 正常执行路径：几乎无开销（仅函数调用）
- Panic 路径：需要记录堆栈，但这是异常情况

The performance overhead of the recovery mechanism is minimal:
- Normal execution path: Almost no overhead (just function calls)
- Panic path: Stack trace logging required, but this is an exceptional case

## 最佳实践 / Best Practices

1. **始终使用 Recovery**: 在所有可能 panic 的边界使用 recovery
2. **记录详细日志**: 确保 panic 信息被完整记录
3. **优雅降级**: Panic 后返回合适的错误响应
4. **测试覆盖**: 为 panic 场景编写测试用例
5. **监控告警**: 监控 panic 频率，及时发现问题

1. **Always Use Recovery**: Use recovery at all boundaries where panic might occur
2. **Log Detailed Information**: Ensure panic information is fully logged
3. **Graceful Degradation**: Return appropriate error responses after panic
4. **Test Coverage**: Write test cases for panic scenarios
5. **Monitoring and Alerting**: Monitor panic frequency to detect issues early

## 故障排查 / Troubleshooting

### 如何查看 Panic 日志 / How to View Panic Logs

Panic 日志会以 ERROR 级别记录，包含：
```
panic recovered
  panic: <panic value>
  stack: <stack trace>
```

### 常见问题 / Common Issues

**Q: Recovery 会隐藏所有错误吗？**
A: 不会。Recovery 只捕获 panic，正常的 error 返回不受影响。

**Q: 性能会受影响吗？**
A: 几乎没有影响。只在实际发生 panic 时才有额外开销。

**Q: 如何测试 panic recovery？**
A: 参考 `recovery_test.go` 中的测试用例。

## 版本历史 / Version History

- **v1.0.1**: 初始实现 Recovery 功能
  - 添加通用 recovery 中间件
  - 在工具注册表层添加保护
  - 在 HTTP/SSE 传输层添加保护
  - 完整的单元测试和集成测试

