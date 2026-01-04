#!/bin/bash
# MCP Server Tools List Script
# Usage: ./list_tools.sh [transport] [detailed]
#   transport: stdio, http, sse (default: http)
#   detailed: show detailed information

TRANSPORT=${1:-http}
DETAILED=$2

# If first arg is "detailed", use http transport
if [ "$TRANSPORT" = "detailed" ]; then
    DETAILED="detailed"
    TRANSPORT="http"
fi

echo "========================================"
echo "MCP Server Tools List"
echo "========================================"
echo "Transport: $TRANSPORT"
echo ""

if [ "$TRANSPORT" = "stdio" ]; then
    # Stdio mode - launch server as subprocess
    echo "Using Stdio transport..."
    echo ""
    if [ "$DETAILED" = "detailed" ]; then
        go run cmd/client/list_tools.go -transport stdio -command ./mcp-toolkit -args "-transport,stdio" -detailed
    else
        go run cmd/client/list_tools.go -transport stdio -command ./mcp-toolkit -args "-transport,stdio"
    fi
elif [ "$TRANSPORT" = "sse" ]; then
    # SSE mode
    echo "Checking if SSE server is running..."
    if pgrep -x "mcp-toolkit" > /dev/null; then
        echo "Server is already running."
    else
        echo "Starting SSE server..."
        ./mcp-toolkit -transport sse -sse-port 8081 &
        echo "Waiting for server to start..."
        sleep 3
    fi
    echo ""
    if [ "$DETAILED" = "detailed" ]; then
        go run cmd/client/list_tools.go -transport sse -port 8081 -path /sse -detailed
    else
        go run cmd/client/list_tools.go -transport sse -port 8081 -path /sse
    fi
else
    # HTTP mode (default)
    echo "Checking if HTTP server is running..."
    if pgrep -x "mcp-toolkit" > /dev/null; then
        echo "Server is already running."
    else
        echo "Starting HTTP server..."
        ./mcp-toolkit -transport http -http-port 8080 &
        echo "Waiting for server to start..."
        sleep 3
    fi
    echo ""
    if [ "$DETAILED" = "detailed" ]; then
        go run cmd/client/list_tools.go -transport http -detailed
    else
        go run cmd/client/list_tools.go -transport http
    fi
fi

echo ""
echo "Done!"

