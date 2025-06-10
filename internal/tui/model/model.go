package model

import (
	"froggit/internal/git"
)

type View int

const (
	FileView View = iota
	CommitView
	BranchView
	RemoteView
	AddRemoteView
	NewBranchView // Nueva vista
	ConfirmDialog // Nuevo tipo de vista para diálogos de confirmación
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
	IsFetching    bool
	IsPulling     bool
	NewBranchName string // Nuevo campo

	DialogType   string // "delete_branch" o "discard_changes"
	DialogTarget string // Nombre del archivo o rama
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
		IsFetching:    false,
		IsPulling:     false,
		NewBranchName: "",
		DialogType:    "",
		DialogTarget:  "",
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
