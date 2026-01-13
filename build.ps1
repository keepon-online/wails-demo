# Wails Demo 构建脚本
# 用法: .\build.ps1 [version]
# 示例: .\build.ps1 v1.0.0

param(
    [string]$Version = "dev"
)

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Wails Demo 构建脚本" -ForegroundColor Cyan
Write-Host "  版本: $Version" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# 创建输出目录
$distDir = "dist"
if (-not (Test-Path $distDir)) {
    New-Item -ItemType Directory -Path $distDir | Out-Null
}

# 步骤 1: 构建 Wails 应用
Write-Host "`n[1/3] 构建 Wails 应用..." -ForegroundColor Yellow
wails build -platform windows/amd64 -ldflags "-X 'wails-demo/internal/updater.Version=$Version'" -o wails-demo.exe

if ($LASTEXITCODE -ne 0) {
    Write-Host "Wails 构建失败!" -ForegroundColor Red
    exit 1
}
Write-Host "Wails 构建成功!" -ForegroundColor Green

# 步骤 2: 创建便携版 ZIP
Write-Host "`n[2/3] 创建便携版 ZIP..." -ForegroundColor Yellow
$zipPath = "$distDir\wails-demo_windows_amd64.zip"
if (Test-Path $zipPath) { Remove-Item $zipPath }
Compress-Archive -Path "build\bin\wails-demo.exe" -DestinationPath $zipPath
Write-Host "便携版创建成功: $zipPath" -ForegroundColor Green

# 步骤 3: 构建 NSIS 安装包
Write-Host "`n[3/3] 构建 NSIS 安装包..." -ForegroundColor Yellow

# 检查 NSIS 是否安装
$nsisPath = "C:\Program Files (x86)\NSIS\makensis.exe"
if (-not (Test-Path $nsisPath)) {
    Write-Host "NSIS 未安装，跳过安装包构建" -ForegroundColor Yellow
    Write-Host "请安装 NSIS: choco install nsis" -ForegroundColor Yellow
} else {
    # 更新 NSIS 脚本中的版本号
    $nsiFile = "build\nsis\installer.nsi"
    $nsiContent = Get-Content $nsiFile -Raw
    $versionNum = $Version -replace '^v', ''
    $nsiContent = $nsiContent -replace '!define PRODUCT_VERSION "[^"]*"', "!define PRODUCT_VERSION `"$versionNum`""
    Set-Content $nsiFile $nsiContent
    
    # 编译 NSIS 脚本
    & $nsisPath $nsiFile
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "NSIS 安装包创建成功!" -ForegroundColor Green
    } else {
        Write-Host "NSIS 构建失败!" -ForegroundColor Red
    }
}

# 显示输出文件
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  构建完成!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "`n输出文件:" -ForegroundColor White
Get-ChildItem $distDir | ForEach-Object {
    Write-Host "  - $($_.Name)" -ForegroundColor Gray
}
