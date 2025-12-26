//go:build sonic && (amd64 || arm64) && (linux || darwin || windows)

// Package json sonic 高性能 JSON 库实现
// sonic 是字节跳动开发的高性能 JSON 库，性能最佳
// 仅支持 amd64/arm64 架构和 Linux/macOS/Windows 系统
package json

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/option"
)

// jsonLibName 当前 JSON 库名称
const jsonLibName = "bytedance/sonic"

// rawMessage 使用标准库的 RawMessage（sonic 兼容）
type rawMessage = json.RawMessage

// decoder 包装 sonic.Decoder 接口，使其行为与标准库一致
// sonic.Decoder 是接口类型，需要包装成结构体以便返回指针
type decoder struct {
	dec sonic.Decoder
}

// Decode 从输入流解码 JSON 值到 v
func (d *decoder) Decode(v any) error {
	return d.dec.Decode(v)
}

// Buffered 返回解码器缓冲区中剩余的数据
func (d *decoder) Buffered() io.Reader {
	return d.dec.Buffered()
}

// DisallowUnknownFields 禁止未知字段
func (d *decoder) DisallowUnknownFields() {
	d.dec.DisallowUnknownFields()
}

// More 报告当前数组或对象中是否还有更多元素
func (d *decoder) More() bool {
	return d.dec.More()
}

// UseNumber 使用 Number 类型而不是 float64 来解码数字
func (d *decoder) UseNumber() {
	d.dec.UseNumber()
}

// encoder 包装 sonic.Encoder 接口，使其行为与标准库一致
// sonic.Encoder 是接口类型，需要包装成结构体以便返回指针
type encoder struct {
	enc sonic.Encoder
}

// Encode 将 v 编码为 JSON 并写入输出流
func (e *encoder) Encode(v any) error {
	return e.enc.Encode(v)
}

// SetEscapeHTML 设置是否转义 HTML 字符
func (e *encoder) SetEscapeHTML(on bool) {
	e.enc.SetEscapeHTML(on)
}

// SetIndent 设置编码器的缩进
func (e *encoder) SetIndent(prefix, indent string) {
	e.enc.SetIndent(prefix, indent)
}

// marshal 使用 sonic 序列化
func marshal(v any) ([]byte, error) {
	return sonic.Marshal(v)
}

// marshalIndent 使用 sonic 序列化（带缩进）
func marshalIndent(v any, prefix, indent string) ([]byte, error) {
	return sonic.MarshalIndent(v, prefix, indent)
}

// unmarshal 使用 sonic 反序列化
func unmarshal(data []byte, v any) error {
	return sonic.Unmarshal(data, v)
}

// valid 使用 sonic 验证 JSON
func valid(data []byte) bool {
	return sonic.Valid(data)
}

// name 返回库名称
func name() string {
	return jsonLibName
}

// newDecoder 创建 sonic 解码器
// 返回包装后的解码器指针，与标准库 API 保持一致
func newDecoder(r io.Reader) *decoder {
	return &decoder{dec: sonic.ConfigDefault.NewDecoder(r)}
}

// newEncoder 创建 sonic 编码器
// 返回包装后的编码器指针，与标准库 API 保持一致
func newEncoder(w io.Writer) *encoder {
	return &encoder{enc: sonic.ConfigDefault.NewEncoder(w)}
}

// pretouch 预热指定类型，提前编译 JIT 代码
// 这可以避免首次序列化/反序列化大型结构体时的延迟
func pretouch(t reflect.Type, _ ...any) error {
	return sonic.Pretouch(t)
}

// pretouchWithOptions 使用自定义选项预热指定类型
func pretouchWithOptions(t reflect.Type, depth int, recursiveDepth int) error {
	return sonic.Pretouch(t,
		option.WithCompileMaxInlineDepth(depth),
		option.WithCompileRecursiveDepth(recursiveDepth),
	)
}
