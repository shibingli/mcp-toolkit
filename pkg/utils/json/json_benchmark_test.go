// Package json JSON 库性能基准测试
package json

import (
	"bytes"
	"testing"
)

// ComplexData 复杂测试数据结构
type ComplexData struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
	Nested      *NestedData            `json:"nested"`
	Items       []ItemData             `json:"items"`
	Config      ConfigData             `json:"config"`
	Stats       StatsData              `json:"stats"`
	Enabled     bool                   `json:"enabled"`
	Priority    int                    `json:"priority"`
	Score       float64                `json:"score"`
}

// NestedData 嵌套数据
type NestedData struct {
	Level1 *Level1Data `json:"level1"`
}

// Level1Data 第一层嵌套
type Level1Data struct {
	Level2 *Level2Data `json:"level2"`
	Values []int       `json:"values"`
}

// Level2Data 第二层嵌套
type Level2Data struct {
	Level3 *Level3Data `json:"level3"`
	Data   string      `json:"data"`
}

// Level3Data 第三层嵌套
type Level3Data struct {
	DeepValue   string            `json:"deep_value"`
	DeepNumbers []float64         `json:"deep_numbers"`
	DeepMap     map[string]string `json:"deep_map"`
}

// ItemData 列表项数据
type ItemData struct {
	ItemID     string            `json:"item_id"`
	ItemName   string            `json:"item_name"`
	Quantity   int               `json:"quantity"`
	Price      float64           `json:"price"`
	Available  bool              `json:"available"`
	Categories []string          `json:"categories"`
	Attributes map[string]string `json:"attributes"`
}

// ConfigData 配置数据
type ConfigData struct {
	MaxConnections int             `json:"max_connections"`
	Timeout        int             `json:"timeout"`
	RetryCount     int             `json:"retry_count"`
	Endpoints      []string        `json:"endpoints"`
	Features       map[string]bool `json:"features"`
	Limits         map[string]int  `json:"limits"`
}

// StatsData 统计数据
type StatsData struct {
	TotalRequests  int64   `json:"total_requests"`
	SuccessCount   int64   `json:"success_count"`
	FailureCount   int64   `json:"failure_count"`
	AverageLatency float64 `json:"average_latency"`
	P99Latency     float64 `json:"p99_latency"`
	Timestamps     []int64 `json:"timestamps"`
	ErrorCodes     []int   `json:"error_codes"`
}

// 预生成测试数据和JSON字节
var (
	benchData     = generateComplexData()
	benchDataJSON []byte
)

func init() {
	var err error
	benchDataJSON, err = Marshal(benchData)
	if err != nil {
		panic("failed to marshal benchmark data: " + err.Error())
	}
}

// generateComplexData 生成复杂测试数据
func generateComplexData() *ComplexData {
	return &ComplexData{
		ID:          "test-id-12345-abcde-67890",
		Name:        "复杂测试数据结构 Complex Test Data",
		Description: "这是一个用于性能测试的复杂数据结构，包含多层嵌套、数组、Map等多种类型",
		Tags:        []string{"performance", "benchmark", "json", "测试", "性能", "基准"},
		Metadata: map[string]interface{}{
			"version": "1.0.0", "author": "test", "created_at": 1702656000,
			"updated_at": 1702742400, "is_active": true, "retry_count": 3, "threshold": 0.95,
		},
		Nested: &NestedData{
			Level1: &Level1Data{
				Level2: &Level2Data{
					Level3: &Level3Data{
						DeepValue:   "深层嵌套值 Deep nested value",
						DeepNumbers: []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9},
						DeepMap:     map[string]string{"key1": "value1", "key2": "value2", "key3": "value3", "键4": "值4"},
					},
					Data: "Level2 data content 第二层数据内容",
				},
				Values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 1000, 10000},
			},
		},
		Items: generateItems(20),
		Config: ConfigData{
			MaxConnections: 1000, Timeout: 30000, RetryCount: 5,
			Endpoints: []string{"http://api1.example.com", "http://api2.example.com", "http://api3.example.com"},
			Features:  map[string]bool{"feature_a": true, "feature_b": false, "feature_c": true},
			Limits:    map[string]int{"rate_limit": 1000, "burst_limit": 100, "max_size": 10485760},
		},
		Stats: StatsData{
			TotalRequests: 1234567890, SuccessCount: 1234000000, FailureCount: 567890,
			AverageLatency: 12.345, P99Latency: 99.999,
			Timestamps: []int64{1702656000, 1702656060, 1702656120, 1702656180, 1702656240},
			ErrorCodes: []int{400, 401, 403, 404, 500, 502, 503, 504},
		},
		Enabled: true, Priority: 100, Score: 98.765,
	}
}

// generateItems 生成测试项目列表
func generateItems(count int) []ItemData {
	items := make([]ItemData, count)
	for i := 0; i < count; i++ {
		items[i] = ItemData{
			ItemID:     "item-" + string(rune('A'+i%26)) + "-" + string(rune('0'+i%10)),
			ItemName:   "商品名称 Product Name " + string(rune('A'+i%26)),
			Quantity:   (i + 1) * 10,
			Price:      float64(i+1) * 9.99,
			Available:  i%2 == 0,
			Categories: []string{"category1", "category2", "分类3"},
			Attributes: map[string]string{"color": "red", "size": "large", "颜色": "红色"},
		}
	}
	return items
}

// BenchmarkMarshal 测试序列化性能
func BenchmarkMarshal(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(benchData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMarshalIndent 测试带缩进序列化性能
func BenchmarkMarshalIndent(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MarshalIndent(benchData, "", "  ")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUnmarshal 测试反序列化性能
func BenchmarkUnmarshal(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var data ComplexData
		err := Unmarshal(benchDataJSON, &data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMarshalToString 测试序列化为字符串性能
func BenchmarkMarshalToString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MarshalToString(benchData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUnmarshalFromString 测试从字符串反序列化性能
func BenchmarkUnmarshalFromString(b *testing.B) {
	jsonStr := string(benchDataJSON)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var data ComplexData
		err := UnmarshalFromString(jsonStr, &data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkValid 测试JSON验证性能
func BenchmarkValid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Valid(benchDataJSON)
	}
}

// BenchmarkEncoder 测试流式编码性能
func BenchmarkEncoder(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		err := enc.Encode(benchData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDecoder 测试流式解码性能
func BenchmarkDecoder(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(benchDataJSON)
		dec := NewDecoder(reader)
		var data ComplexData
		err := dec.Decode(&data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMarshalLargeArray 测试大数组序列化性能
func BenchmarkMarshalLargeArray(b *testing.B) {
	items := generateItems(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(items)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUnmarshalLargeArray 测试大数组反序列化性能
func BenchmarkUnmarshalLargeArray(b *testing.B) {
	items := generateItems(100)
	jsonData, _ := Marshal(items)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []ItemData
		err := Unmarshal(jsonData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMarshalMap 测试Map序列化性能
func BenchmarkMarshalMap(b *testing.B) {
	data := map[string]interface{}{
		"string": "test value 测试值", "number": 12345.6789, "boolean": true, "null": nil,
		"array":  []interface{}{1, 2, 3, "a", "b", "c"},
		"nested": map[string]interface{}{"key1": "value1", "key2": 123, "key3": []int{1, 2, 3, 4, 5}},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUnmarshalMap 测试Map反序列化性能
func BenchmarkUnmarshalMap(b *testing.B) {
	jsonData := []byte(`{"string":"test value 测试值","number":12345.6789,"boolean":true,"null":null,"array":[1,2,3,"a","b","c"],"nested":{"key1":"value1","key2":123,"key3":[1,2,3,4,5]}}`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		err := Unmarshal(jsonData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParallel 测试并行序列化性能
func BenchmarkParallel(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = Marshal(benchData)
		}
	})
}

// BenchmarkParallelUnmarshal 测试并行反序列化性能
func BenchmarkParallelUnmarshal(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var data ComplexData
			_ = Unmarshal(benchDataJSON, &data)
		}
	})
}
