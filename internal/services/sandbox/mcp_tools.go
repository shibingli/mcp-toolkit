package sandbox

import (
	"context"

	"mcp-toolkit/pkg/transport"
	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all filesystem tools to MCP server / 注册所有文件系统工具到MCP服务器
func (s *Service) RegisterTools(mcpServer *mcp.Server) {
	// ==================== File Operation Tools / 文件操作工具 ====================

	// Create file / 创建文件
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "create_file",
		Description: "CREATE A NEW FILE with specified content. Primary tool for file creation - use whenever you need to create or completely replace a file's content. Automatically creates parent directories. Keywords: create, new file, write file, save file, make file. / 创建新文件并写入内容。这是文件创建的主要工具。自动创建父目录。关键词：创建、新建文件、写入文件、保存文件。",
		InputSchema: types.GetToolSchema("create_file"),
	}, s.handleCreateFile)

	// Create directory / 创建目录
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "create_directory",
		Description: "CREATE NEW DIRECTORIES and folder structures. Automatically creates all parent directories (like 'mkdir -p'). Use for creating folders, setting up directory structures, organizing project layout. Keywords: create directory, mkdir, make folder, new folder. / 创建新目录和文件夹结构。自动创建所有父目录（类似 'mkdir -p'）。关键词：创建目录、新建文件夹。",
		InputSchema: types.GetToolSchema("create_directory"),
	}, s.handleCreateDir)

	// Read file / 读取文件
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "read_file",
		Description: "READ AND RETRIEVE complete file content. Primary tool for reading files - use whenever you need to see what's inside a file. Use for viewing, inspecting, checking file contents before modification. Keywords: read, view, show, display, get content, inspect file, open file. / 读取并获取完整文件内容。这是读取文件的主要工具。关键词：读取、查看、显示、获取内容、检查文件。",
		InputSchema: types.GetToolSchema("read_file"),
	}, s.handleReadFile)

	// Write file / 写入文件
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "write_file",
		Description: "WRITE OR UPDATE content to existing file, completely replacing current content. Workflow: read file first, modify content, then write back. For new files use 'create_file'. Keywords: write, update, modify, save, edit file, change file. / 写入或更新现有文件内容，完全替换当前内容。工作流程：先读取文件，修改内容，然后写回。新文件请使用 'create_file'。关键词：写入、更新、修改、保存、编辑文件。",
		InputSchema: types.GetToolSchema("write_file"),
	}, s.handleWriteFile)

	// Delete file or directory (auto-detect) / 删除文件或目录（自动判断）
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "delete",
		Description: "[RECOMMENDED - PRIMARY DELETION TOOL] DELETE any file or directory automatically. MAIN deletion tool - use for ALL deletion operations. Intelligently detects file/directory type and handles appropriately. Performs recursive deletion for directories. Keywords: delete, remove, erase, clean, rm, unlink. / [推荐 - 主要删除工具] 自动删除任何文件或目录。这是主要的删除工具 - 用于所有删除操作。智能检测文件/目录类型并适当处理。对目录执行递归删除。关键词：删除、移除、清理。",
		InputSchema: types.GetToolSchema("delete"),
	}, s.handleDelete)

	// Delete file only / 删除文件（仅限文件）
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "delete_file",
		Description: "Delete a specific file. Only works on files, not directories. Use 'delete' tool if you're unsure whether the path is a file or directory. / 删除指定的文件。仅用于删除文件，不能删除目录。如果不确定是文件还是目录，请使用 delete 工具。",
		InputSchema: types.GetToolSchema("delete_file"),
	}, s.handleDeleteFile)

	// Delete directory only / 删除目录（仅限目录）
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "delete_directory",
		Description: "Delete a directory and optionally all its contents. Only works on directories, not files. Use 'delete' tool if you're unsure whether the path is a file or directory. / 删除目录及其所有内容。仅用于删除目录，不能删除文件。如果不确定是文件还是目录，请使用 delete 工具。",
		InputSchema: types.GetToolSchema("delete_directory"),
	}, s.handleDeleteDirectory)

	// Copy file or directory (auto-detect) / 复制文件或目录（自动判断）
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "copy",
		Description: "[RECOMMENDED - PRIMARY COPY TOOL] COPY files or directories to new location. MAIN copy tool - use for ALL copy operations. Automatically detects file/directory type and handles appropriately. Performs recursive copy for directories. Keywords: copy, duplicate, clone, backup, cp. / [推荐 - 主要复制工具] 复制文件或目录到新位置。这是主要的复制工具 - 用于所有复制操作。自动检测文件/目录类型。对目录执行递归复制。关键词：复制、备份、克隆。",
		InputSchema: types.GetToolSchema("copy"),
	}, s.handleCopy)

	// Copy file only / 复制文件（仅限文件）
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "copy_file",
		Description: "Copy specific file to new location. Only works on files, not directories. Use 'copy' tool if unsure about path type. Keywords: copy file, duplicate file. / 复制指定文件到新位置。仅用于复制文件，不能复制目录。如果不确定路径类型，请使用 copy 工具。关键词：复制文件、备份文件。",
		InputSchema: types.GetToolSchema("copy_file"),
	}, s.handleCopyFile)

	// Copy directory only / 复制目录（仅限目录）
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "copy_directory",
		Description: "Copy directory and all contents to new location. Only works on directories, not files. Use 'copy' tool if unsure about path type. Keywords: copy directory, duplicate folder. / 复制目录及其所有内容到新位置。仅用于复制目录，不能复制文件。如果不确定路径类型，请使用 copy 工具。关键词：复制目录、备份文件夹。",
		InputSchema: types.GetToolSchema("copy_directory"),
	}, s.handleCopyDirectory)

	// Move file or directory (auto-detect) / 移动文件或目录（自动判断）
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "move",
		Description: "[RECOMMENDED - PRIMARY MOVE/RENAME TOOL] MOVE or RENAME files and directories. MAIN move/rename tool - use for ALL move and rename operations. Automatically detects file/directory type and handles appropriately. Keywords: move, rename, mv, relocate. / [推荐 - 主要移动/重命名工具] 移动或重命名文件和目录。这是主要的移动/重命名工具 - 用于所有移动和重命名操作。自动检测文件/目录类型。关键词：移动、重命名、重定位。",
		InputSchema: types.GetToolSchema("move"),
	}, s.handleMove)

	// Move file only / 移动文件（仅限文件）
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "move_file",
		Description: "Move or rename specific file. Only works on files, not directories. Use 'move' tool if unsure about path type. Keywords: move file, rename file. / 移动或重命名指定文件。仅用于移动文件，不能移动目录。如果不确定路径类型，请使用 move 工具。关键词：移动文件、重命名文件。",
		InputSchema: types.GetToolSchema("move_file"),
	}, s.handleMoveFile)

	// Move directory only / 移动目录（仅限目录）
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "move_directory",
		Description: "Move or rename directory. Only works on directories, not files. Use 'move' tool if unsure about path type. Keywords: move directory, rename folder. / 移动或重命名目录。仅用于移动目录，不能移动文件。如果不确定路径类型，请使用 move 工具。关键词：移动目录、重命名文件夹。",
		InputSchema: types.GetToolSchema("move_directory"),
	}, s.handleMoveDirectory)

	// List directory / 列出目录
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_directory",
		Description: "LIST AND EXPLORE directory contents. Shows all files and subdirectories with detailed information. Primary tool for directory exploration. Use for seeing what files are in a directory, exploring project structure. Keywords: list, ls, dir, show files, browse. / 列出和探索目录内容。显示所有文件和子目录的详细信息。这是目录探索的主要工具。关键词：列出、显示文件、浏览。",
		InputSchema: types.GetToolSchema("list_directory"),
	}, s.handleListDir)

	// Search files / 搜索文件
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "search_files",
		Description: "SEARCH AND FIND files matching patterns. Recursively searches through all subdirectories using glob patterns. Use for finding files by extension, locating specific files by name pattern. Keywords: search, find, locate, grep files, filter. / 搜索和查找匹配模式的文件。使用glob模式递归搜索所有子目录。用于按扩展名查找文件、按名称模式定位特定文件。关键词：搜索、查找、定位、过滤。",
		InputSchema: types.GetToolSchema("search_files"),
	}, s.handleSearch)

	// Batch delete / 批量删除
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "batch_delete",
		Description: "Delete multiple files or directories in a single operation. Each path is processed independently, and the tool will report success/failure for each item. / 批量删除多个文件或目录。每个路径独立处理，工具会报告每个项目的成功/失败状态。",
		InputSchema: types.GetToolSchema("batch_delete"),
	}, s.handleBatchDelete)

	// File stat / 文件状态
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "file_stat",
		Description: "Get detailed information about a file or directory, including size, permissions, modification time, and type (file/directory/symlink). / 获取文件或目录的详细信息，包括大小、权限、修改时间和类型（文件/目录/符号链接）。",
		InputSchema: types.GetToolSchema("file_stat"),
	}, s.handleFileStat)

	// File exists / 文件是否存在
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "file_exists",
		Description: "Check if a file or directory exists at the specified path. Returns true if exists, false otherwise. Useful for conditional operations. / 检查指定路径的文件或目录是否存在。存在返回 true，否则返回 false。用于条件操作。",
		InputSchema: types.GetToolSchema("file_exists"),
	}, s.handleFileExists)

	// ==================== Time Tools / 时间工具 ====================

	// Get current time / 获取当前时间
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_current_time",
		Description: "GET CURRENT DATE AND TIME with timezone support. Returns formatted datetime, timezone info, and Unix timestamp. Use for getting current time, checking time in different timezones, timestamps for logging. Keywords: time, date, now, current time, timestamp, timezone, clock. / 获取当前日期和时间，支持时区。返回格式化的日期时间、时区信息和Unix时间戳。用于获取当前时间、检查不同时区的时间、日志时间戳。关键词：时间、日期、当前时间、时间戳、时区。",
		InputSchema: types.GetToolSchema("get_current_time"),
	}, s.handleGetCurrentTime)

	// ==================== System Info Tools / 系统信息工具 ====================

	// Get system info / 获取系统信息
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_system_info",
		Description: "GET COMPREHENSIVE SYSTEM INFORMATION about current machine. Returns OS details, CPU, memory, GPU, network interfaces. Use for checking system specs, verifying hardware capabilities, getting OS info, checking available resources. Keywords: system info, hardware, specs, os info, cpu, memory, ram, gpu. / 获取当前机器的全面系统信息。返回操作系统详情、CPU、内存、显卡、网卡信息。用于检查系统规格、验证硬件能力、获取操作系统信息、检查可用资源。关键词：系统信息、硬件、规格、操作系统信息、CPU、内存、显卡。",
		InputSchema: types.GetToolSchema("get_system_info"),
	}, s.handleGetSystemInfo)

	// ==================== Command Execution Tools / 命令执行工具 ====================

	// Execute command / 执行命令
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "execute_command",
		Description: "EXECUTE SHELL COMMANDS synchronously and get output. Use for running CLI tools, build scripts, version control, package managers. Common uses: git commands, npm/pip/go commands, Python/Node.js scripts, build tools. IMPORTANT: Do NOT use for file operations - use dedicated file tools instead. Keywords: execute, run, command, shell, bash, terminal, cli, script, git, npm, python. / 同步执行shell命令并获取输出。用于运行CLI工具、构建脚本、版本控制、包管理器。常见用途：git命令、npm/pip/go命令、Python/Node.js脚本、构建工具。重要：不要用于文件操作 - 请使用专用的文件工具。关键词：执行、运行、命令、shell、终端、脚本。",
		InputSchema: types.GetToolSchema("execute_command"),
	}, s.handleExecuteCommand)

	// Get command blacklist / 获取命令黑名单
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_command_blacklist",
		Description: "Get the current command and directory blacklist. Returns lists of blocked commands and directories that cannot be executed or accessed. / 获取当前命令和目录黑名单。返回被阻止执行或访问的命令和目录列表。",
		InputSchema: types.GetToolSchema("get_command_blacklist"),
	}, s.handleGetCommandBlacklist)

	// Update command blacklist / 更新命令黑名单
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "update_command_blacklist",
		Description: "Update the command and directory blacklist. Add commands or directories that should be blocked from execution or access. / 更新命令和目录黑名单。添加应被阻止执行或访问的命令或目录。",
		InputSchema: types.GetToolSchema("update_command_blacklist"),
	}, s.handleUpdateCommandBlacklist)

	// Get working directory / 获取当前工作目录
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_working_directory",
		Description: "Get the current working directory path. Returns the absolute path of the current working directory. / 获取当前工作目录路径。返回当前工作目录的绝对路径。",
		InputSchema: types.GetToolSchema("get_working_directory"),
	}, s.handleGetWorkingDirectory)

	// Change directory / 切换工作目录
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "change_directory",
		Description: "Change the current working directory. Similar to 'cd' command. The new directory must exist and be within the sandbox. / 切换当前工作目录。类似 'cd' 命令。新目录必须存在且在沙箱范围内。",
		InputSchema: types.GetToolSchema("change_directory"),
	}, s.handleChangeDirectory)

	// Execute command asynchronously / 异步执行命令
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "execute_command_async",
		Description: "Execute a command asynchronously in the background. Returns a task ID immediately that can be used to check status, get output, or cancel the command. Use this for long-running commands. / 在后台异步执行命令。立即返回任务ID，可用于检查状态、获取输出或取消命令。用于长时间运行的命令。",
		InputSchema: types.GetToolSchema("execute_command_async"),
	}, s.handleExecuteCommandAsync)

	// Get command task / 获取命令任务
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_command_task",
		Description: "Get detailed information about a specific asynchronous command task, including its status, output, start time, and duration. / 获取特定异步命令任务的详细信息，包括状态、输出、开始时间和持续时间。",
		InputSchema: types.GetToolSchema("get_command_task"),
	}, s.handleGetCommandTask)

	// Cancel command task / 取消命令任务
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "cancel_command_task",
		Description: "Cancel a running asynchronous command task. The task will be terminated if it's still running. / 取消正在运行的异步命令任务。如果任务仍在运行，将被终止。",
		InputSchema: types.GetToolSchema("cancel_command_task"),
	}, s.handleCancelCommandTask)

	// Get command history / 获取命令历史
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_command_history",
		Description: "Get the history of executed commands. Returns a list of previously executed commands with their results. / 获取命令执行历史记录。返回之前执行的命令及其结果列表。",
		InputSchema: types.GetToolSchema("get_command_history"),
	}, s.handleGetCommandHistory)

	// Clear command history / 清空命令历史
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "clear_command_history",
		Description: "Clear all command execution history records. This action cannot be undone. / 清空所有命令执行历史记录。此操作不可撤销。",
		InputSchema: types.GetToolSchema("clear_command_history"),
	}, s.handleClearCommandHistory)

	// Set permission level / 设置权限级别
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "set_permission_level",
		Description: "Set the command execution permission level. Higher levels allow more privileged operations. Level 0 is most restrictive, level 3 is least restrictive. / 设置命令执行权限级别。级别越高允许的操作越多。级别0最严格，级别3最宽松。",
		InputSchema: types.GetToolSchema("set_permission_level"),
	}, s.handleSetPermissionLevel)

	// Get permission level / 获取权限级别
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_permission_level",
		Description: "Get the current command execution permission level. Returns the current level (0-3) and its description. / 获取当前命令执行权限级别。返回当前级别（0-3）及其描述。",
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

// handleGetSystemInfo 处理获取系统信息请求 / Handle get system info request
func (s *Service) handleGetSystemInfo(_ context.Context, _ *mcp.CallToolRequest, args types.GetSystemInfoRequest) (*mcp.CallToolResult, *types.GetSystemInfoResponse, error) {
	resp, err := s.GetSystemInfo(&args)
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
	// ==================== File Operation Tools / 文件操作工具 ====================

	// Create file / 创建文件
	registry.RegisterTool(&mcp.Tool{
		Name:        "create_file",
		Description: "Create a new file with the specified content. If the file already exists, it will be overwritten. Parent directories will be created automatically if they don't exist. / 创建新文件并写入内容。如果文件已存在则覆盖。父目录不存在时会自动创建。",
		InputSchema: types.GetToolSchema("create_file"),
	}, s.wrapCreateFile)

	// Create directory / 创建目录
	registry.RegisterTool(&mcp.Tool{
		Name:        "create_directory",
		Description: "Create a new directory. Parent directories will be created automatically if they don't exist (similar to 'mkdir -p'). / 创建新目录。父目录不存在时会自动创建（类似 'mkdir -p'）。",
		InputSchema: types.GetToolSchema("create_directory"),
	}, s.wrapCreateDir)

	// Read file / 读取文件
	registry.RegisterTool(&mcp.Tool{
		Name:        "read_file",
		Description: "Read and return the content of a file. Use this to view file contents before making modifications or to understand existing code/configuration. / 读取并返回文件内容。用于在修改前查看文件内容或理解现有代码/配置。",
		InputSchema: types.GetToolSchema("read_file"),
	}, s.wrapReadFile)

	// Write file / 写入文件
	registry.RegisterTool(&mcp.Tool{
		Name:        "write_file",
		Description: "Write content to an existing file, completely replacing its current content. Use 'create_file' for new files. For partial modifications, read the file first, modify the content, then write back. / 写入内容到现有文件，完全替换当前内容。新文件请使用 'create_file'。部分修改请先读取文件，修改后再写回。",
		InputSchema: types.GetToolSchema("write_file"),
	}, s.wrapWriteFile)

	// Delete file or directory (auto-detect) / 删除文件或目录（自动判断）
	registry.RegisterTool(&mcp.Tool{
		Name:        "delete",
		Description: "[RECOMMENDED] Delete a file or directory. This tool automatically detects whether the path is a file or directory and handles it appropriately. For directories, it performs recursive deletion. This is the preferred tool for all deletion operations. / [推荐] 删除文件或目录，自动判断类型。对于目录会递归删除。这是所有删除操作的首选工具。",
		InputSchema: types.GetToolSchema("delete"),
	}, s.wrapDelete)

	// Delete file only / 删除文件（仅限文件）
	registry.RegisterTool(&mcp.Tool{
		Name:        "delete_file",
		Description: "Delete a specific file. Only works on files, not directories. Use 'delete' tool if you're unsure whether the path is a file or directory. / 删除指定的文件。仅用于删除文件，不能删除目录。如果不确定是文件还是目录，请使用 delete 工具。",
		InputSchema: types.GetToolSchema("delete_file"),
	}, s.wrapDeleteFile)

	// Delete directory only / 删除目录（仅限目录）
	registry.RegisterTool(&mcp.Tool{
		Name:        "delete_directory",
		Description: "Delete a directory and optionally all its contents. Only works on directories, not files. Use 'delete' tool if you're unsure whether the path is a file or directory. / 删除目录及其所有内容。仅用于删除目录，不能删除文件。如果不确定是文件还是目录，请使用 delete 工具。",
		InputSchema: types.GetToolSchema("delete_directory"),
	}, s.wrapDeleteDirectory)

	// Copy file or directory (auto-detect) / 复制文件或目录（自动判断）
	registry.RegisterTool(&mcp.Tool{
		Name:        "copy",
		Description: "[RECOMMENDED] Copy a file or directory to a new location. This tool automatically detects whether the source is a file or directory and handles it appropriately. For directories, it performs recursive copy. This is the preferred tool for all copy operations. / [推荐] 复制文件或目录到新位置，自动判断类型。对于目录会递归复制。这是所有复制操作的首选工具。",
		InputSchema: types.GetToolSchema("copy"),
	}, s.wrapCopy)

	// Copy file only / 复制文件（仅限文件）
	registry.RegisterTool(&mcp.Tool{
		Name:        "copy_file",
		Description: "Copy a specific file to a new location. Only works on files, not directories. Use 'copy' tool if you're unsure whether the path is a file or directory. / 复制指定的文件到新位置。仅用于复制文件，不能复制目录。如果不确定是文件还是目录，请使用 copy 工具。",
		InputSchema: types.GetToolSchema("copy_file"),
	}, s.wrapCopyFile)

	// Copy directory only / 复制目录（仅限目录）
	registry.RegisterTool(&mcp.Tool{
		Name:        "copy_directory",
		Description: "Copy a directory and all its contents to a new location. Only works on directories, not files. Use 'copy' tool if you're unsure whether the path is a file or directory. / 复制目录及其所有内容到新位置。仅用于复制目录，不能复制文件。如果不确定是文件还是目录，请使用 copy 工具。",
		InputSchema: types.GetToolSchema("copy_directory"),
	}, s.wrapCopyDirectory)

	// Move file or directory (auto-detect) / 移动文件或目录（自动判断）
	registry.RegisterTool(&mcp.Tool{
		Name:        "move",
		Description: "[RECOMMENDED] Move or rename a file or directory. This tool automatically detects whether the source is a file or directory and handles it appropriately. This is the preferred tool for all move/rename operations. / [推荐] 移动或重命名文件或目录，自动判断类型。这是所有移动/重命名操作的首选工具。",
		InputSchema: types.GetToolSchema("move"),
	}, s.wrapMove)

	// Move file only / 移动文件（仅限文件）
	registry.RegisterTool(&mcp.Tool{
		Name:        "move_file",
		Description: "Move or rename a specific file. Only works on files, not directories. Use 'move' tool if you're unsure whether the path is a file or directory. / 移动或重命名指定的文件。仅用于移动文件，不能移动目录。如果不确定是文件还是目录，请使用 move 工具。",
		InputSchema: types.GetToolSchema("move_file"),
	}, s.wrapMoveFile)

	// Move directory only / 移动目录（仅限目录）
	registry.RegisterTool(&mcp.Tool{
		Name:        "move_directory",
		Description: "Move or rename a directory. Only works on directories, not files. Use 'move' tool if you're unsure whether the path is a file or directory. / 移动或重命名指定的目录。仅用于移动目录，不能移动文件。如果不确定是文件还是目录，请使用 move 工具。",
		InputSchema: types.GetToolSchema("move_directory"),
	}, s.wrapMoveDirectory)

	// List directory / 列出目录
	registry.RegisterTool(&mcp.Tool{
		Name:        "list_directory",
		Description: "List all files and subdirectories in a directory. Returns file names, types (file/directory), sizes, and modification times. Useful for exploring directory structure. / 列出目录中的所有文件和子目录。返回文件名、类型（文件/目录）、大小和修改时间。用于探索目录结构。",
		InputSchema: types.GetToolSchema("list_directory"),
	}, s.wrapListDir)

	// Search files / 搜索文件
	registry.RegisterTool(&mcp.Tool{
		Name:        "search_files",
		Description: "Search for files matching a pattern within a directory. Supports glob patterns with wildcards (* for any characters, ? for single character). Searches recursively through subdirectories. / 在目录中搜索匹配模式的文件。支持通配符（* 匹配任意字符，? 匹配单个字符）。递归搜索子目录。",
		InputSchema: types.GetToolSchema("search_files"),
	}, s.wrapSearch)

	// Batch delete / 批量删除
	registry.RegisterTool(&mcp.Tool{
		Name:        "batch_delete",
		Description: "Delete multiple files or directories in a single operation. Each path is processed independently, and the tool will report success/failure for each item. / 批量删除多个文件或目录。每个路径独立处理，工具会报告每个项目的成功/失败状态。",
		InputSchema: types.GetToolSchema("batch_delete"),
	}, s.wrapBatchDelete)

	// File stat / 文件状态
	registry.RegisterTool(&mcp.Tool{
		Name:        "file_stat",
		Description: "Get detailed information about a file or directory, including size, permissions, modification time, and type (file/directory/symlink). / 获取文件或目录的详细信息，包括大小、权限、修改时间和类型（文件/目录/符号链接）。",
		InputSchema: types.GetToolSchema("file_stat"),
	}, s.wrapFileStat)

	// File exists / 文件是否存在
	registry.RegisterTool(&mcp.Tool{
		Name:        "file_exists",
		Description: "Check if a file or directory exists at the specified path. Returns true if exists, false otherwise. Useful for conditional operations. / 检查指定路径的文件或目录是否存在。存在返回 true，否则返回 false。用于条件操作。",
		InputSchema: types.GetToolSchema("file_exists"),
	}, s.wrapFileExists)

	// ==================== Time Tools / 时间工具 ====================

	// Get current time / 获取当前时间
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_current_time",
		Description: "Get the current system time. Can optionally specify a timezone (e.g. 'Asia/Shanghai', 'America/New_York'). Returns formatted datetime, timezone information, and Unix timestamp. / 获取当前系统时间。可指定时区（如 'Asia/Shanghai'、'America/New_York'），不指定则使用系统本地时区。返回格式化的日期时间、时区信息和Unix时间戳。",
		InputSchema: types.GetToolSchema("get_current_time"),
	}, s.wrapGetCurrentTime)

	// ==================== System Info Tools / 系统信息工具 ====================

	// Get system info / 获取系统信息
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_system_info",
		Description: "Get comprehensive system information including OS details, CPU, memory, GPU, and network interfaces. Returns detailed hardware and software information about the current system. / 获取全面的系统信息，包括操作系统详情、CPU、内存、显卡和网卡信息。返回当前系统的详细硬件和软件信息。",
		InputSchema: types.GetToolSchema("get_system_info"),
	}, s.wrapGetSystemInfo)

	// ==================== Command Execution Tools / 命令执行工具 ====================

	// Execute command / 执行命令
	registry.RegisterTool(&mcp.Tool{
		Name:        "execute_command",
		Description: "Execute a shell command synchronously and return the output. Use this for running CLI tools like git, npm, python, etc. WARNING: Do NOT use this for file operations (create, delete, copy, move) - use the dedicated file tools instead. / 同步执行shell命令并返回输出。用于运行 git、npm、python 等CLI工具。警告：不要用于文件操作（创建、删除、复制、移动），请使用专门的文件工具。",
		InputSchema: types.GetToolSchema("execute_command"),
	}, s.wrapExecuteCommand)

	// Get command blacklist / 获取命令黑名单
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_command_blacklist",
		Description: "Get the current command and directory blacklist. Returns lists of blocked commands and directories that cannot be executed or accessed. / 获取当前命令和目录黑名单。返回被阻止执行或访问的命令和目录列表。",
		InputSchema: types.GetToolSchema("get_command_blacklist"),
	}, s.wrapGetCommandBlacklist)

	// Update command blacklist / 更新命令黑名单
	registry.RegisterTool(&mcp.Tool{
		Name:        "update_command_blacklist",
		Description: "Update the command and directory blacklist. Add commands or directories that should be blocked from execution or access. / 更新命令和目录黑名单。添加应被阻止执行或访问的命令或目录。",
		InputSchema: types.GetToolSchema("update_command_blacklist"),
	}, s.wrapUpdateCommandBlacklist)

	// Get working directory / 获取当前工作目录
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_working_directory",
		Description: "Get the current working directory path. Returns the absolute path of the current working directory. / 获取当前工作目录路径。返回当前工作目录的绝对路径。",
		InputSchema: types.GetToolSchema("get_working_directory"),
	}, s.wrapGetWorkingDirectory)

	// Change directory / 切换工作目录
	registry.RegisterTool(&mcp.Tool{
		Name:        "change_directory",
		Description: "Change the current working directory. Similar to 'cd' command. The new directory must exist and be within the sandbox. / 切换当前工作目录。类似 'cd' 命令。新目录必须存在且在沙箱范围内。",
		InputSchema: types.GetToolSchema("change_directory"),
	}, s.wrapChangeDirectory)

	// Execute command asynchronously / 异步执行命令
	registry.RegisterTool(&mcp.Tool{
		Name:        "execute_command_async",
		Description: "Execute a command asynchronously in the background. Returns a task ID immediately that can be used to check status, get output, or cancel the command. Use this for long-running commands. / 在后台异步执行命令。立即返回任务ID，可用于检查状态、获取输出或取消命令。用于长时间运行的命令。",
		InputSchema: types.GetToolSchema("execute_command_async"),
	}, s.wrapExecuteCommandAsync)

	// Get command task / 获取命令任务
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_command_task",
		Description: "Get detailed information about a specific asynchronous command task, including its status, output, start time, and duration. / 获取特定异步命令任务的详细信息，包括状态、输出、开始时间和持续时间。",
		InputSchema: types.GetToolSchema("get_command_task"),
	}, s.wrapGetCommandTask)

	// Cancel command task / 取消命令任务
	registry.RegisterTool(&mcp.Tool{
		Name:        "cancel_command_task",
		Description: "Cancel a running asynchronous command task. The task will be terminated if it's still running. / 取消正在运行的异步命令任务。如果任务仍在运行，将被终止。",
		InputSchema: types.GetToolSchema("cancel_command_task"),
	}, s.wrapCancelCommandTask)

	// Get command history / 获取命令历史
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_command_history",
		Description: "Get the history of executed commands. Returns a list of previously executed commands with their results. / 获取命令执行历史记录。返回之前执行的命令及其结果列表。",
		InputSchema: types.GetToolSchema("get_command_history"),
	}, s.wrapGetCommandHistory)

	// Clear command history / 清空命令历史
	registry.RegisterTool(&mcp.Tool{
		Name:        "clear_command_history",
		Description: "Clear all command execution history records. This action cannot be undone. / 清空所有命令执行历史记录。此操作不可撤销。",
		InputSchema: types.GetToolSchema("clear_command_history"),
	}, s.wrapClearCommandHistory)

	// Set permission level / 设置权限级别
	registry.RegisterTool(&mcp.Tool{
		Name:        "set_permission_level",
		Description: "Set the command execution permission level. Higher levels allow more privileged operations. Level 0 is most restrictive, level 3 is least restrictive. / 设置命令执行权限级别。级别越高允许的操作越多。级别0最严格，级别3最宽松。",
		InputSchema: types.GetToolSchema("set_permission_level"),
	}, s.wrapSetPermissionLevel)

	// Get permission level / 获取权限级别
	registry.RegisterTool(&mcp.Tool{
		Name:        "get_permission_level",
		Description: "Get the current command execution permission level. Returns the current level (0-3) and its description. / 获取当前命令执行权限级别。返回当前级别（0-3）及其描述。",
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

func (s *Service) wrapGetSystemInfo(ctx context.Context, arguments interface{}) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	var args types.GetSystemInfoRequest
	if err = json.Unmarshal(argsJSON, &args); err != nil {
		return nil, err
	}
	result, _, err := s.handleGetSystemInfo(ctx, nil, args)
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
