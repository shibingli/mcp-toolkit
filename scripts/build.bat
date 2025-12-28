@echo off
REM 跨平台构建脚本 (Windows) / Cross-platform build script (Windows)
REM 用于编译 Windows、Linux、macOS 的二进制文件

setlocal enabledelayedexpansion

REM 项目信息 / Project information
set PROJECT_NAME=mcp-toolkit
set VERSION=%VERSION%
if "%VERSION%"=="" (
    for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%i
    if "!VERSION!"=="" set VERSION=dev
)

for /f "tokens=*" %%i in ('powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH:mm:ss' -AsUTC"') do set BUILD_TIME=%%i
for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%i
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown

REM 构建标签 / Build tags
set BUILD_TAGS=sonic

REM 输出目录 / Output directory
set OUTPUT_DIR=dist
if exist %OUTPUT_DIR% rmdir /s /q %OUTPUT_DIR%
mkdir %OUTPUT_DIR%

REM 构建信息 / Build information
set LDFLAGS=-s -w -X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME% -X main.GitCommit=%GIT_COMMIT%

echo =========================================
echo Building %PROJECT_NAME% %VERSION%
echo Build Time: %BUILD_TIME%
echo Git Commit: %GIT_COMMIT%
echo =========================================
echo.

REM 构建 Windows amd64
call :build windows amd64

REM 构建 Windows arm64
call :build windows arm64

REM 构建 Linux amd64
call :build linux amd64

REM 构建 Linux arm64
call :build linux arm64

REM 构建 macOS amd64
call :build darwin amd64

REM 构建 macOS arm64
call :build darwin arm64

echo.
echo =========================================
echo Build completed successfully!
echo Output directory: %OUTPUT_DIR%
echo =========================================
echo.
echo Generated files:
dir /b %OUTPUT_DIR%\*.zip %OUTPUT_DIR%\*.tar.gz 2>nul

goto :eof

:build
set GOOS=%1
set GOARCH=%2
set OUTPUT_NAME=%PROJECT_NAME%

if "%GOOS%"=="windows" set OUTPUT_NAME=%OUTPUT_NAME%.exe

set OUTPUT_PATH=%OUTPUT_DIR%\%PROJECT_NAME%-%GOOS%-%GOARCH%
mkdir %OUTPUT_PATH%

echo Building for %GOOS%/%GOARCH%...

REM 构建 / Build
set CGO_ENABLED=0
go build -tags="%BUILD_TAGS%" -ldflags="%LDFLAGS%" -o "%OUTPUT_PATH%\%OUTPUT_NAME%" .

REM 复制文件 / Copy files
copy README.md %OUTPUT_PATH%\ >nul
copy LICENSE %OUTPUT_PATH%\ >nul

REM 创建压缩包 / Create archive
cd %OUTPUT_DIR%
if "%GOOS%"=="windows" (
    powershell -Command "Compress-Archive -Path '%PROJECT_NAME%-%GOOS%-%GOARCH%' -DestinationPath '%PROJECT_NAME%-%VERSION%-%GOOS%-%GOARCH%.zip' -Force"
) else (
    tar -czf "%PROJECT_NAME%-%VERSION%-%GOOS%-%GOARCH%.tar.gz" "%PROJECT_NAME%-%GOOS%-%GOARCH%"
)
cd ..

echo √ Built %GOOS%/%GOARCH%
echo.

goto :eof

