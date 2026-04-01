//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"testing"
)

func TestParseAddedDecisions(t *testing.T) {
	diff := `diff --git a/.context/DECISIONS.md b/.context/DECISIONS.md
index abc1234..def5678 100644
--- a/.context/DECISIONS.md
+++ b/.context/DECISIONS.md
@@ -1,3 +1,11 @@
 # Decisions

+## [2026-01-28-051426] Use SQLite for storage
+
+Rationale: lightweight and embedded.
+
+## [2026-01-29-120000] Prefer JSONL over CSV
+
+Rationale: easier streaming.
+
 ## [2026-01-01-000000] Existing decision

 Already there.`

	refs := parseAddedEntries(diff, "decision")
	if len(refs) != 2 {
		t.Fatalf("parseAddedEntries() returned %d refs, want 2: %v", len(refs), refs)
	}
	if refs[0] != "decision:1" {
		t.Errorf("refs[0] = %q, want %q", refs[0], "decision:1")
	}
	if refs[1] != "decision:2" {
		t.Errorf("refs[1] = %q, want %q", refs[1], "decision:2")
	}
}

func TestParseAddedLearnings(t *testing.T) {
	diff := `diff --git a/.context/LEARNINGS.md b/.context/LEARNINGS.md
index abc1234..def5678 100644
--- a/.context/LEARNINGS.md
+++ b/.context/LEARNINGS.md
@@ -1,5 +1,9 @@
 # Learnings

+## [2026-03-01-090000] Always quote shell paths
+
+Unquoted paths with spaces cause subtle failures.
+
 ## [2026-01-15-083000] Existing learning

 Already there.`

	refs := parseAddedEntries(diff, "learning")
	if len(refs) != 1 {
		t.Fatalf("parseAddedEntries() returned %d refs, want 1: %v", len(refs), refs)
	}
	if refs[0] != "learning:1" {
		t.Errorf("refs[0] = %q, want %q", refs[0], "learning:1")
	}
}

func TestParseAddedTasks(t *testing.T) {
	diff := `diff --git a/.context/TASKS.md b/.context/TASKS.md
index abc1234..def5678 100644
--- a/.context/TASKS.md
+++ b/.context/TASKS.md
@@ -1,7 +1,9 @@
 # Tasks

+- [x] Implement staged file analysis
+- [x] Write unit tests for trace package
 - [ ] Pending task one
 - [ ] Pending task two`

	refs := parseCompletedTasks(diff)
	if len(refs) != 2 {
		t.Fatalf("parseCompletedTasks() returned %d refs, want 2: %v", len(refs), refs)
	}
	if refs[0] != "task:1" {
		t.Errorf("refs[0] = %q, want %q", refs[0], "task:1")
	}
	if refs[1] != "task:2" {
		t.Errorf("refs[1] = %q, want %q", refs[1], "task:2")
	}
}

func TestParseNoAdditions(t *testing.T) {
	diff := `diff --git a/.context/DECISIONS.md b/.context/DECISIONS.md
index abc1234..def5678 100644
--- a/.context/DECISIONS.md
+++ b/.context/DECISIONS.md
@@ -1,5 +1,3 @@
 # Decisions

-## [2026-01-01-000000] Removed decision
-
-Old rationale.`

	refs := parseAddedEntries(diff, "decision")
	if len(refs) != 0 {
		t.Errorf("parseAddedEntries() returned %d refs, want 0: %v", len(refs), refs)
	}
}
