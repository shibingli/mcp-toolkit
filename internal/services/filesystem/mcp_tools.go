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
		InputSchema: types.GetToolSchema("create_file"),
	}, s.handleCreateFile)

	// 创建目录 / Create directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "create_directory",
		Description: "创建新目录 / Create a new directory",
		InputSchema: types.GetToolSchema("create_directory"),
	}, s.handleCreateDir)

	// 读取文件 / Read file
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "read_file",
		Description: "读取文件内容 / Read file content",
		InputSchema: types.GetToolSchema("read_file"),
	}, s.handleReadFile)

	// 写入文件 / Write file
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "write_file",
		Description: "写入或覆盖文件内容 / Write or overwrite file content",
		InputSchema: types.GetToolSchema("write_file"),
	}, s.handleWriteFile)

	// 删除文件或目录（自动判断）/ Delete file or directory (auto-detect)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "delete",
		Description: "[推荐] 删除文件或目录，自动判断类型。这是删除操作的首选工具。/ [Recommended] Delete file or directory, auto-detect type. This is the preferred tool for deletion.",
		InputSchema: types.GetToolSchema("delete"),
	}, s.handleDelete)

	// 删除文件（仅限文件）/ Delete file only
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "delete_file",
		Description: "删除指定的文件。仅用于删除文件，不能删除目录。如果不确定是文件还是目录，请使用 delete 工具。/ Delete specified file. Only for files, cannot delete directories. Use 'delete' tool if unsure.",
		InputSchema: types.GetToolSchema("delete_file"),
	}, s.handleDeleteFile)

	// 删除目录（仅限目录）/ Delete directory only
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "delete_directory",
		Description: "删除指定的目录及其所有内容。仅用于删除目录，不能删除文件。如果不确定是文件还是目录，请使用 delete 工具。/ Delete specified directory and all its contents. Only for directories, cannot delete files. Use 'delete' tool if unsure.",
		InputSchema: types.GetToolSchema("delete_directory"),
	}, s.handleDeleteDirectory)

	// 复制文件或目录 / Copy file or directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "copy",
		Description: "复制文件或目录到新位置 / Copy file or directory to new location",
		InputSchema: types.GetToolSchema("copy"),
	}, s.handleCopy)

	// 复制文件（仅限文件）/ Copy file only
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "copy_file",
		Description: "复制指定的文件。仅用于复制文件，不能复制目录。如果不确定是文件还是目录，请使用 copy 工具。/ Copy specified file. Only for files, cannot copy directories. Use 'copy' tool if unsure.",
		InputSchema: types.GetToolSchema("copy_file"),
	}, s.handleCopyFile)

	// 复制目录（仅限目录）/ Copy directory only
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "copy_directory",
		Description: "复制指定的目录及其所有内容。仅用于复制目录，不能复制文件。如果不确定是文件还是目录，请使用 copy 工具。/ Copy specified directory and all its contents. Only for directories, cannot copy files. Use 'copy' tool if unsure.",
		InputSchema: types.GetToolSchema("copy_directory"),
	}, s.handleCopyDirectory)

	// 移动文件或目录（自动判断）/ Move file or directory (auto-detect)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "move",
		Description: "[推荐] 移动或重命名文件或目录，自动判断类型。这是移动操作的首选工具。/ [Recommended] Move or rename file or directory, auto-detect type. This is the preferred tool for moving.",
		InputSchema: types.GetToolSchema("move"),
	}, s.handleMove)

	// 移动文件（仅限文件）/ Move file only
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "move_file",
		Description: "移动或重命名指定的文件。仅用于移动文件，不能移动目录。如果不确定是文件还是目录，请使用 move 工具。/ Move or rename specified file. Only for files, cannot move directories. Use 'move' tool if unsure.",
		InputSchema: types.GetToolSchema("move_file"),
	}, s.handleMoveFile)

	// 移动目录（仅限目录）/ Move directory only
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "move_directory",
		Description: "移动或重命名指定的目录。仅用于移动目录，不能移动文件。如果不确定是文件还是目录，请使用 move 工具。/ Move or rename specified directory. Only for directories, cannot move files. Use 'move' tool if unsure.",
		InputSchema: types.GetToolSchema("move_directory"),
	}, s.handleMoveDirectory)

	// 列出目录 / List directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_directory",
		Description: "列出目录中的文件和子目录 / List files and subdirectories in a directory",
		InputSchema: types.GetToolSchema("list_directory"),
	}, s.handleListDir)

	// 搜索文件 / Search files
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "search_files",
		Description: "根据文件名模式搜索文件 / Search files by filename pattern",
		InputSchema: types.GetToolSchema("search_files"),
	}, s.handleSearch)

	// 批量删除 / Batch delete
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "batch_delete",
		Description: "批量删除多个文件或目录 / Batch delete multiple files or directories",
		InputSchema: types.GetToolSchema("batch_delete"),
	}, s.handleBatchDelete)

	// 文件状态 / File stat
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "file_stat",
		Description: "获取文件或目录的详细信息 / Get detailed information about a file or directory",
		InputSchema: types.GetToolSchema("file_stat"),
	}, s.handleFileStat)

	// 文件是否存在 / File exists
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "file_exists",
		Description: "检查文件或目录是否存在 / Check if a file or directory exists",
		InputSchema: types.GetToolSchema("file_exists"),
	}, s.handleFileExists)

	// 获取当前时间 / Get current time
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_current_time",
		Description: "获取当前系统时间。可指定时区（如 'Asia/Shanghai'、'America/New_York'），不指定则使用系统本地时区。返回格式化的日期时间、时区信息和Unix时间戳。/ Get current system time. Can specify timezone (e.g. 'Asia/Shanghai', 'America/New_York'), uses system local timezone if not specified. Returns formatted datetime, timezone info and Unix timestamp.",
		InputSchema: types.GetToolSchema("get_current_time"),
	}, s.handleGetCurrentTime)

	// 执行命令 / Execute command
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "execute_command",
		Description: "执行shell命令（如git、npm、python等）。注意：不要用于文件操作（创建、删除、复制、移动文件/目录），请使用专门的文件工具如 create_file、delete、delete_file、delete_directory、copy、move 等。/ Execute shell commands (git, npm, python, etc). NOTE: Do NOT use for file operations (create, delete, copy, move files/directories), use dedicated file tools instead: create_file, delete, delete_file, delete_directory, copy, move, etc.",
		InputSchema: types.GetToolSchema("execute_command"),
	}, s.handleExecuteCommand)

	// 获取命令黑名单 / Get command blacklist
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_command_blacklist",
		Description: "获取命令和目录黑名单 / Get command and directory blacklist",
		InputSchema: types.GetToolSchema("get_command_blacklist"),
	}, s.handleGetCommandBlacklist)

	// 更新命令黑名单 / Update command blacklist
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "update_command_blacklist",
		Description: "更新命令和目录黑名单 / Update command and directory blacklist",
		InputSchema: types.GetToolSchema("update_command_blacklist"),
	}, s.handleUpdateCommandBlacklist)

	// 获取当前工作目录 / Get working directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_working_directory",
		Description: "获取当前工作目录 / Get current working directory",
		InputSchema: types.GetToolSchema("get_working_directory"),
	}, s.handleGetWorkingDirectory)

	// 切换工作目录 / Change directory
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "change_directory",
		Description: "切换当前工作目录 / Change current working directory",
		InputSchema: types.GetToolSchema("change_directory"),
	}, s.handleChangeDirectory)

	// 异步执行命令 / Execute command asynchronously
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "execute_command_async",
		Description: "异步执行命令,返回任务ID / Execute command asynchronously, returns task ID",
		InputSchema: types.GetToolSchema("execute_command_async"),
	}, s.handleExecuteCommandAsync)

	// 获取命令任务 / Get command task
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_command_task",
		Description: "获取异步命令任务状态 / Get async command task status",
		InputSchema: types.GetToolSchema("get_command_task_status"),
	}, s.handleGetCommandTask)

	// 取消命令任务 / Cancel command task
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "cancel_command_task",
		Description: "取消正在执行的命令任务 / Cancel running command task",
		InputSchema: types.GetToolSchema("cancel_command_task"),
	}, s.handleCancelCommandTask)

	// 获取命令历史 / Get command history
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_command_history",
		Description: "获取命令执行历史记录 / Get command execution history",
		InputSchema: types.GetToolSchema("get_command_history"),
	}, s.handleGetCommandHistory)

	// 清空命令历史 / Clear command history
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "clear_command_history",
		Description: "清空命令执行历史记录 / Clear command execution history",
		InputSchema: types.GetToolSchema("clear_command_history"),
	}, s.handleClearCommandHistory)

	// 设置权限级别 / Set permission level
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "set_permission_level",
		Description: "设置命令执行权限级别 / Set command execution permission level",
		InputSchema: types.GetToolSchema("set_permission_level"),
	}, s.handleSetPermissionLevel)

	// 获取权限级别 / Get permission level
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_permission_level",
		Description: "获取当前权限级别 / Get current permission level",
		InputSchema: types.GetToolSchema("get_permission_level"),
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

// handleDelete 处理删除请求（自动判断文件或目录）/ Handle delete request (auto-detect file or directory)
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

// handleDeleteFile 处理删除文件请求 / Handle delete file request
func (s *Service) handleDeleteFile(_ context.Context, _ *mcp.CallToolRequest, args types.DeleteFileRequest) (*mcp.CallToolResult, *types.DeleteFileResponse, error) {
	resp, err := s.DeleteFile(&args)
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

// handleDeleteDirectory 处理删除目录请求 / Handle delete directory request
func (s *Service) handleDeleteDirectory(_ context.Context, _ *mcp.CallToolRequest, args types.DeleteDirectoryRequest) (*mcp.CallToolResult, *types.DeleteDirectoryResponse, error) {
	// 默认递归删除 / Default to recursive delete
	if !args.Recursive {
		args.Recursive = true
	}
	resp, err := s.DeleteDirectory(&args)
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

// handleCopy 处理复制请求（自动判断文件或目录）/ Handle copy request (auto-detect file or directory)
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

// handleCopyFile 处理复制文件请求 / Handle copy file request
func (s *Service) handleCopyFile(_ context.Context, _ *mcp.CallToolRequest, args types.CopyFileRequest) (*mcp.CallToolResult, *types.CopyFileResponse, error) {
	resp, err := s.CopyFile(&args)
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

// handleCopyDirectory 处理复制目录请求 / Handle copy directory request
func (s *Service) handleCopyDirectory(_ context.Context, _ *mcp.CallToolRequest, args types.CopyDirectoryRequest) (*mcp.CallToolResult, *types.CopyDirectoryResponse, error) {
	resp, err := s.CopyDirectory(&args)
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

// handleMove 处理移动请求（自动判断文件或目录）/ Handle move request (auto-detect file or directory)
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

// handleMoveFile 处理移动文件请求 / Handle move file request
func (s *Service) handleMoveFile(_ context.Context, _ *mcp.CallToolRequest, args types.MoveFileRequest) (*mcp.CallToolResult, *types.MoveFileResponse, error) {
	resp, err := s.MoveFile(&args)
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

// handleMoveDirectory 处理移动目录请求 / Handle move directory request
func (s *Service) handleMoveDirectory(_ context.Context, _ *mcp.CallToolRequest, args types.MoveDirectoryRequest) (*mcp.CallToolResult, *types.MoveDirectoryResponse, error) {
	resp, err := s.MoveDirectory(&args)
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
func (s *Service) handleGetCurrentTime(_ context.Context, _ *mcp.CallToolRequest, args types.GetCurrentTimeRequest) (*mcp.CallToolResult, *types.GetCurrentTimeResponse, error) {
	resp, err := s.GetCurrentTime(&args)
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
		InputSchema: types.GetToolSchema("create_file"),
	}, s.wrapCreateFile)

	// 创建目录 / Create directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "create_directory",
		Description: "创建新目录 / Create a new directory",
		InputSchema: types.GetToolSchema("create_directory"),
	}, s.wrapCreateDir)

	// 读取文件 / Read file
	registry.RegisterTool(&mcp.Tool{
		Name:        "read_file",
		Description: "读取文件内容 / Read file content",
		InputSchema: types.GetToolSchema("read_file"),
	}, s.wrapReadFile)

	// 写入文件 / Write file
	registry.RegisterTool(&mcp.Tool{
		Name:        "write_file",
		Description: "写入或覆盖文件内容 / Write or overwrite file content",
		InputSchema: types.GetToolSchema("write_file"),
	}, s.wrapWriteFile)

	// 删除文件或目录（自动判断）/ Delete file or directory (auto-detect)
	registry.RegisterTool(&mcp.Tool{
		Name:        "delete",
		Description: "[推荐] 删除文件或目录，自动判断类型。这是删除操作的首选工具。/ [Recommended] Delete file or directory, auto-detect type. This is the preferred tool for deletion.",
		InputSchema: types.GetToolSchema("delete"),
	}, s.wrapDelete)

	// 删除文件（仅限文件）/ Delete file only
	registry.RegisterTool(&mcp.Tool{
		Name:        "delete_file",
		Description: "删除指定的文件。仅用于删除文件，不能删除目录。如果不确定是文件还是目录，请使用 delete 工具。/ Delete specified file. Only for files, cannot delete directories. Use 'delete' tool if unsure.",
		InputSchema: types.GetToolSchema("delete_file"),
	}, s.wrapDeleteFile)

	// 删除目录（仅限目录）/ Delete directory only
	registry.RegisterTool(&mcp.Tool{
		Name:        "delete_directory",
		Description: "删除指定的目录及其所有内容。仅用于删除目录，不能删除文件。如果不确定是文件还是目录，请使用 delete 工具。/ Delete specified directory and all its contents. Only for directories, cannot delete files. Use 'delete' tool if unsure.",
		InputSchema: types.GetToolSchema("delete_directory"),
	}, s.wrapDeleteDirectory)

	// 复制文件或目录（自动判断）/ Copy file or directory (auto-detect)
	registry.RegisterTool(&mcp.Tool{
		Name:        "copy",
		Description: "[推荐] 复制文件或目录到新位置，自动判断类型。这是复制操作的首选工具。/ [Recommended] Copy file or directory to new location, auto-detect type. This is the preferred tool for copying.",
		InputSchema: types.GetToolSchema("copy"),
	}, s.wrapCopy)

	// 复制文件（仅限文件）/ Copy file only
	registry.RegisterTool(&mcp.Tool{
		Name:        "copy_file",
		Description: "复制指定的文件。仅用于复制文件，不能复制目录。如果不确定是文件还是目录，请使用 copy 工具。/ Copy specified file. Only for files, cannot copy directories. Use 'copy' tool if unsure.",
		InputSchema: types.GetToolSchema("copy_file"),
	}, s.wrapCopyFile)

	// 复制目录（仅限目录）/ Copy directory only
	registry.RegisterTool(&mcp.Tool{
		Name:        "copy_directory",
		Description: "复制指定的目录及其所有内容。仅用于复制目录，不能复制文件。如果不确定是文件还是目录，请使用 copy 工具。/ Copy specified directory and all its contents. Only for directories, cannot copy files. Use 'copy' tool if unsure.",
		InputSchema: types.GetToolSchema("copy_directory"),
	}, s.wrapCopyDirectory)

	// 移动文件或目录（自动判断）/ Move file or directory (auto-detect)
	registry.RegisterTool(&mcp.Tool{
		Name:        "move",
		Description: "[推荐] 移动或重命名文件或目录，自动判断类型。这是移动操作的首选工具。/ [Recommended] Move or rename file or directory, auto-detect type. This is the preferred tool for moving.",
		InputSchema: types.GetToolSchema("move"),
	}, s.wrapMove)

	// 移动文件（仅限文件）/ Move file only
	registry.RegisterTool(&mcp.Tool{
		Name:        "move_file",
		Description: "移动或重命名指定的文件。仅用于移动文件，不能移动目录。如果不确定是文件还是目录，请使用 move 工具。/ Move or rename specified file. Only for files, cannot move directories. Use 'move' tool if unsure.",
		InputSchema: types.GetToolSchema("move_file"),
	}, s.wrapMoveFile)

	// 移动目录（仅限目录）/ Move directory only
	registry.RegisterTool(&mcp.Tool{
		Name:        "move_directory",
		Description: "移动或重命名指定的目录。仅用于移动目录，不能移动文件。如果不确定是文件还是目录，请使用 move 工具。/ Move or rename specified directory. Only for directories, cannot move files. Use 'move' tool if unsure.",
		InputSchema: types.GetToolSchema("move_directory"),
	}, s.wrapMoveDirectory)

	// 列出目录 / List directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "list_directory",
		Description: "列出目录中的文件和子目录 / List files and subdirectories in a directory",
		InputSchema: types.GetToolSchema("list_directory"),
	}, s.wrapListDir)

	// 搜索文件 / Search files
	registry.RegisterTool(&mcp.Tool{
		Name:        "search_files",
		Description: "根据文件名模式搜索文件 / Search files by filename pattern",
		InputSchema: types.GetToolSchema("search_files"),
	}, s.wrapSearch)

	// 批量删除 / Batch delete
	registry.RegisterTool(&mcp.Tool{
		Name:        "batch_delete",
		Description: "批量删除多个文件或目录 / Batch delete multiple files or directories",
		InputSchema: types.GetToolSchema("batch_delete"),
	}, s.wrapBatchDelete)

	// 文件状态 / File stat
	registry.RegisterTool(&mcp.Tool{
		Name:        "file_stat",
		Description: "获取文件或目录的详细信息 / Get detailed information about a file or directory",
		InputSchema: types.GetToolSchema("file_stat"),
	}, s.wrapFileStat)

	// 文件是否存在 / File exists
	registry.RegisterTool(&mcp.Tool{
		Name:        "file_exists",
		Description: "检查文件或目录是否存在 / Check if a file or directory exists",
		InputSchema: types.GetToolSchema("file_exists"),
	}, s.wrapFileExists)

	// 获取当前时间 / Get current time
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_current_time",
		Description: "获取当前系统时间。可指定时区（如 'Asia/Shanghai'、'America/New_York'），不指定则使用系统本地时区。返回格式化的日期时间、时区信息和Unix时间戳。/ Get current system time. Can specify timezone (e.g. 'Asia/Shanghai', 'America/New_York'), uses system local timezone if not specified. Returns formatted datetime, timezone info and Unix timestamp.",
		InputSchema: types.GetToolSchema("get_current_time"),
	}, s.wrapGetCurrentTime)

	// 执行命令 / Execute command
	registry.RegisterTool(&mcp.Tool{
		Name:        "execute_command",
		Description: "执行shell命令（如git、npm、python等）。注意：不要用于文件操作（创建、删除、复制、移动文件/目录），请使用专门的文件工具如 create_file、delete、delete_file、delete_directory、copy、move 等。/ Execute shell commands (git, npm, python, etc). NOTE: Do NOT use for file operations (create, delete, copy, move files/directories), use dedicated file tools instead: create_file, delete, delete_file, delete_directory, copy, move, etc.",
		InputSchema: types.GetToolSchema("execute_command"),
	}, s.wrapExecuteCommand)

	// 获取命令黑名单 / Get command blacklist
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_command_blacklist",
		Description: "获取命令和目录黑名单 / Get command and directory blacklist",
		InputSchema: types.GetToolSchema("get_command_blacklist"),
	}, s.wrapGetCommandBlacklist)

	// 更新命令黑名单 / Update command blacklist
	registry.RegisterTool(&mcp.Tool{
		Name:        "update_command_blacklist",
		Description: "更新命令和目录黑名单 / Update command and directory blacklist",
		InputSchema: types.GetToolSchema("update_command_blacklist"),
	}, s.wrapUpdateCommandBlacklist)

	// 获取当前工作目录 / Get working directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_working_directory",
		Description: "获取当前工作目录 / Get current working directory",
		InputSchema: types.GetToolSchema("get_working_directory"),
	}, s.wrapGetWorkingDirectory)

	// 切换工作目录 / Change directory
	registry.RegisterTool(&mcp.Tool{
		Name:        "change_directory",
		Description: "切换当前工作目录 / Change current working directory",
		InputSchema: types.GetToolSchema("change_directory"),
	}, s.wrapChangeDirectory)

	// 异步执行命令 / Execute command asynchronously
	registry.RegisterTool(&mcp.Tool{
		Name:        "execute_command_async",
		Description: "异步执行命令,返回任务ID / Execute command asynchronously, returns task ID",
		InputSchema: types.GetToolSchema("execute_command_async"),
	}, s.wrapExecuteCommandAsync)

	// 获取命令任务 / Get command task
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_command_task",
		Description: "获取异步命令任务状态 / Get async command task status",
		InputSchema: types.GetToolSchema("get_command_task"),
	}, s.wrapGetCommandTask)

	// 取消命令任务 / Cancel command task
	registry.RegisterTool(&mcp.Tool{
		Name:        "cancel_command_task",
		Description: "取消正在执行的命令任务 / Cancel running command task",
		InputSchema: types.GetToolSchema("cancel_command_task"),
	}, s.wrapCancelCommandTask)

	// 获取命令历史 / Get command history
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_command_history",
		Description: "获取命令执行历史记录 / Get command execution history",
		InputSchema: types.GetToolSchema("get_command_history"),
	}, s.wrapGetCommandHistory)

	// 清空命令历史 / Clear command history
	registry.RegisterTool(&mcp.Tool{
		Name:        "clear_command_history",
		Description: "清空命令执行历史记录 / Clear command execution history",
		InputSchema: types.GetToolSchema("clear_command_history"),
	}, s.wrapClearCommandHistory)

	// 设置权限级别 / Set permission level
	registry.RegisterTool(&mcp.Tool{
		Name:        "set_permission_level",
		Description: "设置命令执行权限级别 / Set command execution permission level",
		InputSchema: types.GetToolSchema("set_permission_level"),
	}, s.wrapSetPermissionLevel)

	// 获取权限级别 / Get permission level
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_permission_level",
		Description: "获取当前权限级别 / Get current permission level",
		InputSchema: types.GetToolSchema("get_permission_level"),
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

func (s *Service) wrapDeleteFile(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.DeleteFileRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleDeleteFile(ctx, nil, args)
	return result, err
}

func (s *Service) wrapDeleteDirectory(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.DeleteDirectoryRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleDeleteDirectory(ctx, nil, args)
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

func (s *Service) wrapCopyFile(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.CopyFileRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleCopyFile(ctx, nil, args)
	return result, err
}

func (s *Service) wrapCopyDirectory(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.CopyDirectoryRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleCopyDirectory(ctx, nil, args)
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

func (s *Service) wrapMoveFile(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.MoveFileRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleMoveFile(ctx, nil, args)
	return result, err
}

func (s *Service) wrapMoveDirectory(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.MoveDirectoryRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleMoveDirectory(ctx, nil, args)
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
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.GetCurrentTimeRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleGetCurrentTime(ctx, nil, args)
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
