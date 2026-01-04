package main

import (
	"context"
	"fmt"
	"time"

	"mcp-toolkit/pkg/client"
	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"
)

// ComprehensiveTestSuite 全面测试套件 / Comprehensive test suite
type ComprehensiveTestSuite struct {
	client  client.Client
	ctx     context.Context
	results []TestResult
}

// NewComprehensiveTestSuite 创建全面测试套件 / Create comprehensive test suite
func NewComprehensiveTestSuite(mcpClient client.Client) *ComprehensiveTestSuite {
	return &ComprehensiveTestSuite{
		client:  mcpClient,
		ctx:     context.Background(),
		results: []TestResult{},
	}
}

// RunAllTests 运行所有测试 / Run all tests
func (s *ComprehensiveTestSuite) RunAllTests() []TestResult {
	fmt.Println("\n=== 开始全面测试 MCP Toolkit 所有功能 ===")

	// 1. 文件操作工具测试 / File operation tools tests
	s.testFileOperations()

	// 2. 目录操作工具测试 / Directory operation tools tests
	s.testDirectoryOperations()

	// 3. 文件信息工具测试 / File info tools tests
	s.testFileInfoTools()

	// 4. 系统信息工具测试 / System info tools tests
	s.testSystemInfoTools()

	// 5. 命令执行工具测试 / Command execution tools tests
	s.testCommandExecutionTools()

	// 6. 异步命令工具测试 / Async command tools tests
	s.testAsyncCommandTools()

	// 7. 命令管理工具测试 / Command management tools tests
	s.testCommandManagementTools()

	// 8. 权限管理工具测试 / Permission management tools tests
	s.testPermissionManagementTools()

	// 9. 时间工具测试 / Time tools tests
	s.testTimeTools()

	// 10. 下载工具测试 / Download tools tests
	s.testDownloadTools()

	// 11. 清理测试 / Cleanup tests
	s.testCleanup()

	return s.results
}

// testFileOperations 测试文件操作工具 / Test file operation tools
func (s *ComprehensiveTestSuite) testFileOperations() {
	fmt.Println("=== 1. 文件操作工具测试 ===")

	// create_file
	s.addResult(s.testCreateFile())

	// read_file
	s.addResult(s.testReadFile())

	// write_file
	s.addResult(s.testWriteFile())

	// copy (auto-detect)
	s.addResult(s.testCopy())

	// copy_file
	s.addResult(s.testCopyFile())

	// move (auto-detect)
	s.addResult(s.testMove())

	// move_file
	s.addResult(s.testMoveFile())

	// delete (auto-detect)
	s.addResult(s.testDelete())

	// delete_file
	s.addResult(s.testDeleteFile())
}

// testDirectoryOperations 测试目录操作工具 / Test directory operation tools
func (s *ComprehensiveTestSuite) testDirectoryOperations() {
	fmt.Println("\n=== 2. 目录操作工具测试 ===")

	// create_directory
	s.addResult(s.testCreateDirectory())

	// list_directory
	s.addResult(s.testListDirectory())

	// copy_directory
	s.addResult(s.testCopyDirectory())

	// move_directory
	s.addResult(s.testMoveDirectory())

	// delete_directory
	s.addResult(s.testDeleteDirectory())

	// batch_delete
	s.addResult(s.testBatchDelete())

	// get_working_directory
	s.addResult(s.testGetWorkingDirectory())

	// change_directory
	s.addResult(s.testChangeDirectory())
}

// testFileInfoTools 测试文件信息工具 / Test file info tools
func (s *ComprehensiveTestSuite) testFileInfoTools() {
	fmt.Println("\n=== 3. 文件信息工具测试 ===")

	// file_stat
	s.addResult(s.testFileStat())

	// file_exists
	s.addResult(s.testFileExists())

	// search_files
	s.addResult(s.testSearchFiles())
}

// testSystemInfoTools 测试系统信息工具 / Test system info tools
func (s *ComprehensiveTestSuite) testSystemInfoTools() {
	fmt.Println("\n=== 4. 系统信息工具测试 ===")

	// get_system_info
	s.addResult(s.testGetSystemInfo())
}

// testAsyncCommandTools 测试异步命令工具 / Test async command tools
func (s *ComprehensiveTestSuite) testAsyncCommandTools() {
	fmt.Println("\n=== 6. 异步命令工具测试 ===")

	// execute_command_async
	s.addResult(s.testExecuteCommandAsync())

	// get_command_task
	s.addResult(s.testGetCommandTask())

	// cancel_command_task
	s.addResult(s.testCancelCommandTask())
}

// testCommandManagementTools 测试命令管理工具 / Test command management tools
func (s *ComprehensiveTestSuite) testCommandManagementTools() {
	fmt.Println("\n=== 7. 命令管理工具测试 ===")

	// get_command_blacklist
	s.addResult(s.testGetCommandBlacklist())

	// update_command_blacklist
	s.addResult(s.testUpdateCommandBlacklist())

	// get_command_history
	s.addResult(s.testGetCommandHistory())

	// clear_command_history
	s.addResult(s.testClearCommandHistory())
}

// testPermissionManagementTools 测试权限管理工具 / Test permission management tools
func (s *ComprehensiveTestSuite) testPermissionManagementTools() {
	fmt.Println("\n=== 8. 权限管理工具测试 ===")

	// get_permission_level
	s.addResult(s.testGetPermissionLevel())

	// set_permission_level
	s.addResult(s.testSetPermissionLevel())
}

// testTimeTools 测试时间工具 / Test time tools
func (s *ComprehensiveTestSuite) testTimeTools() {
	fmt.Println("\n=== 9. 时间工具测试 ===")

	// get_current_time
	s.addResult(s.testGetCurrentTime())
}

// testDownloadTools 测试下载工具 / Test download tools
func (s *ComprehensiveTestSuite) testDownloadTools() {
	fmt.Println("\n=== 10. 下载工具测试 ===")

	// download_file
	s.addResult(s.testDownloadFile())
}

// testCleanup 清理测试 / Cleanup tests
func (s *ComprehensiveTestSuite) testCleanup() {
	fmt.Println("\n=== 11. 清理测试文件 ===")

	// 清理所有测试文件 / Clean up all test files
	s.addResult(s.testFinalCleanup())
}

// addResult 添加测试结果 / Add test result
func (s *ComprehensiveTestSuite) addResult(result TestResult) {
	s.results = append(s.results, result)
}

// ==================== 具体测试方法实现 / Specific test method implementations ====================

// testCreateFile 测试创建文件 / Test create file
func (s *ComprehensiveTestSuite) testCreateFile() TestResult {
	fmt.Println("  测试 create_file...")

	resp, err := s.client.CallTool(s.ctx, "create_file", types.CreateFileRequest{
		Path:    "test_file.txt",
		Content: "Hello, MCP Toolkit!",
	})

	if err != nil {
		return TestResult{Name: "create_file", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件创建成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "create_file", Success: true}
}

// testReadFile 测试读取文件 / Test read file
func (s *ComprehensiveTestSuite) testReadFile() TestResult {
	fmt.Println("  测试 read_file...")

	resp, err := s.client.CallTool(s.ctx, "read_file", types.ReadFileRequest{
		Path: "test_file.txt",
	})

	if err != nil {
		return TestResult{Name: "read_file", Success: false, Error: err.Error()}
	}

	var readResp types.ReadFileResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &readResp); err != nil {
		return TestResult{Name: "read_file", Success: false, Error: err.Error()}
	}

	if readResp.Content != "Hello, MCP Toolkit!" {
		return TestResult{Name: "read_file", Success: false, Error: "内容不匹配"}
	}

	fmt.Printf("    ✅ 文件读取成功，内容: %s\n", readResp.Content)
	return TestResult{Name: "read_file", Success: true}
}

// testWriteFile 测试写入文件 / Test write file
func (s *ComprehensiveTestSuite) testWriteFile() TestResult {
	fmt.Println("  测试 write_file...")

	resp, err := s.client.CallTool(s.ctx, "write_file", types.WriteFileRequest{
		Path:    "test_file.txt",
		Content: "Updated content by write_file",
	})

	if err != nil {
		return TestResult{Name: "write_file", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件写入成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "write_file", Success: true}
}

// testCopy 测试复制（自动检测）/ Test copy (auto-detect)
func (s *ComprehensiveTestSuite) testCopy() TestResult {
	fmt.Println("  测试 copy (auto-detect)...")

	resp, err := s.client.CallTool(s.ctx, "copy", types.CopyRequest{
		Source:      "test_file.txt",
		Destination: "test_file_copy.txt",
	})

	if err != nil {
		return TestResult{Name: "copy", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件复制成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "copy", Success: true}
}

// testCopyFile 测试复制文件 / Test copy file
func (s *ComprehensiveTestSuite) testCopyFile() TestResult {
	fmt.Println("  测试 copy_file...")

	resp, err := s.client.CallTool(s.ctx, "copy_file", types.CopyFileRequest{
		Source:      "test_file.txt",
		Destination: "test_file_copy2.txt",
	})

	if err != nil {
		return TestResult{Name: "copy_file", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件复制成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "copy_file", Success: true}
}

// testMove 测试移动（自动检测）/ Test move (auto-detect)
func (s *ComprehensiveTestSuite) testMove() TestResult {
	fmt.Println("  测试 move (auto-detect)...")

	resp, err := s.client.CallTool(s.ctx, "move", types.MoveRequest{
		Source:      "test_file_copy.txt",
		Destination: "test_file_moved.txt",
	})

	if err != nil {
		return TestResult{Name: "move", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件移动成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "move", Success: true}
}

// testMoveFile 测试移动文件 / Test move file
func (s *ComprehensiveTestSuite) testMoveFile() TestResult {
	fmt.Println("  测试 move_file...")

	resp, err := s.client.CallTool(s.ctx, "move_file", types.MoveFileRequest{
		Source:      "test_file_copy2.txt",
		Destination: "test_file_moved2.txt",
	})

	if err != nil {
		return TestResult{Name: "move_file", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件移动成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "move_file", Success: true}
}

// testDelete 测试删除（自动检测）/ Test delete (auto-detect)
func (s *ComprehensiveTestSuite) testDelete() TestResult {
	fmt.Println("  测试 delete (auto-detect)...")

	resp, err := s.client.CallTool(s.ctx, "delete", types.DeleteRequest{
		Path: "test_file_moved.txt",
	})

	if err != nil {
		return TestResult{Name: "delete", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件删除成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "delete", Success: true}
}

// testDeleteFile 测试删除文件 / Test delete file
func (s *ComprehensiveTestSuite) testDeleteFile() TestResult {
	fmt.Println("  测试 delete_file...")

	resp, err := s.client.CallTool(s.ctx, "delete_file", types.DeleteFileRequest{
		Path: "test_file_moved2.txt",
	})

	if err != nil {
		return TestResult{Name: "delete_file", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件删除成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "delete_file", Success: true}
}

// testCreateDirectory 测试创建目录 / Test create directory
func (s *ComprehensiveTestSuite) testCreateDirectory() TestResult {
	fmt.Println("  测试 create_directory...")

	resp, err := s.client.CallTool(s.ctx, "create_directory", types.CreateDirRequest{
		Path: "test_dir",
	})

	if err != nil {
		return TestResult{Name: "create_directory", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 目录创建成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "create_directory", Success: true}
}

// testListDirectory 测试列出目录 / Test list directory
func (s *ComprehensiveTestSuite) testListDirectory() TestResult {
	fmt.Println("  测试 list_directory...")

	resp, err := s.client.CallTool(s.ctx, "list_directory", types.ListDirRequest{
		Path: ".",
	})

	if err != nil {
		return TestResult{Name: "list_directory", Success: false, Error: err.Error()}
	}

	var listResp types.ListDirResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &listResp); err != nil {
		return TestResult{Name: "list_directory", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 目录列出成功，找到 %d 个条目\n", len(listResp.Files))
	return TestResult{Name: "list_directory", Success: true}
}

// testCopyDirectory 测试复制目录 / Test copy directory
func (s *ComprehensiveTestSuite) testCopyDirectory() TestResult {
	fmt.Println("  测试 copy_directory...")

	resp, err := s.client.CallTool(s.ctx, "copy_directory", types.CopyDirectoryRequest{
		Source:      "test_dir",
		Destination: "test_dir_copy",
	})

	if err != nil {
		return TestResult{Name: "copy_directory", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 目录复制成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "copy_directory", Success: true}
}

// testMoveDirectory 测试移动目录 / Test move directory
func (s *ComprehensiveTestSuite) testMoveDirectory() TestResult {
	fmt.Println("  测试 move_directory...")

	resp, err := s.client.CallTool(s.ctx, "move_directory", types.MoveDirectoryRequest{
		Source:      "test_dir_copy",
		Destination: "test_dir_moved",
	})

	if err != nil {
		return TestResult{Name: "move_directory", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 目录移动成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "move_directory", Success: true}
}

// testDeleteDirectory 测试删除目录 / Test delete directory
func (s *ComprehensiveTestSuite) testDeleteDirectory() TestResult {
	fmt.Println("  测试 delete_directory...")

	resp, err := s.client.CallTool(s.ctx, "delete_directory", types.DeleteDirectoryRequest{
		Path: "test_dir_moved",
	})

	if err != nil {
		return TestResult{Name: "delete_directory", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 目录删除成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "delete_directory", Success: true}
}

// testBatchDelete 测试批量删除 / Test batch delete
func (s *ComprehensiveTestSuite) testBatchDelete() TestResult {
	fmt.Println("  测试 batch_delete...")

	// 先创建一些测试文件 / Create some test files first
	_, _ = s.client.CallTool(s.ctx, "create_file", types.CreateFileRequest{
		Path:    "batch_test1.txt",
		Content: "test1",
	})
	_, _ = s.client.CallTool(s.ctx, "create_file", types.CreateFileRequest{
		Path:    "batch_test2.txt",
		Content: "test2",
	})

	resp, err := s.client.CallTool(s.ctx, "batch_delete", types.BatchDeleteRequest{
		Paths: []string{"batch_test1.txt", "batch_test2.txt"},
	})

	if err != nil {
		return TestResult{Name: "batch_delete", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 批量删除成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "batch_delete", Success: true}
}

// testGetWorkingDirectory 测试获取工作目录 / Test get working directory
func (s *ComprehensiveTestSuite) testGetWorkingDirectory() TestResult {
	fmt.Println("  测试 get_working_directory...")

	resp, err := s.client.CallTool(s.ctx, "get_working_directory", types.GetWorkingDirectoryRequest{})

	if err != nil {
		return TestResult{Name: "get_working_directory", Success: false, Error: err.Error()}
	}

	var wdResp types.GetWorkingDirectoryResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &wdResp); err != nil {
		return TestResult{Name: "get_working_directory", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 获取工作目录成功: %s\n", wdResp.WorkDir)
	return TestResult{Name: "get_working_directory", Success: true}
}

// testChangeDirectory 测试更改目录 / Test change directory
func (s *ComprehensiveTestSuite) testChangeDirectory() TestResult {
	fmt.Println("  测试 change_directory...")

	// 先获取当前目录 / Get current directory first
	resp1, _ := s.client.CallTool(s.ctx, "get_working_directory", types.GetWorkingDirectoryRequest{})
	var wdResp types.GetWorkingDirectoryResponse
	_ = json.UnmarshalFromString(resp1.Content[0].Text, &wdResp)
	originalDir := wdResp.WorkDir

	// 更改到test_dir / Change to test_dir
	resp, err := s.client.CallTool(s.ctx, "change_directory", types.ChangeDirectoryRequest{
		Path: "test_dir",
	})

	if err != nil {
		return TestResult{Name: "change_directory", Success: false, Error: err.Error()}
	}

	// 恢复原目录 / Restore original directory
	_, _ = s.client.CallTool(s.ctx, "change_directory", types.ChangeDirectoryRequest{
		Path: originalDir,
	})

	fmt.Printf("    ✅ 更改目录成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "change_directory", Success: true}
}

// testFileStat 测试文件状态 / Test file stat
func (s *ComprehensiveTestSuite) testFileStat() TestResult {
	fmt.Println("  测试 file_stat...")

	resp, err := s.client.CallTool(s.ctx, "file_stat", types.FileStatRequest{
		Path: "test_file.txt",
	})

	if err != nil {
		return TestResult{Name: "file_stat", Success: false, Error: err.Error()}
	}

	var statResp types.FileStatResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &statResp); err != nil {
		return TestResult{Name: "file_stat", Success: false, Error: err.Error()}
	}

	fileType := "file"
	if statResp.IsDir {
		fileType = "directory"
	}
	fmt.Printf("    ✅ 文件状态获取成功: 大小=%d, 类型=%s\n", statResp.Size, fileType)
	return TestResult{Name: "file_stat", Success: true}
}

// testFileExists 测试文件存在性 / Test file exists
func (s *ComprehensiveTestSuite) testFileExists() TestResult {
	fmt.Println("  测试 file_exists...")

	resp, err := s.client.CallTool(s.ctx, "file_exists", types.FileExistsRequest{
		Path: "test_file.txt",
	})

	if err != nil {
		return TestResult{Name: "file_exists", Success: false, Error: err.Error()}
	}

	var existsResp types.FileExistsResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &existsResp); err != nil {
		return TestResult{Name: "file_exists", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件存在性检查成功: exists=%v\n", existsResp.Exists)
	return TestResult{Name: "file_exists", Success: true}
}

// testSearchFiles 测试搜索文件 / Test search files
func (s *ComprehensiveTestSuite) testSearchFiles() TestResult {
	fmt.Println("  测试 search_files...")

	resp, err := s.client.CallTool(s.ctx, "search_files", types.SearchRequest{
		Path:    ".",
		Pattern: "*.txt",
	})

	if err != nil {
		return TestResult{Name: "search_files", Success: false, Error: err.Error()}
	}

	var searchResp types.SearchResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &searchResp); err != nil {
		return TestResult{Name: "search_files", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件搜索成功，找到 %d 个文件\n", len(searchResp.Files))
	return TestResult{Name: "search_files", Success: true}
}

// testGetSystemInfo 测试获取系统信息 / Test get system info
func (s *ComprehensiveTestSuite) testGetSystemInfo() TestResult {
	fmt.Println("  测试 get_system_info...")

	resp, err := s.client.CallTool(s.ctx, "get_system_info", types.GetSystemInfoRequest{})

	if err != nil {
		return TestResult{Name: "get_system_info", Success: false, Error: err.Error()}
	}

	var sysResp types.GetSystemInfoResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &sysResp); err != nil {
		return TestResult{Name: "get_system_info", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 系统信息获取成功:\n")
	fmt.Printf("       OS: %s %s\n", sysResp.OS.Platform, sysResp.OS.Architecture)
	fmt.Printf("       CPU: %d 核心\n", sysResp.CPU.LogicalCores)
	fmt.Printf("       内存: %d MB\n", sysResp.Memory.Total/(1024*1024))
	return TestResult{Name: "get_system_info", Success: true}
}

// testExecuteCommand 测试执行命令 / Test execute command
func (s *ComprehensiveTestSuite) testExecuteCommand() TestResult {
	fmt.Println("  测试 execute_command...")

	var command string
	var args []string
	// 根据操作系统选择命令 / Choose command based on OS
	if s.isWindows() {
		command = "cmd"
		args = []string{"/c", "echo", "Hello from MCP"}
	} else {
		command = "echo"
		args = []string{"Hello from MCP"}
	}

	resp, err := s.client.CallTool(s.ctx, "execute_command", types.ExecuteCommandRequest{
		Command: command,
		Args:    args,
		WorkDir: ".",
	})

	if err != nil {
		return TestResult{Name: "execute_command", Success: false, Error: err.Error()}
	}

	var cmdResp types.ExecuteCommandResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &cmdResp); err != nil {
		return TestResult{Name: "execute_command", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 命令执行成功，退出码: %d\n", cmdResp.ExitCode)
	return TestResult{Name: "execute_command", Success: true}
}

// testExecuteCommandAsync 测试异步执行命令 / Test execute command async
func (s *ComprehensiveTestSuite) testExecuteCommandAsync() TestResult {
	fmt.Println("  测试 execute_command_async...")

	var command string
	var args []string
	if s.isWindows() {
		command = "cmd"
		args = []string{"/c", "ping", "127.0.0.1", "-n", "3"}
	} else {
		command = "sleep"
		args = []string{"2"}
	}

	resp, err := s.client.CallTool(s.ctx, "execute_command_async", types.ExecuteCommandAsyncRequest{
		Command: command,
		Args:    args,
		WorkDir: ".",
	})

	if err != nil {
		return TestResult{Name: "execute_command_async", Success: false, Error: err.Error()}
	}

	var asyncResp types.ExecuteCommandAsyncResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &asyncResp); err != nil {
		return TestResult{Name: "execute_command_async", Success: false, Error: err.Error()}
	}

	lastAsyncTaskID = asyncResp.TaskID
	fmt.Printf("    ✅ 异步命令启动成功，任务ID: %s\n", asyncResp.TaskID)

	// 等待一下让命令执行 / Wait a bit for command execution
	time.Sleep(1 * time.Second)

	return TestResult{Name: "execute_command_async", Success: true}
}

// testGetCommandTask 测试获取命令任务 / Test get command task
func (s *ComprehensiveTestSuite) testGetCommandTask() TestResult {
	fmt.Println("  测试 get_command_task...")

	if lastAsyncTaskID == "" {
		return TestResult{Name: "get_command_task", Success: false, Error: "没有可用的任务ID"}
	}

	resp, err := s.client.CallTool(s.ctx, "get_command_task", types.GetCommandTaskRequest{
		TaskID: lastAsyncTaskID,
	})

	if err != nil {
		return TestResult{Name: "get_command_task", Success: false, Error: err.Error()}
	}

	var taskResp types.GetCommandTaskResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &taskResp); err != nil {
		return TestResult{Name: "get_command_task", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 任务信息获取成功，状态: %s\n", taskResp.Task.Status)
	return TestResult{Name: "get_command_task", Success: true}
}

// testCancelCommandTask 测试取消命令任务 / Test cancel command task
func (s *ComprehensiveTestSuite) testCancelCommandTask() TestResult {
	fmt.Println("  测试 cancel_command_task...")

	// 启动一个新的长时间运行的任务 / Start a new long-running task
	var command string
	var args []string
	if s.isWindows() {
		command = "cmd"
		args = []string{"/c", "ping", "127.0.0.1", "-n", "10"}
	} else {
		command = "sleep"
		args = []string{"10"}
	}

	resp1, err := s.client.CallTool(s.ctx, "execute_command_async", types.ExecuteCommandAsyncRequest{
		Command: command,
		Args:    args,
		WorkDir: ".",
	})

	if err != nil {
		return TestResult{Name: "cancel_command_task", Success: false, Error: err.Error()}
	}

	var asyncResp types.ExecuteCommandAsyncResponse
	if err = json.UnmarshalFromString(resp1.Content[0].Text, &asyncResp); err != nil {
		return TestResult{Name: "cancel_command_task", Success: false, Error: err.Error()}
	}

	// 取消任务 / Cancel the task
	resp, err := s.client.CallTool(s.ctx, "cancel_command_task", types.CancelCommandTaskRequest{
		TaskID: asyncResp.TaskID,
	})

	if err != nil {
		return TestResult{Name: "cancel_command_task", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 任务取消成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "cancel_command_task", Success: true}
}

// testGetCommandBlacklist 测试获取命令黑名单 / Test get command blacklist
func (s *ComprehensiveTestSuite) testGetCommandBlacklist() TestResult {
	fmt.Println("  测试 get_command_blacklist...")

	resp, err := s.client.CallTool(s.ctx, "get_command_blacklist", types.GetCommandBlacklistRequest{})

	if err != nil {
		return TestResult{Name: "get_command_blacklist", Success: false, Error: err.Error()}
	}

	var blacklistResp types.GetCommandBlacklistResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &blacklistResp); err != nil {
		return TestResult{Name: "get_command_blacklist", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 黑名单获取成功，命令数: %d, 目录数: %d\n", len(blacklistResp.Commands), len(blacklistResp.Directories))
	return TestResult{Name: "get_command_blacklist", Success: true}
}

// testUpdateCommandBlacklist 测试更新命令黑名单 / Test update command blacklist
func (s *ComprehensiveTestSuite) testUpdateCommandBlacklist() TestResult {
	fmt.Println("  测试 update_command_blacklist...")

	resp, err := s.client.CallTool(s.ctx, "update_command_blacklist", types.UpdateCommandBlacklistRequest{
		Commands:    []string{"test_blocked_cmd"},
		Directories: []string{"/test/blocked/dir"},
	})

	if err != nil {
		return TestResult{Name: "update_command_blacklist", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 黑名单更新成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "update_command_blacklist", Success: true}
}

// testGetCommandHistory 测试获取命令历史 / Test get command history
func (s *ComprehensiveTestSuite) testGetCommandHistory() TestResult {
	fmt.Println("  测试 get_command_history...")

	resp, err := s.client.CallTool(s.ctx, "get_command_history", types.GetCommandHistoryRequest{})

	if err != nil {
		return TestResult{Name: "get_command_history", Success: false, Error: err.Error()}
	}

	var historyResp types.GetCommandHistoryResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &historyResp); err != nil {
		return TestResult{Name: "get_command_history", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 命令历史获取成功，记录数: %d\n", len(historyResp.History))
	return TestResult{Name: "get_command_history", Success: true}
}

// testClearCommandHistory 测试清空命令历史 / Test clear command history
func (s *ComprehensiveTestSuite) testClearCommandHistory() TestResult {
	fmt.Println("  测试 clear_command_history...")

	resp, err := s.client.CallTool(s.ctx, "clear_command_history", types.ClearCommandHistoryRequest{})

	if err != nil {
		return TestResult{Name: "clear_command_history", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 命令历史清空成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "clear_command_history", Success: true}
}

// testGetPermissionLevel 测试获取权限级别 / Test get permission level
func (s *ComprehensiveTestSuite) testGetPermissionLevel() TestResult {
	fmt.Println("  测试 get_permission_level...")

	resp, err := s.client.CallTool(s.ctx, "get_permission_level", types.GetPermissionLevelRequest{})

	if err != nil {
		return TestResult{Name: "get_permission_level", Success: false, Error: err.Error()}
	}

	var permResp types.GetPermissionLevelResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &permResp); err != nil {
		return TestResult{Name: "get_permission_level", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 权限级别获取成功: %v\n", permResp.Level)
	return TestResult{Name: "get_permission_level", Success: true}
}

// testSetPermissionLevel 测试设置权限级别 / Test set permission level
func (s *ComprehensiveTestSuite) testSetPermissionLevel() TestResult {
	fmt.Println("  测试 set_permission_level...")

	// 设置为只读权限 / Set to read-only permission
	resp, err := s.client.CallTool(s.ctx, "set_permission_level", types.SetPermissionLevelRequest{
		Level: types.PermissionLevelReadOnly,
	})

	if err != nil {
		return TestResult{Name: "set_permission_level", Success: false, Error: err.Error()}
	}

	// 恢复为标准权限 / Restore to standard permission
	_, _ = s.client.CallTool(s.ctx, "set_permission_level", types.SetPermissionLevelRequest{
		Level: types.PermissionLevelStandard,
	})

	fmt.Printf("    ✅ 权限级别设置成功: %s\n", resp.Content[0].Text)
	return TestResult{Name: "set_permission_level", Success: true}
}

// testGetCurrentTime 测试获取当前时间 / Test get current time
func (s *ComprehensiveTestSuite) testGetCurrentTime() TestResult {
	fmt.Println("  测试 get_current_time...")

	resp, err := s.client.CallTool(s.ctx, "get_current_time", map[string]interface{}{})

	if err != nil {
		return TestResult{Name: "get_current_time", Success: false, Error: err.Error()}
	}

	var timeResp types.GetCurrentTimeResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &timeResp); err != nil {
		return TestResult{Name: "get_current_time", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 当前时间获取成功: %s (%s)\n", timeResp.Time, timeResp.TimeZone)
	return TestResult{Name: "get_current_time", Success: true}
}

// testDownloadFile 测试下载文件 / Test download file
func (s *ComprehensiveTestSuite) testDownloadFile() TestResult {
	fmt.Println("  测试 download_file...")

	// 使用一个公开的测试URL / Use a public test URL
	resp, err := s.client.CallTool(s.ctx, "download_file", types.DownloadFileRequest{
		URL:     "https://www.example.com/",
		Path:    "downloads/example.html",
		Method:  "GET",
		Timeout: 30,
	})

	if err != nil {
		return TestResult{Name: "download_file", Success: false, Error: err.Error()}
	}

	var downloadResp types.DownloadFileResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &downloadResp); err != nil {
		return TestResult{Name: "download_file", Success: false, Error: err.Error()}
	}

	fmt.Printf("    ✅ 文件下载成功: %s (大小: %d 字节)\n", downloadResp.SandboxPath, downloadResp.Size)
	return TestResult{Name: "download_file", Success: true}
}

// testFinalCleanup 最终清理 / Final cleanup
func (s *ComprehensiveTestSuite) testFinalCleanup() TestResult {
	fmt.Println("  清理所有测试文件...")

	// 清理所有测试文件和目录 / Clean up all test files and directories
	testPaths := []string{
		"test_file.txt",
		"test_dir",
		"downloads",
	}

	for _, path := range testPaths {
		_, _ = s.client.CallTool(s.ctx, "delete", types.DeleteRequest{
			Path: path,
		})
	}

	fmt.Println("    ✅ 清理完成")
	return TestResult{Name: "cleanup", Success: true}
}

// isWindows 检查是否为Windows系统 / Check if running on Windows
func (s *ComprehensiveTestSuite) isWindows() bool {
	// 通过获取系统信息来判断 / Determine by getting system info
	resp, err := s.client.CallTool(s.ctx, "get_system_info", types.GetSystemInfoRequest{})
	if err != nil {
		return false
	}

	var sysResp types.GetSystemInfoResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &sysResp); err != nil {
		return false
	}

	return sysResp.OS.Platform == "windows"
}

// testCommandExecutionTools 测试命令执行工具 / Test command execution tools
func (s *ComprehensiveTestSuite) testCommandExecutionTools() {
	fmt.Println("\n=== 5. 命令执行工具测试 ===")

	// execute_command
	s.addResult(s.testExecuteCommand())
}
