//go:build !sonic && !gojson && !jsoniter

// Package json 标准库实现
package json

import (
	"encoding/json"
	"io"
	"reflect"
)

// jsonLibName 当前 JSON 库名称
const jsonLibName = "encoding/json"

// rawMessage 是 json.RawMessage 的别名
type rawMessage = json.RawMessage

// decoder 是 json.Decoder 的别名
type decoder = json.Decoder

// encoder 是 json.Encoder 的别名
type encoder = json.Encoder

// marshal 使用标准库序列化
func marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// marshalIndent 使用标准库序列化（带缩进）
func marshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// unmarshal 使用标准库反序列化
func unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// valid 使用标准库验证 JSON
func valid(data []byte) bool {
	return json.Valid(data)
}

// name 返回库名称
func name() string {
	return jsonLibName
}

// newDecoder 创建标准库解码器
func newDecoder(r io.Reader) *decoder {
	return json.NewDecoder(r)
}

// newEncoder 创建标准库编码器
func newEncoder(w io.Writer) *encoder {
	return json.NewEncoder(w)
}

// pretouch 标准库不需要预热，直接返回 nil
func pretouch(_ reflect.Type, _ ...any) error {
	return nil
}

// pretouchWithOptions 标准库不需要预热，直接返回 nil
func pretouchWithOptions(_ reflect.Type, _ int, _ int) error {
	return nil
}
