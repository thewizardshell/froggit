#!/usr/bin/env bash

set -e

REPO="thewizardshell/froggit"
BIN_NAME="froggit"

# ASCII Banner
echo -e "\033[1;32m"
cat <<"EOF"
 __      __   _                    ___                   _ _   
 \ \    / /__| |__ ___ _ __  ___  | __| _ ___  __ _ __ _(_) |_ 
  \ \/\/ / -_) / _/ _ \ '  \/ -_) | _| '_/ _ \/ _` / _` | |  _|
   \_/\_/\___|_\__\___/_|_|_\___| |_||_| \___/\__, \__, |_|\__|
                                              |___/|___/
EOF
echo -e "\033[0m"

# Ask user before continuing (especially important for Windows users)
read -r -p "Do you want to continue? [y/N] " answer < /dev/tty
if [[ ! "$answer" =~ ^[Yy]$ ]]; then
  echo "Installation cancelled by user."
  exit 0
fi

# Loading animation function
loading() {
  local msg="$1"
  echo -n "$msg"
  for i in {1..3}; do
    echo -n "."
    sleep 0.3
  done
  echo
}

loading "🔍 Detecting system"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
linux*) OS="linux" ;;
darwin*) OS="darwin" ;;
msys* | cygwin* | mingw*) OS="windows" ;;
*) echo "❌ Unsupported OS: $OS" && exit 1 ;;
esac

# Detect ARCH
ARCH=$(uname -m)
case "$ARCH" in
x86_64) ARCH="amd64" ;;
aarch64 | arm64) ARCH="arm64" ;;
*) echo "❌ Unsupported architecture: $ARCH" && exit 1 ;;
esac

if [[ "$OS" == "windows" ]]; then
  echo -n "🪟 Windows detected. Do you want to run the PowerShell installer? (y/n): "
  read -r answer
  if [[ "$answer" =~ ^[Yy]$ ]]; then
    loading "Running PowerShell installer"
    powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "iwr -useb https://raw.githubusercontent.com/thewizardshell/froggit/master/scripts/install.ps1 | iex"
  else
    echo "Installation cancelled by user."
  fi
  exit 0
fi

ZIP_NAME="${OS}-${ARCH}.zip"
URL="https://github.com/${REPO}/releases/latest/download/${ZIP_NAME}"

loading "⬇️  Downloading Froggit ${OS}-${ARCH}"
curl -L -o froggit.zip "$URL"

loading "📦 Unzipping"
unzip -o froggit.zip

# The extracted binary has the OS-ARCH suffix
EXTRACTED_NAME="${BIN_NAME}-${OS}-${ARCH}"
FINAL_NAME="${BIN_NAME}"

loading "🚚 Installing to /usr/local/bin"
chmod +x "$EXTRACTED_NAME"
sudo mv "$EXTRACTED_NAME" /usr/local/bin/froggit

echo -e "\n✅ \033[1;32mFroggit installed successfully!\033[0m"
echo "👉 Run 'froggit' to get started 🐸"

# Cleanup
rm froggit.zip
