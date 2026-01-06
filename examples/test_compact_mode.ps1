# 测试精简模式的效果 / Test compact mode effect
# 此脚本用于对比普通模式和精简模式的token消耗差异

Write-Host "=== MCP工具列表精简模式测试 ===" -ForegroundColor Cyan
Write-Host ""

# 启动MCP服务器（HTTP模式）
Write-Host "启动MCP服务器..." -ForegroundColor Yellow
$serverProcess = Start-Process -FilePath ".\mcp-toolkit.exe" -ArgumentList "-transport", "http", "-http-port", "8080" -PassThru -WindowStyle Hidden
Start-Sleep -Seconds 2

try
{
    # 测试普通模式
    Write-Host "`n1. 测试普通模式（完整schema）..." -ForegroundColor Green
    $normalRequest = @{
        jsonrpc = "2.0"
        id = 1
        method = "tools/list"
        params = @{ }
    } | ConvertTo-Json -Depth 10

    $normalResponse = Invoke-RestMethod -Uri "http://127.0.0.1:8080/mcp" -Method Post -Body $normalRequest -ContentType "application/json"
    $normalJson = $normalResponse | ConvertTo-Json -Depth 20
    $normalSize = [System.Text.Encoding]::UTF8.GetByteCount($normalJson)
    $normalTokens = [Math]::Ceiling($normalSize / 4)  # 粗略估算：1 token ≈ 4 bytes

    Write-Host "  - 响应大小: $normalSize bytes" -ForegroundColor White
    Write-Host "  - 估算tokens: ~$normalTokens tokens" -ForegroundColor White
    Write-Host "  - 工具数量: $( $normalResponse.result.tools.Count )" -ForegroundColor White

    # 测试精简模式
    Write-Host "`n2. 测试精简模式（compact=true）..." -ForegroundColor Green
    $compactRequest = @{
        jsonrpc = "2.0"
        id = 2
        method = "tools/list"
        params = @{
            compact = $true
        }
    } | ConvertTo-Json -Depth 10

    $compactResponse = Invoke-RestMethod -Uri "http://127.0.0.1:8080/mcp" -Method Post -Body $compactRequest -ContentType "application/json"
    $compactJson = $compactResponse | ConvertTo-Json -Depth 20
    $compactSize = [System.Text.Encoding]::UTF8.GetByteCount($compactJson)
    $compactTokens = [Math]::Ceiling($compactSize / 4)

    Write-Host "  - 响应大小: $compactSize bytes" -ForegroundColor White
    Write-Host "  - 估算tokens: ~$compactTokens tokens" -ForegroundColor White
    Write-Host "  - 工具数量: $( $compactResponse.result.tools.Count )" -ForegroundColor White

    # 计算优化效果
    Write-Host "`n3. 优化效果对比..." -ForegroundColor Green
    $reduction = $normalSize - $compactSize
    $reductionPercent = [Math]::Round(($reduction / $normalSize) * 100, 2)
    $tokenReduction = $normalTokens - $compactTokens

    Write-Host "  - 大小减少: $reduction bytes ($reductionPercent%)" -ForegroundColor Cyan
    Write-Host "  - Token减少: ~$tokenReduction tokens" -ForegroundColor Cyan

    if ($compactTokens -lt 3072)
    {
        Write-Host "  - ✓ 精简模式可以在3072 token限制内使用" -ForegroundColor Green
    }
    else
    {
        Write-Host "  - ✗ 精简模式仍超过3072 token限制" -ForegroundColor Red
    }

    # 保存示例到文件
    Write-Host "`n4. 保存示例到文件..." -ForegroundColor Green
    $normalJson | Out-File -FilePath "examples\tools_list_normal.json" -Encoding UTF8
    $compactJson | Out-File -FilePath "examples\tools_list_compact.json" -Encoding UTF8
    Write-Host "  - 普通模式: examples\tools_list_normal.json" -ForegroundColor White
    Write-Host "  - 精简模式: examples\tools_list_compact.json" -ForegroundColor White

    Write-Host "`n=== 测试完成 ===" -ForegroundColor Cyan

}
finally
{
    # Stop server
    Write-Host "`nStopping MCP server..." -ForegroundColor Yellow
    Stop-Process -Id $serverProcess.Id -Force
}

