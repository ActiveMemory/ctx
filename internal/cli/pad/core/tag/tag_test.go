//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tag

import (
	"testing"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		name  string
		entry string
		want  []string
	}{
		{"empty", "", nil},
		{"no tags", "just plain text", nil},
		{"single tag", "fix test #later", []string{"later"}},
		{"tag at start", "#urgent deploy hotfix", []string{"urgent"}},
		{"multiple tags", "fix test #later #ci", []string{"ci", "later"}},
		{"duplicate tags", "#a some text #a", []string{"a"}},
		{"hyphenated tag", "task #high-priority done", []string{"high-priority"}},
		{"numeric tag", "issue #42 is broken", []string{"42"}},
		{"underscore tag", "check #my_tag here", []string{"my_tag"}},
		{"tag at end", "some text #done", []string{"done"}},
		{"only tag", "#solo", []string{"solo"}},
		{
			"mid-word hash ignored",
			"foo#bar baz",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Extract(tt.entry)
			if !stringSliceEqual(got, tt.want) {
				t.Errorf("Extract(%q) = %v, want %v", tt.entry, got, tt.want)
			}
		})
	}
}

func TestHas(t *testing.T) {
	tests := []struct {
		name  string
		entry string
		tag   string
		want  bool
	}{
		{"present", "fix test #later", "later", true},
		{"absent", "fix test #later", "urgent", false},
		{"empty entry", "", "later", false},
		{"at start", "#urgent deploy", "urgent", true},
		{"at end", "deploy #urgent", "urgent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Has(tt.entry, tt.tag); got != tt.want {
				t.Errorf("Has(%q, %q) = %v, want %v",
					tt.entry, tt.tag, got, tt.want)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name   string
		entry  string
		filter string
		want   bool
	}{
		{"positive match", "fix #later", "later", true},
		{"positive no match", "fix #later", "urgent", false},
		{"negated match", "fix #later", "~urgent", true},
		{"negated no match", "fix #later", "~later", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Match(tt.entry, tt.filter); got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v",
					tt.entry, tt.filter, got, tt.want)
			}
		})
	}
}

func TestMatchAll(t *testing.T) {
	tests := []struct {
		name    string
		entry   string
		filters []string
		want    bool
	}{
		{"all match", "fix #later #ci", []string{"later", "ci"}, true},
		{"one fails", "fix #later", []string{"later", "ci"}, false},
		{"empty filters", "fix #later", nil, true},
		{"negation combo", "fix #later #ci", []string{"later", "~urgent"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchAll(tt.entry, tt.filters); got != tt.want {
				t.Errorf("MatchAll(%q, %v) = %v, want %v",
					tt.entry, tt.filters, got, tt.want)
			}
		})
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
