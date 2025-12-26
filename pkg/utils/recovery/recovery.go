// Copyright 2024 MCP Toolkit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package recovery

import (
	"fmt"
	"runtime/debug"

	"go.uber.org/zap"
)

// RecoveryHandler panic恢复处理器 / Panic recovery handler
type RecoveryHandler struct {
	logger *zap.Logger
}

// NewRecoveryHandler 创建恢复处理器 / Create recovery handler
func NewRecoveryHandler(logger *zap.Logger) *RecoveryHandler {
	return &RecoveryHandler{
		logger: logger,
	}
}

// Recover 执行panic恢复并返回错误 / Execute panic recovery and return error
// 用于包装可能发生panic的函数调用
// Usage: if err := handler.Recover(func() error { ... }); err != nil { ... }
func (h *RecoveryHandler) Recover(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// 记录panic信息和堆栈 / Log panic info and stack trace
			stack := debug.Stack()
			h.logger.Error("panic recovered",
				zap.Any("panic", r),
				zap.String("stack", string(stack)))

			// 将panic转换为error / Convert panic to error
			err = h.panicToError(r)
		}
	}()

	// 执行函数 / Execute function
	err = fn()
	return err
}

// RecoverWithValue 执行panic恢复并返回值和错误 / Execute panic recovery and return value and error
// 用于包装返回值的函数
// Usage: result, err := handler.RecoverWithValue(func() (interface{}, error) { ... })
func (h *RecoveryHandler) RecoverWithValue(fn func() (interface{}, error)) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			// 记录panic信息和堆栈 / Log panic info and stack trace
			stack := debug.Stack()
			h.logger.Error("panic recovered",
				zap.Any("panic", r),
				zap.String("stack", string(stack)))

			// 将panic转换为error / Convert panic to error
			result = nil
			err = h.panicToError(r)
		}
	}()

	// 执行函数 / Execute function
	result, err = fn()
	return result, err
}

// panicToError 将panic值转换为error / Convert panic value to error
func (h *RecoveryHandler) panicToError(r interface{}) error {
	switch v := r.(type) {
	case error:
		return fmt.Errorf("panic: %w", v)
	case string:
		return fmt.Errorf("panic: %s", v)
	default:
		return fmt.Errorf("panic: %v", v)
	}
}

// SafeGo 安全地启动goroutine，自动恢复panic / Safely start goroutine with automatic panic recovery
// 用于启动后台任务，防止panic导致程序崩溃
// Usage: handler.SafeGo(func() { ... })
func (h *RecoveryHandler) SafeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// 记录panic信息和堆栈 / Log panic info and stack trace
				stack := debug.Stack()
				h.logger.Error("goroutine panic recovered",
					zap.Any("panic", r),
					zap.String("stack", string(stack)))
			}
		}()

		// 执行函数 / Execute function
		fn()
	}()
}

// WrapHandler 包装处理器函数，添加panic恢复 / Wrap handler function with panic recovery
// 返回一个新的处理器函数，该函数会自动捕获并转换panic
// Usage: wrappedHandler := handler.WrapHandler(originalHandler)
func (h *RecoveryHandler) WrapHandler(handler func() error) func() error {
	return func() error {
		return h.Recover(handler)
	}
}

// WrapHandlerWithValue 包装带返回值的处理器函数 / Wrap handler function with value and panic recovery
// 返回一个新的处理器函数，该函数会自动捕获并转换panic
// Usage: wrappedHandler := handler.WrapHandlerWithValue(originalHandler)
func (h *RecoveryHandler) WrapHandlerWithValue(handler func() (interface{}, error)) func() (interface{}, error) {
	return func() (interface{}, error) {
		return h.RecoverWithValue(handler)
	}
}
