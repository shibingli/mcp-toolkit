package filesystem

import (
	"context"

	"mcp-toolkit/pkg/transport"
	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools 注册所有文件系统工具到MCP服务器 / Register all filesystem tools to MCP server
func (s *Service) RegisterTools(mcpServer *mcp.Server) {
	// 创建文件 / Create file
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "create_file",
		Description: "创建新文件并写入内容 / Create a new file and write content",
	}, s.handleCreateFile)

	// 创建目录 / Create directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "create_directory",
		Description: "创建新目录 / Create a new directory",
	}, s.handleCreateDir)

	// 读取文件 / Read file
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "read_file",
		Description: "读取文件内容 / Read file content",
	}, s.handleReadFile)

	// 写入文件 / Write file
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "write_file",
		Description: "写入或覆盖文件内容 / Write or overwrite file content",
	}, s.handleWriteFile)

	// 删除文件或目录 / Delete file or directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "delete",
		Description: "删除文件或目录 / Delete file or directory",
	}, s.handleDelete)

	// 复制文件或目录 / Copy file or directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "copy",
		Description: "复制文件或目录 / Copy file or directory",
	}, s.handleCopy)

	// 移动文件或目录 / Move file or directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "move",
		Description: "移动或重命名文件或目录 / Move or rename file or directory",
	}, s.handleMove)

	// 列出目录 / List directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_directory",
		Description: "列出目录中的文件和子目录 / List files and subdirectories in a directory",
	}, s.handleListDir)

	// 搜索文件 / Search files
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "search_files",
		Description: "根据文件名模式搜索文件 / Search files by filename pattern",
	}, s.handleSearch)

	// 批量删除 / Batch delete
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "batch_delete",
		Description: "批量删除多个文件或目录 / Batch delete multiple files or directories",
	}, s.handleBatchDelete)

	// 文件状态 / File stat
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "file_stat",
		Description: "获取文件或目录的详细信息 / Get detailed information about a file or directory",
	}, s.handleFileStat)

	// 文件是否存在 / File exists
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "file_exists",
		Description: "检查文件或目录是否存在 / Check if a file or directory exists",
	}, s.handleFileExists)

	// 获取当前时间 / Get current time
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_current_time",
		Description: "获取当前系统时间 / Get current system time",
	}, s.handleGetCurrentTime)

	// 执行命令 / Execute command
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "execute_command",
		Description: "在沙箱目录内执行命令行命令 / Execute command line command within sandbox directory",
	}, s.handleExecuteCommand)

	// 获取命令黑名单 / Get command blacklist
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_command_blacklist",
		Description: "获取命令和目录黑名单 / Get command and directory blacklist",
	}, s.handleGetCommandBlacklist)

	// 更新命令黑名单 / Update command blacklist
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "update_command_blacklist",
		Description: "更新命令和目录黑名单 / Update command and directory blacklist",
	}, s.handleUpdateCommandBlacklist)

	// 获取当前工作目录 / Get working directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_working_directory",
		Description: "获取当前工作目录 / Get current working directory",
	}, s.handleGetWorkingDirectory)

	// 切换工作目录 / Change directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "change_directory",
		Description: "切换当前工作目录 / Change current working directory",
	}, s.handleChangeDirectory)

	// 异步执行命令 / Execute command asynchronously
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "execute_command_async",
		Description: "异步执行命令,返回任务ID / Execute command asynchronously, returns task ID",
	}, s.handleExecuteCommandAsync)

	// 获取命令任务 / Get command task
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_command_task",
		Description: "获取异步命令任务状态 / Get async command task status",
	}, s.handleGetCommandTask)

	// 取消命令任务 / Cancel command task
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "cancel_command_task",
		Description: "取消正在执行的命令任务 / Cancel running command task",
	}, s.handleCancelCommandTask)

	// 获取命令历史 / Get command history
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_command_history",
		Description: "获取命令执行历史记录 / Get command execution history",
	}, s.handleGetCommandHistory)

	// 清空命令历史 / Clear command history
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "clear_command_history",
		Description: "清空命令执行历史记录 / Clear command execution history",
	}, s.handleClearCommandHistory)

	// 设置权限级别 / Set permission level
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "set_permission_level",
		Description: "设置命令执行权限级别 / Set command execution permission level",
	}, s.handleSetPermissionLevel)

	// 获取权限级别 / Get permission level
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_permission_level",
		Description: "获取当前权限级别 / Get current permission level",
	}, s.handleGetPermissionLevel)
}

// handleCreateFile 处理创建文件请求 / Handle create file request
func (s *Service) handleCreateFile(_ context.Context, _ *mcp.CallToolRequest, args types.CreateFileRequest) (*mcp.CallToolResult, *types.CreateFileResponse, error) {
	resp, err := s.CreateFile(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleCreateDir 处理创建目录请求 / Handle create directory request
func (s *Service) handleCreateDir(_ context.Context, _ *mcp.CallToolRequest, args types.CreateDirRequest) (*mcp.CallToolResult, *types.CreateDirResponse, error) {
	resp, err := s.CreateDir(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleReadFile 处理读取文件请求 / Handle read file request
func (s *Service) handleReadFile(_ context.Context, _ *mcp.CallToolRequest, args types.ReadFileRequest) (*mcp.CallToolResult, *types.ReadFileResponse, error) {
	resp, err := s.ReadFile(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleWriteFile 处理写入文件请求 / Handle write file request
func (s *Service) handleWriteFile(_ context.Context, _ *mcp.CallToolRequest, args types.WriteFileRequest) (*mcp.CallToolResult, *types.WriteFileResponse, error) {
	resp, err := s.WriteFile(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleDelete 处理删除请求 / Handle delete request
func (s *Service) handleDelete(_ context.Context, _ *mcp.CallToolRequest, args types.DeleteRequest) (*mcp.CallToolResult, *types.DeleteResponse, error) {
	resp, err := s.Delete(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleCopy 处理复制请求 / Handle copy request
func (s *Service) handleCopy(_ context.Context, _ *mcp.CallToolRequest, args types.CopyRequest) (*mcp.CallToolResult, *types.CopyResponse, error) {
	resp, err := s.Copy(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleMove 处理移动请求 / Handle move request
func (s *Service) handleMove(_ context.Context, _ *mcp.CallToolRequest, args types.MoveRequest) (*mcp.CallToolResult, *types.MoveResponse, error) {
	resp, err := s.Move(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleListDir 处理列出目录请求 / Handle list directory request
func (s *Service) handleListDir(_ context.Context, _ *mcp.CallToolRequest, args types.ListDirRequest) (*mcp.CallToolResult, *types.ListDirResponse, error) {
	resp, err := s.ListDir(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleSearch 处理搜索请求 / Handle search request
func (s *Service) handleSearch(_ context.Context, _ *mcp.CallToolRequest, args types.SearchRequest) (*mcp.CallToolResult, *types.SearchResponse, error) {
	resp, err := s.Search(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleBatchDelete 处理批量删除请求 / Handle batch delete request
func (s *Service) handleBatchDelete(_ context.Context, _ *mcp.CallToolRequest, args types.BatchDeleteRequest) (*mcp.CallToolResult, *types.BatchDeleteResponse, error) {
	resp, err := s.BatchDelete(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleFileStat 处理文件状态请求 / Handle file stat request
func (s *Service) handleFileStat(_ context.Context, _ *mcp.CallToolRequest, args types.FileStatRequest) (*mcp.CallToolResult, *types.FileStatResponse, error) {
	resp, err := s.FileStat(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleFileExists 处理文件存在检查请求 / Handle file exists request
func (s *Service) handleFileExists(_ context.Context, _ *mcp.CallToolRequest, args types.FileExistsRequest) (*mcp.CallToolResult, *types.FileExistsResponse, error) {
	resp, err := s.FileExists(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleGetCurrentTime 处理获取当前时间请求 / Handle get current time request
func (s *Service) handleGetCurrentTime(_ context.Context, _ *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, *types.GetCurrentTimeResponse, error) {
	resp, err := s.GetCurrentTime()
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleExecuteCommand 处理执行命令请求 / Handle execute command request
func (s *Service) handleExecuteCommand(_ context.Context, _ *mcp.CallToolRequest, args types.ExecuteCommandRequest) (*mcp.CallToolResult, *types.ExecuteCommandResponse, error) {
	resp, err := s.ExecuteCommand(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleGetCommandBlacklist 处理获取命令黑名单请求 / Handle get command blacklist request
func (s *Service) handleGetCommandBlacklist(_ context.Context, _ *mcp.CallToolRequest, args types.GetCommandBlacklistRequest) (*mcp.CallToolResult, *types.GetCommandBlacklistResponse, error) {
	resp, err := s.GetCommandBlacklist(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleUpdateCommandBlacklist 处理更新命令黑名单请求 / Handle update command blacklist request
func (s *Service) handleUpdateCommandBlacklist(_ context.Context, _ *mcp.CallToolRequest, args types.UpdateCommandBlacklistRequest) (*mcp.CallToolResult, *types.UpdateCommandBlacklistResponse, error) {
	resp, err := s.UpdateCommandBlacklist(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleGetWorkingDirectory 处理获取工作目录请求 / Handle get working directory request
func (s *Service) handleGetWorkingDirectory(_ context.Context, _ *mcp.CallToolRequest, args types.GetWorkingDirectoryRequest) (*mcp.CallToolResult, *types.GetWorkingDirectoryResponse, error) {
	resp, err := s.GetWorkingDirectory(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleChangeDirectory 处理切换目录请求 / Handle change directory request
func (s *Service) handleChangeDirectory(_ context.Context, _ *mcp.CallToolRequest, args types.ChangeDirectoryRequest) (*mcp.CallToolResult, *types.ChangeDirectoryResponse, error) {
	resp, err := s.ChangeDirectory(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// RegisterToolsToRegistry 注册所有文件系统工具到工具注册表 / Register all filesystem tools to tool registry
func (s *Service) RegisterToolsToRegistry(registry *transport.ToolRegistry) {
	// 创建文件 / Create file
	registry.RegisterTool(&mcp.Tool{
		Name:        "create_file",
		Description: "创建新文件并写入内容 / Create a new file and write content",
	}, s.wrapCreateFile)

	// 创建目录 / Create directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "create_directory",
		Description: "创建新目录 / Create a new directory",
	}, s.wrapCreateDir)

	// 读取文件 / Read file
	registry.RegisterTool(&mcp.Tool{
		Name:        "read_file",
		Description: "读取文件内容 / Read file content",
	}, s.wrapReadFile)

	// 写入文件 / Write file
	registry.RegisterTool(&mcp.Tool{
		Name:        "write_file",
		Description: "写入或覆盖文件内容 / Write or overwrite file content",
	}, s.wrapWriteFile)

	// 删除文件或目录 / Delete file or directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "delete",
		Description: "删除文件或目录 / Delete file or directory",
	}, s.wrapDelete)

	// 复制文件或目录 / Copy file or directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "copy",
		Description: "复制文件或目录 / Copy file or directory",
	}, s.wrapCopy)

	// 移动文件或目录 / Move file or directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "move",
		Description: "移动或重命名文件或目录 / Move or rename file or directory",
	}, s.wrapMove)

	// 列出目录 / List directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "list_directory",
		Description: "列出目录中的文件和子目录 / List files and subdirectories in a directory",
	}, s.wrapListDir)

	// 搜索文件 / Search files
	registry.RegisterTool(&mcp.Tool{
		Name:        "search_files",
		Description: "根据文件名模式搜索文件 / Search files by filename pattern",
	}, s.wrapSearch)

	// 批量删除 / Batch delete
	registry.RegisterTool(&mcp.Tool{
		Name:        "batch_delete",
		Description: "批量删除多个文件或目录 / Batch delete multiple files or directories",
	}, s.wrapBatchDelete)

	// 文件状态 / File stat
	registry.RegisterTool(&mcp.Tool{
		Name:        "file_stat",
		Description: "获取文件或目录的详细信息 / Get detailed information about a file or directory",
	}, s.wrapFileStat)

	// 文件是否存在 / File exists
	registry.RegisterTool(&mcp.Tool{
		Name:        "file_exists",
		Description: "检查文件或目录是否存在 / Check if a file or directory exists",
	}, s.wrapFileExists)

	// 获取当前时间 / Get current time
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_current_time",
		Description: "获取当前系统时间 / Get current system time",
	}, s.wrapGetCurrentTime)

	// 执行命令 / Execute command
	registry.RegisterTool(&mcp.Tool{
		Name:        "execute_command",
		Description: "在沙箱目录内执行命令行命令 / Execute command line command within sandbox directory",
	}, s.wrapExecuteCommand)

	// 获取命令黑名单 / Get command blacklist
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_command_blacklist",
		Description: "获取命令和目录黑名单 / Get command and directory blacklist",
	}, s.wrapGetCommandBlacklist)

	// 更新命令黑名单 / Update command blacklist
	registry.RegisterTool(&mcp.Tool{
		Name:        "update_command_blacklist",
		Description: "更新命令和目录黑名单 / Update command and directory blacklist",
	}, s.wrapUpdateCommandBlacklist)

	// 获取当前工作目录 / Get working directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_working_directory",
		Description: "获取当前工作目录 / Get current working directory",
	}, s.wrapGetWorkingDirectory)

	// 切换工作目录 / Change directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "change_directory",
		Description: "切换当前工作目录 / Change current working directory",
	}, s.wrapChangeDirectory)

	// 异步执行命令 / Execute command asynchronously
	registry.RegisterTool(&mcp.Tool{
		Name:        "execute_command_async",
		Description: "异步执行命令,返回任务ID / Execute command asynchronously, returns task ID",
	}, s.wrapExecuteCommandAsync)

	// 获取命令任务 / Get command task
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_command_task",
		Description: "获取异步命令任务状态 / Get async command task status",
	}, s.wrapGetCommandTask)

	// 取消命令任务 / Cancel command task
	registry.RegisterTool(&mcp.Tool{
		Name:        "cancel_command_task",
		Description: "取消正在执行的命令任务 / Cancel running command task",
	}, s.wrapCancelCommandTask)

	// 获取命令历史 / Get command history
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_command_history",
		Description: "获取命令执行历史记录 / Get command execution history",
	}, s.wrapGetCommandHistory)

	// 清空命令历史 / Clear command history
	registry.RegisterTool(&mcp.Tool{
		Name:        "clear_command_history",
		Description: "清空命令执行历史记录 / Clear command execution history",
	}, s.wrapClearCommandHistory)

	// 设置权限级别 / Set permission level
	registry.RegisterTool(&mcp.Tool{
		Name:        "set_permission_level",
		Description: "设置命令执行权限级别 / Set command execution permission level",
	}, s.wrapSetPermissionLevel)

	// 获取权限级别 / Get permission level
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_permission_level",
		Description: "获取当前权限级别 / Get current permission level",
	}, s.wrapGetPermissionLevel)
}

// 包装函数，将MCP处理器转换为ToolHandler / Wrapper functions to convert MCP handlers to ToolHandler

func (s *Service) wrapCreateFile(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.CreateFileRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleCreateFile(ctx, nil, args)
	return result, err
}

func (s *Service) wrapCreateDir(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.CreateDirRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleCreateDir(ctx, nil, args)
	return result, err
}

func (s *Service) wrapReadFile(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.ReadFileRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleReadFile(ctx, nil, args)
	return result, err
}

func (s *Service) wrapWriteFile(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.WriteFileRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleWriteFile(ctx, nil, args)
	return result, err
}

func (s *Service) wrapDelete(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.DeleteRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleDelete(ctx, nil, args)
	return result, err
}

func (s *Service) wrapCopy(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.CopyRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleCopy(ctx, nil, args)
	return result, err
}

func (s *Service) wrapMove(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.MoveRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleMove(ctx, nil, args)
	return result, err
}

func (s *Service) wrapListDir(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.ListDirRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleListDir(ctx, nil, args)
	return result, err
}

func (s *Service) wrapSearch(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.SearchRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleSearch(ctx, nil, args)
	return result, err
}

func (s *Service) wrapBatchDelete(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.BatchDeleteRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleBatchDelete(ctx, nil, args)
	return result, err
}

func (s *Service) wrapFileStat(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.FileStatRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleFileStat(ctx, nil, args)
	return result, err
}

func (s *Service) wrapFileExists(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.FileExistsRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleFileExists(ctx, nil, args)
	return result, err
}

func (s *Service) wrapGetCurrentTime(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	result, _, err := s.handleGetCurrentTime(ctx, nil, struct{}{})
	return result, err
}

func (s *Service) wrapExecuteCommand(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.ExecuteCommandRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleExecuteCommand(ctx, nil, args)
	return result, err
}

func (s *Service) wrapGetCommandBlacklist(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.GetCommandBlacklistRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleGetCommandBlacklist(ctx, nil, args)
	return result, err
}

func (s *Service) wrapUpdateCommandBlacklist(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.UpdateCommandBlacklistRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleUpdateCommandBlacklist(ctx, nil, args)
	return result, err
}

func (s *Service) wrapGetWorkingDirectory(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.GetWorkingDirectoryRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleGetWorkingDirectory(ctx, nil, args)
	return result, err
}

func (s *Service) wrapChangeDirectory(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.ChangeDirectoryRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleChangeDirectory(ctx, nil, args)
	return result, err
}

func (s *Service) wrapExecuteCommandAsync(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.ExecuteCommandAsyncRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleExecuteCommandAsync(ctx, nil, args)
	return result, err
}

func (s *Service) wrapGetCommandTask(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.GetCommandTaskRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleGetCommandTask(ctx, nil, args)
	return result, err
}

func (s *Service) wrapCancelCommandTask(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.CancelCommandTaskRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleCancelCommandTask(ctx, nil, args)
	return result, err
}

func (s *Service) wrapGetCommandHistory(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.GetCommandHistoryRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleGetCommandHistory(ctx, nil, args)
	return result, err
}

func (s *Service) wrapClearCommandHistory(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.ClearCommandHistoryRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleClearCommandHistory(ctx, nil, args)
	return result, err
}

func (s *Service) wrapSetPermissionLevel(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.SetPermissionLevelRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleSetPermissionLevel(ctx, nil, args)
	return result, err
}

func (s *Service) wrapGetPermissionLevel(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.GetPermissionLevelRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleGetPermissionLevel(ctx, nil, args)
	return result, err
}

// handleExecuteCommandAsync 处理异步执行命令请求 / Handle execute command async request
func (s *Service) handleExecuteCommandAsync(_ context.Context, _ *mcp.CallToolRequest, args types.ExecuteCommandAsyncRequest) (*mcp.CallToolResult, *types.ExecuteCommandAsyncResponse, error) {
	resp, err := s.ExecuteCommandAsync(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleGetCommandTask 处理获取命令任务请求 / Handle get command task request
func (s *Service) handleGetCommandTask(_ context.Context, _ *mcp.CallToolRequest, args types.GetCommandTaskRequest) (*mcp.CallToolResult, *types.GetCommandTaskResponse, error) {
	resp, err := s.GetCommandTask(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleCancelCommandTask 处理取消命令任务请求 / Handle cancel command task request
func (s *Service) handleCancelCommandTask(_ context.Context, _ *mcp.CallToolRequest, args types.CancelCommandTaskRequest) (*mcp.CallToolResult, *types.CancelCommandTaskResponse, error) {
	resp, err := s.CancelCommandTask(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleGetCommandHistory 处理获取命令历史请求 / Handle get command history request
func (s *Service) handleGetCommandHistory(_ context.Context, _ *mcp.CallToolRequest, args types.GetCommandHistoryRequest) (*mcp.CallToolResult, *types.GetCommandHistoryResponse, error) {
	resp, err := s.GetCommandHistory(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleClearCommandHistory 处理清空命令历史请求 / Handle clear command history request
func (s *Service) handleClearCommandHistory(_ context.Context, _ *mcp.CallToolRequest, args types.ClearCommandHistoryRequest) (*mcp.CallToolResult, *types.ClearCommandHistoryResponse, error) {
	resp, err := s.ClearCommandHistory(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleSetPermissionLevel 处理设置权限级别请求 / Handle set permission level request
func (s *Service) handleSetPermissionLevel(_ context.Context, _ *mcp.CallToolRequest, args types.SetPermissionLevelRequest) (*mcp.CallToolResult, *types.SetPermissionLevelResponse, error) {
	resp, err := s.SetPermissionLevel(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}

// handleGetPermissionLevel 处理获取权限级别请求 / Handle get permission level request
func (s *Service) handleGetPermissionLevel(_ context.Context, _ *mcp.CallToolRequest, args types.GetPermissionLevelRequest) (*mcp.CallToolResult, *types.GetPermissionLevelResponse, error) {
	resp, err := s.GetPermissionLevel(&args)
	if err != nil {
		return nil, nil, err
	}

	resultJSON, _ := json.MarshalToString(resp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultJSON},
		},
	}, resp, nil
}
