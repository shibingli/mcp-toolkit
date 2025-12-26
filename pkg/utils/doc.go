// Package utils 提供了 MCP Toolkit 的工具函数和辅助功能
//
// 本包包含多个子包，提供各种通用功能。
//
// # 子包
//
// json 包（pkg/utils/json）：
//   - 高性能 JSON 序列化/反序列化
//   - 支持多种 JSON 库（Sonic、go-json、jsoniter、标准库）
//   - 编译时选择 JSON 库
//   - JSON 结构体预热机制（Sonic）
//
// recovery 包（pkg/utils/recovery）：
//   - Panic recovery 机制
//   - 自动捕获和恢复 panic
//   - 记录完整堆栈信息
//   - 转换 panic 为 error
//
// # JSON 包
//
// JSON 包提供了统一的 JSON 操作接口，支持多种 JSON 库：
//
// 使用 Sonic（推荐，最高性能）：
//
//	go build -tags="sonic"
//
// 使用 go-json：
//
//	go build -tags="gojson"
//
// 使用 jsoniter：
//
//	go build -tags="jsoniter"
//
// 使用标准库（默认）：
//
//	go build
//
// 使用示例：
//
//	import "mcp-toolkit/pkg/utils/json"
//
//	// 序列化
//	data, err := json.Marshal(obj)
//
//	// 反序列化
//	err := json.Unmarshal(data, &obj)
//
//	// 字符串序列化
//	str, err := json.MarshalToString(obj)
//
//	// 字符串反序列化
//	err := json.UnmarshalFromString(str, &obj)
//
//	// 获取当前使用的 JSON 库名称
//	name := json.Name()
//
// JSON 预热（仅 Sonic 有效）：
//
//	// 预热所有结构体
//	failedCount, err := json.PretouchAll()
//
//	// 预热单个结构体
//	err := json.Pretouch(reflect.TypeOf(MyStruct{}))
//
// # Recovery 包
//
// Recovery 包提供了 panic 恢复机制，保护系统稳定性：
//
// 使用示例：
//
//	import "mcp-toolkit/pkg/utils/recovery"
//
//	// 创建 recovery handler
//	handler := recovery.NewRecoveryHandler(logger)
//
//	// 恢复 panic
//	err := handler.Recover(func() error {
//	    // 可能 panic 的代码
//	    return doSomething()
//	})
//
//	// 恢复 panic 并返回值
//	result, err := handler.RecoverWithValue(func() (interface{}, error) {
//	    // 可能 panic 的代码
//	    return doSomethingWithResult()
//	})
//
// Recovery 特性：
//
//   - 自动捕获 panic
//   - 记录完整堆栈信息
//   - 转换为 error 返回
//   - 支持返回值的函数
//   - 集成 zap 日志
//
// # 性能优化
//
// JSON 性能优化：
//
//   - Sonic：使用 JIT 编译，性能最高
//   - go-json：纯 Go 实现，性能优秀
//   - jsoniter：兼容标准库，性能良好
//   - 结构体预热：消除首次序列化延迟
//
// Recovery 性能优化：
//
//   - 最小化性能开销
//   - 仅在 panic 时记录日志
//   - 复用 error 对象
//
// # 测试
//
// JSON 包测试：
//
//	go test -v ./pkg/utils/json/...
//
// Recovery 包测试：
//
//	go test -v ./pkg/utils/recovery/...
//
// # 基准测试
//
// JSON 包包含完整的基准测试：
//
//	go test -bench=. ./pkg/utils/json/...
//
// 比较不同 JSON 库的性能：
//
//	# Sonic
//	go test -tags="sonic" -bench=. ./pkg/utils/json/...
//
//	# go-json
//	go test -tags="gojson" -bench=. ./pkg/utils/json/...
//
//	# jsoniter
//	go test -tags="jsoniter" -bench=. ./pkg/utils/json/...
//
//	# 标准库
//	go test -bench=. ./pkg/utils/json/...
//
// # 设计原则
//
//   - 统一接口：所有 JSON 库使用相同接口
//   - 编译时选择：通过 build tags 选择实现
//   - 零依赖切换：切换 JSON 库不需要修改代码
//   - 性能优先：优先使用高性能实现
//   - 安全保护：Recovery 机制保护系统稳定性
package utils
