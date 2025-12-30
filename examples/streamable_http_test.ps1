# Streamable HTTP 功能测试脚本 (PowerShell)
# Streamable HTTP Feature Test Script (PowerShell)

$ErrorActionPreference = "Stop"

# 配置 / Configuration
$Host_Addr = "localhost"
$Port = "8080"
$BaseUrl = "http://${Host_Addr}:${Port}/mcp"
$ProtocolVersion = "2025-06-18"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Streamable HTTP 功能测试" -ForegroundColor Cyan
Write-Host "Streamable HTTP Feature Test" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# 1. 初始化会话 / Initialize session
Write-Host "1. 初始化会话 / Initialize session" -ForegroundColor Yellow
Write-Host "---"

$InitBody = @{
    jsonrpc = "2.0"
    id = 1
    method = "initialize"
    params = @{
        protocolVersion = $ProtocolVersion
        capabilities = @{ }
        clientInfo = @{
            name = "test-client"
            version = "1.0.0"
        }
    }
} | ConvertTo-Json -Depth 10

$InitHeaders = @{
    "Content-Type" = "application/json"
    "Accept" = "application/json"
    "MCP-Protocol-Version" = $ProtocolVersion
}

try
{
    $InitResponse = Invoke-WebRequest -Uri $BaseUrl -Method POST -Headers $InitHeaders -Body $InitBody
    $SessionId = $InitResponse.Headers["Mcp-Session-Id"]

    if ( [string]::IsNullOrEmpty($SessionId))
    {
        Write-Host "❌ 失败: 未获取到会话ID / Failed: No session ID received" -ForegroundColor Red
        exit 1
    }

    Write-Host "✅ 成功: 会话ID = $SessionId" -ForegroundColor Green
    Write-Host "✅ Success: Session ID = $SessionId" -ForegroundColor Green
}
catch
{
    Write-Host "❌ 失败: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 2. 列出工具 / List tools
Write-Host "2. 列出工具 / List tools" -ForegroundColor Yellow
Write-Host "---"

$ToolsBody = @{
    jsonrpc = "2.0"
    id = 2
    method = "tools/list"
} | ConvertTo-Json

$ToolsHeaders = @{
    "Content-Type" = "application/json"
    "Accept" = "application/json"
    "Mcp-Session-Id" = $SessionId
}

try
{
    $ToolsResponse = Invoke-RestMethod -Uri $BaseUrl -Method POST -Headers $ToolsHeaders -Body $ToolsBody
    $ToolsCount = $ToolsResponse.result.tools.Count

    if ($ToolsCount -gt 0)
    {
        Write-Host "✅ 成功: 获取到 $ToolsCount 个工具" -ForegroundColor Green
        Write-Host "✅ Success: Retrieved $ToolsCount tools" -ForegroundColor Green
    }
    else
    {
        Write-Host "❌ 失败: 未获取到工具列表" -ForegroundColor Red
        exit 1
    }
}
catch
{
    Write-Host "❌ 失败: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 3. 调用工具 / Call tool
Write-Host "3. 调用工具 / Call tool (get_current_time)" -ForegroundColor Yellow
Write-Host "---"

$CallBody = @{
    jsonrpc = "2.0"
    id = 3
    method = "tools/call"
    params = @{
        name = "get_current_time"
    }
} | ConvertTo-Json

try
{
    $CallResponse = Invoke-RestMethod -Uri $BaseUrl -Method POST -Headers $ToolsHeaders -Body $CallBody

    if ($CallResponse.result)
    {
        Write-Host "✅ 成功: 工具调用成功" -ForegroundColor Green
        Write-Host "✅ Success: Tool call succeeded" -ForegroundColor Green
        Write-Host "响应 / Response:"
        $CallResponse | ConvertTo-Json -Depth 10
    }
    else
    {
        Write-Host "❌ 失败: 工具调用失败" -ForegroundColor Red
        exit 1
    }
}
catch
{
    Write-Host "❌ 失败: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 4. 测试 SSE 流 / Test SSE streaming
Write-Host "4. 测试 SSE 流 / Test SSE streaming" -ForegroundColor Yellow
Write-Host "---"

$SSEHeaders = @{
    "Content-Type" = "application/json"
    "Accept" = "text/event-stream"
    "Mcp-Session-Id" = $SessionId
}

try
{
    $SSEResponse = Invoke-WebRequest -Uri $BaseUrl -Method POST -Headers $SSEHeaders -Body $ToolsBody

    if ($SSEResponse.Content -match "data:")
    {
        Write-Host "✅ 成功: SSE 流响应正常" -ForegroundColor Green
        Write-Host "✅ Success: SSE stream response OK" -ForegroundColor Green
        Write-Host "SSE 响应 / SSE Response:"
        Write-Host $SSEResponse.Content.Substring(0,[Math]::Min(200, $SSEResponse.Content.Length))
    }
    else
    {
        Write-Host "❌ 失败: SSE 流响应异常" -ForegroundColor Red
        exit 1
    }
}
catch
{
    Write-Host "❌ 失败: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 5. 终止会话 / Terminate session
Write-Host "5. 终止会话 / Terminate session" -ForegroundColor Yellow
Write-Host "---"

$DeleteHeaders = @{
    "Mcp-Session-Id" = $SessionId
}

try
{
    $DeleteResponse = Invoke-WebRequest -Uri $BaseUrl -Method DELETE -Headers $DeleteHeaders

    if ($DeleteResponse.StatusCode -eq 200)
    {
        Write-Host "✅ 成功: 会话已终止" -ForegroundColor Green
        Write-Host "✅ Success: Session terminated" -ForegroundColor Green
    }
    else
    {
        Write-Host "❌ 失败: 会话终止失败 (HTTP $( $DeleteResponse.StatusCode ))" -ForegroundColor Red
        exit 1
    }
}
catch
{
    Write-Host "❌ 失败: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "✅ 所有测试通过!" -ForegroundColor Green
Write-Host "✅ All tests passed!" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan

