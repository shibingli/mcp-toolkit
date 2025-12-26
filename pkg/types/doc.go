// Package types 定义了 MCP Toolkit 中使用的所有数据类型和常量
//
// 本包包含以下主要内容：
//
// # 常量定义
//
// 全局常量（constants.go）：
//   - ProtocolVersion：MCP 协议版本
//   - ServerName：服务器名称
//   - ServerVersion：服务器版本
//   - DefaultSandboxDir：默认沙箱目录
//
// # 配置类型
//
// 传输配置（config.go）：
//   - TransportType：传输方式类型（Stdio/HTTP/SSE）
//   - HTTPConfig：HTTP 传输配置
//   - SSEConfig：SSE 传输配置
//
// # 文件系统类型
//
// 请求类型（filesystem.go）：
//   - CreateFileRequest：创建文件请求
//   - CreateDirRequest：创建目录请求
//   - ReadFileRequest：读取文件请求
//   - WriteFileRequest：写入文件请求
//   - DeleteRequest：删除请求
//   - CopyRequest：复制请求
//   - MoveRequest：移动请求
//   - ListDirRequest：列出目录请求
//   - SearchRequest：搜索请求
//   - GetFileInfoRequest：获取文件信息请求
//   - FileExistsRequest：文件存在检查请求
//   - ExecuteCommandRequest：执行命令请求
//   - ExecuteCommandAsyncRequest：异步执行命令请求
//   - GetCommandTaskRequest：获取命令任务请求
//   - CancelCommandTaskRequest：取消命令任务请求
//   - UpdateCommandBlacklistRequest：更新命令黑名单请求
//   - ChangeDirRequest：切换目录请求
//   - SetPermissionLevelRequest：设置权限级别请求
//
// 响应类型（filesystem.go）：
//   - CreateFileResponse：创建文件响应
//   - CreateDirResponse：创建目录响应
//   - ReadFileResponse：读取文件响应
//   - WriteFileResponse：写入文件响应
//   - DeleteResponse：删除响应
//   - CopyResponse：复制响应
//   - MoveResponse：移动响应
//   - ListDirResponse：列出目录响应
//   - SearchResponse：搜索响应
//   - GetFileInfoResponse：获取文件信息响应
//   - FileExistsResponse：文件存在检查响应
//   - ExecuteCommandResponse：执行命令响应
//   - ExecuteCommandAsyncResponse：异步执行命令响应
//   - GetCommandTaskResponse：获取命令任务响应
//   - CancelCommandTaskResponse：取消命令任务响应
//   - GetCommandBlacklistResponse：获取命令黑名单响应
//   - UpdateCommandBlacklistResponse：更新命令黑名单响应
//   - GetWorkingDirectoryResponse：获取工作目录响应
//   - ChangeDirResponse：切换目录响应
//   - GetCurrentTimeResponse：获取当前时间响应
//   - SetPermissionLevelResponse：设置权限级别响应
//   - GetPermissionLevelResponse：获取权限级别响应
//
// 数据模型（filesystem.go）：
//   - FileInfo：文件信息
//   - CommandTask：命令任务
//   - CommandBlacklist：命令黑名单
//   - PermissionLevel：权限级别
//
// # MCP 协议类型
//
// MCP 请求/响应（mcp.go）：
//   - MCPRequest：MCP 请求
//   - MCPResponse：MCP 响应
//   - MCPError：MCP 错误
//   - MCPToolsCallRequest：工具调用请求
//
// SSE 消息（mcp.go）：
//   - SSEMessage：SSE 消息
//
// # 使用示例
//
// 创建文件请求：
//
//	req := &types.CreateFileRequest{
//	    Path:    "test.txt",
//	    Content: "Hello, World!",
//	}
//
// 配置 HTTP 传输：
//
//	config := &types.HTTPConfig{
//	    Host:           "127.0.0.1",
//	    Port:           8080,
//	    Path:           "/mcp",
//	    EnableCORS:     true,
//	    AllowedOrigins: []string{"*"},
//	    Timeout:        30,
//	}
//
// 使用全局常量：
//
//	fmt.Println(types.ServerName)     // "mcp-toolkit"
//	fmt.Println(types.ServerVersion)  // "1.0.1"
//	fmt.Println(types.ProtocolVersion) // "2024-11-05"
//
// # 设计原则
//
//   - 所有结构体都定义在 pkg/types 包中（符合项目规范）
//   - 所有字段都有中英文注释
//   - 使用 JSON 标签支持序列化/反序列化
//   - 常量集中管理，避免硬编码
//   - 类型安全，避免使用 interface{} 和 map[string]interface{}
//
// # JSON 预热
//
// 所有结构体都在 pkg/utils/json/pretouch.go 中进行 JSON 预热，
// 以提高首次序列化/反序列化的性能（仅对 Sonic 有效）。
package types
