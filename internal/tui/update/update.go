package update

import (
	"fmt"
	"froggit/internal/config"
	"froggit/internal/git"
	"froggit/internal/tui/model"
	"froggit/internal/tui/update/actions"
	"froggit/internal/tui/update/async"
	"froggit/internal/tui/update/handlers"
	"froggit/internal/tui/update/messages"
	"froggit/internal/utils"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func Update(m model.Model, cfg config.Config, msg tea.Msg) (model.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if cfg.Git.AutoFetch && m.CurrentView == model.FileView && !m.IsFetching && !m.AutoFetchDone {
		m.AutoFetchDone = true
		m.IsFetching = true
		cmds = append(cmds, async.PerformAutoFetch(), async.Spinner())
	}

	if !cfg.Git.AutoFetch && m.CurrentView == model.FileView && m.CurrentBranch != "" && !m.HasRemoteChanges {
		cmds = append(cmds, async.PerformRemoteChangesCheck(m.CurrentBranch))
	}

	if m.CurrentView == model.ConfirmCloneRepoView {
		if key, ok := msg.(tea.KeyMsg); ok {
			var handled bool
			m, handled = actions.HandleConfirmCloneRepo(m, key.String())
			if handled {
				return m, tea.Batch(cmds...)
			}
		}
	}

	switch msg := msg.(type) {

	case messages.SwitchBranchMsg:
		if msg.Err != nil {
			m.Message = fmt.Sprintf("✗ Error switching to branch %s: %s", msg.TargetBranch, msg.Err)
			m.MessageType = "error"
			return m, nil
		}

		m.CurrentBranch = msg.TargetBranch
		m.RefreshData()
		if msg.NextAction == "merge" {
			err := git.Merge(msg.SourceBranch)

			conflicts, _ := git.GetConflictFiles()

			if len(conflicts) > 0 {
				m.LogLines = conflicts
				m.Message = fmt.Sprintf("Conflicts detected while merging %s into %s. Please resolve them and use [P] Proceed or [X] Cancel.", msg.SourceBranch, msg.TargetBranch)
				m.MessageType = "warning"
				return m, nil
			} else if err != nil {
				m.Message = fmt.Sprintf("✗ Error merging %s into %s: %s", msg.SourceBranch, msg.TargetBranch, err)
				m.MessageType = "error"
				return m, nil
			} else {
				m.Message = fmt.Sprintf("✓ Successfully merged %s into %s. Press [P] to push to remote.", msg.SourceBranch, msg.TargetBranch)
				m.MessageType = "success"
				m.CurrentView = model.MergeView
				m.AwaitingPush = true
			}
			return m, nil
		} else if msg.NextAction == "rebase" {
			err := git.Rebase(msg.TargetBranch)

			conflicts, _ := git.GetConflictFiles()

			if len(conflicts) > 0 {
				m.LogLines = conflicts
				m.Message = fmt.Sprintf("Conflicts detected while rebasing %s onto %s. Please resolve them and use [P] Proceed or [X] Cancel.", msg.SourceBranch, msg.TargetBranch)
				m.MessageType = "warning"
				return m, nil
			} else if err != nil {
				m.Message = fmt.Sprintf("✗ Error rebasing %s onto %s: %s", msg.SourceBranch, msg.TargetBranch, err)
				m.MessageType = "error"
				return m, nil
			} else {
				m.Message = fmt.Sprintf("✓ Successfully rebased %s onto %s", msg.SourceBranch, msg.TargetBranch)
				m.MessageType = "success"
				m.CurrentView = model.FileView
				m.DialogTarget = ""
				m.LogLines = nil
				m.RefreshData()
			}
			return m, nil
		}

	case tea.KeyMsg:

		if m.CurrentView == model.QuickStartView {
			switch msg.String() {
			case "up":
				if m.Cursor > 0 {
					m.Cursor--
				}
				return m, nil
			case "down":
				maxOptions := 3
				if m.Cursor < maxOptions-1 {
					m.Cursor++
				}
				return m, nil
			case "enter":
				if m.Cursor == 0 || (m.Cursor == 1 && m.HasGitHubCLI) || (m.Cursor == 2 && m.HasGitHubCLI) {
					return m, tea.Quit
				}
				return m, nil
			case "1":
				m.Cursor = 0
				return m, tea.Quit
			case "2":
				if m.HasGitHubCLI {
					m.Cursor = 1
					return m, tea.Quit
				}
				return m, nil
			case "3":
				if m.HasGitHubCLI {
					m.Cursor = 2
					return m, tea.Quit
				}
				return m, nil
			case "q":
				return m, tea.Quit
			}
		}

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

		if m.CurrentView == model.MergeView {
			return handlers.HandleMergeView(m, msg)
		}

		if m.CurrentView == model.RebaseView {
			return handlers.HandleRebaseView(m, msg)
		}

		if m.CurrentView == model.StashView {
			return handlers.HandleStashView(m, msg)
		}

		if m.CurrentView == model.StashMessageView {
			return handlers.HandleStashMessageView(m, msg)
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
				case "delete_remote":
					if err := git.RemoveRemote(m.DialogTarget); err != nil {
						m.Message = fmt.Sprintf("✗ Error deleting remote: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = "✓ Remote deleted successfully"
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
				case "drop_stash":
					if err := git.StashDrop(m.DialogTarget); err != nil {
						m.Message = fmt.Sprintf("✗ Error dropping stash: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = "✓ Stash dropped successfully"
						m.MessageType = "success"
						m.RefreshData()
					}
					m.CurrentView = model.StashView
					return m, nil
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
		case "M":
			if m.CurrentView == model.FileView && m.AdvancedMode {
				m.CurrentView = model.MergeView
				m.Cursor = 0
				m.DialogTarget = ""
				m.LogLines = nil
				m.Message = fmt.Sprintf("Current branch: %s - Select target branch to merge INTO", m.CurrentBranch)
				m.MessageType = "info"
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
				if len(m.Branches) == 0 {
					m.RefreshData()
				}
				return m, nil
			}
		case "S":
			if m.CurrentView == model.FileView && m.AdvancedMode {
				m.CurrentView = model.StashView
				m.Cursor = 0
				m.Message = ""
				m.MessageType = ""
				m.RefreshData()
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

		case "tab":
			if m.CurrentView == model.AddRemoteView {
				if m.InputField == "name" {
					m.InputField = "url"
				} else if m.InputField == "url" {
					m.InputField = "name"
				} else {
					m.InputField = "name"
				}
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

			case model.AddRemoteView:
				if m.InputField == "name" && m.RemoteName != "" {
					m.InputField = "url"
				} else if m.InputField == "url" && m.RemoteURL != "" {
					if err := git.AddRemote(m.RemoteName, m.RemoteURL); err != nil {
						m.Message = fmt.Sprintf("✗ Error adding remote: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = fmt.Sprintf("✓ Remote %s added successfully", m.RemoteName)
						m.MessageType = "success"
						m.CurrentView = model.RemoteView
						m.RemoteName = ""
						m.RemoteURL = ""
						m.InputField = ""
						m.RefreshData()
					}
				} else if m.InputField == "" {
					m.InputField = "name"
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
			case model.AddRemoteView:
				if m.InputField == "name" && len(m.RemoteName) > 0 {
					m.RemoteName = m.RemoteName[:len(m.RemoteName)-1]
				} else if m.InputField == "url" && len(m.RemoteURL) > 0 {
					m.RemoteURL = m.RemoteURL[:len(m.RemoteURL)-1]
				}
			}
			return m, nil

		case "up":
			if m.CurrentView != model.CommitView && m.CurrentView != model.NewBranchView && m.CurrentView != model.AddRemoteView && m.Cursor > 0 {
				m.Cursor--
				if m.CurrentView == model.FileView && m.Cursor < m.FileViewOffset {
					m.FileViewOffset = m.Cursor
				}
			}
			return m, nil

		case "down":
			if m.CurrentView != model.CommitView && m.CurrentView != model.NewBranchView && m.CurrentView != model.AddRemoteView {
				switch m.CurrentView {
				case model.FileView:
					if m.Cursor < len(m.Files)-1 {
						m.Cursor++
						if m.Cursor >= m.FileViewOffset+m.FileViewHeight {
							m.FileViewOffset = m.Cursor - m.FileViewHeight + 1
						}
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
			switch msg.String() {
			case "tab":
				if !m.IsGeneratingAI && m.CopilotAvailable {
					m.IsGeneratingAI = true
					m.Message = "Generating AI commit message..."
					m.MessageType = "info"
					return m, tea.Batch(async.PerformAICommitGeneration(), async.Spinner())
				}
				return m, nil
			}
			if len(msg.Runes) == 1 && utils.IsPrintableChar(msg.Runes[0]) {
				m.CommitMsg += string(msg.Runes)
				return m, nil
			}
		}

		if m.CurrentView == model.NewBranchView {
			if len(msg.Runes) == 1 && utils.IsPrintableChar(msg.Runes[0]) {
				m.NewBranchName += string(msg.Runes)
				return m, nil
			}
		}

		if m.CurrentView == model.AddRemoteView {
			if len(msg.Runes) == 1 && utils.IsPrintableChar(msg.Runes[0]) {
				if m.InputField == "name" {
					m.RemoteName += string(msg.Runes)
				} else if m.InputField == "url" {
					m.RemoteURL += string(msg.Runes)
				}
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

		if m.CurrentView == model.RemoteView {
			switch msg.String() {
			case "n":
				m.CurrentView = model.AddRemoteView
				m.RemoteName = ""
				m.RemoteURL = ""
				m.InputField = "name"
				return m, nil
			case "d":
				if len(m.Remotes) > 0 && m.Cursor < len(m.Remotes) {
					toDel := m.Remotes[m.Cursor]
					remoteName := strings.Split(toDel, " -> ")[0]
					m.DialogType = "delete_remote"
					m.DialogTarget = remoteName
					m.CurrentView = model.ConfirmDialog
				}
				return m, nil
			}
		}

		if m.CurrentView == model.FileView {
			switch msg.String() {
			case " ":
				if len(m.Files) > 0 && m.Cursor < len(m.Files) {
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
					// Refresh files to update status (?, A, M, etc.)
					files, _ := git.GetModifiedFiles()
					m.Files = files
					// Restore cursor position
					if m.Cursor >= len(m.Files) && len(m.Files) > 0 {
						m.Cursor = len(m.Files) - 1
					}
				} else {
					if len(m.Files) > 0 {
						m.Cursor = 0
					} else {
						m.Message = "⚠ No files to stage"
						m.MessageType = "warning"
					}
				}
			case "a":
				if len(m.Files) > 0 {
					for i := range m.Files {
						if !m.Files[i].Staged {
							m.Files[i].Staged = true
							git.Add(m.Files[i].Name)
						}
					}
					m.Message = "✓ All files added to stage"
					m.MessageType = "success"
					// Refresh files to update status (?, A, M, etc.)
					files, _ := git.GetModifiedFiles()
					m.Files = files
					if m.Cursor >= len(m.Files) && len(m.Files) > 0 {
						m.Cursor = len(m.Files) - 1
					}
				} else {
					m.Message = "⚠ No files to stage"
					m.MessageType = "warning"
				}
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
				utils.ValidateCursor(&m)
				m.Message = "✓ Status updated"
				m.MessageType = "success"
			case "p":
				hasCommits, err := git.HasCommitsToush()
				if err != nil {
					m.Message = fmt.Sprintf("✗ Error checking commits: %s", err)
					m.MessageType = "error"
					return m, nil
				}

				if !hasCommits {
					m.Message = "⚠ No commits to push"
					m.MessageType = "error"
					return m, nil
				}

				if !m.IsPushing {
					m.IsPushing = true
					m.Message = "Pushing..."
					m.MessageType = "info"
					return m, tea.Batch(async.PerformPushWithConfig(cfg.Git.DefaultBranch), async.Spinner())
				}
			case "f":
				if !m.IsFetching {
					m.IsFetching = true
					m.Message = "Fetching..."
					m.MessageType = "info"
					return m, tea.Batch(async.PerformFetch(), async.Spinner())
				}
			case "l":
				if !m.IsPulling {
					m.IsPulling = true
					m.Message = "Pulling..."
					m.MessageType = "info"
					return m, tea.Batch(async.PerformPull(), async.Spinner())
				}
			case "L":
				m, cmd := OpenLogGraphView(m)
				return m, cmd

			case "A":
				m.MessageType = "info"
				return m, nil
			case "x":
				if len(m.Files) > 0 && m.Cursor < len(m.Files) {
					m.DialogType = "discard_changes"
					m.DialogTarget = m.Files[m.Cursor].Name
					m.CurrentView = model.ConfirmDialog
				} else {
					m.Message = "⚠ No file selected or no files available"
					m.MessageType = "warning"
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

	case async.SpinnerTickMsg:
		if m.IsPushing || m.IsFetching || m.IsPulling || m.IsGeneratingAI {
			m.SpinnerIndex = (m.SpinnerIndex + 1) % len(m.SpinnerFrames)
			return m, async.Spinner()
		}
		return m, nil

	case async.AICommitMsg:
		m.IsGeneratingAI = false
		if msg.Err != nil {
			m.Message = fmt.Sprintf("✗ AI error: %s", msg.Err)
			m.MessageType = "error"
		} else {
			m.CommitMsg = msg.Message
			m.Message = "✓ AI generated commit message"
			m.MessageType = "success"
		}
		return m, nil

	case async.PushMsg:
		m.IsPushing = false
		if msg.Err != nil {
			m.Message = fmt.Sprintf("✗ Error pushing changes: %s", msg.Err)
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes pushed to remote successfully"
			m.MessageType = "success"
			m.CurrentView = model.FileView
			m.DialogTarget = ""
			m.LogLines = nil
			m.RefreshData()
			utils.ValidateCursor(&m)
		}

	case async.FetchMsg:
		m.IsFetching = false
		if msg.Err != nil {
			m.Message = fmt.Sprintf("✗ Error fetching changes: %s", msg.Err)
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes fetched successfully"
			m.MessageType = "success"
			m.RefreshData()
			utils.ValidateCursor(&m)
		}

	case async.PullMsg:
		m.IsPulling = false
		if msg.Err != nil {
			m.Message = fmt.Sprintf("✗ Error pulling changes: %s", msg.Err)
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes pulled successfully"
			m.MessageType = "success"
			m.RefreshData()
			utils.ValidateCursor(&m)
		}

	case async.RemoteChangesCheckMsg:
		if msg.Err == nil {
			m.HasRemoteChanges = msg.HasChanges
		}
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}
