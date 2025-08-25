package utils

import (
	"testing"

	"froggit/internal/tui/model"
)

func TestIsPrintableChar(t *testing.T) {
	cases := []struct {
		r    rune
		want bool
	}{
		{31, false},
		{32, true},
		{65, true},
		{126, true},
		{127, false},
		{128, true},
		{255, true},
		{256, false},
	}

	for _, c := range cases {
		got := IsPrintableChar(c.r)
		if got != c.want {
			t.Fatalf("IsPrintableChar(%d) = %v; want %v", c.r, got, c.want)
		}
	}
}

func TestValidateCursor_FileBranchRemoteViews(t *testing.T) {
	// FileView: empty files should clamp cursor to 0
	m := &model.Model{
		CurrentView: model.FileView,
		Files:       nil,
		Cursor:      5,
	}
	ValidateCursor(m)
	if m.Cursor != 0 {
		t.Fatalf("expected cursor 0 for empty Files, got %d", m.Cursor)
	}

	// BranchView: cursor greater than len(branches)-1 should clamp
	m = &model.Model{
		CurrentView: model.BranchView,
		Branches:    []string{"a", "b", "c"},
		Cursor:      5,
	}
	ValidateCursor(m)
	if m.Cursor != 2 {
		t.Fatalf("expected cursor 2 for BranchView, got %d", m.Cursor)
	}

	// Negative cursor clamps to 0
	m = &model.Model{
		CurrentView: model.RemoteView,
		Remotes:     []string{"r1"},
		Cursor:      -3,
	}
	ValidateCursor(m)
	if m.Cursor != 0 {
		t.Fatalf("expected cursor 0 for negative cursor, got %d", m.Cursor)
	}
}
