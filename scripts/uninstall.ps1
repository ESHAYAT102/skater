param(
    [switch]$WithSkate
)

$ErrorActionPreference = "Stop"

$BinaryName = "skater.exe"
$SkateBinaryName = "skate.exe"
$InstallDir = Join-Path $HOME ".local\bin"

function Remove-Binary {
    param(
        [string]$Path,
        [string]$Name
    )

    if (Test-Path $Path) {
        Remove-Item -Force $Path
        Write-Host "removed $Name from $Path"
    }
    else {
        Write-Host "$Name was not found at $Path"
    }
}

Remove-Binary -Path (Join-Path $InstallDir $BinaryName) -Name $BinaryName

if ($WithSkate) {
    Remove-Binary -Path (Join-Path $InstallDir $SkateBinaryName) -Name $SkateBinaryName
}
else {
    Write-Host "left $SkateBinaryName installed; rerun with -WithSkate to remove $(Join-Path $InstallDir $SkateBinaryName)"
}
