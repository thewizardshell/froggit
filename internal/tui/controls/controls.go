package controls

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type Control struct {
	Key         string
	Description string
	Group       string
}

type ControlSet struct {
	controls []Control
	width    int
}

func NewControlSet() *ControlSet {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 80
	}
	return &ControlSet{
		controls: []Control{},
		width:    width,
	}
}

func (cs *ControlSet) Add(key, description, group string) *ControlSet {
	cs.controls = append(cs.controls, Control{
		Key:         key,
		Description: description,
		Group:       group,
	})
	return cs
}

func (cs *ControlSet) AddMultiple(controls []Control) *ControlSet {
	cs.controls = append(cs.controls, controls...)
	return cs
}

func (cs *ControlSet) Render() string {
	if len(cs.controls) == 0 {
		return ""
	}

	return cs.renderSimple()
}

func (cs *ControlSet) renderSimple() string {
	var controlParts []string

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("76")).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("248"))

	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	for _, control := range cs.controls {
		part := keyStyle.Render(control.Key) + " " + descStyle.Render(control.Description)
		controlParts = append(controlParts, part)
	}

	content := strings.Join(controlParts, separatorStyle.Render(" │ "))

	if lipgloss.Width(content) > cs.width - 4 {
		var lines []string
		currentLine := ""
		currentWidth := 0
		maxWidth := cs.width - 4

		for i, part := range controlParts {
			partWidth := lipgloss.Width(part)
			separatorWidth := 3 // " │ "

			totalWidth := currentWidth
			if currentWidth > 0 {
				totalWidth += separatorWidth
			}
			totalWidth += partWidth

			if totalWidth <= maxWidth {
				if currentLine != "" {
					currentLine += separatorStyle.Render(" │ ") + part
				} else {
					currentLine = part
				}
				currentWidth = totalWidth
			} else {
				if currentLine != "" {
					lines = append(lines, currentLine)
				}
				currentLine = part
				currentWidth = partWidth
			}

			if i == len(controlParts)-1 && currentLine != "" {
				lines = append(lines, currentLine)
			}
		}

		content = strings.Join(lines, "\n")
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("76")).   // Green border
		Foreground(lipgloss.Color("250")).
		Padding(0, 1).
		MarginTop(1)

	return boxStyle.Render(content)
}

func NewFileViewControls(staged bool, hasFiles bool, advancedMode bool) *ControlSet {
	cs := NewControlSet()

	cs.Add("↑/↓", "navigate", "navigation")

	if !advancedMode {
		if hasFiles {
			cs.Add("space", "stage/unstage", "files")
			cs.Add("x", "discard changes", "files")
		}
		if staged {
			cs.Add("c", "commit", "files")
		}
		cs.Add("a", "stage all", "files")
		cs.Add("r", "refresh", "files")
		cs.Add("f", "fetch", "git")
		cs.Add("l", "pull", "git")
		cs.Add("p", "push", "git")
		cs.Add("b", "branches", "nav")
		cs.Add("m", "remotes", "nav")
		cs.Add("?", "help", "general")
	} else {
		cs.Add("L", "log graph", "advanced")
		cs.Add("M", "merge", "advanced")
		cs.Add("R", "rebase", "advanced")
		cs.Add("S", "stash", "advanced")
		cs.Add("esc", "exit advanced", "mode")
		cs.Add("?", "help", "general")
	}

	return cs
}

func NewBranchViewControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("↑/↓", "navigate", "navigation")
	cs.Add("enter", "switch branch", "actions")
	cs.Add("n", "new branch", "actions")
	cs.Add("d", "delete branch", "actions")
	cs.Add("esc", "back", "navigation")
	cs.Add("q", "quit", "general")
	return cs
}

func NewRemoteViewControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("↑/↓", "navigate", "navigation")
	cs.Add("n", "add remote", "actions")
	cs.Add("d", "delete remote", "actions")
	cs.Add("esc", "back", "navigation")
	cs.Add("q", "quit", "general")
	return cs
}

func NewCommitViewControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("enter", "commit changes", "actions")
	cs.Add("backspace", "delete char", "edit")
	cs.Add("esc", "cancel", "navigation")
	return cs
}

func NewMergeViewControls(hasSelection bool) *ControlSet {
	cs := NewControlSet()
	cs.Add("↑/↓", "navigate", "navigation")
	cs.Add("enter", "select branch", "actions")
	if hasSelection {
		cs.Add("M", "merge", "actions")
	}
	cs.Add("esc", "cancel", "navigation")
	return cs
}

func NewRebaseViewControls(hasSelection bool) *ControlSet {
	cs := NewControlSet()
	cs.Add("↑/↓", "navigate", "navigation")
	cs.Add("enter", "select branch", "actions")
	if hasSelection {
		cs.Add("R", "rebase", "actions")
	}
	cs.Add("esc", "cancel", "navigation")
	return cs
}

func NewStashViewControls(hasChanges bool, hasStashes bool) *ControlSet {
	cs := NewControlSet()

	if hasStashes {
		cs.Add("↑/↓", "navigate", "navigation")
		cs.Add("enter", "apply stash", "actions")
		cs.Add("p", "pop stash", "actions")
		cs.Add("d", "drop stash", "actions")
		cs.Add("v", "view stash", "actions")
	}

	if hasChanges {
		cs.Add("s", "save stash", "actions")
	}

	cs.Add("esc", "back", "navigation")
	cs.Add("?", "help", "general")
	return cs
}

func NewStashMessageViewControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("enter", "save stash", "actions")
	cs.Add("backspace", "delete char", "edit")
	cs.Add("esc", "cancel", "navigation")
	return cs
}

func NewNewBranchViewControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("enter", "create branch", "actions")
	cs.Add("backspace", "delete char", "edit")
	cs.Add("esc", "cancel", "navigation")
	return cs
}

func NewAddRemoteViewControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("tab", "switch field", "navigation")
	cs.Add("enter", "confirm/next", "actions")
	cs.Add("backspace", "delete char", "edit")
	cs.Add("esc", "cancel", "navigation")
	return cs
}

func NewGitHubControlsViewControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("r", "list repositories", "github")
	cs.Add("esc", "back", "navigation")
	return cs
}

func NewConfirmDialogControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("y", "yes", "actions")
	cs.Add("n", "no", "actions")
	cs.Add("esc", "cancel", "navigation")
	return cs
}

func NewRepositoryListControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("↑/↓", "navigate", "navigation")
	cs.Add("c", "clone repository", "actions")
	cs.Add("esc", "back", "navigation")
	return cs
}

func NewQuickStartControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("↑/↓", "navigate", "navigation")
	cs.Add("enter", "select option", "actions")
	cs.Add("1/2/3", "quick select", "actions")
	cs.Add("q", "quit", "general")
	return cs
}

func NewHelpViewControls() *ControlSet {
	cs := NewControlSet()
	cs.Add("esc", "back", "navigation")
	cs.Add("?", "close help", "navigation")
	cs.Add("q", "quit", "general")
	return cs
}