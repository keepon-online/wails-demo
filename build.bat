@echo off
setlocal enabledelayedexpansion

:: Wails Demo 构建脚本
:: 用法: build.bat [version]
:: 示例: build.bat v1.0.0

set VERSION=%1
if "%VERSION%"=="" set VERSION=dev

echo ========================================
echo   Wails Demo Build Script
echo   Version: %VERSION%
echo ========================================

:: 创建输出目录
if not exist "dist" mkdir dist

:: 步骤 1: 构建 Wails 应用
echo.
echo [1/3] Building Wails application...
wails build -platform windows/amd64 -ldflags "-X 'wails-demo/internal/updater.Version=%VERSION%'" -o wails-demo.exe

if %ERRORLEVEL% neq 0 (
    echo Wails build failed!
    exit /b 1
)
echo Wails build successful!

:: 步骤 2: 创建便携版 ZIP
echo.
echo [2/3] Creating portable ZIP...
if exist "dist\wails-demo_windows_amd64.zip" del "dist\wails-demo_windows_amd64.zip"
powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "Compress-Archive -Path 'build\bin\wails-demo.exe' -DestinationPath 'dist\wails-demo_windows_amd64.zip'"
echo Portable version created: dist\wails-demo_windows_amd64.zip

:: 步骤 3: 构建 NSIS 安装包
echo.
echo [3/3] Building NSIS installer...

:: 检查 NSIS 路径
set "NSIS_PATH=C:\Program Files (x86)\NSIS\makensis.exe"
if not exist "%NSIS_PATH%" (
    set "NSIS_PATH=C:\Program Files\NSIS\makensis.exe"
)
if not exist "%NSIS_PATH%" (
    echo NSIS not installed, skipping installer build
    echo Please install NSIS: choco install nsis
    goto :done
)

:: 编译 NSIS 脚本（使用命令行定义版本号覆盖脚本中的默认值）
set VERSION_NUM=%VERSION:v=%
echo Compiling NSIS installer...
"%NSIS_PATH%" /DPRODUCT_VERSION=%VERSION_NUM% build\nsis\installer.nsi

if %ERRORLEVEL% equ 0 (
    echo NSIS installer created successfully!
) else (
    echo NSIS build failed!
)

:done
echo.
echo ========================================
echo   Build Complete!
echo ========================================
echo.
echo Output files:
dir /b dist 2>nul

endlocal
