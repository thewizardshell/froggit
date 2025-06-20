# Froggit PowerShell Installer
Write-Host "
 __      __   _                    ___                   _ _   
 \ \    / /__| |__ ___ _ __  ___  | __| _ ___  __ _ __ _(_) |_ 
  \ \/\/ / -_) / _/ _ \ '  \/ -_) | _| '_/ _ \/ _` / _` | |  _|
   \_/\_/\___|_\__\___/_|_|_\___| |_||_| \___/\__, \__, |_|\__|
                                              |___/|___/        
" -ForegroundColor Green

# Ask user if they want to continue
$answer = Read-Host "Do you want to continue with the Froggit installation? (y/n)"
if ($answer -notmatch '^[Yy]$') {
    Write-Host "Installation cancelled by user." -ForegroundColor Red
    exit 0
}

# Detect architecture
$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$downloadBase = "https://github.com/thewizardshell/froggit/releases/latest/download"
$zipUrl = "$downloadBase/windows-$arch.zip"

$tempPath = "$env:TEMP\froggit"
$installDir = "C:\tools\froggit"
$exePath = "$installDir\froggit.exe"

# Create temp and install directories
New-Item -ItemType Directory -Force -Path $tempPath | Out-Null
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

Write-Host "📦 Downloading Froggit for Windows $arch..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $zipUrl -OutFile "$tempPath\froggit.zip"

Write-Host "📂 Extracting..." -ForegroundColor Cyan
Add-Type -AssemblyName System.IO.Compression.FileSystem
try {
    [System.IO.Compression.ZipFile]::ExtractToDirectory("$tempPath\froggit.zip", $tempPath, $true)
} catch {
    Write-Error "Failed to extract archive: $_"
    exit 1
}

if (-Not (Test-Path "$tempPath\froggit.exe")) {
    Write-Error "froggit.exe not found after extraction!"
    exit 1
}

Write-Host "🚚 Installing executable..." -ForegroundColor Cyan
Move-Item -Force -Path "$tempPath\froggit.exe" -Destination $exePath

# Add install directory to system PATH if not present
$envPath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
$pathChanged = $false
if ($envPath -notlike "*$installDir*") {
    Write-Host "🛠 Adding $installDir to system PATH..." -ForegroundColor Yellow
    $newPath = "$envPath;$installDir"
    try {
        [System.Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
        $pathChanged = $true
    } catch {
        Write-Warning "Failed to update system PATH. Please run this script as Administrator."
    }
}

Write-Host "`n✅ Froggit installed at: $exePath" -ForegroundColor Green
Write-Host "👉 You can now run 'froggit' from any terminal." -ForegroundColor Green

if ($pathChanged) {
    Write-Host "🔄 Please restart your terminal or log out/in to apply the updated PATH." -ForegroundColor Yellow
}

# Cleanup
Remove-Item -Force "$tempPath\froggit.zip"
Remove-Item -Recurse -Force $tempPath

