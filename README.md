<h1 align="center">Froggit 🐸</h1>

<p align="center">
  <img src="assets/logo.png" alt="Froggit Logo" width="450" />
</p>

<p align="center">
  <strong>A modern, intuitive Terminal User Interface (TUI) for Git</strong><br>
  Simplify your Git workflow with visual feedback and keyboard-driven controls
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#features">Features</a> •
  <a href="#usage">Usage</a> •
  <a href="#keyboard-shortcuts">Shortcuts</a> •
  <a href="docs/">Documentation</a>
</p>

## ⚡️ Quick Start

```bash
# Install
go install github.com/thewizardshell/froggit@latest

# Or build from source
git clone https://github.com/thewizardshell/froggit.git
cd froggit
go build

# Run
froggit
```

## 🎯 Features

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

## 📝 License

MIT License - see [LICENSE](LICENSE) for details

## 👤 Author

**Vicente Roa** ([@thewizardshell](https://github.com/thewizardshell))

---

<p align="center">
  Made with ❤️ and Go
</p>