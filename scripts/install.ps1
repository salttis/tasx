param(
    [string] $Version = 'dev',
    [switch] $Rollback
)

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
$dist = Join-Path $repoRoot 'dist'
$targetDir = Join-Path $HOME '.personal\scripts'
$ldflagsVariable = 'github.com/salttis/tasx/internal/cli.version'
$commands = @('tasx', 'tx')

if ($Rollback) {
    $rolledBack = 0
    foreach ($command in $commands) {
        $target = Join-Path $targetDir "$command.exe"
        $backups = @(Get-ChildItem -LiteralPath $targetDir -Filter "$command.exe.previous*" -File -ErrorAction SilentlyContinue |
            Sort-Object LastWriteTimeUtc -Descending)
        if ($backups.Count -eq 0) {
            Write-Output "No previous $command.exe backup exists; leaving it unchanged."
            continue
        }

        $backup = $backups[0].FullName
        if ([IO.Path]::GetFileName($backup) -match '\.previous-absent(?:-|$)') {
            if (Test-Path -LiteralPath $target) {
                Remove-Item -LiteralPath $target -Force
            }
            Remove-Item -LiteralPath $backup -Force
            Write-Output "Rolled back $command command to its previously absent state."
            $rolledBack++
            continue
        }

        $staged = Join-Path $targetDir ".$command-$PID.rollback"
        $hadTarget = Test-Path -LiteralPath $target
        try {
            if ($hadTarget) {
                Move-Item -LiteralPath $target -Destination $staged
            }
            Move-Item -LiteralPath $backup -Destination $target
            if ($hadTarget) {
                Move-Item -LiteralPath $staged -Destination $backup
                [IO.File]::SetLastWriteTimeUtc($backup, [DateTime]::UtcNow)
            }
        } catch {
            if (Test-Path -LiteralPath $staged) {
                if (-not (Test-Path -LiteralPath $target)) {
                    Move-Item -LiteralPath $staged -Destination $target
                } else {
                    $recovery = Join-Path $targetDir ("$command.exe.rollback-recovery-{0}" -f (Get-Date -Format 'yyyyMMddHHmmssfff'))
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
        Write-Output "Rolled back $command command: $target"
        $rolledBack++
    }

    if ($rolledBack -eq 0) {
        throw "No previous Tasx command binaries exist in $targetDir."
    }
    return
}

if ($Version -notmatch '^(dev|v?[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?)$') {
    throw "Version must be 'dev' or a semantic version (for example 1.0.0)."
}
if ($Version.StartsWith('v')) {
    $Version = $Version.Substring(1)
}

$distBinary = Join-Path $dist 'tasx.exe'
$installations = @()
foreach ($command in $commands) {
    $target = Join-Path $targetDir "$command.exe"
    $installations += [pscustomobject]@{
        Name          = $command
        Target        = $target
        Staged        = Join-Path $targetDir ".$command-$PID.new.exe"
        HadTarget     = Test-Path -LiteralPath $target
        Backup        = $null
        BackupCreated = $false
        Installed     = $false
    }
}

New-Item -ItemType Directory -Path $dist -Force | Out-Null
Push-Location $repoRoot
try {
    & go build -trimpath -ldflags "-X $ldflagsVariable=$Version" -o $distBinary .\cmd\tasx
    if ($LASTEXITCODE -ne 0) { throw 'Tasx build failed.' }
} finally {
    Pop-Location
}

Copy-Item -LiteralPath $distBinary -Destination (Join-Path $dist 'tx.exe')
New-Item -ItemType Directory -Path $targetDir -Force | Out-Null
try {
    foreach ($installation in $installations) {
        Copy-Item -LiteralPath $distBinary -Destination $installation.Staged
        if ($installation.HadTarget) {
            $installation.Backup = "$($installation.Target).previous"
            if (Test-Path -LiteralPath $installation.Backup) {
                $stamp = "{0}-{1}" -f (Get-Date -Format 'yyyyMMddHHmmssfff'), [guid]::NewGuid().ToString('N')
                $installation.Backup = Join-Path $targetDir "$($installation.Name).exe.previous-$stamp"
            }
        } else {
            $stamp = "{0}-{1}" -f (Get-Date -Format 'yyyyMMddHHmmssfff'), [guid]::NewGuid().ToString('N')
            $installation.Backup = Join-Path $targetDir "$($installation.Name).exe.previous-absent-$stamp"
        }
    }
} catch {
    foreach ($installation in $installations) {
        if (Test-Path -LiteralPath $installation.Staged) {
            Remove-Item -LiteralPath $installation.Staged -Force
        }
    }
    throw
}

try {
    foreach ($installation in $installations) {
        if ($installation.HadTarget) {
            Move-Item -LiteralPath $installation.Target -Destination $installation.Backup
            $installation.BackupCreated = $true
            [IO.File]::SetLastWriteTimeUtc($installation.Backup, [DateTime]::UtcNow)
        } else {
            New-Item -ItemType File -Path $installation.Backup | Out-Null
            $installation.BackupCreated = $true
        }
        Move-Item -LiteralPath $installation.Staged -Destination $installation.Target
        $installation.Installed = $true
    }
} catch {
    $installError = $_
    $recoveryErrors = @()
    $rollbackInstallations = @($installations | Where-Object { $_.BackupCreated })
    [array]::Reverse($rollbackInstallations)
    foreach ($installation in $rollbackInstallations) {
        try {
            if ($installation.Installed -and (Test-Path -LiteralPath $installation.Target)) {
                Remove-Item -LiteralPath $installation.Target -Force
            }
            if ($installation.HadTarget -and (Test-Path -LiteralPath $installation.Backup)) {
                Move-Item -LiteralPath $installation.Backup -Destination $installation.Target
            } elseif (-not $installation.HadTarget -and (Test-Path -LiteralPath $installation.Backup)) {
                Remove-Item -LiteralPath $installation.Backup -Force
            }
        } catch {
            $recoveryErrors += "$($installation.Name): $($_.Exception.Message)"
        }
    }
    if ($recoveryErrors.Count -gt 0) {
        throw "Installation failed: $($installError.Exception.Message) Rollback was incomplete: $($recoveryErrors -join '; ')"
    }
    throw $installError
} finally {
    foreach ($installation in $installations) {
        if (Test-Path -LiteralPath $installation.Staged) {
            Remove-Item -LiteralPath $installation.Staged -Force
        }
    }
}

foreach ($installation in $installations) {
    Write-Output "Installed $($installation.Name) ${Version}: $($installation.Target)"
    if ($installation.HadTarget) {
        Write-Output "Previous binary saved to: $($installation.Backup)"
    }
}
Write-Output "Rollback with: pwsh -NoProfile -File .\scripts\install.ps1 -Rollback"
