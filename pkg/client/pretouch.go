package client

import (
	"reflect"

	"mcp-toolkit/pkg/utils/json"
)

// pretouchTypes 需要预热的客户端类型列表 / Client types that need to be preheated
var pretouchTypes = []any{
	InitializeResponse{},
	ServerCapabilities{},
	ToolsCapability{},
	ResourcesCapability{},
	PromptsCapability{},
	ServerInfo{},
	ListToolsResponse{},
	Tool{},
	CallToolResponse{},
	Content{},
}

// PretouchAll 预热所有客户端类型 / Preheat all client types
// 建议在使用客户端前调用此函数 / Recommended to call before using client
func PretouchAll() (failedCount int, firstErr error) {
	for _, v := range pretouchTypes {
		t := reflect.TypeOf(v)
		if err := json.Pretouch(t); err != nil {
			failedCount++
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return failedCount, firstErr
}
