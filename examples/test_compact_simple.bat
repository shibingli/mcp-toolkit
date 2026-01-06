@echo off
chcp 65001 >nul
echo === Testing Compact Mode ===
echo.

echo Starting MCP server without session management...
start /B mcp-toolkit.exe -transport http -http-port 8080 -http-disable-session
timeout /t 3 /nobreak >nul

echo.
echo 1. Testing normal mode...
curl -X POST http://127.0.0.1:8080/mcp ^
  -H "Content-Type: application/json" ^
  -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\",\"params\":{}}" ^
  -o examples\tools_list_normal.json

echo.
echo 2. Testing compact mode...
curl -X POST http://127.0.0.1:8080/mcp ^
  -H "Content-Type: application/json" ^
  -d "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/list\",\"params\":{\"compact\":true}}" ^
  -o examples\tools_list_compact.json

echo.
echo 3. Comparing file sizes...
for %%A in (examples\tools_list_normal.json) do set normal_size=%%~zA
for %%A in (examples\tools_list_compact.json) do set compact_size=%%~zA

echo Normal mode: %normal_size% bytes
echo Compact mode: %compact_size% bytes

echo.
echo Stopping server...
taskkill /F /IM mcp-toolkit.exe >nul 2>&1

echo.
echo === Test Complete ===
echo Check examples\tools_list_normal.json and examples\tools_list_compact.json
pause

