package update

import (
	"froggit/internal/gh"
	"froggit/internal/tui/model"
)

// ShowRepositoryList loads the user's repositories and updates the model to show the list.
func ShowRepositoryList(m model.Model, ghClient *gh.GhClient) model.Model {
	repos, err := gh.ListUserRepositories(ghClient)
	if err != nil {
		m.Message = "Error loading repositories: " + err.Error()
		m.MessageType = "error"
		return m
	}
	m.Repositories = repos
	m.SelectedRepoIndex = 0
	m.CurrentView = model.RepositoryListView
	m.Message = ""
	return m
}
