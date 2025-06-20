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

$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$version = "latest"
$downloadBase = "https://github.com/thewizardshell/froggit/releases/latest/download"
$zipUrl = "$downloadBase/windows-$arch.zip"

$tempPath = "$env:TEMP\froggit"
$installDir = "C:\tools\froggit"
$exePath = "$installDir\froggit.exe"

New-Item -ItemType Directory -Force -Path $tempPath | Out-Null
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

Write-Host "📦 Downloading Froggit for Windows $arch..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $zipUrl -OutFile "$tempPath\froggit.zip"

Write-Host "📂 Extracting..." -ForegroundColor Cyan
Add-Type -AssemblyName System.IO.Compression.FileSystem
[System.IO.Compression.ZipFile]::ExtractToDirectory("$tempPath\froggit.zip", $tempPath, $true)

Move-Item -Force -Path "$tempPath\froggit.exe" -Destination $exePath

$envPath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
$pathChanged = $false
if ($envPath -notlike "*C:\tools\froggit*") {
    Write-Host "🛠 Adding C:\tools\froggit to system PATH..." -ForegroundColor Yellow
    $newPath = "$envPath;C:\tools\froggit"
    [System.Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
    $pathChanged = $true
}

Write-Host "`n✅ Froggit installed at: $exePath" -ForegroundColor Green
Write-Host "👉 You can now use 'froggit' from any terminal." -ForegroundColor Green

if ($pathChanged) {
    Write-Host "🔄 Please restart your terminal or log out/in to apply the updated PATH." -ForegroundColor Yellow
}

