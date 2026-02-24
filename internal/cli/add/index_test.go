//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

// TestDelegation verifies that the wrapper functions correctly delegate
// to the internal/index package. The full logic is tested in internal/index.
func TestDelegation(t *testing.T) {
	content := `# Decisions

## [2026-01-28-051426] Test decision

**Status**: Accepted
`

	// Test ParseEntryHeaders delegation
	entries := ParseEntryHeaders(content)
	if len(entries) != 1 {
		t.Errorf("ParseEntryHeaders() got %d entries, want 1", len(entries))
	}
	if entries[0].Date != "2026-01-28" {
		t.Errorf("ParseEntryHeaders() entry.Date = %q, want %q", entries[0].Date, "2026-01-28")
	}

	// Test ParseDecisionHeaders alias
	decisionEntries := ParseDecisionHeaders(content)
	if len(decisionEntries) != 1 {
		t.Errorf("ParseDecisionHeaders() got %d entries, want 1", len(decisionEntries))
	}

	// Test GenerateIndexTable delegation
	table := GenerateIndexTable(entries, "Decision")
	if !strings.Contains(table, "| Date | Decision |") {
		t.Error("GenerateIndexTable() missing header")
	}

	// Test GenerateIndex convenience function
	indexTable := GenerateIndex(entries)
	if !strings.Contains(indexTable, "| Date | Decision |") {
		t.Error("GenerateIndex() missing header")
	}

	// Test UpdateIndex delegation
	updated := UpdateIndex(content)
	if !strings.Contains(updated, config.IndexStart) {
		t.Error("UpdateIndex() missing INDEX:START marker")
	}

	// Test UpdateLearningsIndex delegation
	learningContent := `# Learnings

## [2026-01-28-191951] Test learning

**Context**: Test

**Lesson**: Test

**Application**: Test
`
	updatedLearning := UpdateLearningsIndex(learningContent)
	if !strings.Contains(updatedLearning, config.IndexStart) {
		t.Error("UpdateLearningsIndex() missing INDEX:START marker")
	}
	if !strings.Contains(updatedLearning, "| Date | Learning |") {
		t.Error("UpdateLearningsIndex() missing Learning column header")
	}
}
