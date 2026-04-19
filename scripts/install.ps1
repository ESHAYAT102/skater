$ErrorActionPreference = "Stop"

$BinaryName = "skater.exe"
$SkateBinaryName = "skate.exe"
$InstallDir = Join-Path $HOME ".local\bin"

function Require-Command {
    param([string]$Name)

    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "$Name is required"
    }
}

Require-Command go

function Install-SkateIfMissing {
    $InstalledSkate = Join-Path $InstallDir $SkateBinaryName
    if ((Get-Command skate -ErrorAction SilentlyContinue) -or (Test-Path $InstalledSkate)) {
        return
    }

    Write-Host "skate not found; installing skate"
    $OldGoBin = $env:GOBIN
    try {
        $env:GOBIN = $InstallDir
        go install github.com/charmbracelet/skate@latest
    }
    finally {
        $env:GOBIN = $OldGoBin
    }
}

try {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

    Install-SkateIfMissing

    $OldGoBin = $env:GOBIN
    try {
        $env:GOBIN = $InstallDir
        go install github.com/ESHAYAT102/skater@latest
    }
    finally {
        $env:GOBIN = $OldGoBin
    }

    Write-Host "installed $BinaryName to $(Join-Path $InstallDir $BinaryName)"

    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $PathEntries = @()
    if ($UserPath) {
        $PathEntries = $UserPath -split ";"
    }

    if ($PathEntries -notcontains $InstallDir) {
        $NewPath = if ($UserPath) { "$UserPath;$InstallDir" } else { $InstallDir }
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        $env:Path = "$env:Path;$InstallDir"
        Write-Host "added $InstallDir to your user PATH"
        Write-Host "restart your terminal before running skater"
    }
}
