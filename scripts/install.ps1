param(
    [string] $Version = 'dev',
    [switch] $Rollback
)

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
$dist = Join-Path $repoRoot 'dist'
$targetDir = Join-Path $HOME '.personal\scripts'
$target = Join-Path $targetDir 'tx.exe'
$ldflagsVariable = 'github.com/salttis/tasx/internal/cli.version'

if ($Rollback) {
    $backups = @(Get-ChildItem -LiteralPath $targetDir -Filter 'tx.exe.previous*' -File -ErrorAction SilentlyContinue |
        Sort-Object LastWriteTimeUtc -Descending)
    if ($backups.Count -eq 0) {
        throw "No previous tx.exe backup exists in $targetDir."
    }

    $backup = $backups[0].FullName
    $staged = Join-Path $targetDir ".tx-$PID.rollback"
    $hadTarget = Test-Path -LiteralPath $target
    try {
        if ($hadTarget) {
            Move-Item -LiteralPath $target -Destination $staged
        }
        Move-Item -LiteralPath $backup -Destination $target
        if ($hadTarget) {
            Move-Item -LiteralPath $staged -Destination $backup
        }
    } catch {
        if (Test-Path -LiteralPath $staged) {
            if (-not (Test-Path -LiteralPath $target)) {
                Move-Item -LiteralPath $staged -Destination $target
            } else {
                $recovery = Join-Path $targetDir ("tx.exe.rollback-recovery-{0}" -f (Get-Date -Format 'yyyyMMddHHmmssfff'))
                Move-Item -LiteralPath $staged -Destination $recovery
                throw "Rollback was incomplete; the previous current binary was preserved at $recovery. $($_.Exception.Message)"
            }
        }
        throw
    } finally {
        if (Test-Path -LiteralPath $staged) {
            Remove-Item -LiteralPath $staged -Force
        }
    }

    Write-Output "Rolled back tx command: $target"
    return
}

if ($Version -notmatch '^(dev|v?[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?)$') {
    throw "Version must be 'dev' or a semantic version (for example 1.0.0)."
}
if ($Version.StartsWith('v')) {
    $Version = $Version.Substring(1)
}

$distBinary = Join-Path $dist 'tx.exe'
$staged = Join-Path $targetDir ".tx-$PID.new.exe"
$backup = Join-Path $targetDir 'tx.exe.previous'
if (Test-Path -LiteralPath $backup) {
    $backup = Join-Path $targetDir ("tx.exe.previous-{0}" -f (Get-Date -Format 'yyyyMMddHHmmssfff'))
    if (Test-Path -LiteralPath $backup) {
        throw "Backup path already exists: $backup"
    }
}

New-Item -ItemType Directory -Path $dist -Force | Out-Null
New-Item -ItemType Directory -Path $targetDir -Force | Out-Null
Push-Location $repoRoot
try {
    & go build -trimpath -ldflags "-X $ldflagsVariable=$Version" -o $distBinary .\cmd\tasx
    if ($LASTEXITCODE -ne 0) { throw 'Tasx build failed.' }
} finally {
    Pop-Location
}

Copy-Item -LiteralPath $distBinary -Destination $staged
$hadTarget = Test-Path -LiteralPath $target
try {
    if ($hadTarget) {
        Move-Item -LiteralPath $target -Destination $backup
    }
    Move-Item -LiteralPath $staged -Destination $target
} catch {
    if ($hadTarget -and (Test-Path -LiteralPath $backup) -and -not (Test-Path -LiteralPath $target)) {
        Move-Item -LiteralPath $backup -Destination $target
    }
    throw
} finally {
    if (Test-Path -LiteralPath $staged) {
        Remove-Item -LiteralPath $staged -Force
    }
}

Write-Output "Installed tx ${Version}: $target"
if ($hadTarget) {
    Write-Output "Previous binary saved to: $backup"
    Write-Output "Rollback with: pwsh -NoProfile -File .\scripts\install.ps1 -Rollback"
}
