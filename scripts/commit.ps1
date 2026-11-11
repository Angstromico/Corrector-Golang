<#
.SYNOPSIS
    Interactive commit script enforcing tests, linting, and Conventional Commits.
.DESCRIPTION
    Runs `go vet` and `go test`. If checks pass, prompts for commit Type, Scope, and Message,
    then executes git commit.
.PARAMETER Help
    Show this help message.
.PARAMETER SkipTests
    Skip running tests before commit (use with caution).
#>
param (
    [switch]$Help,
    [switch]$SkipTests
)

if ($Help) {
    Get-Help $MyInvocation.MyCommand.Path
    exit 0
}

$ErrorActionPreference = "Stop"

Write-Host "==> Verifying repository state..." -ForegroundColor Cyan

if (-not $SkipTests) {
    Write-Host "--> Running go vet..." -ForegroundColor Yellow
    go vet ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Error "go vet failed! Aborting commit."
        exit $LASTEXITCODE
    }

    Write-Host "--> Running tests..." -ForegroundColor Yellow
    go test -v ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Error "go test failed! Aborting commit."
        exit $LASTEXITCODE
    }
}

$validTypes = @("feat", "fix", "refactor", "test", "docs", "chore", "ci", "perf", "style")
$validScopes = @("core", "parser", "spellcheck", "grammar", "statistics", "gui", "repo", "build", "e2e")

Write-Host "`n=== Conventional Commit Assistant ===" -ForegroundColor Green
Write-Host "Allowed types: $($validTypes -join ', ')"
Write-Host "Common scopes: $($validScopes -join ', ')"

$type = ""
while ($validTypes -notcontains $type) {
    $type = (Read-Host -Prompt "Enter commit type").Trim().ToLower()
    if ($validTypes -notcontains $type) {
        Write-Host "Invalid type. Please choose from: $($validTypes -join ', ')" -ForegroundColor Red
    }
}

$scope = (Read-Host -Prompt "Enter scope (e.g., core, parser, repo) [optional]").Trim().ToLower()
$message = ""
while ([string]::IsNullOrWhiteSpace($message)) {
    $message = (Read-Host -Prompt "Enter commit message (imperative, lower-case start)").Trim()
    if ([string]::IsNullOrWhiteSpace($message)) {
        Write-Host "Message cannot be empty." -ForegroundColor Red
    }
}

$commitMsg = if ($scope) { "$($type)($($scope)): $message" } else { "$($type): $message" }

Write-Host "`nCommit message: '$commitMsg'" -ForegroundColor Cyan
$confirm = (Read-Host -Prompt "Proceed with commit? (y/n)").Trim().ToLower()
if ($confirm -eq "y" -or $confirm -eq "yes") {
    git commit -m "$commitMsg"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Successfully committed: $commitMsg" -ForegroundColor Green
    }
} else {
    Write-Host "Commit cancelled." -ForegroundColor Yellow
}
