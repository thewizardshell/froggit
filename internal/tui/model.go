package tui

import (
	"giteasy/internal/git"
)

type View int

const (
	FileView View = iota
	CommitView
	BranchView
	RemoteView
	AddRemoteView
)

type Model struct {
	Files         []git.FileItem
	Branches      []string
	Remotes       []string
	CurrentBranch string
	Cursor        int
	CurrentView   View
	CommitMsg     string
	RemoteName    string
	RemoteURL     string
	InputField    string // Para determinar qué campo estamos editando
	Message       string
	MessageType   string // "error", "success", ""
	IsPushing     bool
	SpinnerIndex  int
	SpinnerFrames []string
}

func InitialModel() Model {
	files, _ := git.GetModifiedFiles()
	branches, current := git.GetBranches()
	remotes, _ := git.GetRemotes()

	return Model{
		Files:         files,
		Branches:      branches,
		Remotes:       remotes,
		CurrentBranch: current,
		Cursor:        0,
		CurrentView:   FileView,
		CommitMsg:     "",
		RemoteName:    "",
		RemoteURL:     "",
		InputField:    "",
		Message:       "",
		MessageType:   "",
		SpinnerFrames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		IsPushing:     false,
		SpinnerIndex:  0,
	}
}

func (m *Model) RefreshData() {
	files, _ := git.GetModifiedFiles()
	branches, current := git.GetBranches()
	remotes, _ := git.GetRemotes()

	m.Files = files
	m.Branches = branches
	m.Remotes = remotes
	m.CurrentBranch = current
}
