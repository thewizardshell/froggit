package model

import (
	"froggit/internal/copilot"
	"froggit/internal/gh"
	"froggit/internal/git"
	"strings"
)

type View int

const (
	QuickStartView View = iota
	FileView
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

	LogLines []string

	Repositories      []gh.Repository
	SelectedRepoIndex int
	RepoToClone       *gh.Repository

	SelectedBranch      string   // branch selected for merge/rebase
	MergeConflictFiles  []string // files in conflict during merge
	RebaseConflictFiles []string // files in conflict during rebase
	IsMerging           bool
	IsRebasing          bool
	MergeStep           string // "select", "confirm", "conflict"
	RebaseStep          string // "select", "confirm", "conflict"

	Stashes       []string
	StashMessage  string
	SelectedStash int
	IsStashing    bool

	FileViewOffset int
	FileViewHeight int

	IsGeneratingAI   bool
	CopilotAvailable bool

	DialogType   string
	DialogTarget string

	AwaitingPush bool

	QuickStartOptions []string
	HasGitHubCLI      bool
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
		CopilotAvailable: copilot.IsAvailable(),
		FileViewOffset:   0,
		FileViewHeight:   5,
	}
}

func (m *Model) RefreshData() {
	filesCh := make(chan []git.FileItem)
	branchesCh := make(chan []string)
	currentCh := make(chan string)
	remotesCh := make(chan []string)
	hasRemoteChangesCh := make(chan bool)

	go func() {
		files, _ := git.GetModifiedFiles()
		filesCh <- files
	}()
	go func() {
		branches, current := git.GetBranches()
		branchesCh <- branches
		currentCh <- current
	}()
	go func() {
		remotes, _ := git.GetRemotes()
		remotesCh <- remotes
	}()

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

	stashesCh := make(chan []string)
	go func() {
		stashOutput, _ := git.StashList()
		stashes := parseStashList(stashOutput)
		stashesCh <- stashes
	}()
	m.Stashes = <-stashesCh
}

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
