package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"mcp-toolkit/pkg/client"
	"mcp-toolkit/pkg/types"

	"go.uber.org/zap"
)

// ToolSchema 工具输入模式 / Tool input schema
type ToolSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

// main 主函数 / Main function
func main() {
	// 解析命令行参数 / Parse command line arguments
	transport := flag.String("transport", "http", "Transport mode: stdio, http, sse")
	host := flag.String("host", "127.0.0.1", "MCP server host (for http/sse)")
	port := flag.Int("port", 8080, "MCP server port (for http/sse)")
	path := flag.String("path", "/mcp", "MCP server path (for http/sse)")
	detailed := flag.Bool("detailed", false, "Show detailed tool information including input schema")
	command := flag.String("command", "", "Command to launch MCP server (for stdio mode)")
	args := flag.String("args", "", "Arguments for MCP server command (for stdio mode, comma-separated)")
	flag.Parse()

	// 创建静默日志记录器 / Create silent logger
	logger := zap.NewNop()
	defer func() { _ = logger.Sync() }()

	ctx := context.Background()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║           MCP Server 工具集查询 / MCP Tools List              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")

	// 根据传输模式创建客户端 / Create client based on transport mode
	var mcpClient client.Client
	var err error

	switch strings.ToLower(*transport) {
	case "stdio":
		fmt.Printf("\n传输模式 / Transport: Stdio\n")
		if *command == "" {
			fmt.Println("❌ 错误 / Error: -command is required for stdio mode")
			fmt.Println("   示例 / Example: -command ./mcp-toolkit.exe -args \"-transport,stdio\"")
			os.Exit(1)
		}

		// 解析参数 / Parse arguments
		var cmdArgs []string
		if *args != "" {
			cmdArgs = strings.Split(*args, ",")
		}

		fmt.Printf("命令 / Command: %s %v\n\n", *command, cmdArgs)

		mcpClient, err = client.NewStdioClient(*command, cmdArgs, logger)
		if err != nil {
			fmt.Printf("❌ 创建Stdio客户端失败 / Failed to create Stdio client: %v\n", err)
			os.Exit(1)
		}

	case "sse":
		fmt.Printf("\n传输模式 / Transport: SSE\n")
		fmt.Printf("连接地址 / Server URL: http://%s:%d%s\n\n", *host, *port, *path)
		mcpClient = client.NewSSEClient(*host, *port, *path, logger)

	case "http":
		fallthrough
	default:
		fmt.Printf("\n传输模式 / Transport: HTTP\n")
		fmt.Printf("连接地址 / Server URL: http://%s:%d%s\n\n", *host, *port, *path)
		mcpClient = client.NewHTTPClient(*host, *port, *path, logger)
	}

	defer func() { _ = mcpClient.Close() }()

	// 初始化连接 / Initialize connection
	fmt.Println("正在连接服务器... / Connecting to server...")
	initResp, err := mcpClient.Initialize(ctx, types.ProtocolVersion)
	if err != nil {
		fmt.Printf("❌ 连接失败 / Connection failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 连接成功 / Connected successfully\n")
	fmt.Printf("   服务器 / Server: %s v%s\n", initResp.ServerInfo.Name, initResp.ServerInfo.Version)
	fmt.Printf("   协议版本 / Protocol: %s\n\n", initResp.ProtocolVersion)

	// 获取工具列表 / Get tools list
	fmt.Println("正在获取工具列表... / Fetching tools list...")
	toolsResp, err := mcpClient.ListTools(ctx)
	if err != nil {
		fmt.Printf("❌ 获取工具列表失败 / Failed to get tools list: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 成功获取 %d 个工具 / Successfully retrieved %d tools\n\n", len(toolsResp.Tools), len(toolsResp.Tools))

	// 按类别分组工具 / Group tools by category
	categories := map[string][]client.Tool{
		"文件操作 / File Operations":      {},
		"目录操作 / Directory Operations": {},
		"文件信息 / File Information":     {},
		"系统信息 / System Information":   {},
		"命令执行 / Command Execution":    {},
		"命令管理 / Command Management":   {},
		"权限管理 / Permission":           {},
		"时间工具 / Time":                 {},
		"下载工具 / Download":             {},
	}

	// 分类工具 / Categorize tools
	for _, tool := range toolsResp.Tools {
		name := tool.Name
		switch {
		case strings.Contains(name, "file") && !strings.Contains(name, "download"):
			categories["文件操作 / File Operations"] = append(categories["文件操作 / File Operations"], tool)
		case strings.Contains(name, "directory") || strings.Contains(name, "dir"):
			categories["目录操作 / Directory Operations"] = append(categories["目录操作 / Directory Operations"], tool)
		case strings.Contains(name, "stat") || strings.Contains(name, "exists") || strings.Contains(name, "search"):
			categories["文件信息 / File Information"] = append(categories["文件信息 / File Information"], tool)
		case strings.Contains(name, "system_info"):
			categories["系统信息 / System Information"] = append(categories["系统信息 / System Information"], tool)
		case strings.Contains(name, "execute_command"):
			categories["命令执行 / Command Execution"] = append(categories["命令执行 / Command Execution"], tool)
		case strings.Contains(name, "command") && !strings.Contains(name, "execute"):
			categories["命令管理 / Command Management"] = append(categories["命令管理 / Command Management"], tool)
		case strings.Contains(name, "permission"):
			categories["权限管理 / Permission"] = append(categories["权限管理 / Permission"], tool)
		case strings.Contains(name, "time"):
			categories["时间工具 / Time"] = append(categories["时间工具 / Time"], tool)
		case strings.Contains(name, "download"):
			categories["下载工具 / Download"] = append(categories["下载工具 / Download"], tool)
		}
	}

	// 打印分类工具列表 / Print categorized tools list
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("                        工具列表 / Tools List                   ")
	fmt.Println("════════════════════════════════════════════════════════════════")

	totalCount := 0
	for category, tools := range categories {
		if len(tools) == 0 {
			continue
		}
		fmt.Printf("\n【%s】(%d 个工具)\n", category, len(tools))
		fmt.Println(strings.Repeat("─", 64))
		for i, tool := range tools {
			totalCount++
			fmt.Printf("%d. %s\n", i+1, tool.Name)
			if *detailed {
				// 显示详细描述 / Show detailed description
				desc := tool.Description
				if len(desc) > 100 {
					desc = desc[:100] + "..."
				}
				fmt.Printf("   描述 / Description: %s\n", desc)

				// 显示输入参数 / Show input parameters
				if tool.InputSchema != nil {
					if props, ok := tool.InputSchema["properties"].(map[string]interface{}); ok && len(props) > 0 {
						fmt.Println("   参数 / Parameters:")

						// 获取必填字段列表 / Get required fields list
						requiredFields := make(map[string]bool)
						if reqList, ok := tool.InputSchema["required"].([]interface{}); ok {
							for _, req := range reqList {
								if reqStr, ok := req.(string); ok {
									requiredFields[reqStr] = true
								}
							}
						}

						// 显示每个参数 / Show each parameter
						for paramName, paramValue := range props {
							paramType := "unknown"
							if paramMap, ok := paramValue.(map[string]interface{}); ok {
								if t, ok := paramMap["type"].(string); ok {
									paramType = t
								}
							}

							required := ""
							if requiredFields[paramName] {
								required = " [必填 / Required]"
							}
							fmt.Printf("     - %s (%s)%s\n", paramName, paramType, required)
						}
					}
				}
				fmt.Println()
			} else {
				// 简短描述 / Short description
				desc := tool.Description
				if len(desc) > 80 {
					desc = desc[:80] + "..."
				}
				fmt.Printf("   %s\n", desc)
			}
		}
	}

	fmt.Println("\n════════════════════════════════════════════════════════════════")
	fmt.Printf("总计 / Total: %d 个工具 / tools\n", totalCount)
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("\n提示 / Tips:")
	fmt.Println("  - 使用 -detailed 参数查看详细信息 / Use -detailed flag for detailed info")
	fmt.Println("  - 使用 -host, -port, -path 参数指定服务器地址 / Use -host, -port, -path to specify server")
	fmt.Println()
}
