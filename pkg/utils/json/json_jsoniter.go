//go:build jsoniter

// Package json json-iterator 高性能 JSON 库实现
// jsoniter 是滴滴开发的高性能 JSON 库，广泛使用，兼容性好
package json

import (
	stdjson "encoding/json"
	"io"
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

// jsonLibName 当前 JSON 库名称
const jsonLibName = "json-iterator/go"

// jsonAPI 使用与标准库兼容的配置
var jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary

// rawMessage 使用标准库的 RawMessage，确保与其他包的类型兼容
type rawMessage = stdjson.RawMessage

// decoder 使用 jsoniter 的解码器
type decoder = jsoniter.Decoder

// encoder 使用 jsoniter 的编码器
type encoder = jsoniter.Encoder

// marshal 使用 jsoniter 序列化
func marshal(v any) ([]byte, error) {
	return jsonAPI.Marshal(v)
}

// marshalIndent 使用 jsoniter 序列化（带缩进）
func marshalIndent(v any, prefix, indent string) ([]byte, error) {
	return jsonAPI.MarshalIndent(v, prefix, indent)
}

// unmarshal 使用 jsoniter 反序列化
func unmarshal(data []byte, v any) error {
	return jsonAPI.Unmarshal(data, v)
}

// valid 使用标准库验证 JSON
// 注意：jsoniter 的 Valid 函数对于简单的 JSON 原始值（如数字、布尔值）返回 false
// 这与标准库行为不一致，因此这里使用标准库的 Valid 函数
func valid(data []byte) bool {
	return stdjson.Valid(data)
}

// name 返回库名称
func name() string {
	return jsonLibName
}

// newDecoder 创建 jsoniter 解码器
func newDecoder(r io.Reader) *decoder {
	return jsonAPI.NewDecoder(r)
}

// newEncoder 创建 jsoniter 编码器
func newEncoder(w io.Writer) *encoder {
	return jsonAPI.NewEncoder(w)
}

// pretouch jsoniter 不需要预热，直接返回 nil
func pretouch(_ reflect.Type, _ ...any) error {
	return nil
}

// pretouchWithOptions jsoniter 不需要预热，直接返回 nil
func pretouchWithOptions(_ reflect.Type, _ int, _ int) error {
	return nil
}
