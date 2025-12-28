@echo off
REM Cross-platform build script for Windows
REM Builds binaries for Windows, Linux, and macOS

setlocal enabledelayedexpansion

REM Project information
set PROJECT_NAME=mcp-toolkit
set VERSION=%VERSION%
if "%VERSION%"=="" (
    for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%i
    if "!VERSION!"=="" set VERSION=dev
)

for /f "tokens=*" %%i in ('powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH:mm:ss' -AsUTC"') do set BUILD_TIME=%%i
for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%i
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown

REM Build tags
set BUILD_TAGS=sonic

REM Output directory
set OUTPUT_DIR=dist
if exist %OUTPUT_DIR% rmdir /s /q %OUTPUT_DIR%
mkdir %OUTPUT_DIR%

REM Build information
set LDFLAGS=-s -w -X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME% -X main.GitCommit=%GIT_COMMIT%

echo =========================================
echo Building %PROJECT_NAME% v%VERSION%
echo Build Time: %BUILD_TIME%
echo Git Commit: %GIT_COMMIT%
echo =========================================
echo.
echo Building for all platforms...
echo.

REM Build Windows amd64
call :build windows amd64

REM Build Windows arm64
call :build windows arm64

REM Build Linux amd64
call :build linux amd64

REM Build Linux arm64
call :build linux arm64

REM Build macOS amd64
call :build darwin amd64

REM Build macOS arm64
call :build darwin arm64

REM Generate checksums
echo.
echo Generating checksums...
cd %OUTPUT_DIR%
powershell -Command "Get-ChildItem -Include *.zip,*.tar.gz -Recurse | Get-FileHash -Algorithm SHA256 | ForEach-Object { $_.Hash.ToLower() + '  ' + $_.Path.Split('\')[-1] } | Out-File -Encoding ASCII checksums.txt"
cd ..

echo.
echo =========================================
echo Build completed successfully!
echo Output directory: %OUTPUT_DIR%
echo =========================================
echo.
echo Generated files:
echo.
echo Windows packages (zip):
dir /b %OUTPUT_DIR%\*windows*.zip 2>nul
echo.
echo Linux/macOS packages (tar.gz):
dir /b %OUTPUT_DIR%\*linux*.tar.gz %OUTPUT_DIR%\*darwin*.tar.gz 2>nul
echo.
echo Checksums:
dir /b %OUTPUT_DIR%\checksums.txt 2>nul

goto :eof

:build
set GOOS=%1
set GOARCH=%2
set OUTPUT_NAME=%PROJECT_NAME%

if "%GOOS%"=="windows" set OUTPUT_NAME=%OUTPUT_NAME%.exe

set OUTPUT_PATH=%OUTPUT_DIR%\%PROJECT_NAME%-%GOOS%-%GOARCH%
mkdir %OUTPUT_PATH%

echo Building for %GOOS%/%GOARCH%...

REM Build
set CGO_ENABLED=0
go build -tags="%BUILD_TAGS%" -ldflags="%LDFLAGS%" -o "%OUTPUT_PATH%\%OUTPUT_NAME%" .

REM Copy files
copy README.md %OUTPUT_PATH%\ >nul
copy LICENSE %OUTPUT_PATH%\ >nul

REM Create archive
cd %OUTPUT_DIR%
if "%GOOS%"=="windows" (
    powershell -Command "Compress-Archive -Path '%PROJECT_NAME%-%GOOS%-%GOARCH%' -DestinationPath '%PROJECT_NAME%-%VERSION%-%GOOS%-%GOARCH%.zip' -Force"
) else (
    tar -czf "%PROJECT_NAME%-%VERSION%-%GOOS%-%GOARCH%.tar.gz" "%PROJECT_NAME%-%GOOS%-%GOARCH%"
)
cd ..

echo   Built %GOOS%/%GOARCH%
goto :eof

