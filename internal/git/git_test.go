package git

import "testing"

func TestGetStashRef(t *testing.T) {
	cases := map[string]string{
		"stash@{0}: WIP on main": "stash@{0}",
		"stash@{1}: On feature":  "stash@{1}",
		"":                       "stash@{0}",
	}
	for in, want := range cases {
		got := GetStashRef(in)
		if got != want {
			t.Fatalf("GetStashRef(%q) = %q; want %q", in, got, want)
		}
	}
}
