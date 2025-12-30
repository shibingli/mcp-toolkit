#!/bin/bash

# Streamable HTTP 功能测试脚本
# Streamable HTTP Feature Test Script

set -e

# 配置 / Configuration
HOST="localhost"
PORT="8080"
BASE_URL="http://${HOST}:${PORT}/mcp"
PROTOCOL_VERSION="2025-06-18"

echo "=========================================="
echo "Streamable HTTP 功能测试"
echo "Streamable HTTP Feature Test"
echo "=========================================="
echo ""

# 1. 初始化会话 / Initialize session
echo "1. 初始化会话 / Initialize session"
echo "---"

INIT_RESPONSE=$(curl -s -i -X POST "${BASE_URL}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "MCP-Protocol-Version: ${PROTOCOL_VERSION}" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "'"${PROTOCOL_VERSION}"'",
      "capabilities": {},
      "clientInfo": {
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  }')

# 提取会话ID / Extract session ID
SESSION_ID=$(echo "$INIT_RESPONSE" | grep -i "Mcp-Session-Id:" | cut -d' ' -f2 | tr -d '\r')

if [ -z "$SESSION_ID" ]; then
  echo "❌ 失败: 未获取到会话ID / Failed: No session ID received"
  exit 1
fi

echo "✅ 成功: 会话ID = $SESSION_ID"
echo "✅ Success: Session ID = $SESSION_ID"
echo ""

# 2. 列出工具 / List tools
echo "2. 列出工具 / List tools"
echo "---"

TOOLS_RESPONSE=$(curl -s -X POST "${BASE_URL}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "Mcp-Session-Id: ${SESSION_ID}" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }')

TOOLS_COUNT=$(echo "$TOOLS_RESPONSE" | grep -o '"name"' | wc -l)

if [ "$TOOLS_COUNT" -gt 0 ]; then
  echo "✅ 成功: 获取到 $TOOLS_COUNT 个工具"
  echo "✅ Success: Retrieved $TOOLS_COUNT tools"
else
  echo "❌ 失败: 未获取到工具列表"
  echo "❌ Failed: No tools retrieved"
  exit 1
fi
echo ""

# 3. 调用工具 / Call tool
echo "3. 调用工具 / Call tool (get_current_time)"
echo "---"

TOOL_RESPONSE=$(curl -s -X POST "${BASE_URL}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "Mcp-Session-Id: ${SESSION_ID}" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "get_current_time"
    }
  }')

if echo "$TOOL_RESPONSE" | grep -q '"result"'; then
  echo "✅ 成功: 工具调用成功"
  echo "✅ Success: Tool call succeeded"
  echo "响应 / Response:"
  echo "$TOOL_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$TOOL_RESPONSE"
else
  echo "❌ 失败: 工具调用失败"
  echo "❌ Failed: Tool call failed"
  exit 1
fi
echo ""

# 4. 测试 SSE 流 / Test SSE streaming
echo "4. 测试 SSE 流 / Test SSE streaming"
echo "---"

SSE_RESPONSE=$(curl -s -X POST "${BASE_URL}" \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -H "Mcp-Session-Id: ${SESSION_ID}" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/list"
  }' | head -n 5)

if echo "$SSE_RESPONSE" | grep -q "data:"; then
  echo "✅ 成功: SSE 流响应正常"
  echo "✅ Success: SSE stream response OK"
  echo "SSE 响应 / SSE Response:"
  echo "$SSE_RESPONSE"
else
  echo "❌ 失败: SSE 流响应异常"
  echo "❌ Failed: SSE stream response error"
  exit 1
fi
echo ""

# 5. 终止会话 / Terminate session
echo "5. 终止会话 / Terminate session"
echo "---"

DELETE_RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "${BASE_URL}" \
  -H "Mcp-Session-Id: ${SESSION_ID}")

HTTP_CODE=$(echo "$DELETE_RESPONSE" | tail -n 1)

if [ "$HTTP_CODE" = "200" ]; then
  echo "✅ 成功: 会话已终止"
  echo "✅ Success: Session terminated"
else
  echo "❌ 失败: 会话终止失败 (HTTP $HTTP_CODE)"
  echo "❌ Failed: Session termination failed (HTTP $HTTP_CODE)"
  exit 1
fi
echo ""

echo "=========================================="
echo "✅ 所有测试通过!"
echo "✅ All tests passed!"
echo "=========================================="

