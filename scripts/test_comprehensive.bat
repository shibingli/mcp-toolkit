@echo off
REM 全面测试MCP Toolkit所有功能的脚本 / Comprehensive test script for all MCP Toolkit features
REM 使用方法 / Usage: test_comprehensive.bat

echo ========================================
echo MCP Toolkit 全面功能测试
echo Comprehensive MCP Toolkit Feature Test
echo ========================================
echo.

REM 检查是否已编译 / Check if compiled
if not exist "mcp-toolkit.exe" (
    echo [错误] 未找到 mcp-toolkit.exe，请先编译项目
    echo [ERROR] mcp-toolkit.exe not found, please build the project first
    echo.
    echo 运行编译命令 / Run build command:
    echo   go build -tags="sonic" -o mcp-toolkit.exe
    exit /b 1
)

REM 启动MCP服务器 / Start MCP server
echo [1/3] 启动 MCP 服务器...
echo [1/3] Starting MCP server...
start /B mcp-toolkit.exe
timeout /t 3 /nobreak >nul

REM 等待服务器启动 / Wait for server to start
echo [2/3] 等待服务器就绪...
echo [2/3] Waiting for server to be ready...
timeout /t 2 /nobreak >nul

REM 运行测试客户端 / Run test client
echo [3/3] 运行全面测试...
echo [3/3] Running comprehensive tests...
echo.

go run cmd/client/main.go cmd/client/comprehensive_suite.go -verbose

REM 保存退出码 / Save exit code
set TEST_EXIT_CODE=%ERRORLEVEL%

REM 停止服务器 / Stop server
echo.
echo 停止 MCP 服务器...
echo Stopping MCP server...
taskkill /F /IM mcp-toolkit.exe >nul 2>&1

echo.
if %TEST_EXIT_CODE% EQU 0 (
    echo ========================================
    echo ✅ 所有测试通过！
    echo ✅ All tests passed!
    echo ========================================
) else (
    echo ========================================
    echo ❌ 测试失败，退出码: %TEST_EXIT_CODE%
    echo ❌ Tests failed with exit code: %TEST_EXIT_CODE%
    echo ========================================
)

exit /b %TEST_EXIT_CODE%

