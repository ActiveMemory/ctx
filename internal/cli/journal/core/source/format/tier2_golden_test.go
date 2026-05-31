//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/entity"
)

// TestJournalEntryPartMatchesLegacy asserts JournalEntryPart still
// produces, byte-for-byte, the pre-migration output for the metadata
// table, plan, and tool-result <details> paths now rendered via the
// metaTable and details templates. Fixtures were captured from the
// legacy fmt.Sprintf/paired-tag code path (see git history).
func TestJournalEntryPartMatchesLegacy(t *testing.T) {
	t.Setenv("TZ", "UTC")

	base := func() *entity.Session {
		return &entity.Session{
			ID: "abc12345-session-id", Slug: "test-slug",
			Tool: "claude-code", Project: "myproject",
			StartTime: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
			EndTime:   time.Date(2026, 1, 15, 11, 0, 0, 0, time.UTC),
			Duration:  30 * time.Minute, TurnCount: 2,
			TotalTokens: 15000, TotalTokensIn: 10000, TotalTokensOut: 5000,
		}
	}

	single := base()
	single.Messages = []entity.Message{
		{Role: "user", Text: "Hello",
			Timestamp: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)},
		{Role: "assistant", Text: "Hi there!",
			Timestamp: time.Date(2026, 1, 15, 10, 30, 5, 0, time.UTC)},
	}

	metafull := base()
	metafull.GitBranch = "main"
	metafull.Model = "claude-opus"
	metafull.Messages = single.Messages

	plan := base()
	plan.Messages = []entity.Message{
		{Role: "assistant", PlanContent: "step one\nstep two",
			Timestamp: time.Date(2026, 1, 15, 10, 31, 0, 0, time.UTC)},
	}

	tooluse := base()
	tooluse.Messages = []entity.Message{
		{Role: "assistant",
			Timestamp: time.Date(2026, 1, 15, 10, 32, 0, 0, time.UTC),
			ToolUses: []entity.ToolUse{
				{ID: "t1", Name: "Read", Input: `{"file_path":"/tmp/x.go"}`},
			}},
		{Role: "user",
			Timestamp: time.Date(2026, 1, 15, 10, 32, 1, 0, time.UTC),
			ToolResults: []entity.ToolResult{
				{ToolUseID: "t1", Content: "package main\nfunc main() {}"},
				{ToolUseID: "t2", Content: "boom", IsError: true},
				{ToolUseID: "t3", Content: strings.Repeat("line\n", 15)},
			}},
	}

	cases := map[string]struct {
		s           *entity.Session
		part, total int
	}{
		"single":   {single, 1, 1},
		"metafull": {metafull, 1, 2},
		"plan":     {plan, 1, 1},
		"tooluse":  {tooluse, 1, 1},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			want, readErr := os.ReadFile(
				filepath.Join("testdata", name+".golden"),
			)
			if readErr != nil {
				t.Fatal(readErr)
			}
			got := JournalEntryPart(
				c.s, c.s.Messages, 0, c.part, c.total, "b", "",
			)
			if got != string(want) {
				t.Errorf(
					"drift:\n--- want ---\n%q\n--- got ---\n%q",
					string(want), got,
				)
			}
		})
	}
}
