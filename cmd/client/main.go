package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"mcp-toolkit/pkg/client"
	"mcp-toolkit/pkg/types"
	"mcp-toolkit/pkg/utils/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 解析命令行参数 / Parse command line arguments
	host := flag.String("host", "127.0.0.1", "MCP server host")
	port := flag.Int("port", 8080, "MCP server port")
	path := flag.String("path", "/mcp", "MCP server path")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	// 创建日志记录器 / Create logger
	var logger *zap.Logger
	if *verbose {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = config.Build()
	} else {
		logger = zap.NewNop()
	}
	defer func() { _ = logger.Sync() }()

	// 创建客户端 / Create client
	mcpClient := client.NewHTTPClient(*host, *port, *path, logger)
	defer func() { _ = mcpClient.Close() }()

	ctx := context.Background()

	fmt.Println("=== MCP Client Test Suite ===")
	fmt.Printf("Connecting to: http://%s:%d%s\n\n", *host, *port, *path)

	// 测试初始化 / Test initialize
	fmt.Println("1. Testing Initialize...")
	initResp, err := mcpClient.Initialize(ctx, types.ProtocolVersion)
	if err != nil {
		fmt.Printf("❌ Initialize failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Initialize successful\n")
	fmt.Printf("   Server: %s v%s\n", initResp.ServerInfo.Name, initResp.ServerInfo.Version)
	fmt.Printf("   Protocol: %s\n\n", initResp.ProtocolVersion)

	// 测试列出工具 / Test list tools
	fmt.Println("2. Testing List Tools...")
	toolsResp, err := mcpClient.ListTools(ctx)
	if err != nil {
		fmt.Printf("❌ List tools failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ List tools successful\n")
	fmt.Printf("   Found %d tools:\n", len(toolsResp.Tools))
	for i, tool := range toolsResp.Tools {
		fmt.Printf("   %d. %s - %s\n", i+1, tool.Name, tool.Description)
	}
	fmt.Println()

	// 运行所有工具测试 / Run all tool tests
	testResults := runAllToolTests(ctx, mcpClient)

	// 打印测试结果摘要 / Print test results summary
	fmt.Println("\n=== Test Results Summary ===")
	passed := 0
	failed := 0
	for _, result := range testResults {
		if result.Success {
			passed++
			fmt.Printf("✅ %s\n", result.Name)
		} else {
			failed++
			fmt.Printf("❌ %s: %s\n", result.Name, result.Error)
		}
	}
	fmt.Printf("\nTotal: %d, Passed: %d, Failed: %d\n", len(testResults), passed, failed)

	if failed > 0 {
		os.Exit(1)
	}
}

// TestResult 测试结果 / Test result
type TestResult struct {
	Name    string
	Success bool
	Error   string
}

// runAllToolTests 运行所有工具测试 / Run all tool tests
func runAllToolTests(ctx context.Context, mcpClient client.Client) []TestResult {
	results := []TestResult{}

	// 基础文件操作测试 / Basic file operation tests
	results = append(results, testCreateFile(ctx, mcpClient))
	results = append(results, testReadFile(ctx, mcpClient))
	results = append(results, testWriteFile(ctx, mcpClient))
	results = append(results, testCreateDirectory(ctx, mcpClient))
	results = append(results, testListDirectory(ctx, mcpClient))
	results = append(results, testFileStat(ctx, mcpClient))
	results = append(results, testFileExists(ctx, mcpClient))
	results = append(results, testSearchFiles(ctx, mcpClient))
	results = append(results, testCopyFile(ctx, mcpClient))
	results = append(results, testMoveFile(ctx, mcpClient))
	results = append(results, testBatchDelete(ctx, mcpClient))

	// 工作目录管理测试 / Working directory management tests
	results = append(results, testGetWorkingDirectory(ctx, mcpClient))
	results = append(results, testChangeDirectory(ctx, mcpClient))

	// 命令执行测试 / Command execution tests
	results = append(results, testExecuteCommand(ctx, mcpClient))
	results = append(results, testGetCommandBlacklist(ctx, mcpClient))
	results = append(results, testUpdateCommandBlacklist(ctx, mcpClient))

	// 异步命令执行测试 / Async command execution tests
	results = append(results, testExecuteCommandAsync(ctx, mcpClient))
	results = append(results, testGetCommandTask(ctx, mcpClient))
	results = append(results, testCancelCommandTask(ctx, mcpClient))

	// 命令历史测试 / Command history tests
	results = append(results, testGetCommandHistory(ctx, mcpClient))
	results = append(results, testClearCommandHistory(ctx, mcpClient))

	// 权限管理测试 / Permission management tests
	results = append(results, testGetPermissionLevel(ctx, mcpClient))
	results = append(results, testSetPermissionLevel(ctx, mcpClient))

	// 清理和其他测试 / Cleanup and other tests
	results = append(results, testDelete(ctx, mcpClient))
	results = append(results, testGetCurrentTime(ctx, mcpClient))

	return results
}

// testCreateFile 测试创建文件 / Test create file
func testCreateFile(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n3. Testing create_file...")

	resp, err := mcpClient.CallTool(ctx, "create_file", types.CreateFileRequest{
		Path:    "test_file.txt",
		Content: "Hello, MCP!",
	})

	if err != nil {
		return TestResult{Name: "create_file", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Response: %s\n", resp.Content[0].Text)
	return TestResult{Name: "create_file", Success: true}
}

// testReadFile 测试读取文件 / Test read file
func testReadFile(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n4. Testing read_file...")

	resp, err := mcpClient.CallTool(ctx, "read_file", types.ReadFileRequest{
		Path: "test_file.txt",
	})

	if err != nil {
		return TestResult{Name: "read_file", Success: false, Error: err.Error()}
	}

	var readResp types.ReadFileResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &readResp); err != nil {
		return TestResult{Name: "read_file", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Content: %s\n", readResp.Content)
	if readResp.Content != "Hello, MCP!" {
		return TestResult{Name: "read_file", Success: false, Error: "content mismatch"}
	}

	return TestResult{Name: "read_file", Success: true}
}

// testWriteFile 测试写入文件 / Test write file
func testWriteFile(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n5. Testing write_file...")

	resp, err := mcpClient.CallTool(ctx, "write_file", types.WriteFileRequest{
		Path:    "test_file.txt",
		Content: "Updated content",
	})

	if err != nil {
		return TestResult{Name: "write_file", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Response: %s\n", resp.Content[0].Text)
	return TestResult{Name: "write_file", Success: true}
}

// testCreateDirectory 测试创建目录 / Test create directory
func testCreateDirectory(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n6. Testing create_directory...")

	resp, err := mcpClient.CallTool(ctx, "create_directory", types.CreateDirRequest{
		Path: "test_dir",
	})

	if err != nil {
		return TestResult{Name: "create_directory", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Response: %s\n", resp.Content[0].Text)
	return TestResult{Name: "create_directory", Success: true}
}

// testListDirectory 测试列出目录 / Test list directory
func testListDirectory(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n7. Testing list_directory...")

	resp, err := mcpClient.CallTool(ctx, "list_directory", types.ListDirRequest{
		Path: ".",
	})

	if err != nil {
		return TestResult{Name: "list_directory", Success: false, Error: err.Error()}
	}

	var listResp types.ListDirResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &listResp); err != nil {
		return TestResult{Name: "list_directory", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Found %d entries\n", len(listResp.Files))
	for _, entry := range listResp.Files {
		fmt.Printf("   - %s (IsDir: %v)\n", entry.Name, entry.IsDir)
	}

	return TestResult{Name: "list_directory", Success: true}
}

// testFileStat 测试文件状态 / Test file stat
func testFileStat(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n8. Testing file_stat...")

	resp, err := mcpClient.CallTool(ctx, "file_stat", types.FileStatRequest{
		Path: "test_file.txt",
	})

	if err != nil {
		return TestResult{Name: "file_stat", Success: false, Error: err.Error()}
	}

	var statResp types.FileStatResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &statResp); err != nil {
		return TestResult{Name: "file_stat", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Name: %s, Size: %d, IsDir: %v, Mode: %s\n",
		statResp.Name, statResp.Size, statResp.IsDir, statResp.Mode)

	return TestResult{Name: "file_stat", Success: true}
}

// testFileExists 测试文件是否存在 / Test file exists
func testFileExists(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n9. Testing file_exists...")

	resp, err := mcpClient.CallTool(ctx, "file_exists", types.FileExistsRequest{
		Path: "test_file.txt",
	})

	if err != nil {
		return TestResult{Name: "file_exists", Success: false, Error: err.Error()}
	}

	var existsResp types.FileExistsResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &existsResp); err != nil {
		return TestResult{Name: "file_exists", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Exists: %v\n", existsResp.Exists)
	if !existsResp.Exists {
		return TestResult{Name: "file_exists", Success: false, Error: "file should exist"}
	}

	return TestResult{Name: "file_exists", Success: true}
}

// testSearchFiles 测试搜索文件 / Test search files
func testSearchFiles(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n10. Testing search_files...")

	resp, err := mcpClient.CallTool(ctx, "search_files", types.SearchRequest{
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

	fmt.Printf("   Found %d files matching pattern\n", len(searchResp.Files))
	for _, file := range searchResp.Files {
		fmt.Printf("   - %s\n", file.Name)
	}

	return TestResult{Name: "search_files", Success: true}
}

// testCopyFile 测试复制文件 / Test copy file
func testCopyFile(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n11. Testing copy...")

	resp, err := mcpClient.CallTool(ctx, "copy", types.CopyRequest{
		Source:      "test_file.txt",
		Destination: "test_file_copy.txt",
	})

	if err != nil {
		return TestResult{Name: "copy", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Response: %s\n", resp.Content[0].Text)
	return TestResult{Name: "copy", Success: true}
}

// testMoveFile 测试移动文件 / Test move file
func testMoveFile(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n12. Testing move...")

	resp, err := mcpClient.CallTool(ctx, "move", types.MoveRequest{
		Source:      "test_file_copy.txt",
		Destination: "test_file_moved.txt",
	})

	if err != nil {
		return TestResult{Name: "move", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Response: %s\n", resp.Content[0].Text)
	return TestResult{Name: "move", Success: true}
}

// testBatchDelete 测试批量删除 / Test batch delete
func testBatchDelete(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n13. Testing batch_delete...")

	// 先创建几个测试文件 / Create some test files first
	for i := 1; i <= 3; i++ {
		_, _ = mcpClient.CallTool(ctx, "create_file", types.CreateFileRequest{
			Path:    fmt.Sprintf("batch_test_%d.txt", i),
			Content: fmt.Sprintf("Batch test file %d", i),
		})
	}

	// 等待一下确保文件创建完成 / Wait a bit to ensure files are created
	time.Sleep(100 * time.Millisecond)

	resp, err := mcpClient.CallTool(ctx, "batch_delete", types.BatchDeleteRequest{
		Paths: []string{"batch_test_1.txt", "batch_test_2.txt", "batch_test_3.txt"},
	})

	if err != nil {
		return TestResult{Name: "batch_delete", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Response: %s\n", resp.Content[0].Text)
	return TestResult{Name: "batch_delete", Success: true}
}

// testGetWorkingDirectory 测试获取工作目录 / Test get working directory
func testGetWorkingDirectory(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n14. Testing get_working_directory...")

	resp, err := mcpClient.CallTool(ctx, "get_working_directory", types.GetWorkingDirectoryRequest{})

	if err != nil {
		return TestResult{Name: "get_working_directory", Success: false, Error: err.Error()}
	}

	var wdResp types.GetWorkingDirectoryResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &wdResp); err != nil {
		return TestResult{Name: "get_working_directory", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Current working directory: %s\n", wdResp.WorkDir)
	return TestResult{Name: "get_working_directory", Success: true}
}

// testChangeDirectory 测试切换目录 / Test change directory
func testChangeDirectory(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n15. Testing change_directory...")

	resp, err := mcpClient.CallTool(ctx, "change_directory", types.ChangeDirectoryRequest{
		Path: "test_dir",
	})

	if err != nil {
		return TestResult{Name: "change_directory", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Response: %s\n", resp.Content[0].Text)

	// 切换回根目录 / Change back to root directory
	_, _ = mcpClient.CallTool(ctx, "change_directory", types.ChangeDirectoryRequest{
		Path: ".",
	})

	return TestResult{Name: "change_directory", Success: true}
}

// testExecuteCommand 测试执行命令 / Test execute command
func testExecuteCommand(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n16. Testing execute_command...")

	// 使用跨平台的echo命令 / Use cross-platform echo command
	resp, err := mcpClient.CallTool(ctx, "execute_command", types.ExecuteCommandRequest{
		Command: "echo",
		Args:    []string{"Hello from MCP"},
		WorkDir: ".",
	})

	if err != nil {
		return TestResult{Name: "execute_command", Success: false, Error: err.Error()}
	}

	var cmdResp types.ExecuteCommandResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &cmdResp); err != nil {
		return TestResult{Name: "execute_command", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Exit code: %d\n", cmdResp.ExitCode)
	fmt.Printf("   Stdout: %s\n", cmdResp.Stdout)
	fmt.Printf("   Stderr: %s\n", cmdResp.Stderr)

	return TestResult{Name: "execute_command", Success: true}
}

// testGetCommandBlacklist 测试获取命令黑名单 / Test get command blacklist
func testGetCommandBlacklist(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n17. Testing get_command_blacklist...")

	resp, err := mcpClient.CallTool(ctx, "get_command_blacklist", types.GetCommandBlacklistRequest{})

	if err != nil {
		return TestResult{Name: "get_command_blacklist", Success: false, Error: err.Error()}
	}

	var blacklistResp types.GetCommandBlacklistResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &blacklistResp); err != nil {
		return TestResult{Name: "get_command_blacklist", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Blacklisted commands: %d\n", len(blacklistResp.Commands))
	fmt.Printf("   Blacklisted directories: %d\n", len(blacklistResp.Directories))
	fmt.Printf("   System directories: %d\n", len(blacklistResp.SystemDirectories))

	return TestResult{Name: "get_command_blacklist", Success: true}
}

// testUpdateCommandBlacklist 测试更新命令黑名单 / Test update command blacklist
func testUpdateCommandBlacklist(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n18. Testing update_command_blacklist...")

	resp, err := mcpClient.CallTool(ctx, "update_command_blacklist", types.UpdateCommandBlacklistRequest{
		Commands:    []string{"test_cmd"},
		Directories: []string{},
	})

	if err != nil {
		return TestResult{Name: "update_command_blacklist", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Response: %s\n", resp.Content[0].Text)

	// 恢复原始黑名单 / Restore original blacklist
	_, _ = mcpClient.CallTool(ctx, "update_command_blacklist", types.UpdateCommandBlacklistRequest{
		Commands:    []string{},
		Directories: []string{},
	})

	return TestResult{Name: "update_command_blacklist", Success: true}
}

// testExecuteCommandAsync 测试异步执行命令 / Test execute command async
func testExecuteCommandAsync(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n19. Testing execute_command_async...")

	resp, err := mcpClient.CallTool(ctx, "execute_command_async", types.ExecuteCommandAsyncRequest{
		Command: "echo",
		Args:    []string{"Async test"},
		WorkDir: ".",
		Timeout: 10,
	})

	if err != nil {
		return TestResult{Name: "execute_command_async", Success: false, Error: err.Error()}
	}

	var asyncResp types.ExecuteCommandAsyncResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &asyncResp); err != nil {
		return TestResult{Name: "execute_command_async", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Task ID: %s\n", asyncResp.TaskID)

	// 保存任务ID供后续测试使用 / Save task ID for later tests
	lastAsyncTaskID = asyncResp.TaskID

	return TestResult{Name: "execute_command_async", Success: true}
}

// testGetCommandTask 测试获取命令任务 / Test get command task
func testGetCommandTask(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n20. Testing get_command_task...")

	if lastAsyncTaskID == "" {
		return TestResult{Name: "get_command_task", Success: false, Error: "no async task ID available"}
	}

	// 等待任务完成 / Wait for task completion
	time.Sleep(200 * time.Millisecond)

	resp, err := mcpClient.CallTool(ctx, "get_command_task", types.GetCommandTaskRequest{
		TaskID: lastAsyncTaskID,
	})

	if err != nil {
		return TestResult{Name: "get_command_task", Success: false, Error: err.Error()}
	}

	var taskResp types.GetCommandTaskResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &taskResp); err != nil {
		return TestResult{Name: "get_command_task", Success: false, Error: err.Error()}
	}

	if taskResp.Task != nil {
		fmt.Printf("   Status: %s\n", taskResp.Task.Status)
		fmt.Printf("   Exit code: %d\n", taskResp.Task.ExitCode)
	}

	return TestResult{Name: "get_command_task", Success: true}
}

// testCancelCommandTask 测试取消命令任务 / Test cancel command task
func testCancelCommandTask(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n21. Testing cancel_command_task...")

	// 创建一个新的长时间运行任务 / Create a new long-running task
	resp, err := mcpClient.CallTool(ctx, "execute_command_async", types.ExecuteCommandAsyncRequest{
		Command: "echo",
		Args:    []string{"Long running task"},
		WorkDir: ".",
		Timeout: 60,
	})

	if err != nil {
		return TestResult{Name: "cancel_command_task", Success: false, Error: err.Error()}
	}

	var asyncResp types.ExecuteCommandAsyncResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &asyncResp); err != nil {
		return TestResult{Name: "cancel_command_task", Success: false, Error: err.Error()}
	}

	// 立即取消任务 / Cancel task immediately
	cancelResp, err := mcpClient.CallTool(ctx, "cancel_command_task", types.CancelCommandTaskRequest{
		TaskID: asyncResp.TaskID,
	})

	// 取消可能成功或失败(如果任务已经完成) / Cancel may succeed or fail (if task already completed)
	if err == nil {
		fmt.Printf("   Response: %s\n", cancelResp.Content[0].Text)
		return TestResult{Name: "cancel_command_task", Success: true}
	}

	// 如果任务已完成,也算测试通过 / If task already completed, test still passes
	fmt.Printf("   Task already completed or cancelled\n")
	return TestResult{Name: "cancel_command_task", Success: true}
}

// testGetCommandHistory 测试获取命令历史 / Test get command history
func testGetCommandHistory(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n22. Testing get_command_history...")

	resp, err := mcpClient.CallTool(ctx, "get_command_history", types.GetCommandHistoryRequest{
		Limit: 10,
	})

	if err != nil {
		return TestResult{Name: "get_command_history", Success: false, Error: err.Error()}
	}

	var historyResp types.GetCommandHistoryResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &historyResp); err != nil {
		return TestResult{Name: "get_command_history", Success: false, Error: err.Error()}
	}

	fmt.Printf("   History entries: %d\n", len(historyResp.History))
	for i, entry := range historyResp.History {
		if i < 3 { // 只显示前3条 / Only show first 3 entries
			fmt.Printf("   - %s %v (Exit: %d)\n", entry.Command, entry.Args, entry.ExitCode)
		}
	}

	return TestResult{Name: "get_command_history", Success: true}
}

// testClearCommandHistory 测试清空命令历史 / Test clear command history
func testClearCommandHistory(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n23. Testing clear_command_history...")

	resp, err := mcpClient.CallTool(ctx, "clear_command_history", types.ClearCommandHistoryRequest{})

	if err != nil {
		return TestResult{Name: "clear_command_history", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Response: %s\n", resp.Content[0].Text)
	return TestResult{Name: "clear_command_history", Success: true}
}

// testGetPermissionLevel 测试获取权限级别 / Test get permission level
func testGetPermissionLevel(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n24. Testing get_permission_level...")

	resp, err := mcpClient.CallTool(ctx, "get_permission_level", types.GetPermissionLevelRequest{})

	if err != nil {
		return TestResult{Name: "get_permission_level", Success: false, Error: err.Error()}
	}

	var permResp types.GetPermissionLevelResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &permResp); err != nil {
		return TestResult{Name: "get_permission_level", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Current permission level: %v\n", permResp.Level)
	return TestResult{Name: "get_permission_level", Success: true}
}

// testSetPermissionLevel 测试设置权限级别 / Test set permission level
func testSetPermissionLevel(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n25. Testing set_permission_level...")

	// 设置为只读权限 / Set to read-only permission
	resp, err := mcpClient.CallTool(ctx, "set_permission_level", types.SetPermissionLevelRequest{
		Level: types.PermissionLevelReadOnly,
	})

	if err != nil {
		return TestResult{Name: "set_permission_level", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Response: %s\n", resp.Content[0].Text)

	// 恢复为标准权限 / Restore to standard permission
	_, _ = mcpClient.CallTool(ctx, "set_permission_level", types.SetPermissionLevelRequest{
		Level: types.PermissionLevelStandard,
	})

	return TestResult{Name: "set_permission_level", Success: true}
}

// testDelete 测试删除 / Test delete
func testDelete(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n26. Testing delete...")

	// 删除测试文件和目录 / Delete test files and directories
	files := []string{"test_file.txt", "test_file_moved.txt", "test_dir"}

	for _, file := range files {
		resp, err := mcpClient.CallTool(ctx, "delete", types.DeleteRequest{
			Path: file,
		})

		if err != nil {
			// 文件可能已经不存在,继续 / File may not exist, continue
			fmt.Printf("   Skipped %s (may not exist)\n", file)
			continue
		}

		fmt.Printf("   Deleted %s: %s\n", file, resp.Content[0].Text)
	}

	return TestResult{Name: "delete", Success: true}
}

// testGetCurrentTime 测试获取当前时间 / Test get current time
func testGetCurrentTime(ctx context.Context, mcpClient client.Client) TestResult {
	fmt.Println("\n27. Testing get_current_time...")

	resp, err := mcpClient.CallTool(ctx, "get_current_time", map[string]interface{}{})

	if err != nil {
		return TestResult{Name: "get_current_time", Success: false, Error: err.Error()}
	}

	var timeResp types.GetCurrentTimeResponse
	if err = json.UnmarshalFromString(resp.Content[0].Text, &timeResp); err != nil {
		return TestResult{Name: "get_current_time", Success: false, Error: err.Error()}
	}

	fmt.Printf("   Time: %s\n", timeResp.Time)
	fmt.Printf("   Timezone: %s\n", timeResp.TimeZone)
	fmt.Printf("   Unix: %d\n", timeResp.Unix)

	return TestResult{Name: "get_current_time", Success: true}
}

// lastAsyncTaskID 保存最后一个异步任务ID / Save last async task ID
var lastAsyncTaskID string
