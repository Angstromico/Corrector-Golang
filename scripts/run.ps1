<#
.SYNOPSIS
    Run GramatiGo in development mode.
#>
param (
    [switch]$Help
)

if ($Help) {
    Get-Help $MyInvocation.MyCommand.Path
    exit 0
}

$ErrorActionPreference = "Stop"

if (Get-Command "wails" -ErrorAction SilentlyContinue) {
    Write-Host "==> Starting Wails development server..." -ForegroundColor Cyan
    wails dev
} else {
    Write-Host "==> Wails CLI not found on PATH. Running standard Go app..." -ForegroundColor Yellow
    go run main.go
}
