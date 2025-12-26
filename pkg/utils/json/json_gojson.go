//go:build gojson

// Package json goccy/go-json 高性能 JSON 库实现
// go-json 是一个高性能的 JSON 库，具有良好的跨平台兼容性
package json

import (
	stdjson "encoding/json"
	"io"
	"reflect"

	gojson "github.com/goccy/go-json"
)

// jsonLibName 当前 JSON 库名称
const jsonLibName = "goccy/go-json"

// rawMessage 使用标准库的 RawMessage，确保与其他包的类型兼容
type rawMessage = stdjson.RawMessage

// decoder 使用 go-json 的解码器
type decoder = gojson.Decoder

// encoder 使用 go-json 的编码器
type encoder = gojson.Encoder

// marshal 使用 go-json 序列化
func marshal(v any) ([]byte, error) {
	return gojson.Marshal(v)
}

// marshalIndent 使用 go-json 序列化（带缩进）
func marshalIndent(v any, prefix, indent string) ([]byte, error) {
	return gojson.MarshalIndent(v, prefix, indent)
}

// unmarshal 使用 go-json 反序列化
func unmarshal(data []byte, v any) error {
	return gojson.Unmarshal(data, v)
}

// valid 使用 go-json 验证 JSON
func valid(data []byte) bool {
	return gojson.Valid(data)
}

// name 返回库名称
func name() string {
	return jsonLibName
}

// newDecoder 创建 go-json 解码器
func newDecoder(r io.Reader) *decoder {
	return gojson.NewDecoder(r)
}

// newEncoder 创建 go-json 编码器
func newEncoder(w io.Writer) *encoder {
	return gojson.NewEncoder(w)
}

// pretouch go-json 不需要预热，直接返回 nil
func pretouch(_ reflect.Type, _ ...any) error {
	return nil
}

// pretouchWithOptions go-json 不需要预热，直接返回 nil
func pretouchWithOptions(_ reflect.Type, _ int, _ int) error {
	return nil
}
