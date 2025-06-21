package update

import (
	"fmt"
	"strings"

	"froggit/internal/git"
	"froggit/internal/tui/model"

	tea "github.com/charmbracelet/bubbletea"
)

// HandleLogGraphKey handles key messages when in the LogGraphView
func HandleLogGraphKey(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.CurrentView = model.FileView
		return m, nil

	case "up":
		if m.Cursor > 0 {
			m.Cursor--
		}
		return m, nil

	case "down":
		if m.Cursor < len(m.LogLines)-1 {
			m.Cursor++
		}
		return m, nil

	default:
		return m, nil
	}
}

// OpenLogGraphView loads the git log graph and updates the model to enter LogGraphView
func OpenLogGraphView(m model.Model) (model.Model, tea.Cmd) {
	graph, err := git.LogsGraph()
	if err != nil {
		m.Message = fmt.Sprintf("✗ Error retrieving log graph: %s", err)
		m.MessageType = "error"
		return m, nil
	}

	m.LogLines = strings.Split(strings.TrimSpace(graph), "\n")
	m.Cursor = 0
	m.CurrentView = model.LogGraphView
	m.Message = ""

	return m, nil
}
