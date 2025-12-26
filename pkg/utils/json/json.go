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

// Package json 提供统一的 JSON 编解码接口
// 通过 Go 编译标识 (build tags) 切换不同的 JSON 库实现：
//   - 默认: 使用标准库 encoding/json
//   - sonic: 使用 bytedance/sonic (最高性能，仅支持 amd64/arm64 + Linux/macOS/Windows)
//   - gojson: 使用 goccy/go-json (高性能，跨平台兼容)
//   - jsoniter: 使用 json-iterator/go (高性能，广泛使用)
//
// 使用示例:
//
//	go build -tags=sonic ./...
//	go build -tags=gojson ./...
//	go build -tags=jsoniter ./...
//	go build ./...  // 使用标准库
package json

import (
	"io"
	"reflect"
)

// RawMessage 是原始编码的 JSON 值
// 它实现了 Marshaler 和 Unmarshaler 接口，可以用于延迟 JSON 解码或预计算 JSON 编码
type RawMessage = rawMessage

// Marshal 将对象序列化为 JSON 字节数组
// 具体实现由编译标识决定
func Marshal(v any) ([]byte, error) {
	return marshal(v)
}

// MarshalIndent 将对象序列化为带缩进的 JSON 字节数组
// prefix 为每行前缀，indent 为缩进字符串
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return marshalIndent(v, prefix, indent)
}

// Unmarshal 将 JSON 字节数组反序列化为对象
func Unmarshal(data []byte, v any) error {
	return unmarshal(data, v)
}

// MarshalToString 将对象序列化为 JSON 字符串
func MarshalToString(v any) (string, error) {
	data, err := marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalFromString 将 JSON 字符串反序列化为对象
func UnmarshalFromString(s string, v any) error {
	return unmarshal([]byte(s), v)
}

// Valid 检查字节数组是否为有效的 JSON
func Valid(data []byte) bool {
	return valid(data)
}

// Name 返回当前使用的 JSON 库名称
func Name() string {
	return name()
}

// NewDecoder 创建一个从 r 读取的新 JSON 解码器
// 返回的解码器支持 Decode、Buffered、More、UseNumber、DisallowUnknownFields 等方法
func NewDecoder(r io.Reader) *decoder {
	return newDecoder(r)
}

// NewEncoder 创建一个写入 w 的新 JSON 编码器
// 返回的编码器支持 Encode、SetEscapeHTML、SetIndent 等方法
func NewEncoder(w io.Writer) *encoder {
	return newEncoder(w)
}

// Pretouch 预热指定类型，提前编译 JIT 代码（仅 sonic 有效）
// 对于 sonic 库，这可以避免首次序列化/反序列化大型结构体时的延迟
// 对于其他 JSON 库，此函数直接返回 nil
//
// 使用示例:
//
//	func init() {
//	    var v YourStruct
//	    _ = json.Pretouch(reflect.TypeOf(v))
//	}
func Pretouch(t reflect.Type) error {
	return pretouch(t)
}

// PretouchWithDepth 使用自定义深度预热指定类型（仅 sonic 有效）
// maxInlineDepth: 最大内联深度，用于控制编译时间
// recursiveDepth: 递归深度，用于处理深层嵌套结构
//
// 对于大型嵌套结构体，可以设置较小的 maxInlineDepth 来减少编译时间
// 对于深层嵌套结构体（嵌套深度 > 默认值），可以设置 recursiveDepth 来提高稳定性
func PretouchWithDepth(t reflect.Type, maxInlineDepth, recursiveDepth int) error {
	return pretouchWithOptions(t, maxInlineDepth, recursiveDepth)
}
