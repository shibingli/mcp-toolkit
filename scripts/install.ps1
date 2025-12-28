# MCP Toolkit 安装脚本 (Windows) / MCP Toolkit installation script (Windows)

param(
    [string]$Version = "",
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\mcp-toolkit"
)

$ErrorActionPreference = "Stop"

# 配置 / Configuration
$Repo = "shibingli/mcp-toolkit"  # 请替换为你的 GitHub 仓库 / Replace with your GitHub repository
$BinaryName = "mcp-toolkit.exe"

# 颜色输出函数 / Color output functions
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
    exit 1
}

# 检测架构 / Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { Write-Error-Custom "Unsupported architecture: $arch" }
    }
}

# 获取最新版本 / Get latest version
function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        return $response.tag_name
    }
    catch {
        Write-Error-Custom "Failed to get latest version: $_"
    }
}

# 下载并安装 / Download and install
function Install-McpToolkit {
    Write-Info "Detecting platform..."
    $arch = Get-Architecture
    Write-Info "Architecture: $arch"
    
    if ([string]::IsNullOrEmpty($Version)) {
        Write-Info "Getting latest version..."
        $Version = Get-LatestVersion
    }
    Write-Info "Version: $Version"
    
    $downloadUrl = "https://github.com/$Repo/releases/download/$Version/mcp-toolkit-$Version-windows-$arch.zip"
    $tempDir = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_ }
    $zipFile = Join-Path $tempDir "mcp-toolkit.zip"
    
    Write-Info "Downloading from $downloadUrl..."
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipFile
    }
    catch {
        Write-Error-Custom "Failed to download: $_"
    }
    
    Write-Info "Extracting..."
    Expand-Archive -Path $zipFile -DestinationPath $tempDir -Force
    
    Write-Info "Installing to $InstallDir..."
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    $extractedDir = Get-ChildItem -Path $tempDir -Directory | Where-Object { $_.Name -like "mcp-toolkit-*" } | Select-Object -First 1
    $sourceBinary = Join-Path $extractedDir.FullName $BinaryName
    $destBinary = Join-Path $InstallDir $BinaryName
    
    Copy-Item -Path $sourceBinary -Destination $destBinary -Force
    
    Remove-Item -Path $tempDir -Recurse -Force
    
    Write-Info "Installation completed!"
    Write-Info "Binary installed to: $destBinary"
    
    # 检查是否在 PATH 中 / Check if in PATH
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$InstallDir*") {
        Write-Warn "$InstallDir is not in your PATH"
        Write-Host ""
        Write-Host "Do you want to add it to your PATH? (Y/N)" -ForegroundColor Yellow
        $response = Read-Host
        if ($response -eq "Y" -or $response -eq "y") {
            $newPath = "$currentPath;$InstallDir"
            [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
            Write-Info "Added to PATH. Please restart your terminal."
        }
        else {
            Write-Warn "You can manually add the following to your PATH:"
            Write-Host "    $InstallDir"
        }
    }
    
    Write-Info "Run '$BinaryName --help' to get started"
}

# 卸载 / Uninstall
function Uninstall-McpToolkit {
    $binaryPath = Join-Path $InstallDir $BinaryName
    if (Test-Path $binaryPath) {
        Remove-Item -Path $binaryPath -Force
        Write-Info "Uninstalled successfully"
        
        # 如果目录为空，删除目录 / Remove directory if empty
        if ((Get-ChildItem -Path $InstallDir).Count -eq 0) {
            Remove-Item -Path $InstallDir -Force
        }
    }
    else {
        Write-Warn "Binary not found at $binaryPath"
    }
}

# 主函数 / Main function
if ($args[0] -eq "uninstall") {
    Uninstall-McpToolkit
}
else {
    Install-McpToolkit
}

