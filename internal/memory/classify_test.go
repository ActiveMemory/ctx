//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		target string
	}{
		{
			name:   "convention: always use",
			text:   "always use bun for this project",
			target: config.EntryConvention,
		},
		{
			name:   "convention: prefer",
			text:   "prefer filepath.Join over string concatenation",
			target: config.EntryConvention,
		},
		{
			name:   "convention: never use",
			text:   "never use global state in handlers",
			target: config.EntryConvention,
		},
		{
			name:   "decision: decided",
			text:   "decided to use SQLite over Postgres for local storage",
			target: config.EntryDecision,
		},
		{
			name:   "decision: chose",
			text:   "chose marker-based merge for bidirectional sync",
			target: config.EntryDecision,
		},
		{
			name:   "learning: learned",
			text:   "learned that golangci-lint v2 ignores inline nolint",
			target: config.EntryLearning,
		},
		{
			name:   "learning: gotcha",
			text:   "gotcha: symlinks in project path produce different slugs",
			target: config.EntryLearning,
		},
		{
			name:   "learning: watch out",
			text:   "watch out for race conditions in concurrent map access",
			target: config.EntryLearning,
		},
		{
			name:   "task: need to",
			text:   "need to add tests for the import command",
			target: config.EntryTask,
		},
		{
			name:   "task: todo",
			text:   "todo: wire up the publish command",
			target: config.EntryTask,
		},
		{
			name:   "skip: session notes",
			text:   "Session 2026-03-05: Memory Bridge Design",
			target: TargetSkip,
		},
		{
			name:   "skip: generic paragraph",
			text:   "Worked on the memory bridge today. Made good progress.",
			target: TargetSkip,
		},
		{
			name:   "case insensitive",
			text:   "ALWAYS USE ctx from PATH",
			target: config.EntryConvention,
		},
		{
			name:   "convention wins over decision (priority order)",
			text:   "always use the approach we decided on",
			target: config.EntryConvention,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := Entry{Text: tt.text}
			got := Classify(entry)
			if got.Target != tt.target {
				t.Errorf("Classify(%q) = %q, want %q (keywords: %v)",
					tt.text, got.Target, tt.target, got.Keywords)
			}
		})
	}
}
