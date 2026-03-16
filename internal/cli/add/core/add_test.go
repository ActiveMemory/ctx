//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"strings"
	"testing"
)

// TestAppendEntry tests the AppendEntry function directly.
func TestAppendEntry(t *testing.T) {
	t.Run("decision prepends after header", func(t *testing.T) {
		// Use timestamp format "## [" to match what FormatDecision produces
		existing := []byte("# Decisions\n\n## [2026-01-01] Old Decision\n\nContent\n")
		entry := "## [2026-01-02] New Decision\n\nNew content\n"

		result := AppendEntry(existing, entry, "decision", "")

		resultStr := string(result)
		newIdx := strings.Index(resultStr, "New Decision")
		oldIdx := strings.Index(resultStr, "Old Decision")

		if newIdx == -1 || oldIdx == -1 {
			t.Fatalf("decisions not found in result: %s", resultStr)
		}
		if newIdx >= oldIdx {
			t.Errorf("new decision should appear before old, but new at %d, old at %d", newIdx, oldIdx)
		}
	})

	t.Run("learning prepends after separator", func(t *testing.T) {
		// Use section format "## [" to match what FormatLearning produces
		existing := []byte("# Learnings\n\n<!-- comment -->\n\n## [2026-01-01] Old Learning\n\nContent\n")
		entry := "## [2026-01-02] New Learning\n\nContent\n"

		result := AppendEntry(existing, entry, "learning", "")

		resultStr := string(result)
		newIdx := strings.Index(resultStr, "New Learning")
		oldIdx := strings.Index(resultStr, "Old Learning")

		if newIdx == -1 || oldIdx == -1 {
			t.Fatalf("learnings not found in result: %s", resultStr)
		}
		if newIdx >= oldIdx {
			t.Errorf("new learning should appear before old, but new at %d, old at %d", newIdx, oldIdx)
		}
	})

	t.Run("convention appends at end", func(t *testing.T) {
		existing := []byte("# Conventions\n\n- Old convention\n")
		entry := "- New convention\n"

		result := AppendEntry(existing, entry, "convention", "")

		resultStr := string(result)
		newIdx := strings.Index(resultStr, "New convention")
		oldIdx := strings.Index(resultStr, "Old convention")

		if newIdx == -1 || oldIdx == -1 {
			t.Fatal("conventions not found in result")
		}
		if newIdx <= oldIdx {
			t.Errorf("new convention should appear after old (appended), but new at %d, old at %d", newIdx, oldIdx)
		}
	})

	t.Run("decision on empty file", func(t *testing.T) {
		existing := []byte("# Decisions\n\n<!-- Add decisions here -->\n")
		entry := "## [2026-01-01] First Decision\n\nContent\n"

		result := AppendEntry(existing, entry, "decision", "")

		if !strings.Contains(string(result), "First Decision") {
			t.Errorf("decision not found in result: %s", result)
		}
	})

	t.Run("learning on empty file", func(t *testing.T) {
		existing := []byte("# Learnings\n\n<!-- Add gotchas here -->\n")
		entry := "## [2026-01-01] First Learning\n\nContent\n"

		result := AppendEntry(existing, entry, "learning", "")

		if !strings.Contains(string(result), "First Learning") {
			t.Errorf("learning not found in result: %s", result)
		}
	})
}

// TestInsertDecisionNotInsideComment is a regression test for the bug where
// ctx add decision inserts the entry inside the <!-- DECISION FORMATS -->
// HTML comment block on a fresh DECISIONS.md, making it invisible when rendered.
func TestInsertDecisionNotInsideComment(t *testing.T) {
	// This is the exact structure of a freshly scaffolded DECISIONS.md.
	freshDecisions := "# Decisions\n\n" +
		"<!-- INDEX:START -->\n" +
		"<!-- INDEX:END -->\n\n" +
		"<!-- DECISION FORMATS\n\n" +
		"## Quick Format (Y-Statement)\n\n" +
		"For lightweight decisions, a single statement suffices.\n\n" +
		"## Full Format\n\n" +
		"For significant decisions:\n\n" +
		"## [YYYY-MM-DD] Decision Title\n\n" +
		"**Status**: Accepted\n\n" +
		"-->\n"

	entry := "## [2026-02-18] My Important Decision\n\n" +
		"**Status**: Accepted\n\n" +
		"**Context**: What prompted it\n\n" +
		"**Rationale**: Why this choice\n\n" +
		"**Consequence**: What changes\n"

	result := AppendEntry([]byte(freshDecisions), entry, "decision", "")
	resultStr := string(result)

	// The entry must appear in the output.
	if !strings.Contains(resultStr, "My Important Decision") {
		t.Fatalf("decision not found in result:\n%s", resultStr)
	}

	formatBlockCloseIdx := strings.LastIndex(freshDecisions, "-->")
	entryIdx := strings.Index(resultStr, "My Important Decision")
	if formatBlockCloseIdx == -1 {
		t.Fatal("closing --> not found in template")
	}
	if entryIdx <= formatBlockCloseIdx {
		t.Errorf(
			"decision was inserted inside the HTML comment block: "+
				"entry at index %d, DECISION FORMATS block closes at index %d\n\nFull result:\n%s",
			entryIdx, formatBlockCloseIdx, resultStr,
		)
	}
}

// TestInsertTaskDefaultPlacement tests task insertion without --section.
func TestInsertTaskDefaultPlacement(t *testing.T) {
	t.Run("inserts before first unchecked task", func(t *testing.T) {
		existing := []byte("# Tasks\n\n### Phase 1\n\n- [x] Done task\n- [ ] Pending task\n")
		entry := "- [ ] New task\n"

		result := AppendEntry(existing, entry, "task", "")
		resultStr := string(result)

		newIdx := strings.Index(resultStr, "New task")
		pendingIdx := strings.Index(resultStr, "Pending task")
		doneIdx := strings.Index(resultStr, "Done task")

		if newIdx == -1 || pendingIdx == -1 {
			t.Fatalf("tasks not found in result:\n%s", resultStr)
		}
		if newIdx >= pendingIdx {
			t.Errorf("new task should appear before existing pending task, but new at %d, pending at %d", newIdx, pendingIdx)
		}
		if newIdx <= doneIdx {
			t.Errorf("new task should appear after completed task, but new at %d, done at %d", newIdx, doneIdx)
		}
	})

	t.Run("appends at end when all tasks checked", func(t *testing.T) {
		existing := []byte("# Tasks\n\n- [x] Done one\n- [x] Done two\n")
		entry := "- [ ] New task\n"

		result := AppendEntry(existing, entry, "task", "")
		resultStr := string(result)

		newIdx := strings.Index(resultStr, "New task")
		lastDoneIdx := strings.LastIndex(resultStr, "Done two")

		if newIdx == -1 {
			t.Fatalf("new task not found in result:\n%s", resultStr)
		}
		if newIdx <= lastDoneIdx {
			t.Errorf("new task should appear after completed tasks, but new at %d, last done at %d", newIdx, lastDoneIdx)
		}
	})

	t.Run("explicit section overrides auto-placement", func(t *testing.T) {
		existing := []byte("# Tasks\n\n### Phase 1\n\n- [ ] Phase 1 task\n\n### Maintenance\n\n- [ ] Maint task\n")
		entry := "- [ ] New maint task\n"

		result := AppendEntry(existing, entry, "task", "Maintenance")
		resultStr := string(result)

		newIdx := strings.Index(resultStr, "New maint task")
		maintHeaderIdx := strings.Index(resultStr, "### Maintenance")
		maintTaskIdx := strings.Index(resultStr, "Maint task")

		if newIdx == -1 {
			t.Fatalf("new task not found in result:\n%s", resultStr)
		}
		if newIdx <= maintHeaderIdx {
			t.Errorf("new task should appear after Maintenance header, but new at %d, header at %d", newIdx, maintHeaderIdx)
		}
		if newIdx >= maintTaskIdx {
			t.Errorf("new task should appear before existing maint task, but new at %d, existing at %d", newIdx, maintTaskIdx)
		}
	})

	t.Run("phased file without Next Up section", func(t *testing.T) {
		existing := []byte("# Tasks\n\n### Phase 1\n\n- [x] Old done\n\n### Phase 2\n\n- [ ] Phase 2 pending\n")
		entry := "- [ ] New task\n"

		result := AppendEntry(existing, entry, "task", "")
		resultStr := string(result)

		newIdx := strings.Index(resultStr, "New task")
		phase2Idx := strings.Index(resultStr, "Phase 2 pending")

		if newIdx == -1 || phase2Idx == -1 {
			t.Fatalf("tasks not found in result:\n%s", resultStr)
		}
		if newIdx >= phase2Idx {
			t.Errorf("new task should appear before Phase 2 pending task, but new at %d, phase2 at %d", newIdx, phase2Idx)
		}
	})
}

// TestIsInsideHTMLComment unit-tests the IsInsideHTMLComment helper directly.
func TestIsInsideHTMLComment(t *testing.T) {
	cases := []struct {
		name    string
		content string
		idx     int
		want    bool
	}{
		{
			name:    "position before any comment",
			content: "hello <!-- comment --> world",
			idx:     3, // inside "hello"
			want:    false,
		},
		{
			name:    "position inside comment",
			content: "hello <!-- comment --> world",
			idx:     10, // inside " comment "
			want:    true,
		},
		{
			name:    "position after comment close",
			content: "hello <!-- comment --> world",
			idx:     23, // inside " world"
			want:    false,
		},
		{
			name:    "inline single-line comment: position after close",
			content: "<!-- INDEX:START -->\n## [real entry]",
			idx:     20, // start of "## ["
			want:    false,
		},
		{
			name:    "multi-line comment: heading inside",
			content: "<!-- FORMATS\n## [YYYY-MM-DD] Template\n-->\n## [real]",
			idx:     13, // start of "## [YYYY"
			want:    true,
		},
		{
			name:    "multi-line comment: heading after close",
			content: "<!-- FORMATS\n## [YYYY-MM-DD] Template\n-->\n## [real]",
			idx:     41, // start of "## [real]"
			want:    false,
		},
		{
			name:    "unclosed comment treated as inside",
			content: "<!-- unclosed\n## [heading]",
			idx:     14,
			want:    true,
		},
		{
			name:    "no comment at all",
			content: "# Decisions\n\n## [2026-01-01] Entry\n",
			idx:     13, // start of "## ["
			want:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsInsideHTMLComment(tc.content, tc.idx)
			if got != tc.want {
				t.Errorf("IsInsideHTMLComment(%q, %d) = %v, want %v",
					tc.content, tc.idx, got, tc.want)
			}
		})
	}
}
