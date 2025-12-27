// Package json 提供 JSON 预热功能
package json

import (
	"reflect"

	"mcp-toolkit/pkg/types"
)

// pretouchTypes 需要预热的类型列表
// 这些类型在程序启动时会被预热，以避免首次序列化/反序列化时的延迟
var pretouchTypes = []any{
	// ========== 文件系统相关结构体 / Filesystem related structures ==========

	// 基础类型 / Basic types
	types.FileInfo{},

	// 请求类型 / Request types
	types.CreateFileRequest{},
	types.CreateDirRequest{},
	types.ReadFileRequest{},
	types.WriteFileRequest{},
	types.DeleteRequest{},
	types.DeleteFileRequest{},
	types.DeleteDirectoryRequest{},
	types.CopyRequest{},
	types.CopyFileRequest{},
	types.CopyDirectoryRequest{},
	types.MoveRequest{},
	types.MoveFileRequest{},
	types.MoveDirectoryRequest{},
	types.GetCurrentTimeRequest{},
	types.ListDirRequest{},
	types.SearchRequest{},
	types.BatchDeleteRequest{},
	types.FileStatRequest{},
	types.FileExistsRequest{},
	types.ExecuteCommandRequest{},
	types.GetCommandBlacklistRequest{},
	types.UpdateCommandBlacklistRequest{},
	types.GetWorkingDirectoryRequest{},
	types.ChangeDirectoryRequest{},
	types.ExecuteCommandAsyncRequest{},
	types.GetCommandTaskRequest{},
	types.CancelCommandTaskRequest{},
	types.GetCommandHistoryRequest{},
	types.ClearCommandHistoryRequest{},
	types.SetPermissionLevelRequest{},
	types.GetPermissionLevelRequest{},
	types.GetSystemInfoRequest{},

	// 响应类型 / Response types
	types.FileExistsResponse{},
	types.ReadFileResponse{},
	types.ListDirResponse{},
	types.SearchResponse{},
	types.OperationResponse{},
	types.GetTimeResponse{},
	types.ExecuteCommandResponse{},
	types.GetCommandBlacklistResponse{},
	types.GetWorkingDirectoryResponse{},
	types.ExecuteCommandAsyncResponse{},
	types.GetCommandTaskResponse{},
	types.GetCommandHistoryResponse{},
	types.GetPermissionLevelResponse{},
	types.GetSystemInfoResponse{},
	types.CommandHistoryEntry{},
	types.CommandTask{},

	// 系统信息相关结构体 / System info related structures
	types.OSInfo{},
	types.CPUInfo{},
	types.MemoryInfo{},
	types.GPUInfo{},
	types.NetworkInfo{},
	types.IPAddress{},

	// ========== 配置相关结构体 / Configuration related structures ==========
	types.ServerConfig{},
	types.HTTPConfig{},
	types.SSEConfig{},

	// ========== MCP协议相关结构体 / MCP protocol related structures ==========
	types.MCPRequest{},
	types.MCPResponse{},
	types.MCPError{},
	types.MCPToolsListRequest{},
	types.MCPToolsCallRequest{},
	types.MCPInitializeRequest{},
	types.SSEMessage{},

	// ========== JSON Schema 相关结构体 / JSON Schema related structures ==========
	types.JSONSchema{},
	types.Property{},
	types.Items{},
}

// PretouchAll 预热所有已注册的类型
// 建议在程序启动时调用此函数，以避免首次请求时的延迟
// 返回预热失败的类型数量和第一个错误
func PretouchAll() (failedCount int, firstErr error) {
	for _, v := range pretouchTypes {
		t := reflect.TypeOf(v)
		if err := Pretouch(t); err != nil {
			failedCount++
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return failedCount, firstErr
}

// PretouchTypes 预热指定的类型列表
// 参数 values 是需要预热的类型实例（零值即可）
func PretouchTypes(values ...any) error {
	for _, v := range values {
		t := reflect.TypeOf(v)
		if err := Pretouch(t); err != nil {
			return err
		}
	}
	return nil
}

// RegisterPretouchType 注册需要预热的类型
// 可以在 init() 函数中调用此方法注册自定义类型
func RegisterPretouchType(v any) {
	pretouchTypes = append(pretouchTypes, v)
}
