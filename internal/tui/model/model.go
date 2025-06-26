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
	AdvancedMode     bool

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
		AdvancedMode:     false,
	}
}

func (m *Model) RefreshData() {
	// Use goroutines to fetch data in parallel
	filesCh := make(chan []git.FileItem)
	branchesCh := make(chan []string)
	currentCh := make(chan string)
	remotesCh := make(chan []string)
	hasRemoteChangesCh := make(chan bool)

	// Files
	go func() {
		files, _ := git.GetModifiedFiles()
		filesCh <- files
	}()
	// Branches and current branch
	go func() {
		branches, current := git.GetBranches()
		branchesCh <- branches
		currentCh <- current
	}()
	// Remotes
	go func() {
		remotes, _ := git.GetRemotes()
		remotesCh <- remotes
	}()

	// Wait for branches to get current branch, then check remote changes
	branches := <-branchesCh
	current := <-currentCh
	go func() {
		hasRemoteChanges, _ := git.HasRemoteChanges(current)
		hasRemoteChangesCh <- hasRemoteChanges
	}()

	m.Files = <-filesCh
	m.Branches = branches
	m.Remotes = <-remotesCh
	m.CurrentBranch = current
	m.HasRemoteChanges = <-hasRemoteChangesCh
}
