@echo off
REM 验证JSON Schema Pattern修复的脚本 / Script to verify JSON Schema Pattern fix

echo ========================================
echo 验证 JSON Schema Pattern 修复
echo Verify JSON Schema Pattern Fix
echo ========================================
echo.

echo [1/2] 检查修复内容...
echo [1/2] Checking fix content...
findstr /C:"Pattern:     \"^https?://.*$\"" pkg\types\schema.go >nul
if %ERRORLEVEL% EQU 0 (
    echo ✅ Pattern 修复正确
    echo ✅ Pattern fix is correct
) else (
    echo ❌ Pattern 修复不正确
    echo ❌ Pattern fix is incorrect
    exit /b 1
)

echo.
echo [2/2] 编译测试...
echo [2/2] Compile test...
go build -tags="sonic" -o mcp-toolkit-test.exe
if %ERRORLEVEL% EQU 0 (
    echo ✅ 编译成功
    echo ✅ Compile successful
    del mcp-toolkit-test.exe
) else (
    echo ❌ 编译失败
    echo ❌ Compile failed
    exit /b 1
)

echo.
echo ========================================
echo ✅ 所有验证通过！
echo ✅ All verifications passed!
echo ========================================
echo.
echo 修复内容：
echo Fix content:
echo   文件: pkg/types/schema.go
echo   File: pkg/types/schema.go
echo   行号: 661
echo   Line: 661
echo   修复前: Pattern: "^https?://.*",
echo   Before: Pattern: "^https?://.*",
echo   修复后: Pattern: "^https?://.*$",
echo   After:  Pattern: "^https?://.*$",
echo.

