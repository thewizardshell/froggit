<h1 align="center"> Froggit </h1>

<div align="center">
  <img src="https://github.com/user-attachments/assets/d4194260-341d-425c-872d-ae623c1ec189" alt="Froggit Logo" width="400" />
  
  <p><strong>A modern, minimalist Git TUI</strong></p>
  <p>Designed for clarity, speed, and smooth integration with your terminal workflow.</p>

  [![Release](https://img.shields.io/github/v/release/thewizardshell/froggit?label=Release&color=27ae60)](https://github.com/thewizardshell/froggit/releases)
  [![Installer Tests](https://github.com/thewizardshell/froggit/actions/workflows/test-installers.yml/badge.svg)](https://github.com/thewizardshell/froggit/actions/workflows/test-installers.yml)
  [![Platforms](https://img.shields.io/badge/Platforms-Windows%20%7C%20Linux%20%7C%20macOS-2ecc71?style=flat)](#)
  [![Go Version](https://img.shields.io/badge/Go-1.20%2B-145a32?style=flat)](#)

  <img src="https://github.com/user-attachments/assets/7b5a2dd6-fbc3-4035-83fe-a072d2298f57" alt="Froggit preview" width="700" />
</div>



## Installation

### Quick Install (Recommended)

**Linux / macOS:**
```bash
curl -s https://raw.githubusercontent.com/thewizardshell/froggit/master/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
iwr https://raw.githubusercontent.com/thewizardshell/froggit/master/scripts/install.ps1 -UseBasicParsing | iex
```

### Build from Source

```bash
git clone https://github.com/thewizardshell/froggit.git
cd froggit
go mod tidy
go build
./froggit
```

## Configuration

Froggit can be customized using a `froggit.yml` configuration file. The configuration file should be placed in the same directory as the Froggit executable.

### Creating Configuration File

Create a `froggit.yml` file next to your Froggit executable with the following structure:

```yaml
ui:
  branding: true          # Show Froggit branding (default: true)
  position: "center"      # UI position: "left", "center", "right" (default: "left")

git:
  autofetch: true         # Automatically fetch from remote (default: true)
  defaultbranch: "main"   # Default branch for new repositories (default: "main")
```

### Configuration Options

#### UI Settings (`ui`)
| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `branding` | boolean | `true` | Display Froggit branding and visual elements |
| `position` | string | `"left"` | UI positioning: `"left"`, `"center"`, or `"right"` |

#### Git Settings (`git`)
| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `autofetch` | boolean | `true` | Automatically fetch from remote repositories on startup |
| `defaultbranch` | string | `"main"` | Default branch name for new repositories and push operations |

### Example Configurations

**Minimal Configuration:**
```yaml
git:
  defaultbranch: "master"
```

**Full Configuration:**
```yaml
ui:
  branding: false
  position: "center"
git:
  autofetch: false
  defaultbranch: "develop"
```

**Note:** If no configuration file is found, Froggit will use the default values shown above.

## Requirements

- [Git](https://git-scm.com/) installed and accessible in your terminal
- [Go 1.20+](https://go.dev/dl/) (only required if building from source)
- A terminal with [Nerd Fonts](https://www.nerdfonts.com/) support
- [GitHub CLI](https://cli.github.com/) (`gh`) – optional, for GitHub integration

## Feature Support

### Git Operations
| Feature | Status | Description |
|---------|--------|-------------|
| Stage all | 🟢 | Stage all changes |
| Branches | 🟢 | View and manage branches |
| Remotes | 🟢 | Manage remote repositories |
| Push | 🟢 | Push changes to remote |
| Fetch | 🟢 | Fetch from remote |
| Pull | 🟢 | Pull changes (when remote changes available) |
| Commit | 🟢 | Create commits |
| Discard changes | 🟢 | Discard uncommitted changes |
| Refresh | 🟢 | Refresh repository status |
| Advanced mode | 🟢 | Access to logs, merge, stash, rebase |
| Logs | 🟢 | View commit history |
| Merge | 🟢 | Merge branches |
| Stash | 🟡 | Stash changes |
| Rebase | 🟢 | Rebase branches |

### GitHub CLI Integration
| Feature | Status | Description |
|---------|--------|-------------|
| Create repository | 🟢 | Create new GitHub repository |
| Clone repository | 🟢 | Clone from your GitHub repositories |

🟢 Supported &nbsp;&nbsp; 🟡 In Development &nbsp;&nbsp; 🔴 Planned


## GitHub CLI Integration

Froggit integrates seamlessly with [GitHub CLI](https://cli.github.com/) to enhance your workflow.

```bash
gh auth login
```

Once authenticated, Froggit will detect `gh` and enable features like cloning repositories directly from GitHub.

<div align="center">
  <img src="https://github.com/user-attachments/assets/8f3de6e0-16bf-4ac6-bc91-d434512df4d1" alt="GitHub CLI integration in Froggit" width="700" />
</div>

## Key Shortcuts

### File Management
- `↑ / ↓`: Navigate files
- `Space`: Stage/unstage files
- `a`: Stage all changes
- `x`: Discard changes
- `c`: Commit changes

### Branch Operations
- `b`: View branches
- `n`: Create new branch
- `d`: Delete branch
- `Enter`: Switch branch

### Advanced Mode
- `A`: Enter advanced mode
- `M`: Merge (in advanced mode)
- `R`: Rebase (in advanced mode)

### Global
- `q`, `Ctrl+C`: Quit
- `Esc`: Go back
- `?`: Show help

For a complete list of shortcuts, see the [keyboard shortcuts documentation](docs/keyboard-shortcuts.md).

## Documentation

- [Installation Guide](https://froggit-docs.vercel.app/guides/install/)
- [Architecture Overview](docs/architecture.md)
- [Development Guide](docs/development.md)
- [Contributing Guidelines](docs/contributing.md)
- [Keyboard Shortcuts](docs/keyboard-shortcuts.md)

## Related Tools

- [LazyGit](https://github.com/jesseduffield/lazygit) – Git TUI for power users
- [tig](https://github.com/jonas/tig) – Terminal Git history browser
- [Magit](https://github.com/magit/magit) – Git for Emacs




## Learn More About Git

* [**Git Handbook**](https://alias404.github.io/Git/#/) – A practical, visual Git reference



## Contributing

We welcome contributions! Please see our [Contributing Guidelines](docs/contributing.md) for details on how to get started.



## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=thewizardshell/froggit&type=Date)](https://www.star-history.com/#thewizardshell/froggit&Date)


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Vicente Roa**  
GitHub: [@thewizardshell](https://github.com/thewizardshell)

---

<div align="center">
  <img src="https://github.com/user-attachments/assets/123b5ff0-da29-48b1-b5bf-e39d670642d6" alt="Froggit" width="100" />
</div>
