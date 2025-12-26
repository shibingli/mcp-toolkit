package json

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

// testStruct 测试用结构体
type testStruct struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Active  bool     `json:"active"`
	Tags    []string `json:"tags,omitempty"`
	Score   float64  `json:"score"`
	private string   //nolint:unused // 私有字段，不应被序列化
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{"nil value", nil, false},
		{"string", "hello", false},
		{"int", 42, false},
		{"float", 3.14, false},
		{"bool", true, false},
		{"slice", []int{1, 2, 3}, false},
		{"map", map[string]int{"a": 1, "b": 2}, false},
		{"struct", testStruct{Name: "test", Age: 25}, false},
		{"channel (unsupported)", make(chan int), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarshalIndent(t *testing.T) {
	input := testStruct{Name: "test", Age: 25, Active: true}
	data, err := MarshalIndent(input, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent() error = %v", err)
	}
	if len(data) == 0 || data[0] != '{' {
		t.Error("MarshalIndent() returned invalid data")
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  any
		wantErr bool
	}{
		{"string", `"hello"`, new(string), false},
		{"int", `42`, new(int), false},
		{"bool", `true`, new(bool), false},
		{"invalid json", `{invalid}`, new(map[string]any), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnmarshalStruct(t *testing.T) {
	input := `{"name":"test","age":25,"active":true,"tags":["go","json"],"score":95.5}`
	var result testStruct
	if err := Unmarshal([]byte(input), &result); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if result.Name != "test" || result.Age != 25 || !result.Active {
		t.Errorf("Unmarshal() got unexpected result: %+v", result)
	}
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	original := testStruct{Name: "roundtrip", Age: 30, Active: true, Tags: []string{"a", "b"}, Score: 88.8}
	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	var result testStruct
	if err = Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if !reflect.DeepEqual(original, result) {
		t.Errorf("Round trip failed: got %+v, want %+v", result, original)
	}
}

func TestMarshalToString(t *testing.T) {
	s, err := MarshalToString(map[string]int{"a": 1})
	if err != nil {
		t.Fatalf("MarshalToString() error = %v", err)
	}
	if s != `{"a":1}` {
		t.Errorf("MarshalToString() = %v, want %v", s, `{"a":1}`)
	}
}

func TestUnmarshalFromString(t *testing.T) {
	var result map[string]int
	if err := UnmarshalFromString(`{"a":1}`, &result); err != nil {
		t.Fatalf("UnmarshalFromString() error = %v", err)
	}
	if result["a"] != 1 {
		t.Errorf("result[\"a\"] = %v, want 1", result["a"])
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{`{"a":1}`, true},
		{`[1,2,3]`, true},
		{`"hello"`, true},
		{`42`, true},
		{`true`, true},
		{`null`, true},
		{`{"a":1`, false},
		{``, false},
	}
	for _, tt := range tests {
		if got := Valid([]byte(tt.input)); got != tt.want {
			t.Errorf("Valid(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestName(t *testing.T) {
	n := Name()
	if n == "" {
		t.Error("Name() returned empty string")
	}
	t.Logf("Using JSON library: %s", n)
}

func TestMarshalSpecialChars(t *testing.T) {
	tests := []string{"你好世界", "🎉🚀", "line1\nline2", `"quoted"`, `path\to\file`}
	for _, input := range tests {
		data, err := Marshal(input)
		if err != nil {
			t.Fatalf("Marshal(%q) error = %v", input, err)
		}
		var result string
		if err = Unmarshal(data, &result); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if result != input {
			t.Errorf("Round trip failed: got %q, want %q", result, input)
		}
	}
}

func TestRawMessage(t *testing.T) {
	// 测试 RawMessage 类型
	type wrapper struct {
		Type string     `json:"type"`
		Data RawMessage `json:"data"`
	}

	// 测试序列化
	raw := RawMessage(`{"nested":"value"}`)
	w := wrapper{Type: "test", Data: raw}
	data, err := Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// 验证 RawMessage 被正确嵌入
	if !bytes.Contains(data, []byte(`"nested":"value"`)) {
		t.Errorf("RawMessage not properly embedded: %s", data)
	}

	// 测试反序列化
	input := `{"type":"example","data":{"key":"value","num":123}}`
	var result wrapper
	if err = Unmarshal([]byte(input), &result); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if result.Type != "example" {
		t.Errorf("Type = %v, want example", result.Type)
	}
	if !bytes.Contains(result.Data, []byte(`"key":"value"`)) {
		t.Errorf("RawMessage data incorrect: %s", result.Data)
	}
}

func TestNewDecoder(t *testing.T) {
	input := `{"name":"decoder_test","age":42}`
	reader := strings.NewReader(input)
	dec := NewDecoder(reader)

	var result testStruct
	if err := dec.Decode(&result); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if result.Name != "decoder_test" || result.Age != 42 {
		t.Errorf("Decode() got unexpected result: %+v", result)
	}
}

func TestNewEncoder(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	input := testStruct{Name: "encoder_test", Age: 30, Active: true}
	if err := enc.Encode(input); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"name":"encoder_test"`) {
		t.Errorf("Encode() output missing name: %s", output)
	}
	if !strings.Contains(output, `"age":30`) {
		t.Errorf("Encode() output missing age: %s", output)
	}
}

func TestDecoderMultipleObjects(t *testing.T) {
	// 测试解码多个 JSON 对象
	input := `{"name":"first","age":1}
{"name":"second","age":2}
{"name":"third","age":3}`
	reader := strings.NewReader(input)
	dec := NewDecoder(reader)

	expected := []testStruct{
		{Name: "first", Age: 1},
		{Name: "second", Age: 2},
		{Name: "third", Age: 3},
	}

	for i, exp := range expected {
		var result testStruct
		if err := dec.Decode(&result); err != nil {
			t.Fatalf("Decode() object %d error = %v", i, err)
		}
		if result.Name != exp.Name || result.Age != exp.Age {
			t.Errorf("Decode() object %d: got %+v, want %+v", i, result, exp)
		}
	}
}

func TestPretouch(t *testing.T) {
	// 测试预热单个类型
	type SimpleStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	err := Pretouch(reflect.TypeOf(SimpleStruct{}))
	if err != nil {
		t.Errorf("Pretouch() error = %v", err)
	}
}

func TestPretouchWithDepth(t *testing.T) {
	// 测试带深度参数的预热
	type NestedStruct struct {
		Name   string `json:"name"`
		Nested *struct {
			Value int `json:"value"`
		} `json:"nested"`
	}

	err := PretouchWithDepth(reflect.TypeOf(NestedStruct{}), 3, 2)
	if err != nil {
		t.Errorf("PretouchWithDepth() error = %v", err)
	}
}

func TestPretouchAll(t *testing.T) {
	// 测试预热所有已注册类型
	failedCount, err := PretouchAll()
	if err != nil {
		t.Logf("PretouchAll() first error = %v (may be expected for non-sonic builds)", err)
	}
	if failedCount > 0 {
		t.Logf("PretouchAll() failed count = %d (may be expected for non-sonic builds)", failedCount)
	}
}

func TestPretouchTypes(t *testing.T) {
	// 测试预热指定类型列表
	type TypeA struct {
		A string `json:"a"`
	}
	type TypeB struct {
		B int `json:"b"`
	}

	err := PretouchTypes(TypeA{}, TypeB{})
	if err != nil {
		t.Errorf("PretouchTypes() error = %v", err)
	}
}

func TestRegisterPretouchType(t *testing.T) {
	// 测试注册自定义预热类型
	type CustomType struct {
		Custom string `json:"custom"`
	}

	// 记录注册前的类型数量
	initialCount := len(pretouchTypes)

	// 注册新类型
	RegisterPretouchType(CustomType{})

	// 验证类型已注册
	if len(pretouchTypes) != initialCount+1 {
		t.Errorf("RegisterPretouchType() did not add type, count: %d, expected: %d",
			len(pretouchTypes), initialCount+1)
	}
}
