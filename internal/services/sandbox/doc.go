// Package sandbox 实现了 MCP Toolkit 的沙箱服务
//
// 本包提供了安全的文件系统操作、命令执行、系统信息获取等功能，
// 所有操作都在沙箱目录内进行。
//
// # 主要功能
//
// 文件操作：
//   - 创建文件（create_file）
//   - 读取文件（read_file）
//   - 写入文件（write_file）
//   - 删除文件/目录（delete）
//   - 复制文件/目录（copy）
//   - 移动/重命名（move）
//   - 获取文件信息（get_file_info）
//   - 检查文件存在（file_exists）
//
// 目录操作：
//   - 创建目录（create_directory）
//   - 列出目录（list_directory）
//   - 搜索文件（search_files）
//   - 获取工作目录（get_working_directory）
//   - 切换工作目录（change_directory）
//
// 命令执行：
//   - 同步执行命令（execute_command）
//   - 异步执行命令（execute_command_async）
//   - 获取命令任务（get_command_task）
//   - 取消命令任务（cancel_command_task）
//   - 命令黑名单管理（get_command_blacklist、update_command_blacklist）
//
// 系统功能：
//   - 获取当前时间（get_current_time）
//   - 权限级别管理（get_permission_level、set_permission_level）
//
// # 核心组件
//
// Service：
//   - 文件系统服务的主要实现
//   - 管理沙箱目录
//   - 提供所有文件系统操作
//   - 集成日志记录
//
// # 安全机制
//
// 沙箱隔离：
//   - 所有操作限制在沙箱目录内
//   - 路径遍历攻击防护（..）
//   - 绝对路径验证
//   - 符号链接检查
//
// 参数验证：
//   - Nil 检查
//   - 空字符串检查
//   - 路径长度限制（4096 字符）
//   - 文件大小限制（100MB）
//   - 批量操作限制（1000 个文件）
//
// 命令执行安全：
//   - 命令黑名单机制
//   - 目录黑名单机制
//   - 工作目录限制
//   - 超时控制
//   - 权限级别控制
//
// # 使用示例
//
// 创建服务：
//
//	logger := zap.NewNop()
//	service, err := filesystem.NewService("/path/to/sandbox", logger)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// 注册工具到 MCP 服务器：
//
//	service.RegisterTools(mcpServer)
//
// 注册工具到工具注册表：
//
//	service.RegisterToolsToRegistry(toolRegistry)
//
// 创建文件：
//
//	req := types.CreateFileRequest{
//	    Path:    "test.txt",
//	    Content: "Hello, World!",
//	}
//	resp, err := service.CreateFile(ctx, req)
//
// 读取文件：
//
//	req := types.ReadFileRequest{
//	    Path: "test.txt",
//	}
//	resp, err := service.ReadFile(ctx, req)
//	fmt.Println(resp.Content)
//
// 执行命令：
//
//	req := types.ExecuteCommandRequest{
//	    Command: "ls -la",
//	    Timeout: 30,
//	}
//	resp, err := service.ExecuteCommand(ctx, req)
//	fmt.Println(resp.Output)
//
// # 文件结构
//
//   - service.go：核心服务实现
//   - validation.go：参数验证逻辑
//   - constants.go：常量定义
//   - mcp_tools.go：MCP 工具注册
//   - command.go：命令执行功能
//   - command_async.go：异步命令执行
//   - command_blacklist.go：命令黑名单管理
//   - permission.go：权限级别管理
//
// # 常量定义
//
// 文件权限：
//   - DefaultDirPerm：默认目录权限（0755）
//   - DefaultFilePerm：默认文件权限（0644）
//
// 限制值：
//   - MaxFileSize：最大文件大小（100MB）
//   - MaxPathLength：最大路径长度（4096）
//   - MaxBatchDeleteCount：最大批量删除数量（1000）
//   - DefaultCommandTimeout：默认命令超时（30秒）
//   - MaxCommandTimeout：最大命令超时（300秒）
//
// # 错误处理
//
// 服务使用标准的 error 处理：
//
//   - 参数验证错误
//   - 文件系统错误
//   - 权限错误
//   - 命令执行错误
//
// 所有错误都包含详细的上下文信息。
//
// # 测试
//
// 本包包含完整的单元测试：
//
//   - service_test.go：服务功能测试
//   - mcp_tools_test.go：MCP 工具测试
//
// 运行测试：
//
//	go test -v ./internal/services/filesystem/...
//
// 测试覆盖率：
//
//	go test -cover ./internal/services/filesystem/...
//
// # 性能
//
//   - 路径缓存优化
//   - 最小化文件系统调用
//   - 批量操作支持
//   - 异步命令执行
//
// # 日志
//
// 服务使用 zap 进行结构化日志记录：
//
//   - 操作日志（Info 级别）
//   - 错误日志（Error 级别）
//   - 调试日志（Debug 级别）
//
// # 线程安全
//
// Service 是线程安全的，可以在多个 goroutine 中并发使用。
package sandbox
