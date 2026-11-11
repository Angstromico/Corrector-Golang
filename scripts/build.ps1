<#
.SYNOPSIS
    Build GramatiGo optimized binary with -ldflags="-s -w".
#>
param (
    [string]$Platform = "windows/amd64",
    [switch]$NoGui,
    [switch]$Help
)

if ($Help) {
    Get-Help $MyInvocation.MyCommand.Path
    exit 0
}

$ErrorActionPreference = "Stop"

$outputDir = "build/bin"
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
}

$ldflags = "-s -w"
if (-not $NoGui -and $Platform -like "windows*") {
    $ldflags = "-H windowsgui -s -w"
}

Write-Host "==> Building GramatiGo for $Platform..." -ForegroundColor Cyan

if (Get-Command "wails" -ErrorAction SilentlyContinue) {
    wails build -clean -platform $Platform -ldflags "$ldflags"
} else {
    Write-Host "==> Wails CLI not found. Running go build..." -ForegroundColor Yellow
    $ext = if ($Platform -like "windows*") { ".exe" } else { "" }
    go build -ldflags="$ldflags" -o "$outputDir/corrector$ext" .
}

Write-Host "==> Build finished successfully." -ForegroundColor Green
