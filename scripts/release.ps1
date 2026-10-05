param(
    [Parameter(Mandatory, Position = 0)]
    [ValidateSet('minor', 'major')]
    [string] $Bump
)

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
$repository = 'salttis/tasx'
$releaseWorkflow = 'release.yml'
$runnerQueueTimeout = [TimeSpan]::FromMinutes(5)
$workflowTimeout = [TimeSpan]::FromHours(1)
$commands = @('git', 'go', 'gh', 'tar')

foreach ($name in $commands) {
    if (-not (Get-Command $name -ErrorAction SilentlyContinue)) {
        throw "Required command '$name' was not found on PATH."
    }
}

function Invoke-Checked {
    param(
        [Parameter(Mandatory)][string] $Command,
        [Parameter(Mandatory)][string[]] $Arguments
    )

    & $Command @Arguments | Out-Host
    if ($LASTEXITCODE -ne 0) {
        throw "'$Command $($Arguments -join ' ')' failed with exit code $LASTEXITCODE."
    }
}

function Get-Release {
    param([Parameter(Mandatory)][string] $Tag)

    $output = & gh release view $Tag --repo $repository --json name,tagName,isDraft,isPrerelease,assets,url 2>$null
    if ($LASTEXITCODE -ne 0) {
        return $null
    }
    return $output | ConvertFrom-Json
}

function Get-ReleaseRun {
    param(
        [Parameter(Mandatory)][string] $Tag,
        [Parameter(Mandatory)][string] $Commit
    )

    $runs = & gh run list --repo $repository --workflow $releaseWorkflow --commit $Commit --limit 20 --json databaseId,headSha,headBranch,status,conclusion,event
    if ($LASTEXITCODE -ne 0) {
        throw 'Could not query the GitHub release workflow runs.'
    }
    $matchingRuns = @($runs | ConvertFrom-Json | Where-Object {
        $_.headSha -eq $Commit -and $_.headBranch -eq $Tag -and $_.event -eq 'push'
    })
    if ($matchingRuns.Count -gt 1) {
        throw "More than one release workflow run was found for $Tag."
    }
    if ($matchingRuns.Count -eq 0) {
        return $null
    }
    return $matchingRuns[0]
}

function Get-StableTags {
    $tagNames = @(git tag --list 'v[0-9]*.[0-9]*.[0-9]*')
    if ($LASTEXITCODE -ne 0) {
        throw 'Could not list local Git tags.'
    }

    $stableTags = foreach ($tagName in $tagNames) {
        if ($tagName -match '^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$') {
            [pscustomobject]@{
                Name    = $tagName
                Version = [version]::new([int]$Matches[1], [int]$Matches[2], [int]$Matches[3])
            }
        }
    }
    return @($stableTags | Sort-Object Version -Descending)
}

function Test-Archive {
    param(
        [Parameter(Mandatory)][string] $Archive,
        [Parameter(Mandatory)][string] $Package,
        [Parameter(Mandatory)][AllowEmptyString()][string] $Extension
    )

    $entries = @(tar -tzf $Archive)
    if ($LASTEXITCODE -ne 0) {
        throw "Could not read archive $([IO.Path]::GetFileName($Archive))."
    }
    foreach ($requiredPath in @(
        "$Package/tasx$Extension",
        "$Package/tx$Extension",
        "$Package/README.md",
        "$Package/CHANGELOG.md",
        "$Package/LICENSE",
        "$Package/docs/README.md"
    )) {
        if ($entries -notcontains $requiredPath) {
            throw "Archive $([IO.Path]::GetFileName($Archive)) is missing $requiredPath."
        }
    }
    if (@($entries | Where-Object { $_ -like "$Package/third_party_licenses/*" }).Count -eq 0) {
        throw "Archive $([IO.Path]::GetFileName($Archive)) contains no third-party license files."
    }
}

function New-ReleasePackages {
    param(
        [Parameter(Mandatory)][string] $Version,
        [Parameter(Mandatory)][string] $OutputDirectory
    )

    $toolDirectory = Join-Path $OutputDirectory 'tools'
    $distDirectory = Join-Path $OutputDirectory 'dist'
    New-Item -ItemType Directory -Force -Path $toolDirectory, $distDirectory | Out-Null
    $oldGoBin = $env:GOBIN
    $env:GOBIN = $toolDirectory
    try {
        Invoke-Checked go @('install', 'github.com/google/go-licenses/v2@v2.0.1')
    } finally {
        $env:GOBIN = $oldGoBin
        if ($null -eq $oldGoBin) {
            Remove-Item Env:GOBIN -ErrorAction SilentlyContinue
        }
    }

    $targets = @(
        @{ OS = 'windows'; Arch = 'amd64'; Extension = '.exe' },
        @{ OS = 'windows'; Arch = 'arm64'; Extension = '.exe' },
        @{ OS = 'linux'; Arch = 'amd64'; Extension = '' },
        @{ OS = 'linux'; Arch = 'arm64'; Extension = '' },
        @{ OS = 'darwin'; Arch = 'amd64'; Extension = '' },
        @{ OS = 'darwin'; Arch = 'arm64'; Extension = '' }
    )
    $oldGoOS = $env:GOOS
    $oldGoArch = $env:GOARCH
    $archives = [System.Collections.Generic.List[string]]::new()
    try {
        foreach ($target in $targets) {
            $package = "tasx_${Version}_$($target.OS)_$($target.Arch)"
            $packageDirectory = Join-Path $distDirectory $package
            New-Item -ItemType Directory -Path $packageDirectory | Out-Null
            $env:GOOS = $target.OS
            $env:GOARCH = $target.Arch

            Write-Host "Building $($target.OS)/$($target.Arch)..."
            $binary = Join-Path $packageDirectory "tasx$($target.Extension)"
            Invoke-Checked go @(
                'build',
                '-trimpath',
                '-ldflags',
                "-s -w -X github.com/salttis/tasx/internal/cli.version=$Version",
                '-o',
                $binary,
                '.\cmd\tasx'
            )
            Copy-Item -LiteralPath $binary -Destination (Join-Path $packageDirectory "tx$($target.Extension)")

            Write-Host "Collecting licenses for $($target.OS)/$($target.Arch)..."
            $licenseDirectory = Join-Path $packageDirectory 'third_party_licenses'
            Invoke-Checked (Join-Path $toolDirectory 'go-licenses.exe') @(
                'save',
                '.\cmd\tasx',
                "--save_path=$licenseDirectory"
            )
            Copy-Item README.md, CHANGELOG.md, LICENSE -Destination $packageDirectory
            Copy-Item docs -Destination (Join-Path $packageDirectory 'docs') -Recurse

            $archive = Join-Path $distDirectory "$package.tar.gz"
            Invoke-Checked tar @('-czf', $archive, '-C', $distDirectory, $package)
            $archives.Add($archive)
        }
    } finally {
        $env:GOOS = $oldGoOS
        $env:GOARCH = $oldGoArch
        if ($null -eq $oldGoOS) { Remove-Item Env:GOOS -ErrorAction SilentlyContinue }
        if ($null -eq $oldGoArch) { Remove-Item Env:GOARCH -ErrorAction SilentlyContinue }
    }

    $sumLines = foreach ($archive in $archives) {
        $hash = (Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.ToLowerInvariant()
        "$hash  $([IO.Path]::GetFileName($archive))"
    }
    $sumFile = Join-Path $distDirectory 'SHA256SUMS'
    [IO.File]::WriteAllLines($sumFile, $sumLines, [Text.UTF8Encoding]::new($false))
    return @($archives) + @($sumFile)
}

function Complete-Release {
    param(
        [Parameter(Mandatory)][string] $Tag,
        [Parameter(Mandatory)][string] $Version,
        [Parameter(Mandatory)][string] $Commit,
        [Parameter(Mandatory)][string] $OutputDirectory
    )

    $deadline = [DateTime]::UtcNow + $workflowTimeout
    $queueDeadline = [DateTime]::UtcNow + $runnerQueueTimeout
    $run = $null
    $useLocalPackages = $false
    while (-not $run -and [DateTime]::UtcNow -lt $queueDeadline) {
        $run = Get-ReleaseRun -Tag $Tag -Commit $Commit
        if (-not $run) {
            Start-Sleep -Seconds 10
        }
    }
    if (-not $run) {
        Write-Warning "No GitHub release workflow appeared within $runnerQueueTimeout; preparing packages locally."
        $useLocalPackages = $true
    } else {
        while ($run.status -notin @('completed', 'cancelled') -and [DateTime]::UtcNow -lt $deadline) {
            if ($run.status -eq 'queued' -and [DateTime]::UtcNow -ge $queueDeadline) {
                Write-Warning "GitHub has not assigned a runner after $runnerQueueTimeout; cancelling the queued workflow and preparing packages locally."
                Invoke-Checked gh @('run', 'cancel', "$($run.databaseId)", '--repo', $repository)
                do {
                    Start-Sleep -Seconds 5
                    $run = Get-ReleaseRun -Tag $Tag -Commit $Commit
                    if (-not $run) {
                        throw "The GitHub release workflow run for $Tag disappeared while cancelling."
                    }
                } while ($run.status -notin @('completed', 'cancelled') -and [DateTime]::UtcNow -lt $deadline)
                if ($run.conclusion -eq 'cancelled') {
                    $useLocalPackages = $true
                } elseif ($run.conclusion -ne 'success') {
                    throw "GitHub release workflow $($run.databaseId) ended with status '$($run.status)' and conclusion '$($run.conclusion)'."
                }
                break
            }
            Start-Sleep -Seconds 10
            $run = Get-ReleaseRun -Tag $Tag -Commit $Commit
            if (-not $run) {
                throw "The GitHub release workflow run for $Tag disappeared."
            }
        }
        if ([DateTime]::UtcNow -ge $deadline) {
            throw "The GitHub release workflow did not finish within $workflowTimeout."
        }
        if ($run.conclusion -eq 'cancelled') {
            Write-Warning "The GitHub release workflow was cancelled; preparing packages locally."
            $useLocalPackages = $true
        } elseif (-not $useLocalPackages -and ($run.status -ne 'completed' -or $run.conclusion -ne 'success')) {
            throw "GitHub release workflow $($run.databaseId) ended with status '$($run.status)' and conclusion '$($run.conclusion)'."
        }
    }

    $release = Get-Release -Tag $Tag
    if ($useLocalPackages) {
        if ($release -and -not $release.isDraft) {
            throw "Release $Tag is already public; refusing to replace it."
        }
        $files = New-ReleasePackages -Version $Version -OutputDirectory $OutputDirectory
        if ($release) {
            Invoke-Checked gh (@('release', 'upload', $Tag, '--repo', $repository, '--clobber') + $files)
        } else {
            Invoke-Checked gh (@(
                'release', 'create', $Tag,
                '--repo', $repository,
                '--title', "Tasx v$Version",
                '--draft',
                '--generate-notes'
            ) + $files)
        }
        $release = Get-Release -Tag $Tag
    }
    if (-not $release -or -not $release.isDraft -or $release.tagName -ne $Tag) {
        throw "Expected a GitHub draft release for $Tag."
    }

    $expectedNames = @(
        'SHA256SUMS',
        "tasx_${Version}_darwin_amd64.tar.gz",
        "tasx_${Version}_darwin_arm64.tar.gz",
        "tasx_${Version}_linux_amd64.tar.gz",
        "tasx_${Version}_linux_arm64.tar.gz",
        "tasx_${Version}_windows_amd64.tar.gz",
        "tasx_${Version}_windows_arm64.tar.gz"
    )
    $actualNames = @($release.assets | ForEach-Object name | Sort-Object)
    if (Compare-Object ($expectedNames | Sort-Object) $actualNames) {
        throw "Release $Tag does not contain exactly the six platform archives and SHA256SUMS."
    }

    $downloadDirectory = Join-Path $OutputDirectory 'download'
    New-Item -ItemType Directory -Path $downloadDirectory | Out-Null
    Invoke-Checked gh @('release', 'download', $Tag, '--repo', $repository, '--dir', $downloadDirectory, '--pattern', '*')
    $sumFile = Join-Path $downloadDirectory 'SHA256SUMS'
    if (-not (Test-Path -LiteralPath $sumFile)) {
        throw "Release $Tag is missing SHA256SUMS."
    }
    $checksums = @{}
    foreach ($line in Get-Content -LiteralPath $sumFile) {
        if ($line -notmatch '^([0-9a-f]{64})  (tasx_[0-9]+\.[0-9]+\.[0-9]+_[a-z]+_[a-z0-9]+\.tar\.gz)$') {
            throw "Invalid checksum manifest line: $line"
        }
        if ($checksums.ContainsKey($Matches[2])) {
            throw "Duplicate checksum entry for $($Matches[2])."
        }
        $checksums[$Matches[2]] = $Matches[1]
    }
    if ($checksums.Count -ne 6) {
        throw "Expected six checksum entries in $sumFile."
    }

    foreach ($name in $expectedNames | Where-Object { $_ -ne 'SHA256SUMS' }) {
        $archive = Join-Path $downloadDirectory $name
        if (-not (Test-Path -LiteralPath $archive)) {
            throw "Downloaded release asset is missing: $name."
        }
        $hash = (Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.ToLowerInvariant()
        if ($checksums[$name] -ne $hash) {
            throw "SHA-256 verification failed for $name."
        }
        $extension = if ($name -match '_windows_') { '.exe' } else { '' }
        $package = [IO.Path]::GetFileNameWithoutExtension([IO.Path]::GetFileNameWithoutExtension($name))
        Test-Archive -Archive $archive -Package $package -Extension $extension
    }
    if (@(Get-ChildItem -LiteralPath $downloadDirectory -File).Count -ne 7) {
        throw "Downloaded release $Tag has unexpected or missing assets."
    }

    Write-Host "All release archives, licenses, and checksums verified. Publishing $Tag..."
    Invoke-Checked gh @('release', 'edit', $Tag, '--repo', $repository, '--draft=false')
    $published = Get-Release -Tag $Tag
    if (-not $published -or $published.isDraft -or $published.isPrerelease -or $published.tagName -ne $Tag) {
        throw "Could not verify that $Tag is a public stable release."
    }
    return $published
}

Push-Location $repoRoot
$temporaryDirectory = $null
$temporaryDirectoryCreated = $false
try {
    Invoke-Checked gh @('auth', 'status')
    Invoke-Checked git @('fetch', 'origin', '--tags', '--prune')

    $branch = (git branch --show-current).Trim()
    if ($LASTEXITCODE -ne 0 -or $branch -ne 'main') {
        throw "Release from the main branch only; current branch is '$branch'."
    }
    $dirty = @(git status --porcelain)
    if ($LASTEXITCODE -ne 0 -or $dirty.Count -gt 0) {
        throw 'The worktree must be clean before a release. Commit or stash all changes first.'
    }
    $commit = (git rev-parse HEAD).Trim()
    $remoteCommit = (git rev-parse origin/main).Trim()
    if ($LASTEXITCODE -ne 0 -or $commit -ne $remoteCommit) {
        throw 'Local main must match origin/main. Pull or push commits before releasing.'
    }

    $remoteRepository = (gh repo view --json nameWithOwner --jq '.nameWithOwner').Trim()
    if ($LASTEXITCODE -ne 0 -or $remoteRepository -ne $repository) {
        throw "Expected GitHub repository $repository, found '$remoteRepository'."
    }
    if ((gh api repos/$repository --jq '.private').Trim() -ne 'false') {
        throw "Repository $repository is not public."
    }

    $tags = Get-StableTags
    if ($tags.Count -eq 0) {
        throw 'No stable vMAJOR.MINOR.PATCH Git tag exists.'
    }
    $latest = $tags[0]
    $latestCommit = (git rev-parse "$($latest.Name)^{}").Trim()
    if ($LASTEXITCODE -ne 0) {
        throw "Could not resolve tag $($latest.Name)."
    }
    $latestRelease = Get-Release -Tag $latest.Name
    $pendingRelease = $latestRelease
    if ($pendingRelease -and -not $pendingRelease.isDraft) {
        $pendingRelease = $null
    }
    $pendingRun = Get-ReleaseRun -Tag $latest.Name -Commit $latestCommit
    $resumePending = [bool]$pendingRelease -or
        ($pendingRun -and $pendingRun.status -in @('queued', 'in_progress')) -or
        (-not $latestRelease -and $commit -eq $latestCommit)

    if ($resumePending) {
        if ($commit -ne $latestCommit) {
            throw "Pending release $($latest.Name) points to a different commit. Check its GitHub workflow before continuing."
        }
        $tag = $latest.Name
        $version = $latest.Version.ToString()
        Write-Output "Resuming the unfinished release $tag."
    } else {
        if ($tags.Count -gt 1) {
            $previous = $tags[1]
            $previousCommit = (git rev-parse "$($previous.Name)^{}").Trim()
            if ($LASTEXITCODE -ne 0) {
                throw "Could not resolve tag $($previous.Name)."
            }
            git merge-base --is-ancestor $previousCommit $latestCommit
            if ($LASTEXITCODE -ne 0) {
                throw "The latest stable tag $($latest.Name) is not descended from the preceding stable tag $($previous.Name)."
            }
        }
        $commitsSinceRelease = [int](git rev-list --count "$($latest.Name)..HEAD")
        if ($LASTEXITCODE -ne 0 -or $commitsSinceRelease -eq 0) {
            throw "There are no commits after $($latest.Name) to release."
        }
        $major = $latest.Version.Major
        $minor = $latest.Version.Minor
        if ($Bump -eq 'major') {
            $major++
            $minor = 0
            $patch = 0
        } else {
            $minor++
            $patch = 0
        }
        $version = "$major.$minor.$patch"
        $tag = "v$version"
        if (@(git tag --list $tag).Count -gt 0 -or @(git ls-remote --tags origin "refs/tags/$tag").Count -gt 0) {
            throw "Tag $tag already exists."
        }
    }

    Write-Output 'Checking formatting...'
    $unformatted = @(gofmt -l .)
    if ($LASTEXITCODE -ne 0) { throw 'gofmt failed.' }
    if ($unformatted.Count -gt 0) { throw "Go files need gofmt: $($unformatted -join ', ')." }

    Write-Output 'Running Go tests...'
    Invoke-Checked go @('test', './...')
    Write-Output 'Running go vet...'
    Invoke-Checked go @('vet', './...')
    foreach ($target in @(
        @{ OS = 'windows'; Arch = 'amd64' },
        @{ OS = 'windows'; Arch = 'arm64' },
        @{ OS = 'linux'; Arch = 'amd64' },
        @{ OS = 'linux'; Arch = 'arm64' },
        @{ OS = 'darwin'; Arch = 'amd64' },
        @{ OS = 'darwin'; Arch = 'arm64' }
    )) {
        $oldGoOS = $env:GOOS
        $oldGoArch = $env:GOARCH
        try {
            $env:GOOS = $target.OS
            $env:GOARCH = $target.Arch
            Write-Output "Building $($target.OS)/$($target.Arch)..."
            Invoke-Checked go @('build', '-trimpath', './...')
        } finally {
            $env:GOOS = $oldGoOS
            $env:GOARCH = $oldGoArch
            if ($null -eq $oldGoOS) { Remove-Item Env:GOOS -ErrorAction SilentlyContinue }
            if ($null -eq $oldGoArch) { Remove-Item Env:GOARCH -ErrorAction SilentlyContinue }
        }
    }

    $confirmation = Read-Host "Type $tag to create the annotated tag, push it, and publish the public GitHub release"
    if ($confirmation -cne $tag) {
        throw 'Release confirmation did not match; no tag was pushed.'
    }

    if (-not $resumePending) {
        Invoke-Checked git @('tag', '-a', $tag, '-m', "Tasx $tag", $commit)
        try {
            Invoke-Checked git @('push', 'origin', "refs/tags/$tag")
        } catch {
            $remoteTag = @(git ls-remote --tags origin "refs/tags/$tag")
            if ($remoteTag.Count -eq 0) {
                git tag -d $tag | Out-Null
                throw
            }
            $remoteTagObject = ($remoteTag[0] -split '\s+')[0]
            $localTagObject = (git rev-parse "$tag^{tag}").Trim()
            if ($LASTEXITCODE -ne 0 -or $remoteTagObject -ne $localTagObject) {
                git tag -d $tag | Out-Null
                throw "Tag push failed and origin already has a different $tag tag."
            }
            Write-Warning "Tag push returned an error, but $tag exists on origin; continuing with the published tag."
        }
    }

    $temporaryDirectory = Join-Path ([IO.Path]::GetTempPath()) "tasx-release-$version-$([guid]::NewGuid().ToString('N'))"
    if (Test-Path -LiteralPath $temporaryDirectory) {
        throw "Temporary release directory already exists: $temporaryDirectory"
    }
    New-Item -ItemType Directory -Path $temporaryDirectory | Out-Null
    $temporaryDirectoryCreated = $true
    $published = Complete-Release -Tag $tag -Version $version -Commit $commit -OutputDirectory $temporaryDirectory
    Write-Output "Published $($published.name): $($published.url)"
    Write-Output "Verified assets: $($published.assets.Count)"
} finally {
    if ($temporaryDirectoryCreated -and (Test-Path -LiteralPath $temporaryDirectory)) {
        Remove-Item -LiteralPath $temporaryDirectory -Recurse -Force
    }
    Pop-Location
}
