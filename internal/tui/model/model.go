// Package model defines the core application state (Model) for the Froggit TUI.
// It includes view types, the main Model structure, and helper functions
// for initializing and refreshing the application state.
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
	NewBranchView
	ConfirmDialog
	HelpView
	LogGraphView
)

type Model struct {
	Files            []git.FileItem
	Branches         []string
	Remotes          []string
	CurrentBranch    string
	Cursor           int
	CurrentView      View
	CommitMsg        string
	RemoteName       string
	RemoteURL        string
	InputField       string
	Message          string
	MessageType      string
	IsPushing        bool
	SpinnerIndex     int
	SpinnerFrames    []string
	IsFetching       bool
	IsPulling        bool
	NewBranchName    string
	HasRemoteChanges bool
	ShowHelpPanel    bool

	// LogsGraph data
	LogLines []string

	DialogType   string
	DialogTarget string
}

func InitialModel() Model {
	files, _ := git.GetModifiedFiles()
	branches, current := git.GetBranches()
	remotes, _ := git.GetRemotes()
	hasRemoteChanges, _ := git.HasRemoteChanges(current)

	return Model{
		Files:            files,
		Branches:         branches,
		Remotes:          remotes,
		CurrentBranch:    current,
		Cursor:           0,
		CurrentView:      FileView,
		CommitMsg:        "",
		RemoteName:       "",
		RemoteURL:        "",
		InputField:       "",
		Message:          "",
		MessageType:      "",
		SpinnerFrames:    []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		IsPushing:        false,
		SpinnerIndex:     0,
		IsFetching:       false,
		IsPulling:        false,
		NewBranchName:    "",
		HasRemoteChanges: hasRemoteChanges,
		ShowHelpPanel:    false,
		LogLines:         []string{},
		DialogType:       "",
		DialogTarget:     "",
	}
}

func (m *Model) RefreshData() {
	files, _ := git.GetModifiedFiles()
	branches, current := git.GetBranches()
	remotes, _ := git.GetRemotes()
	hasRemoteChanges, _ := git.HasRemoteChanges(current)

	m.Files = files
	m.Branches = branches
	m.Remotes = remotes
	m.CurrentBranch = current
	m.HasRemoteChanges = hasRemoteChanges
}
