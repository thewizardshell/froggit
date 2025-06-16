<h1 align="center">Froggit 🐸</h1>

<p align="center">
  <img src="https://github.com/user-attachments/assets/d4194260-341d-425c-872d-ae623c1ec189" alt="Froggit Logo" width="450" />
</p>

<p align="center">
  <strong>A modern, intuitive Terminal User Interface (TUI) for Git</strong><br>
  Simplify your Git workflow with visual feedback and keyboard-driven controls
</p>

<p align="center">
  <a href="https://froggit-docs.vercel.app/guides/install/">Installation</a> •
  <a href="#features">Features</a> •
  <a href="#usage">Usage</a> •
  <a href="https://github.com/thewizardshell/froggit/blob/master/docs/keyboard-shortcuts.md">Shortcuts</a> •
  <a href="https://github.com/thewizardshell/froggit/tree/master/docs">Documentation</a>
</p>

## ⚡️ Quick Start

![preview_short_froggit](https://github.com/user-attachments/assets/7b5a2dd6-fbc3-4035-83fe-a072d2298f57)

```bash
# Install

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

# 🔗 Related Tools and Alternatives
While Froggit is designed to be a minimal and beginner-friendly Git TUI, there are several excellent tools out there that offer more advanced features or different workflows. Depending on your needs, you might want to explore them as well:

- LazyGit: A powerful and widely-used Git TUI, great for intermediate and advanced users: https://github.com/jesseduffield/lazygit.
- GitKraken: A GUI Git client with a rich graphical interface: https://www.gitkraken.com/.
- tig : Text-mode interface for Git, useful for browsing history and commits: https://github.com/jonas/tig.
- Magit: An interface to Git for Emacs users, highly customizable and powerful: https://github.com/magit/magit.
  
Froggit aims to serve as a stepping stone — especially for those new to Git or who prefer simplicity over complexity. Feel free to try different tools and find what works best for your workflow!
