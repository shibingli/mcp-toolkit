@echo off
echo Verifying JSON Schema Pattern fix...
echo.

echo [1/2] Checking Pattern fix...
findstr /C:"Pattern:     \"^https?://.*$\"" pkg\types\schema.go >nul
if %ERRORLEVEL% EQU 0 (
    echo [OK] Pattern is correct: ^https?://.*$
) else (
    echo [ERROR] Pattern is incorrect
    exit /b 1
)

echo.
echo [2/2] Testing compilation...
go build -tags="sonic" -o test-build.exe
if %ERRORLEVEL% EQU 0 (
    echo [OK] Compilation successful
    del test-build.exe
) else (
    echo [ERROR] Compilation failed
    exit /b 1
)

echo.
echo ========================================
echo All verifications passed!
echo ========================================

