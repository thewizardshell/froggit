// Package update handles the core message loop for user interactions
// in the Froggit TUI application, using the Bubble Tea framework.
package update

import (
	"fmt"
	"time"

	"froggit/internal/gh"
	"froggit/internal/git"
	"froggit/internal/tui/model"

	tea "github.com/charmbracelet/bubbletea"
)

// Internal messages for async operations
type pushMsg struct{ err error }
type fetchMsg struct{ err error }
type pullMsg struct{ err error }
type spinnerTickMsg struct{}
type switchBranchMsg struct {
	err          error
	targetBranch string
	nextAction   string // "merge" or "rebase"
	sourceBranch string
}

// spinner returns a Cmd that emits spinnerTickMsg every 100ms.
func spinner() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

// performPush runs git.Push asynchronously and returns a pushMsg.
func performPush() tea.Cmd {
	return func() tea.Msg {
		return pushMsg{err: git.Push()}
	}
}

// performFetch runs git.Fetch asynchronously and returns a fetchMsg.
func performFetch() tea.Cmd {
	return func() tea.Msg {
		return fetchMsg{err: git.Fetch()}
	}
}

// performPull runs git.Pull asynchronously and returns a pullMsg.
func performPull() tea.Cmd {
	return func() tea.Msg {
		return pullMsg{err: git.Pull()}
	}
}

// performSwitchAndMerge switches to target branch and then merges source branch
func performSwitchAndMerge(targetBranch, sourceBranch string) tea.Cmd {
	return func() tea.Msg {
		err := git.Checkout(targetBranch)
		return switchBranchMsg{
			err:          err,
			targetBranch: targetBranch,
			nextAction:   "merge",
			sourceBranch: sourceBranch,
		}
	}
}

// performSwitchAndRebase switches to source branch and then rebases onto target
func performSwitchAndRebase(sourceBranch, targetBranch string) tea.Cmd {
	return func() tea.Msg {
		err := git.Checkout(sourceBranch)
		return switchBranchMsg{
			err:          err,
			targetBranch: targetBranch,
			nextAction:   "rebase",
			sourceBranch: sourceBranch,
		}
	}
}

var ghClient = gh.NewGhClient() // You can move this elsewhere if you prefer

// GetGhClient returns the shared GitHub client instance
func GetGhClient() *gh.GhClient {
	return ghClient
}

// Update handles Bubble Tea messages and returns the updated model
// along with the next command to execute.
func Update(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	// --- Confirm Clone Repo Dialog ---
	if m.CurrentView == model.ConfirmCloneRepoView {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "y":
				if m.RepoToClone != nil {
					repo := m.RepoToClone
					repoFullName := repo.Owner.Login + "/" + repo.Name
					err := gh.CloneRepository(ghClient, repoFullName, "")
					if err != nil {
						m.Message = "✗ Error cloning: " + err.Error()
						m.MessageType = "error"
					} else {
						m.Message = "✓ Repository cloned successfully"
						m.MessageType = "success"
					}
				}
				m.CurrentView = model.RepositoryListView
				m.RepoToClone = nil
				return m, nil
			case "n", "esc":
				m.CurrentView = model.RepositoryListView
				m.RepoToClone = nil
				return m, nil
			}
		}
	}

	switch msg := msg.(type) {

	case switchBranchMsg:
		if msg.err != nil {
			m.Message = fmt.Sprintf("✗ Error switching to branch %s: %s", msg.targetBranch, msg.err)
			m.MessageType = "error"
			return m, nil
		}

		// Update current branch
		m.CurrentBranch = msg.targetBranch
		m.RefreshData()

		if msg.nextAction == "merge" {
			// Now perform the merge
			err := git.Merge(msg.sourceBranch)
			if err != nil {
				m.Message = fmt.Sprintf("✗ Error merging %s into %s: %s", msg.sourceBranch, msg.targetBranch, err)
				m.MessageType = "error"
				return m, nil
			}

			conflicts, _ := git.GetConflictFiles()
			if len(conflicts) > 0 {
				m.LogLines = conflicts
				m.Message = fmt.Sprintf("Conflicts detected while merging %s into %s. Please resolve them and use [P] Proceed or [X] Cancel.", msg.sourceBranch, msg.targetBranch)
				m.MessageType = "warning"
				return m, nil
			} else {
				m.Message = fmt.Sprintf("✓ Successfully merged %s into %s. Press [P] to push to remote.", msg.sourceBranch, msg.targetBranch)
				m.MessageType = "success"
				m.CurrentView = model.MergeView
				m.AwaitingPush = true
				return m, nil
			}
		} else if msg.nextAction == "rebase" {
			// Now perform the rebase
			err := git.Rebase(msg.targetBranch)
			if err != nil {
				m.Message = fmt.Sprintf("✗ Error rebasing %s onto %s: %s", msg.sourceBranch, msg.targetBranch, err)
				m.MessageType = "error"
				return m, nil
			}

			conflicts, _ := git.GetConflictFiles()
			if len(conflicts) > 0 {
				m.LogLines = conflicts
				m.Message = fmt.Sprintf("Conflicts detected while rebasing %s onto %s. Please resolve them and use [P] Proceed or [X] Cancel.", msg.sourceBranch, msg.targetBranch)
				m.MessageType = "warning"
				return m, nil
			} else {
				m.Message = fmt.Sprintf("✓ Successfully rebased %s onto %s", msg.sourceBranch, msg.targetBranch)
				m.MessageType = "success"
				m.CurrentView = model.FileView
				m.DialogTarget = ""
				m.LogLines = nil
				m.RefreshData()
				return m, nil
			}
		}

	case tea.KeyMsg:
		// --- GitHub Controls removed from main flow ---
		// The repository view is only accessed at startup if there is no git init

		// --- Repository List Navigation ---
		if m.CurrentView == model.RepositoryListView {
			switch msg.String() {
			case "up":
				if m.SelectedRepoIndex > 0 {
					m.SelectedRepoIndex--
				}
				return m, nil
			case "down":
				if m.SelectedRepoIndex < len(m.Repositories)-1 {
					m.SelectedRepoIndex++
				}
				return m, nil
			case "esc":
				m.CurrentView = model.GitHubControlsView
				return m, nil
			case "c":
				if len(m.Repositories) > 0 {
					m.RepoToClone = &m.Repositories[m.SelectedRepoIndex]
					m.CurrentView = model.ConfirmCloneRepoView
				}
				return m, nil
			}
		}

		if m.CurrentView == model.LogGraphView {
			return HandleLogGraphKey(m, msg)
		}

		// --- MergeView controls ---
		if m.CurrentView == model.MergeView {
			switch msg.String() {
			case "up":
				if m.Cursor > 0 {
					m.Cursor--
				}
				return m, nil
			case "down":
				if m.Cursor < len(m.Branches)-1 {
					m.Cursor++
				}
				return m, nil
			case " ", "space":
				if len(m.Branches) > 0 && m.Cursor < len(m.Branches) {
					selected := m.Branches[m.Cursor]
					if selected == m.CurrentBranch {
						m.Message = "⚠ Cannot merge a branch into itself"
						m.MessageType = "warning"
						return m, nil
					}
					if m.DialogTarget == selected {
						m.DialogTarget = "" // unselect if already selected
						m.Message = "Selection unmarked"
						m.MessageType = "info"
					} else {
						m.DialogTarget = selected
						m.Message = fmt.Sprintf("Will merge %s into %s (current: %s)", m.CurrentBranch, selected, m.CurrentBranch)
						m.MessageType = "info"
					}
				}
				return m, nil
			case "M", "m":
				if m.DialogTarget != "" {
					if m.DialogTarget == m.CurrentBranch {
						m.Message = "⚠ Cannot merge a branch into itself"
						m.MessageType = "warning"
						return m, nil
					}

					// Show clear message about what will happen
					m.Message = fmt.Sprintf("Switching to %s and merging %s into it...", m.DialogTarget, m.CurrentBranch)
					m.MessageType = "info"

					// Switch to target branch and then merge current branch into it
					return m, performSwitchAndMerge(m.DialogTarget, m.CurrentBranch)
				} else {
					m.Message = "Select a target branch first by pressing space"
					m.MessageType = "warning"
					return m, nil
				}
			case "P", "p":
				if m.AwaitingPush {
					m.Message = "Pushing..."
					m.MessageType = "info"
					m.AwaitingPush = false
					return m, tea.Batch(performPush(), spinner())
				}
				if len(m.LogLines) > 0 {
					err := git.MergeContinue()
					if err != nil {
						m.Message = fmt.Sprintf("✗ Error continuing merge: %s", err)
						m.MessageType = "error"
						return m, nil
					}
					conflicts, _ := git.GetConflictFiles()
					if len(conflicts) > 0 {
						m.LogLines = conflicts
						m.Message = "Conflicts still present. Please resolve all conflicts."
						m.MessageType = "warning"
						return m, nil
					}
					m.Message = "✓ Merge completed successfully. Press [P] to push to remote."
					m.MessageType = "success"
					m.AwaitingPush = true
					return m, nil
				}
			case "X", "x":
				if len(m.LogLines) > 0 {
					err := git.MergeAbort()
					if err != nil {
						m.Message = fmt.Sprintf("✗ Error aborting merge: %s", err)
						m.MessageType = "error"
						return m, nil
					}
					m.Message = "Merge aborted."
					m.MessageType = "info"
					m.CurrentView = model.FileView
					m.DialogTarget = ""
					m.LogLines = nil
					m.RefreshData()
					return m, nil
				}
			case "esc":
				m.CurrentView = model.FileView
				m.DialogTarget = ""
				m.Message = ""
				m.Cursor = 0
				m.LogLines = nil
				return m, nil
			}
		}

		// --- RebaseView controls ---
		if m.CurrentView == model.RebaseView {
			switch msg.String() {
			case "up":
				if m.Cursor > 0 {
					m.Cursor--
				}
				return m, nil
			case "down":
				if m.Cursor < len(m.Branches)-1 {
					m.Cursor++
				}
				return m, nil
			case " ", "space":
				if len(m.Branches) > 0 && m.Cursor < len(m.Branches) {
					selected := m.Branches[m.Cursor]
					if selected == m.CurrentBranch {
						m.Message = "⚠ Cannot rebase a branch onto itself"
						m.MessageType = "warning"
						return m, nil
					}
					if m.DialogTarget == selected {
						m.DialogTarget = "" // unselect if already selected
						m.Message = "Selection unmarked"
						m.MessageType = "info"
					} else {
						m.DialogTarget = selected
						m.Message = fmt.Sprintf("Will rebase %s onto %s", m.CurrentBranch, selected)
						m.MessageType = "info"
					}
				}
				return m, nil
			case "R", "r":
				if m.DialogTarget != "" {
					if m.DialogTarget == m.CurrentBranch {
						m.Message = "⚠ Cannot rebase a branch onto itself"
						m.MessageType = "warning"
						return m, nil
					}

					// For rebase, we stay on current branch and rebase onto target
					m.Message = fmt.Sprintf("Rebasing %s onto %s...", m.CurrentBranch, m.DialogTarget)
					m.MessageType = "info"

					err := git.Rebase(m.DialogTarget)
					if err != nil {
						m.Message = fmt.Sprintf("✗ Error rebasing: %s", err)
						m.MessageType = "error"
						return m, nil
					}
					conflicts, _ := git.GetConflictFiles()
					if len(conflicts) > 0 {
						m.LogLines = conflicts
						m.Message = "Conflicts detected. Please resolve them and use [P] Proceed or [X] Cancel."
						m.MessageType = "warning"
						return m, nil
					} else {
						m.Message = fmt.Sprintf("✓ Rebase of %s onto %s successful", m.CurrentBranch, m.DialogTarget)
						m.MessageType = "success"
						m.CurrentView = model.FileView
						m.DialogTarget = ""
						m.LogLines = nil
						m.RefreshData()
						return m, nil
					}
				} else {
					m.Message = "Select a target branch first by pressing space"
					m.MessageType = "warning"
					return m, nil
				}
			case "P", "p":
				if len(m.LogLines) > 0 {
					err := git.RebaseContinue()
					if err != nil {
						m.Message = fmt.Sprintf("✗ Error continuing rebase: %s", err)
						m.MessageType = "error"
						return m, nil
					}
					conflicts, _ := git.GetConflictFiles()
					if len(m.LogLines) > 0 {
						m.LogLines = conflicts
						m.Message = "Conflicts still present. Please resolve all conflicts."
						m.MessageType = "warning"
						return m, nil
					}
					m.Message = "✓ Rebase completed successfully"
					m.MessageType = "success"
					m.CurrentView = model.FileView
					m.DialogTarget = ""
					m.LogLines = nil
					m.RefreshData()
					return m, nil
				}
			case "X", "x":
				if len(m.LogLines) > 0 {
					err := git.RebaseAbort()
					if err != nil {
						m.Message = fmt.Sprintf("✗ Error aborting rebase: %s", err)
						m.MessageType = "error"
						return m, nil
					}
					m.Message = "Rebase aborted."
					m.MessageType = "info"
					m.CurrentView = model.FileView
					m.DialogTarget = ""
					m.LogLines = nil
					m.RefreshData()
					return m, nil
				}
			case "esc":
				m.CurrentView = model.FileView
				m.DialogTarget = ""
				m.Message = ""
				m.Cursor = 0
				m.LogLines = nil
				return m, nil
			}
		}

		if m.CurrentView == model.ConfirmDialog {
			switch msg.String() {
			case "y":
				switch m.DialogType {
				case "delete_branch":
					if err := git.DeleteBranch(m.DialogTarget); err != nil {
						m.Message = fmt.Sprintf("✗ Error deleting branch: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = "✓ Branch deleted successfully"
						m.MessageType = "success"
						m.RefreshData()
					}
				case "discard_changes":
					if err := git.DiscardChanges(m.DialogTarget); err != nil {
						m.Message = fmt.Sprintf("✗ Error discarding changes: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = "✓ Changes discarded"
						m.MessageType = "success"
						m.RefreshData()
					}
				}
				m.CurrentView = model.FileView
				return m, nil

			case "n", "esc":
				m.CurrentView = model.FileView
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "A":
			if m.CurrentView == model.FileView && !m.AdvancedMode {
				m.AdvancedMode = true
				return m, nil
			}
		// --- ADVANCED MODE: Merge and Rebase (IMPROVED) ---
		case "M":
			if m.CurrentView == model.FileView && m.AdvancedMode {
				m.CurrentView = model.MergeView
				m.Cursor = 0
				m.DialogTarget = ""
				m.LogLines = nil
				m.Message = fmt.Sprintf("Current branch: %s - Select target branch to merge INTO", m.CurrentBranch)
				m.MessageType = "info"
				// Ensure branches are loaded
				if len(m.Branches) == 0 {
					m.RefreshData()
				}
				return m, nil
			}
		case "R":
			if m.CurrentView == model.FileView && m.AdvancedMode {
				m.CurrentView = model.RebaseView
				m.Cursor = 0
				m.DialogTarget = ""
				m.Message = fmt.Sprintf("Current branch: %s - Select base branch to rebase ONTO", m.CurrentBranch)
				m.MessageType = "info"
				m.LogLines = nil
				// Ensure branches are loaded
				if len(m.Branches) == 0 {
					m.RefreshData()
				}
				return m, nil
			}
		case "a":
			if m.CurrentView == model.FileView && !m.AdvancedMode {
				for i := range m.Files {
					if !m.Files[i].Staged {
						m.Files[i].Staged = true
						git.Add(m.Files[i].Name)
					}
				}
				m.Message = "✓ All files added to stage"
				m.MessageType = "success"
				return m, nil
			}

		case "esc":
			if m.AdvancedMode {
				m.AdvancedMode = false
				if m.CurrentView == model.LogGraphView {
					m.CurrentView = model.FileView
				}
				return m, nil
			}
			if m.CurrentView != model.FileView {
				m.CurrentView = model.FileView
				m.CommitMsg = ""
				m.RemoteName = ""
				m.RemoteURL = ""
				m.NewBranchName = ""
				m.Message = ""
				m.MessageType = ""
			}
			return m, nil

		case "q":
			if m.CurrentView == model.FileView || m.CurrentView == model.BranchView || m.CurrentView == model.RemoteView || m.CurrentView == model.ConfirmDialog || m.CurrentView == model.HelpView {
				return m, tea.Quit
			}

		case "enter":
			switch m.CurrentView {
			case model.CommitView:
				if m.CommitMsg != "" {
					if err := git.Commit(m.CommitMsg); err != nil {
						m.Message = fmt.Sprintf("✗ Error committing: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = "✓ Changes committed successfully"
						m.MessageType = "success"
						m.CurrentView = model.FileView
						m.CommitMsg = ""
						m.RefreshData()
					}
				}
				return m, nil

			case model.NewBranchView:
				if m.NewBranchName != "" {
					if err := git.CreateBranch(m.NewBranchName); err != nil {
						m.Message = fmt.Sprintf("✗ Error creating branch: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = fmt.Sprintf("✓ Branch %s created successfully", m.NewBranchName)
						m.MessageType = "success"
						m.CurrentView = model.BranchView
						m.NewBranchName = ""
						m.RefreshData()
					}
				}
				return m, nil

			case model.BranchView:
				if len(m.Branches) > 0 {
					if m.Cursor >= len(m.Branches) {
						m.Cursor = len(m.Branches) - 1
					}
					selected := m.Branches[m.Cursor]
					if selected != m.CurrentBranch {
						if err := git.Checkout(selected); err != nil {
							m.Message = fmt.Sprintf("✗ Error switching to branch %s: %s", selected, err)
							m.MessageType = "error"
						} else {
							m.Message = fmt.Sprintf("✓ Switched to branch %s", selected)
							m.MessageType = "success"
							m.CurrentBranch = selected
							m.RefreshData()
						}
					} else {
						m.Message = "⚠ You are already on this branch"
						m.MessageType = "info"
					}
				}
				return m, nil
			}

		case "backspace":
			switch m.CurrentView {
			case model.CommitView:
				if len(m.CommitMsg) > 0 {
					m.CommitMsg = m.CommitMsg[:len(m.CommitMsg)-1]
				}
			case model.NewBranchView:
				if len(m.NewBranchName) > 0 {
					m.NewBranchName = m.NewBranchName[:len(m.NewBranchName)-1]
				}
			}
			return m, nil

		case "up":
			if m.CurrentView != model.CommitView && m.CurrentView != model.NewBranchView && m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil

		case "down":
			if m.CurrentView != model.CommitView && m.CurrentView != model.NewBranchView {
				switch m.CurrentView {
				case model.FileView:
					if m.Cursor < len(m.Files)-1 {
						m.Cursor++
					}
				case model.BranchView:
					if m.Cursor < len(m.Branches)-1 {
						m.Cursor++
					}
				case model.RemoteView:
					if m.Cursor < len(m.Remotes)-1 {
						m.Cursor++
					}
				case model.LogGraphView:
				}
			}
			return m, nil
		}

		if m.CurrentView == model.CommitView {
			if len(msg.Runes) == 1 && isPrintableChar(msg.Runes[0]) {
				m.CommitMsg += string(msg.Runes)
				return m, nil
			}
		}

		if m.CurrentView == model.NewBranchView {
			if len(msg.Runes) == 1 && isPrintableChar(msg.Runes[0]) {
				m.NewBranchName += string(msg.Runes)
				return m, nil
			}
		}

		if m.CurrentView == model.BranchView {
			switch msg.String() {
			case "n":
				m.CurrentView = model.NewBranchView
				m.NewBranchName = ""
				return m, nil
			case "d":
				if len(m.Branches) > 0 && m.Cursor < len(m.Branches) {
					toDel := m.Branches[m.Cursor]
					if toDel == m.CurrentBranch {
						m.Message = "✗ Cannot delete current branch"
						m.MessageType = "error"
					} else {
						m.DialogType = "delete_branch"
						m.DialogTarget = toDel
						m.CurrentView = model.ConfirmDialog
					}
				}
				return m, nil
			}
		}

		if m.CurrentView == model.FileView {
			switch msg.String() {
			case " ":
				if len(m.Files) > 0 {
					f := &m.Files[m.Cursor]
					f.Staged = !f.Staged
					if f.Staged {
						git.Add(f.Name)
						m.Message = fmt.Sprintf("✓ File %s added to stage", f.Name)
					} else {
						git.Reset(f.Name)
						m.Message = fmt.Sprintf("✓ File %s removed from stage", f.Name)
					}
					m.MessageType = "success"
				}
			case "a":
				for i := range m.Files {
					if !m.Files[i].Staged {
						m.Files[i].Staged = true
						git.Add(m.Files[i].Name)
					}
				}
				m.Message = "✓ All files added to stage"
				m.MessageType = "success"
			case "c":
				ok := false
				for _, f := range m.Files {
					if f.Staged {
						ok = true
						break
					}
				}
				if ok {
					m.CurrentView = model.CommitView
					m.Message = ""
				} else {
					m.Message = "⚠ No staged files to commit"
					m.MessageType = "error"
				}
			case "b":
				m.CurrentView = model.BranchView
				m.Cursor = 0
				m.Message = ""
			case "m":
				m.CurrentView = model.RemoteView
				m.Cursor = 0
				m.Message = ""
			case "r":
				m.RefreshData()
				m.Message = "✓ Status updated"
				m.MessageType = "success"
			case "p":
				if !m.IsPushing {
					m.IsPushing = true
					m.Message = "Pushing..."
					m.MessageType = "info"
					return m, tea.Batch(performPush(), spinner())
				}
			case "f":
				if !m.IsFetching {
					m.IsFetching = true
					m.Message = "Fetching..."
					m.MessageType = "info"
					return m, tea.Batch(performFetch(), spinner())
				}
			case "l":
				if !m.IsPulling {
					m.IsPulling = true
					m.Message = "Pulling..."
					m.MessageType = "info"
					return m, tea.Batch(performPull(), spinner())
				}
			case "L":
				// Open interactive Git log graph view
				m, cmd := OpenLogGraphView(m)
				return m, cmd

			case "A":
				m.MessageType = "info"
				return m, nil
			case "x":
				if len(m.Files) > 0 {
					m.DialogType = "discard_changes"
					m.DialogTarget = m.Files[m.Cursor].Name
					m.CurrentView = model.ConfirmDialog
				}
			case "?":
				if m.CurrentView == model.FileView {
					m.CurrentView = model.HelpView
					return m, nil
				} else if m.CurrentView == model.HelpView {
					m.CurrentView = model.FileView
					return m, nil
				}
			}
		}

	case spinnerTickMsg:
		if m.IsPushing || m.IsFetching || m.IsPulling {
			m.SpinnerIndex = (m.SpinnerIndex + 1) % len(m.SpinnerFrames)
			return m, spinner()
		}

	case pushMsg:
		m.IsPushing = false
		if msg.err != nil {
			m.Message = fmt.Sprintf("✗ Error pushing changes: %s", msg.err)
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes pushed to remote successfully"
			m.MessageType = "success"
			m.CurrentView = model.FileView
			m.DialogTarget = ""
			m.LogLines = nil
			m.RefreshData()
		}

	case fetchMsg:
		m.IsFetching = false
		if msg.err != nil {
			m.Message = fmt.Sprintf("✗ Error fetching changes: %s", msg.err)
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes fetched successfully"
			m.MessageType = "success"
			m.RefreshData()
		}

	case pullMsg:
		m.IsPulling = false
		if msg.err != nil {
			m.Message = fmt.Sprintf("✗ Error pulling changes: %s", msg.err)
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes pulled successfully"
			m.MessageType = "success"
			m.RefreshData()
		}
	}

	return m, nil
}

// isPrintableChar checks whether a rune is a printable ASCII or extended character.
func isPrintableChar(r rune) bool {
	return (r >= 32 && r <= 126) || (r >= 128 && r <= 255) // Extended ASCII range
}
