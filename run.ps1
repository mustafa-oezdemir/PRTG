param(
    [ValidateSet("linux", "windows")]
    [string]$Target = $(if ($PSVersionTable.PSEdition -eq "Desktop" -or $IsWindows) { "windows" } else { "linux" }),

    [switch]$SkipDocker,
    [switch]$SkipBuild
)

$ErrorActionPreference = "Stop"

function Invoke-Step {
    param(
        [string]$Name,
        [scriptblock]$Command
    )

    Write-Host ""
    Write-Host "==> $Name" -ForegroundColor Cyan
    & $Command
}

if (-not $SkipDocker) {
    Invoke-Step "Starting Grafana with Docker Compose" {
        docker compose up -d --build
    }
}

if (-not $SkipBuild) {
    Invoke-Step "Building backend for $Target" {
        mage -v "build:$Target"
    }
}

Invoke-Step "Starting frontend watcher" {
    npm run dev
}
