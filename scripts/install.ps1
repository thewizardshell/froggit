# Banner green ASCII art

param([switch]$Force)

$banner = @"
 __      __   _                    ___                   _ _   
 \ \    / /__| |__ ___ _ __  ___  | __| _ ___  __ _ __ _(_) |_ 
  \ \/\/ / -_) / _/ _ \ '  \/ -_) | _| '_/ _ \/ _` / _` | |  _|
   \_/\_/\___|_\__\___/_|_|_\___| |_||_| \___/\__, \__, |_|\__|
                                              |___/|___/       
"@

Write-Host $banner -ForegroundColor Green


# Confirmation prompt if not forced
if (-not $Force) {
    $confirm = Read-Host "Do you want to continue with the installation? (y/n)"
    if ($confirm.ToLower() -ne 'y') {
        Write-Host "Installation cancelled." -ForegroundColor Yellow
        exit
    }
}

# Check for admin rights
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin -and -not $Force) {
    Write-Host "WARNING: It is recommended to run as Administrator" -ForegroundColor Yellow
    $continue = Read-Host "Continue anyway? (y/n)"
    if ($continue.ToLower() -ne 'y') { exit }
}

# Configuration
$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$url = "https://github.com/thewizardshell/froggit/releases/latest/download/windows-$arch.zip"
$tempDir = "$env:TEMP\froggit-install"
$installDir = "C:\tools\froggit"
$zipFile = "$tempDir\froggit.zip"

Write-Host "Architecture: $arch" -ForegroundColor Cyan
Write-Host "Downloading from: $url" -ForegroundColor Cyan

# Clean and create directories
if (Test-Path $tempDir) { Remove-Item -Recurse -Force $tempDir }
if (Test-Path $installDir) { Remove-Item -Recurse -Force $installDir }
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

# Download
Write-Host "Downloading..." -ForegroundColor Yellow
try {
    Invoke-WebRequest -Uri $url -OutFile $zipFile -UseBasicParsing
    Write-Host "Download completed" -ForegroundColor Green
} catch {
    Write-Error "Download error: $($_.Exception.Message)"
    exit 1
}

# Verify download
if (-not (Test-Path $zipFile)) {
    Write-Error "File not found after download"
    exit 1
}

$fileSize = (Get-Item $zipFile).Length
Write-Host ("Downloaded file size: {0} KB" -f [math]::Round($fileSize/1KB, 2)) -ForegroundColor Cyan

# Extract
Write-Host "Extracting..." -ForegroundColor Yellow
try {
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    [System.IO.Compression.ZipFile]::ExtractToDirectory($zipFile, $tempDir)
    Write-Host "Extraction completed" -ForegroundColor Green
} catch {
    Write-Error "Extraction error: $($_.Exception.Message)"
    exit 1
}

# Find executable
$exeFiles = Get-ChildItem -Path $tempDir -Filter "*.exe" -Recurse
if ($exeFiles.Count -eq 0) {
    Write-Error "No executable found in the archive"
    Write-Host "Files found:" -ForegroundColor Yellow
    Get-ChildItem -Path $tempDir -Recurse | ForEach-Object { Write-Host "  $($_.Name)" }
    exit 1
}

$sourceExe = $exeFiles[0].FullName
$targetExe = "$installDir\froggit.exe"

# Copy executable
Write-Host "Installing executable..." -ForegroundColor Yellow
try {
    Copy-Item -Path $sourceExe -Destination $targetExe -Force
    Write-Host "Executable installed at: $targetExe" -ForegroundColor Green
} catch {
    Write-Error "Error copying executable: $($_.Exception.Message)"
    exit 1
}

# Update PATH if possible
$currentPath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
if ($currentPath -notlike "*$installDir*") {
    if ($isAdmin) {
        Write-Host "Updating system PATH..." -ForegroundColor Yellow
        try {
            $newPath = $currentPath + ";" + $installDir
            [System.Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
            Write-Host "PATH updated" -ForegroundColor Green
            $pathUpdated = $true
        } catch {
            Write-Warning "Could not update PATH: $($_.Exception.Message)"
            $pathUpdated = $false
        }
    } else {
        Write-Host "To add to PATH, run as administrator or add manually:" -ForegroundColor Yellow
        Write-Host "  $installDir" -ForegroundColor Cyan
        $pathUpdated = $false
    }
}

# Test installation
Write-Host "Testing installation..." -ForegroundColor Yellow
try {
    $version = & $targetExe --version 2>$null
    if ($version) {
        Write-Host "Froggit installed successfully: $version" -ForegroundColor Green
    } else {
        Write-Host "Froggit installed (version not available)" -ForegroundColor Green
    }
} catch {
    Write-Host "Froggit installed (manual test recommended)" -ForegroundColor Yellow
}

# Clean up
Remove-Item -Recurse -Force $tempDir -ErrorAction SilentlyContinue

# Summary
Write-Host "`n=== INSTALLATION COMPLETE ===" -ForegroundColor Green
Write-Host "Location: $targetExe" -ForegroundColor Cyan
if ($pathUpdated) {
    Write-Host "Restart your terminal to use 'froggit' from any location" -ForegroundColor Yellow
} else {
    Write-Host "To use from any location, add to PATH: $installDir" -ForegroundColor Yellow
}
Write-Host "Or run directly: $targetExe" -ForegroundColor Cyan

