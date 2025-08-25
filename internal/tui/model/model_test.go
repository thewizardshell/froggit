package model

import (
	"testing"
)

func TestParseStashList(t *testing.T) {
	out := "stash@{0}: WIP on main: 123abc\nstash@{1}: On feature: 456def\n"
	stashes := parseStashList(out)
	if len(stashes) != 2 {
		t.Fatalf("expected 2 stashes, got %d", len(stashes))
	}
}

func TestSetConflictFiles(t *testing.T) {
	m := &Model{}
	m.SetMergeConflictFiles([]string{"a.go", "b.go"})
	if len(m.MergeConflictFiles) != 2 {
		t.Fatalf("expected 2 merge conflict files, got %d", len(m.MergeConflictFiles))
	}
	m.SetRebaseConflictFiles([]string{"x.go"})
	if len(m.RebaseConflictFiles) != 1 {
		t.Fatalf("expected 1 rebase conflict file, got %d", len(m.RebaseConflictFiles))
	}
}
