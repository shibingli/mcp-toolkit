# MCP Toolkit Installation Script for Windows
# PowerShell script to install, update, or uninstall MCP Toolkit

param(
    [Parameter(Position=0)]
    [ValidateSet('install', 'uninstall', 'update', 'version', 'help')]
    [string]$Command = 'install',
    
    [string]$Version = '',
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\mcp-toolkit",
    [string]$Repo = 'shibingli/mcp-toolkit',
    [switch]$SkipChecksum,
    [switch]$Debug
)

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$BinaryName = 'mcp-toolkit.exe'

# Color output functions
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-ErrorMsg {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
    exit 1
}

function Write-Debug {
    param([string]$Message)
    if ($Debug) {
        Write-Host "[DEBUG] $Message" -ForegroundColor Blue
    }
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    Write-Debug "Detected architecture: $arch"
    switch ($arch) {
        'AMD64' { return 'amd64' }
        'ARM64' { return 'arm64' }
        default { Write-ErrorMsg "Unsupported architecture: $arch" }
    }
}

# Get latest version from GitHub
function Get-LatestVersion {
    Write-Debug "Fetching latest version from GitHub API..."
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -ErrorAction Stop
        $version = $response.tag_name
        Write-Debug "Latest version: $version"
        return $version
    }
    catch {
        Write-ErrorMsg "Failed to get latest version: $_"
    }
}

# Get installed version
function Get-InstalledVersion {
    $binaryPath = Join-Path $InstallDir $BinaryName
    if (Test-Path $binaryPath) {
        try {
            $output = & $binaryPath --version 2>&1 | Select-Object -First 1
            if ($output -match '(\d+\.\d+\.\d+.*)') {
                return $matches[1]
            }
        }
        catch {
            Write-Debug "Failed to get installed version: $_"
        }
    }
    return $null
}

# Verify checksum
function Test-Checksum {
    param(
        [string]$FilePath,
        [string]$ExpectedChecksum
    )
    
    if ($SkipChecksum) {
        Write-Debug "Checksum verification skipped"
        return $true
    }
    
    if ([string]::IsNullOrEmpty($ExpectedChecksum)) {
        Write-Warn "No checksum provided, skipping verification"
        return $true
    }
    
    Write-Info "Verifying checksum..."
    $actualChecksum = (Get-FileHash -Path $FilePath -Algorithm SHA256).Hash.ToLower()
    
    if ($actualChecksum -eq $ExpectedChecksum.ToLower()) {
        Write-Info "Checksum verified successfully"
        return $true
    }
    else {
        Write-ErrorMsg "Checksum verification failed!`nExpected: $ExpectedChecksum`nActual: $actualChecksum"
    }
}

# Install function
function Install-McpToolkit {
    # Check installed version
    $installedVersion = Get-InstalledVersion
    if ($installedVersion) {
        Write-Info "Current installed version: $installedVersion"
    }
    
    Write-Info "Detecting platform..."
    $arch = Get-Architecture
    Write-Info "Architecture: $arch"
    
    if ([string]::IsNullOrEmpty($Version)) {
        Write-Info "Getting latest version..."
        $Version = Get-LatestVersion
    }
    Write-Info "Target version: $Version"
    
    # Check if update is needed
    $versionWithoutV = $Version -replace '^v', ''
    if ($installedVersion -and ($installedVersion -eq $versionWithoutV)) {
        Write-Info "Already up to date (version $installedVersion)"
        return
    }
    
    $archiveName = "mcp-toolkit-$Version-windows-$arch.zip"
    $downloadUrl = "https://github.com/$Repo/releases/download/$Version/$archiveName"
    $checksumsUrl = "https://github.com/$Repo/releases/download/$Version/checksums.txt"
    
    $tempDir = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_ }
    $zipFile = Join-Path $tempDir $archiveName

    try {
        Write-Info "Downloading from $downloadUrl..."
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipFile -ErrorAction Stop

        # Download and verify checksum
        if (-not $SkipChecksum) {
            Write-Debug "Downloading checksums..."
            try {
                $checksumsFile = Join-Path $tempDir 'checksums.txt'
                Invoke-WebRequest -Uri $checksumsUrl -OutFile $checksumsFile -ErrorAction Stop
                $expectedChecksum = (Get-Content $checksumsFile | Select-String $archiveName).Line -split '\s+' | Select-Object -First 1
                Test-Checksum -FilePath $zipFile -ExpectedChecksum $expectedChecksum
            }
            catch {
                Write-Warn "Could not download checksums file, skipping verification"
            }
        }

        Write-Info "Extracting..."
        Expand-Archive -Path $zipFile -DestinationPath $tempDir -Force

        Write-Info "Installing to $InstallDir..."
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        $extractedDir = Get-ChildItem -Path $tempDir -Directory | Where-Object { $_.Name -like 'mcp-toolkit-*' } | Select-Object -First 1
        if (-not $extractedDir) {
            Write-ErrorMsg "Failed to find extracted directory"
        }

        $sourceBinary = Join-Path $extractedDir.FullName $BinaryName
        if (-not (Test-Path $sourceBinary)) {
            Write-ErrorMsg "Binary not found in extracted archive"
        }

        $destBinary = Join-Path $InstallDir $BinaryName

        # Backup old version
        if (Test-Path $destBinary) {
            Write-Info "Backing up old version..."
            Copy-Item -Path $destBinary -Destination "$destBinary.backup" -Force
        }

        # Copy new version
        Copy-Item -Path $sourceBinary -Destination $destBinary -Force

        # Verify installation
        try {
            $null = & $destBinary --version 2>&1
        }
        catch {
            Write-ErrorMsg "Installation verification failed. Binary may be corrupted."
        }

        # Remove backup
        if (Test-Path "$destBinary.backup") {
            Remove-Item "$destBinary.backup" -Force
        }

        Write-Info "Installation completed successfully!"
        Write-Info "Binary installed to: $destBinary"

        # Show version info
        $newVersion = Get-InstalledVersion
        Write-Info "Installed version: $newVersion"

        # Check if in PATH
        $currentPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        if ($currentPath -notlike "*$InstallDir*") {
            Write-Warn "$InstallDir is not in your PATH"
            Write-Host ""
            Write-Host "Do you want to add it to your PATH? (Y/N)" -ForegroundColor Yellow
            $response = Read-Host
            if ($response -eq 'Y' -or $response -eq 'y') {
                $newPath = "$currentPath;$InstallDir"
                [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
                Write-Info "Added to PATH. Please restart your terminal."
            }
            else {
                Write-Warn "You can manually add the following to your PATH:"
                Write-Host "    $InstallDir"
            }
        }
        else {
            Write-Info "Run '$BinaryName --help' to get started"
        }
    }
    catch {
        Write-ErrorMsg "Installation failed: $_"
    }
    finally {
        if (Test-Path $tempDir) {
            Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

# Uninstall function
function Uninstall-McpToolkit {
    $binaryPath = Join-Path $InstallDir $BinaryName
    if (Test-Path $binaryPath) {
        $version = Get-InstalledVersion
        Write-Info "Uninstalling $BinaryName $version..."
        Remove-Item -Path $binaryPath -Force
        Write-Info "Uninstalled successfully"

        # Remove directory if empty
        if ((Get-ChildItem -Path $InstallDir -ErrorAction SilentlyContinue).Count -eq 0) {
            Remove-Item -Path $InstallDir -Force -ErrorAction SilentlyContinue
        }
    }
    else {
        Write-Warn "Binary not found at $binaryPath"
        Write-Warn "Nothing to uninstall"
    }
}

# Update function
function Update-McpToolkit {
    $installedVersion = Get-InstalledVersion
    if (-not $installedVersion) {
        Write-ErrorMsg "$BinaryName is not installed. Run 'install' first."
    }

    Write-Info "Checking for updates..."
    $latestVersion = Get-LatestVersion
    $latestVersionWithoutV = $latestVersion -replace '^v', ''

    if ($installedVersion -eq $latestVersionWithoutV) {
        Write-Info "Already up to date (version $installedVersion)"
    }
    else {
        Write-Info "Update available: $installedVersion -> $latestVersionWithoutV"
        $script:Version = $latestVersion
        Install-McpToolkit
    }
}

# Show version function
function Show-Version {
    $version = Get-InstalledVersion
    if ($version) {
        Write-Host "$BinaryName version $version"
        Write-Host "Installed at: $(Join-Path $InstallDir $BinaryName)"
    }
    else {
        Write-Host "$BinaryName is not installed"
        exit 1
    }
}

# Show help function
function Show-Help {
    @"
MCP Toolkit Installation Script for Windows

Usage: .\install.ps1 [COMMAND] [OPTIONS]

Commands:
    install     Install $BinaryName (default)
    uninstall   Uninstall $BinaryName
    update      Update $BinaryName to the latest version
    version     Show installed version
    help        Show this help message

Parameters:
    -Version <string>       Specify version to install (e.g., v1.0.0)
    -InstallDir <string>    Installation directory (default: %LOCALAPPDATA%\Programs\mcp-toolkit)
    -Repo <string>          GitHub repository (default: shibingli/mcp-toolkit)
    -SkipChecksum           Skip checksum verification
    -Debug                  Enable debug output

Examples:
    # Install latest version
    .\install.ps1 install

    # Install specific version
    .\install.ps1 install -Version v1.0.0

    # Install to custom directory
    .\install.ps1 install -InstallDir "C:\Tools\mcp-toolkit"

    # Update to latest version
    .\install.ps1 update

    # Uninstall
    .\install.ps1 uninstall

    # Show version
    .\install.ps1 version

For more information, visit: https://github.com/$Repo
"@
}

# Main execution
switch ($Command) {
    'install' {
        Install-McpToolkit
    }
    'uninstall' {
        Uninstall-McpToolkit
    }
    'update' {
        Update-McpToolkit
    }
    'version' {
        Show-Version
    }
    'help' {
        Show-Help
    }
    default {
        Write-ErrorMsg "Unknown command: $Command`nRun '.\install.ps1 help' for usage information"
    }
}

