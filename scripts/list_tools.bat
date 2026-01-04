@echo off
REM MCP Server Tools List Script
REM Usage: list_tools.bat [transport] [detailed]
REM   transport: stdio, http, sse (default: http)
REM   detailed: show detailed information

setlocal

set TRANSPORT=%1
set DETAILED=%2

REM Default to http if no transport specified
if "%TRANSPORT%"=="" set TRANSPORT=http
if "%TRANSPORT%"=="detailed" (
    set DETAILED=detailed
    set TRANSPORT=http
)

echo ========================================
echo MCP Server Tools List
echo ========================================
echo Transport: %TRANSPORT%
echo.

if "%TRANSPORT%"=="stdio" (
    REM Stdio mode - launch server as subprocess
    echo Using Stdio transport...
    echo.
    if "%DETAILED%"=="detailed" (
        go run cmd/client/list_tools.go -transport stdio -command .\mcp-toolkit.exe -args "-transport,stdio" -detailed
    ) else (
        go run cmd/client/list_tools.go -transport stdio -command .\mcp-toolkit.exe -args "-transport,stdio"
    )
) else if "%TRANSPORT%"=="sse" (
    REM SSE mode
    echo Checking if SSE server is running...
    tasklist /FI "IMAGENAME eq mcp-toolkit.exe" 2>NUL | find /I /N "mcp-toolkit.exe">NUL
    if "%ERRORLEVEL%"=="0" (
        echo Server is already running.
    ) else (
        echo Starting SSE server...
        start "MCP Server SSE" mcp-toolkit.exe -transport sse -sse-port 8081
        echo Waiting for server to start...
        timeout /t 3 /nobreak >nul
    )
    echo.
    if "%DETAILED%"=="detailed" (
        go run cmd/client/list_tools.go -transport sse -port 8081 -path /sse -detailed
    ) else (
        go run cmd/client/list_tools.go -transport sse -port 8081 -path /sse
    )
) else (
    REM HTTP mode (default)
    echo Checking if HTTP server is running...
    tasklist /FI "IMAGENAME eq mcp-toolkit.exe" 2>NUL | find /I /N "mcp-toolkit.exe">NUL
    if "%ERRORLEVEL%"=="0" (
        echo Server is already running.
    ) else (
        echo Starting HTTP server...
        start "MCP Server HTTP" mcp-toolkit.exe -transport http -http-port 8080
        echo Waiting for server to start...
        timeout /t 3 /nobreak >nul
    )
    echo.
    if "%DETAILED%"=="detailed" (
        go run cmd/client/list_tools.go -transport http -detailed
    ) else (
        go run cmd/client/list_tools.go -transport http
    )
)

echo.
echo Done!
endlocal

