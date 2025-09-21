// Package model defines the core application state (Model) for the Froggit TUI.
// It includes view types, the main Model structure, and helper functions
// for initializing and refreshing the application state.
package model

import (
	"froggit/internal/gh"
	"froggit/internal/git"
	"strings"
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
	RepositoryListView
	ConfirmCloneRepoView
	GitHubControlsView
	MergeView
	RebaseView
	StashView
	StashMessageView
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
	AutoFetchDone    bool
	NewBranchName    string
	HasRemoteChanges bool
	ShowHelpPanel    bool
	AdvancedMode     bool

	// LogsGraph data
	LogLines []string

	// Repositories
	Repositories      []gh.Repository
	SelectedRepoIndex int
	RepoToClone       *gh.Repository

	// Merge/Rebase interactive state
	SelectedBranch      string   // branch selected for merge/rebase
	MergeConflictFiles  []string // files in conflict during merge
	RebaseConflictFiles []string // files in conflict during rebase
	IsMerging           bool
	IsRebasing          bool
	MergeStep           string // "select", "confirm", "conflict"
	RebaseStep          string // "select", "confirm", "conflict"

	// Stash data
	Stashes       []string
	StashMessage  string
	SelectedStash int
	IsStashing    bool

	DialogType   string
	DialogTarget string

	// Awaiting push after merge/rebase
	AwaitingPush bool
}

func InitialModel() Model {
	files, _ := git.GetModifiedFiles()
	branches, current := git.GetBranches()
	remotes, _ := git.GetRemotes()

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
		HasRemoteChanges: false,
		ShowHelpPanel:    false,
		LogLines:         []string{},
		DialogType:       "",
		DialogTarget:     "",
		AdvancedMode:     false,
		AwaitingPush:     false,
		Stashes:          []string{},
		StashMessage:     "",
		SelectedStash:    0,
		IsStashing:       false,
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

	// Wait for branches to get current branch, then check remote changes without fetch
	branches := <-branchesCh
	current := <-currentCh
	go func() {
		hasRemoteChanges, _ := git.HasRemoteChangesWithFetch(current, false)
		hasRemoteChangesCh <- hasRemoteChanges
	}()

	m.Files = <-filesCh
	m.Branches = branches
	m.Remotes = <-remotesCh
	m.CurrentBranch = current
	m.HasRemoteChanges = <-hasRemoteChangesCh

	// Add stash refresh
	stashesCh := make(chan []string)
	go func() {
		stashOutput, _ := git.StashList()
		stashes := parseStashList(stashOutput)
		stashesCh <- stashes
	}()
	m.Stashes = <-stashesCh
}

// parseStashList parses the output of git stash list
func parseStashList(output string) []string {
	if output == "" {
		return []string{}
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var stashes []string
	for _, line := range lines {
		if line != "" {
			stashes = append(stashes, line)
		}
	}
	return stashes
}

func (m *Model) SetMergeConflictFiles(files []string) {
	m.MergeConflictFiles = files
}

func (m *Model) SetRebaseConflictFiles(files []string) {
	m.RebaseConflictFiles = files
}
