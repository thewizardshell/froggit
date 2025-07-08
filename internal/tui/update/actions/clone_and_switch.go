package actions

import (
	"fmt"

	"froggit/internal/gh"
	"froggit/internal/git"
	"froggit/internal/tui/model"
	"froggit/internal/tui/update/messages"
)

var ghClient = gh.NewGhClient()

func HandleConfirmCloneRepo(m model.Model, key string) (model.Model, bool) {
	switch key {
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
		return m, true
	case "n", "esc":
		m.CurrentView = model.RepositoryListView
		m.RepoToClone = nil
		return m, true
	}
	return m, false
}

func HandleSwitchBranchMsg(m model.Model, msg any) (model.Model, bool) {
	switch msg := msg.(type) {
	case messages.SwitchBranchMsg:
		if msg.Err != nil {
			m.Message = fmt.Sprintf("✗ Error switching to branch %s: %s", msg.TargetBranch, msg.Err)
			m.MessageType = "error"
			return m, true
		}

		m.CurrentBranch = msg.TargetBranch
		m.RefreshData()

		if msg.NextAction == "merge" {
			err := git.Merge(msg.SourceBranch)
			if err != nil {
				m.Message = fmt.Sprintf("✗ Error merging %s into %s: %s", msg.SourceBranch, msg.TargetBranch, err)
				m.MessageType = "error"
				return m, true
			}

			conflicts, _ := git.GetConflictFiles()
			if len(conflicts) > 0 {
				m.LogLines = conflicts
				m.Message = fmt.Sprintf("Conflicts detected while merging %s into %s. Please resolve them and use [P] Proceed or [X] Cancel.", msg.SourceBranch, msg.TargetBranch)
				m.MessageType = "warning"
			} else {
				m.Message = fmt.Sprintf("✓ Successfully merged %s into %s. Press [P] to push to remote.", msg.SourceBranch, msg.TargetBranch)
				m.MessageType = "success"
				m.CurrentView = model.MergeView
				m.AwaitingPush = true
			}
			return m, true
		} else if msg.NextAction == "rebase" {
			err := git.Rebase(msg.TargetBranch)
			if err != nil {
				m.Message = fmt.Sprintf("✗ Error rebasing %s onto %s: %s", msg.SourceBranch, msg.TargetBranch, err)
				m.MessageType = "error"
				return m, true
			}

			conflicts, _ := git.GetConflictFiles()
			if len(conflicts) > 0 {
				m.LogLines = conflicts
				m.Message = fmt.Sprintf("Conflicts detected while rebasing %s onto %s. Please resolve them and use [P] Proceed or [X] Cancel.", msg.SourceBranch, msg.TargetBranch)
				m.MessageType = "warning"
			} else {
				m.Message = fmt.Sprintf("✓ Successfully rebased %s onto %s", msg.SourceBranch, msg.TargetBranch)
				m.MessageType = "success"
				m.CurrentView = model.FileView
				m.DialogTarget = ""
				m.LogLines = nil
				m.RefreshData()
			}
			return m, true
		}
	}
	return m, false
}
