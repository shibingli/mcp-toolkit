#!/bin/bash
# 全面测试MCP Toolkit所有功能的脚本 / Comprehensive test script for all MCP Toolkit features
# 使用方法 / Usage: ./test_comprehensive.sh

set -e

echo "========================================"
echo "MCP Toolkit 全面功能测试"
echo "Comprehensive MCP Toolkit Feature Test"
echo "========================================"
echo ""

# 检查是否已编译 / Check if compiled
if [ ! -f "mcp-toolkit" ]; then
    echo "[错误] 未找到 mcp-toolkit，请先编译项目"
    echo "[ERROR] mcp-toolkit not found, please build the project first"
    echo ""
    echo "运行编译命令 / Run build command:"
    echo "  go build -tags=\"sonic\" -o mcp-toolkit"
    exit 1
fi

# 启动MCP服务器 / Start MCP server
echo "[1/3] 启动 MCP 服务器..."
echo "[1/3] Starting MCP server..."
./mcp-toolkit &
MCP_PID=$!
sleep 3

# 等待服务器启动 / Wait for server to start
echo "[2/3] 等待服务器就绪..."
echo "[2/3] Waiting for server to be ready..."
sleep 2

# 运行测试客户端 / Run test client
echo "[3/3] 运行全面测试..."
echo "[3/3] Running comprehensive tests..."
echo ""

set +e
go run cmd/client/main.go cmd/client/comprehensive_suite.go -verbose
TEST_EXIT_CODE=$?
set -e

# 停止服务器 / Stop server
echo ""
echo "停止 MCP 服务器..."
echo "Stopping MCP server..."
kill $MCP_PID 2>/dev/null || true
wait $MCP_PID 2>/dev/null || true

echo ""
if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "========================================"
    echo "✅ 所有测试通过！"
    echo "✅ All tests passed!"
    echo "========================================"
else
    echo "========================================"
    echo "❌ 测试失败，退出码: $TEST_EXIT_CODE"
    echo "❌ Tests failed with exit code: $TEST_EXIT_CODE"
    echo "========================================"
fi

exit $TEST_EXIT_CODE

