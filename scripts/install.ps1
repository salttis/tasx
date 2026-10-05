$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
$dist = Join-Path $repoRoot 'dist'
$targetDir = Join-Path $HOME '.personal\scripts'
$target = Join-Path $targetDir 'tx.exe'

New-Item -ItemType Directory -Path $dist -Force | Out-Null
New-Item -ItemType Directory -Path $targetDir -Force | Out-Null
Push-Location $repoRoot
try {
    & go build -o (Join-Path $dist 'tx.exe') .\cmd\tasx
    if ($LASTEXITCODE -ne 0) { throw 'Tasx build failed.' }
} finally {
    Pop-Location
}
Copy-Item -LiteralPath (Join-Path $dist 'tx.exe') -Destination $target -Force
Write-Output "Installed tx command: $target"
