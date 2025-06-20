<h1 align="center">Froggit 🐸</h1>

<p align="center">
  <img src="https://github.com/user-attachments/assets/d4194260-341d-425c-872d-ae623c1ec189" alt="Froggit Logo" width="450" />
</p>

<p align="center">
  <b>A modern, intuitive Terminal User Interface (TUI) for Git</b><br>
  Simplify your Git workflow with visual feedback and keyboard-driven controls
</p>

<p align="center">
  <img src="https://img.shields.io/badge/💻%20Platforms-Windows%20%7C%20Linux%20%7C%20macOS-27ae60?style=flat-square" alt="Supported Platforms" />
  <img src="https://img.shields.io/badge/⚙️%20Go-1.20%2B-145a32?style=flat-square" alt="Go Version" />
</p>


<p align="center">
  <a href="https://froggit-docs.vercel.app/guides/install/">Installation</a> •
  <a href="#features">Features</a> •
  <a href="#usage">Usage</a> •
  <a href="https://github.com/thewizardshell/froggit/blob/master/docs/keyboard-shortcuts.md">Shortcuts</a> •
  <a href="https://github.com/thewizardshell/froggit/tree/master/docs">Documentation</a>
</p>


---

## ⚡️ Quick Start

![preview_short_froggit](https://github.com/user-attachments/assets/7b5a2dd6-fbc3-4035-83fe-a072d2298f57)


### Linux / macOS

```bash
curl -s https://raw.githubusercontent.com/thewizardshell/froggit/master/scripts/install.sh | bash
```

### Windows (PowerShell)

```powershell
iwr https://raw.githubusercontent.com/thewizardshell/froggit/master/scripts/install.ps1 -UseBasicParsing | iex
```

> ✅ These scripts will:
>
> - Detect your OS and architecture
> - Move it to a folder in your system PATH (e.g. `/usr/local/bin` or `C:\tools\froggit`)

Once installed, run:

```bash
froggit
```

---

### Manual Build

```bash
git clone https://github.com/thewizardshell/froggit.git
cd froggit
go mod tidy
go build -o froggit
sudo mv froggit /usr/local/bin/
```

---

## Features

- 📁 **Visual File Management**

  - Stage/unstage files with visual feedback
  - Bulk actions for multiple files
  - Real-time status updates

- 🌿 **Branch Operations**

  - Create and switch branches
  - Delete branches safely
  - Visual branch selection

- 🔄 **Git Operations**
  - Commit changes with message preview
  - Push/pull with progress indicators
  - Remote repository management

## ⌨️ Keyboard Shortcuts

### Global

- `q` or `Ctrl+C` - Quit
- `Esc` - Go back/cancel
- `?` - Show help

### File View

- `↑/↓` - Navigate files
- `Space` - Stage/unstage file
- `a` - Stage all files
- `c` - Commit staged files
- `x` - Discard changes

### Branch Management

- `b` - Enter branch view
- `n` - Create new branch
- `d` - Delete branch
- `Enter` - Switch to branch

### Repository Actions

- `p` - Push changes
- `f` - Fetch updates
- `l` - Pull changes
- `r` - Refresh status

## 📚 Documentation

Detailed documentation is available in the [docs/](docs/) directory:

- [Architecture Overview](docs/architecture.md)
- [Development Guide](docs/development.md)
- [Contributing Guidelines](docs/contributing.md)
- [Git Commands Reference](docs/git-commands.md)
- [Keyboard Shortcuts](docs/keyboard-shortcuts.md)

## 🛠️ Built With

- [Go](https://golang.org/) - Performance and cross-platform support
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions

## 🔗 Related Tools & Resources

Froggit is built for simplicity, but there are many excellent Git tools worth exploring:

### Alternatives
- [**LazyGit**](https://github.com/jesseduffield/lazygit) – Feature-rich Git TUI for power users.
- [**GitKraken**](https://www.gitkraken.com/) – Visual Git client with a graphical interface.
- [**tig**](https://github.com/jonas/tig) – Terminal-based Git history browser.
- [**Magit**](https://github.com/magit/magit) – Powerful Git interface for Emacs.

### Learn More About Git
- [**Git Handbook**](https://dgamer007.github.io/Git/#/) – A clear and practical guide to mastering Git concepts.

---

## Author

**Vicente Roa**  
GitHub: [@thewizardshell](https://github.com/thewizardshell)

## License

This project is licensed under the [MIT License](LICENSE).
