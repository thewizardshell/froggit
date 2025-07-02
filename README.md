<h1 align="center">Froggit</h1>

<p align="center">
  <img src="https://github.com/user-attachments/assets/d4194260-341d-425c-872d-ae623c1ec189" alt="Froggit Logo" width="450" />
</p>

<p align="center">
  <b>A modern, minimalist Git TUI</b><br>
  Designed for clarity, speed, and smooth integration with your terminal workflow.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Platforms-Windows%20%7C%20Linux%20%7C%20macOS-27ae60?style=flat-square" alt="Supported Platforms" />
  <img src="https://img.shields.io/badge/Go-1.20%2B-145a32?style=flat-square" alt="Go Version" />
</p>

<p align="center">
  <a href="https://froggit-docs.vercel.app/guides/install/">Installation</a> •
  <a href="#features">Features</a> •
  <a href="#github-cli-integration">GitHub CLI Integration</a> •
  <a href="#shortcuts">Shortcuts</a> •
  <a href="#documentation">Documentation</a>
</p>

<p align="center">
  <img src="https://github.com/user-attachments/assets/7b5a2dd6-fbc3-4035-83fe-a072d2298f57" alt="Froggit preview" width="700" />
</p>

---

## Requirements

- [Git](https://git-scm.com/) installed and accessible in your terminal
- [Go 1.20+](https://go.dev/dl/) (only required if building from source)
- A terminal with [Nerd Fonts](https://www.nerdfonts.com/) support for best UI experience
- [GitHub CLI](https://cli.github.com/) (`gh`) – optional, for GitHub integration

---

## Installation

### Option 1: Installer Script (Recommended)

#### Linux / macOS

```bash
curl -s https://raw.githubusercontent.com/thewizardshell/froggit/master/scripts/install.sh | bash
````

#### Windows (PowerShell)

```powershell
iwr https://raw.githubusercontent.com/thewizardshell/froggit/master/scripts/install.ps1 -UseBasicParsing | iex
```

---

### Option 2: Build from Source

```bash
git clone https://github.com/thewizardshell/froggit.git
cd froggit
go mod tidy
go build
./froggit
```

---

## Features

* **Visual Git Interface**

  * Stage, unstage, discard changes
  * View logs and diffs interactively

* **Branch Management**

  * Create, switch, and delete branches

* **GitHub CLI Integration**

  * Clone and explore repositories using `gh` (if installed)

* **Other Tools**

  * Stash support
  * Rebase and merge support
  * Commit previews with feedback

---

## GitHub CLI Integration

Froggit optionally integrates with [GitHub CLI](https://cli.github.com/) for features like cloning and listing repositories.

To set it up:

```bash
gh auth login
```

This will walk you through authenticating with GitHub.

Once logged in, Froggit will detect `gh` and unlock additional options for cloning and working with your GitHub repositories.

<p align="center">
  <!-- Aquí va tu nuevo GIF de integración con GitHub CLI -->
  <img src="https://github.com/user-attachments/assets/8f3de6e0-16bf-4ac6-bc91-d434512df4d1" alt="GitHub CLI integration in Froggit" />
</p>

---

## Shortcuts

(See full list in [keyboard-shortcuts.md](docs/keyboard-shortcuts.md))

### File View

* `↑ / ↓`: Navigate files
* `Space`: Stage/unstage
* `a`: Stage all
* `x`: Discard changes
* `c`: Commit

### Branches

* `b`: Open branch view
* `n`: New branch
* `d`: Delete branch
* `Enter`: Switch

### Global

* `q`, `Ctrl+C`: Quit
* `Esc`: Back
* `?`: Help

---

## Documentation

* [Installation Guide](https://froggit-docs.vercel.app/guides/install/)
* [Architecture Overview](docs/architecture.md)
* [Development Guide](docs/development.md)
* [Contributing](docs/contributing.md)
* [Keyboard Shortcuts](docs/keyboard-shortcuts.md)

---

## Learn More About Git

* [**Git Handbook**](https://dgamer007.github.io/Git/#/) – A clear and practical guide to mastering Git concepts.

---

## Related Tools

* [LazyGit](https://github.com/jesseduffield/lazygit) – Git TUI for power users
* [tig](https://github.com/jonas/tig) – Terminal Git history browser
* [Magit](https://github.com/magit/magit) – Git for Emacs

---

## Author

**Vicente Roa**
GitHub: [@thewizardshell](https://github.com/thewizardshell)

<p align="center">
  <!-- Aquí va tu nuevo GIF de integración con GitHub CLI -->
  <img src="https://github.com/user-attachments/assets/123b5ff0-da29-48b1-b5bf-e39d670642d6" alt="GitHub CLI integration in Froggit" width="100" />
</p>


