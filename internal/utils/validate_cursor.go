// validateCursor ensures the cursor is within valid bounds for the current view
package utils

import "froggit/internal/tui/model"

func ValidateCursor(m *model.Model) {
	switch m.CurrentView {
	case model.FileView:
		if len(m.Files) == 0 {
			m.Cursor = 0
		} else if m.Cursor >= len(m.Files) {
			m.Cursor = len(m.Files) - 1
		} else if m.Cursor < 0 {
			m.Cursor = 0
		}
	case model.BranchView:
		if len(m.Branches) == 0 {
			m.Cursor = 0
		} else if m.Cursor >= len(m.Branches) {
			m.Cursor = len(m.Branches) - 1
		} else if m.Cursor < 0 {
			m.Cursor = 0
		}
	case model.RemoteView:
		if len(m.Remotes) == 0 {
			m.Cursor = 0
		} else if m.Cursor >= len(m.Remotes) {
			m.Cursor = len(m.Remotes) - 1
		} else if m.Cursor < 0 {
			m.Cursor = 0
		}
	}
}
