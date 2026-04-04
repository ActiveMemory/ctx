# Commit Context Tracing Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Link every git commit back to the decisions, tasks, learnings, and sessions that motivated it via `ctx trace`.

**Architecture:** New `internal/trace` package provides the core logic (pending context recording, three-source detection, history/override storage, reference resolution). A new `internal/cli/trace` package wires it into the Cobra CLI as `ctx trace`. Existing commands (`ctx add`, `ctx complete`) gain a one-line `trace.Record()` side-effect. A `ctx trace hook` subcommand generates a prepare-commit-msg shell script that delegates to `ctx trace collect`.

**Tech Stack:** Go, Cobra CLI, JSONL storage, git trailers, prepare-commit-msg hook

---

## File Structure

### New files

| File | Responsibility |
|------|---------------|
| `internal/trace/pending.go` | Record/read/truncate pending context refs in `state/pending-context.jsonl` |
| `internal/trace/pending_test.go` | Tests for pending context operations |
| `internal/trace/staged.go` | Detect context refs from staged `.context/` file diffs |
| `internal/trace/staged_test.go` | Tests for staged file analysis |
| `internal/trace/working.go` | Detect context refs from current working state (in-progress tasks, session env) |
| `internal/trace/working_test.go` | Tests for working state detection |
| `internal/trace/collect.go` | Merge + deduplicate refs from all three sources |
| `internal/trace/collect_test.go` | Tests for collection/merge |
| `internal/trace/history.go` | Read/write `trace/history.jsonl` and `trace/overrides.jsonl` |
| `internal/trace/history_test.go` | Tests for history/override storage |
| `internal/trace/resolve.go` | Resolve ref strings to human-readable context (read DECISIONS.md, etc.) |
| `internal/trace/resolve_test.go` | Tests for reference resolution |
| `internal/trace/types.go` | Shared types: `PendingEntry`, `HistoryEntry`, `OverrideEntry`, `Ref`, `ResolvedRef` |
| `internal/trace/doc.go` | Package documentation |
| `internal/cli/trace/trace.go` | Top-level `Cmd()` that returns the `trace` cobra.Command |
| `internal/cli/trace/cmd/show/cmd.go` | `ctx trace <commit>` and `ctx trace --last N` command definition |
| `internal/cli/trace/cmd/show/run.go` | Execution logic for showing commit context |
| `internal/cli/trace/cmd/file/cmd.go` | `ctx trace file <path>` command definition |
| `internal/cli/trace/cmd/file/run.go` | Execution logic for file tracing |
| `internal/cli/trace/cmd/tag/cmd.go` | `ctx trace tag <commit>` command definition |
| `internal/cli/trace/cmd/tag/run.go` | Execution logic for manual tagging |
| `internal/cli/trace/cmd/collect/cmd.go` | `ctx trace collect` — called by the hook to collect and output trailer |
| `internal/cli/trace/cmd/collect/run.go` | Execution logic for collect |
| `internal/cli/trace/cmd/hook/cmd.go` | `ctx trace hook enable/disable` — manages prepare-commit-msg hook |
| `internal/cli/trace/cmd/hook/run.go` | Hook management logic |
| `internal/config/embed/cmd/trace.go` | Use strings and DescKey constants for trace commands |
| `internal/err/trace/trace.go` | Error constructors for trace operations |
| `internal/err/trace/doc.go` | Package documentation |
| `internal/write/trace/trace.go` | Output formatters for trace results |

### Modified files

| File | Change |
|------|--------|
| `internal/cli/add/cmd/root/run.go` | Add `trace.Record()` call after successful write |
| `internal/cli/task/cmd/complete/run.go` | Add `trace.Record()` call after marking complete |
| `internal/bootstrap/group.go` | Register `trace.Cmd` in the diagnostics group |
| `internal/config/embed/cmd/base.go` | Add `UseTrace` and `DescKeyTrace` constants |
| `internal/config/dir/dir.go` | Add `Trace = "trace"` constant |
| `internal/assets/commands/commands.yaml` | Add trace command descriptions |
| `internal/assets/commands/text/write.yaml` | Add trace output format strings |

---

## Task 1: Core Types and Pending Context Recording

**Files:**
- Create: `internal/trace/doc.go`
- Create: `internal/trace/types.go`
- Create: `internal/trace/pending.go`
- Create: `internal/trace/pending_test.go`
- Modify: `internal/config/dir/dir.go`

### Steps

- [ ] **Step 1: Add Trace directory constant**

In `internal/config/dir/dir.go`, add the `Trace` constant:

```go
// Trace is the subdirectory for commit context tracing within .context/.
Trace = "trace"
```

Add it after the `State` constant in the same `const` block.

- [ ] **Step 2: Create trace package doc**

Create `internal/trace/doc.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace provides commit context tracing — linking git commits
// back to the decisions, tasks, learnings, and sessions that motivated them.
package trace
```

- [ ] **Step 3: Create shared types**

Create `internal/trace/types.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import "time"

// PendingEntry is a single pending context reference accumulated
// between commits.
type PendingEntry struct {
	Ref       string    `json:"ref"`
	Timestamp time.Time `json:"timestamp"`
}

// HistoryEntry is a permanent record of a commit's context references.
type HistoryEntry struct {
	Commit    string    `json:"commit"`
	Refs      []string  `json:"refs"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// OverrideEntry is a manual context tag added to an existing commit.
type OverrideEntry struct {
	Commit    string    `json:"commit"`
	Refs      []string  `json:"refs"`
	Timestamp time.Time `json:"timestamp"`
}

// ResolvedRef holds a resolved context reference with its display text.
type ResolvedRef struct {
	Raw     string // Original ref string (e.g., "decision:12")
	Type    string // "decision", "learning", "task", "convention", "session", "note"
	Number  int    // Entry number (0 for session/note types)
	Title   string // Resolved title or content
	Detail  string // Additional detail (rationale, status, etc.)
	Found   bool   // Whether the reference was resolved
}
```

- [ ] **Step 4: Write failing test for Record**

Create `internal/trace/pending_test.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecord(t *testing.T) {
	tmpDir := t.TempDir()
	stateDir := filepath.Join(tmpDir, "state")
	if err := os.MkdirAll(stateDir, 0750); err != nil {
		t.Fatal(err)
	}

	if err := Record("decision:1", stateDir); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(stateDir, pendingFile))
	if err != nil {
		t.Fatalf("read pending file: %v", err)
	}

	var entry PendingEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if entry.Ref != "decision:1" {
		t.Errorf("got ref %q, want %q", entry.Ref, "decision:1")
	}
	if entry.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestRecordMultiple(t *testing.T) {
	tmpDir := t.TempDir()
	stateDir := filepath.Join(tmpDir, "state")
	if err := os.MkdirAll(stateDir, 0750); err != nil {
		t.Fatal(err)
	}

	_ = Record("decision:1", stateDir)
	_ = Record("task:3", stateDir)
	_ = Record("session:abc123", stateDir)

	data, err := os.ReadFile(filepath.Join(stateDir, pendingFile))
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestReadPending(t *testing.T) {
	tmpDir := t.TempDir()
	stateDir := filepath.Join(tmpDir, "state")
	if err := os.MkdirAll(stateDir, 0750); err != nil {
		t.Fatal(err)
	}

	_ = Record("decision:1", stateDir)
	_ = Record("task:3", stateDir)

	entries, err := ReadPending(stateDir)
	if err != nil {
		t.Fatalf("ReadPending: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Ref != "decision:1" {
		t.Errorf("first ref: got %q, want %q", entries[0].Ref, "decision:1")
	}
	if entries[1].Ref != "task:3" {
		t.Errorf("second ref: got %q, want %q", entries[1].Ref, "task:3")
	}
}

func TestReadPendingEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	stateDir := filepath.Join(tmpDir, "state")
	if err := os.MkdirAll(stateDir, 0750); err != nil {
		t.Fatal(err)
	}

	entries, err := ReadPending(stateDir)
	if err != nil {
		t.Fatalf("ReadPending on missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestTruncatePending(t *testing.T) {
	tmpDir := t.TempDir()
	stateDir := filepath.Join(tmpDir, "state")
	if err := os.MkdirAll(stateDir, 0750); err != nil {
		t.Fatal(err)
	}

	_ = Record("decision:1", stateDir)
	_ = Record("task:3", stateDir)

	if err := TruncatePending(stateDir); err != nil {
		t.Fatalf("TruncatePending: %v", err)
	}

	entries, err := ReadPending(stateDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after truncate, got %d", len(entries))
	}
}
```

- [ ] **Step 5: Run test to verify it fails**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -run TestRecord -v`
Expected: FAIL — functions not defined

- [ ] **Step 6: Implement pending.go**

Create `internal/trace/pending.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/fs"
)

const pendingFile = "pending-context.jsonl"

// Record appends a context reference to the pending context file.
// This is best-effort: errors are returned but callers should treat
// them as non-fatal.
//
// Parameters:
//   - ref: Context reference string (e.g., "decision:12", "task:3")
//   - stateDir: Path to the state directory (.context/state/)
//
// Returns:
//   - error: Non-nil if the file cannot be opened or written
func Record(ref, stateDir string) error {
	if err := os.MkdirAll(stateDir, fs.PermRestrictedDir); err != nil {
		return err
	}

	p := filepath.Join(stateDir, pendingFile)

	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.PermFile)
	if err != nil {
		return err
	}
	defer f.Close()

	entry := PendingEntry{Ref: ref, Timestamp: time.Now().UTC()}
	return json.NewEncoder(f).Encode(entry)
}

// ReadPending reads all pending context entries from the state directory.
// Returns an empty slice if the file does not exist.
//
// Parameters:
//   - stateDir: Path to the state directory (.context/state/)
//
// Returns:
//   - []PendingEntry: Parsed entries
//   - error: Non-nil on read or parse failure
func ReadPending(stateDir string) ([]PendingEntry, error) {
	p := filepath.Join(stateDir, pendingFile)

	f, err := os.Open(filepath.Clean(p))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var entries []PendingEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var entry PendingEntry
		if jsonErr := json.Unmarshal([]byte(line), &entry); jsonErr != nil {
			continue // skip malformed lines
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

// TruncatePending clears the pending context file after a commit.
//
// Parameters:
//   - stateDir: Path to the state directory (.context/state/)
//
// Returns:
//   - error: Non-nil if truncation fails
func TruncatePending(stateDir string) error {
	p := filepath.Join(stateDir, pendingFile)
	return os.Truncate(p, 0)
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -v`
Expected: All PASS

- [ ] **Step 8: Commit**

```bash
git add internal/trace/doc.go internal/trace/types.go internal/trace/pending.go internal/trace/pending_test.go internal/config/dir/dir.go
git commit -m "feat(trace): add pending context recording"
```

---

## Task 2: History and Override Storage

**Files:**
- Create: `internal/trace/history.go`
- Create: `internal/trace/history_test.go`

### Steps

- [ ] **Step 1: Write failing test for history operations**

Create `internal/trace/history_test.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteHistory(t *testing.T) {
	tmpDir := t.TempDir()
	traceDir := filepath.Join(tmpDir, "trace")
	if err := os.MkdirAll(traceDir, 0750); err != nil {
		t.Fatal(err)
	}

	entry := HistoryEntry{
		Commit:  "abc123",
		Refs:    []string{"decision:12", "task:8"},
		Message: "Fix auth token expiry",
	}

	if err := WriteHistory(entry, traceDir); err != nil {
		t.Fatalf("WriteHistory: %v", err)
	}

	entries, err := ReadHistory(traceDir)
	if err != nil {
		t.Fatalf("ReadHistory: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Commit != "abc123" {
		t.Errorf("commit: got %q, want %q", entries[0].Commit, "abc123")
	}
	if len(entries[0].Refs) != 2 {
		t.Errorf("refs count: got %d, want 2", len(entries[0].Refs))
	}
}

func TestReadHistoryForCommit(t *testing.T) {
	tmpDir := t.TempDir()
	traceDir := filepath.Join(tmpDir, "trace")
	if err := os.MkdirAll(traceDir, 0750); err != nil {
		t.Fatal(err)
	}

	_ = WriteHistory(HistoryEntry{
		Commit: "abc123", Refs: []string{"decision:12"}, Message: "First",
	}, traceDir)
	_ = WriteHistory(HistoryEntry{
		Commit: "def456", Refs: []string{"task:3"}, Message: "Second",
	}, traceDir)

	entry, found := ReadHistoryForCommit("abc123", traceDir)
	if !found {
		t.Fatal("expected to find commit abc123")
	}
	if entry.Commit != "abc123" {
		t.Errorf("got commit %q", entry.Commit)
	}

	_, found = ReadHistoryForCommit("missing", traceDir)
	if found {
		t.Error("should not find missing commit")
	}
}

func TestWriteOverride(t *testing.T) {
	tmpDir := t.TempDir()
	traceDir := filepath.Join(tmpDir, "trace")
	if err := os.MkdirAll(traceDir, 0750); err != nil {
		t.Fatal(err)
	}

	entry := OverrideEntry{
		Commit: "abc123",
		Refs:   []string{`"Hotfix for production outage"`},
	}

	if err := WriteOverride(entry, traceDir); err != nil {
		t.Fatalf("WriteOverride: %v", err)
	}

	entries, err := ReadOverrides(traceDir)
	if err != nil {
		t.Fatalf("ReadOverrides: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Commit != "abc123" {
		t.Errorf("commit: got %q", entries[0].Commit)
	}
}

func TestReadOverridesForCommit(t *testing.T) {
	tmpDir := t.TempDir()
	traceDir := filepath.Join(tmpDir, "trace")
	if err := os.MkdirAll(traceDir, 0750); err != nil {
		t.Fatal(err)
	}

	_ = WriteOverride(OverrideEntry{
		Commit: "abc123", Refs: []string{`"Note one"`},
	}, traceDir)
	_ = WriteOverride(OverrideEntry{
		Commit: "abc123", Refs: []string{`"Note two"`},
	}, traceDir)
	_ = WriteOverride(OverrideEntry{
		Commit: "def456", Refs: []string{"decision:5"},
	}, traceDir)

	refs := ReadOverridesForCommit("abc123", traceDir)
	if len(refs) != 2 {
		t.Fatalf("expected 2 override refs for abc123, got %d", len(refs))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -run TestWriteHistory -v`
Expected: FAIL — functions not defined

- [ ] **Step 3: Implement history.go**

Create `internal/trace/history.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/fs"
)

const (
	historyFile  = "history.jsonl"
	overrideFile = "overrides.jsonl"
)

// WriteHistory appends a commit context record to history.jsonl.
//
// Parameters:
//   - entry: The history entry to write
//   - traceDir: Path to the trace directory (.context/trace/)
//
// Returns:
//   - error: Non-nil if the file cannot be opened or written
func WriteHistory(entry HistoryEntry, traceDir string) error {
	if err := os.MkdirAll(traceDir, fs.PermRestrictedDir); err != nil {
		return err
	}

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	p := filepath.Join(traceDir, historyFile)
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.PermFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(entry)
}

// ReadHistory reads all history entries from the trace directory.
//
// Parameters:
//   - traceDir: Path to the trace directory (.context/trace/)
//
// Returns:
//   - []HistoryEntry: Parsed entries (may be empty)
//   - error: Non-nil on read or parse failure
func ReadHistory(traceDir string) ([]HistoryEntry, error) {
	return readJSONL[HistoryEntry](filepath.Join(traceDir, historyFile))
}

// ReadHistoryForCommit finds the history entry for a specific commit.
// Matches by prefix to support short commit hashes.
//
// Parameters:
//   - commitHash: Full or abbreviated commit hash
//   - traceDir: Path to the trace directory
//
// Returns:
//   - HistoryEntry: The matching entry
//   - bool: True if found
func ReadHistoryForCommit(commitHash, traceDir string) (HistoryEntry, bool) {
	entries, err := ReadHistory(traceDir)
	if err != nil {
		return HistoryEntry{}, false
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Commit, commitHash) || strings.HasPrefix(commitHash, e.Commit) {
			return e, true
		}
	}
	return HistoryEntry{}, false
}

// WriteOverride appends a manual tag entry to overrides.jsonl.
//
// Parameters:
//   - entry: The override entry to write
//   - traceDir: Path to the trace directory (.context/trace/)
//
// Returns:
//   - error: Non-nil if the file cannot be opened or written
func WriteOverride(entry OverrideEntry, traceDir string) error {
	if err := os.MkdirAll(traceDir, fs.PermRestrictedDir); err != nil {
		return err
	}

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	p := filepath.Join(traceDir, overrideFile)
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.PermFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(entry)
}

// ReadOverrides reads all override entries from the trace directory.
//
// Parameters:
//   - traceDir: Path to the trace directory (.context/trace/)
//
// Returns:
//   - []OverrideEntry: Parsed entries (may be empty)
//   - error: Non-nil on read or parse failure
func ReadOverrides(traceDir string) ([]OverrideEntry, error) {
	return readJSONL[OverrideEntry](filepath.Join(traceDir, overrideFile))
}

// ReadOverridesForCommit collects all override refs for a specific commit.
//
// Parameters:
//   - commitHash: Full or abbreviated commit hash
//   - traceDir: Path to the trace directory
//
// Returns:
//   - []string: All override refs for this commit
func ReadOverridesForCommit(commitHash, traceDir string) []string {
	entries, err := ReadOverrides(traceDir)
	if err != nil {
		return nil
	}

	var refs []string
	for _, e := range entries {
		if strings.HasPrefix(e.Commit, commitHash) || strings.HasPrefix(commitHash, e.Commit) {
			refs = append(refs, e.Refs...)
		}
	}
	return refs
}

// readJSONL is a generic helper for reading JSONL files.
func readJSONL[T any](path string) ([]T, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var entries []T
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var entry T
		if jsonErr := json.Unmarshal([]byte(line), &entry); jsonErr != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/trace/history.go internal/trace/history_test.go
git commit -m "feat(trace): add history and override storage"
```

---

## Task 3: Staged File Analysis (Source 2)

**Files:**
- Create: `internal/trace/staged.go`
- Create: `internal/trace/staged_test.go`

### Steps

- [ ] **Step 1: Write failing test for staged detection**

Create `internal/trace/staged_test.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import "testing"

func TestParseAddedDecisions(t *testing.T) {
	diff := `+## [2026-03-14-100000] Use short-lived tokens
+
+**Context:** Security review
+
+## [2026-03-14-110000] Rate limiting strategy`

	refs := ParseAddedEntries(diff, "decision")
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
	if refs[0] != "decision:1" || refs[1] != "decision:2" {
		t.Errorf("got refs %v", refs)
	}
}

func TestParseAddedLearnings(t *testing.T) {
	diff := `+## [2026-03-14-100000] Always check for nil
 ## [2026-03-01-090000] Existing learning`

	refs := ParseAddedEntries(diff, "learning")
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0] != "learning:1" {
		t.Errorf("got ref %q, want %q", refs[0], "learning:1")
	}
}

func TestParseAddedTasks(t *testing.T) {
	diff := `+- [x] Implement auth handler #done:2026-03-14-100000
 - [ ] Write tests
+- [x] Add rate limiting #done:2026-03-14-110000`

	refs := ParseCompletedTasks(diff)
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
}

func TestParseNoAdditions(t *testing.T) {
	diff := ` ## [2026-03-01-090000] Existing entry
 - [ ] Existing task`

	refs := ParseAddedEntries(diff, "decision")
	if len(refs) != 0 {
		t.Errorf("expected 0 refs, got %d", len(refs))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -run TestParseAdded -v`
Expected: FAIL — functions not defined

- [ ] **Step 3: Implement staged.go**

Create `internal/trace/staged.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/regex"
)

// StagedRefs detects context references from staged .context/ files
// by examining git diff output.
//
// Parameters:
//   - contextDir: Path to the .context/ directory
//
// Returns:
//   - []string: Detected references (e.g., "decision:1", "task:3")
func StagedRefs(contextDir string) []string {
	var refs []string

	files := []struct {
		name     string
		entryType string
	}{
		{ctx.Decision, "decision"},
		{ctx.Learning, "learning"},
		{ctx.Convention, "convention"},
	}

	for _, f := range files {
		diff := stagedDiff(filepath.Join(contextDir, f.name))
		if diff == "" {
			continue
		}
		refs = append(refs, ParseAddedEntries(diff, f.entryType)...)
	}

	// Check TASKS.md for newly completed tasks
	taskDiff := stagedDiff(filepath.Join(contextDir, ctx.Task))
	if taskDiff != "" {
		refs = append(refs, ParseCompletedTasks(taskDiff)...)
	}

	return refs
}

// ParseAddedEntries extracts entry numbers from added lines in a diff.
// Only lines prefixed with "+" that match the entry header pattern are counted.
//
// Parameters:
//   - diff: Git diff output
//   - entryType: The reference type prefix ("decision", "learning", "convention")
//
// Returns:
//   - []string: Refs like "decision:1", "decision:2"
func ParseAddedEntries(diff, entryType string) []string {
	var refs []string
	count := 0

	for _, line := range strings.Split(diff, "\n") {
		if !strings.HasPrefix(line, "+") {
			continue
		}
		// Remove the leading "+" to match the regex
		content := line[1:]
		if regex.EntryHeader.MatchString(content) {
			count++
			refs = append(refs, fmt.Sprintf("%s:%d", entryType, count))
		}
	}

	return refs
}

// ParseCompletedTasks extracts task refs from newly completed tasks
// in a diff. Lines that are added ("+") and contain "[x]" are counted.
//
// Parameters:
//   - diff: Git diff output for TASKS.md
//
// Returns:
//   - []string: Refs like "task:1", "task:2"
func ParseCompletedTasks(diff string) []string {
	var refs []string
	count := 0

	for _, line := range strings.Split(diff, "\n") {
		if !strings.HasPrefix(line, "+") {
			continue
		}
		content := line[1:]
		match := regex.Task.FindStringSubmatch(content)
		if match != nil && (len(match) > 2 && match[2] == "x") {
			count++
			refs = append(refs, fmt.Sprintf("task:%d", count))
		}
	}

	return refs
}

// stagedDiff returns the staged diff for a specific file.
// Returns empty string if the file is not staged or git is not available.
func stagedDiff(filePath string) string {
	cmd := exec.Command("git", "diff", "--cached", "--", filePath)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -run "TestParseAdded|TestParseCompleted|TestParseNo" -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/trace/staged.go internal/trace/staged_test.go
git commit -m "feat(trace): add staged file analysis for context detection"
```

---

## Task 4: Working State Detection (Source 3)

**Files:**
- Create: `internal/trace/working.go`
- Create: `internal/trace/working_test.go`

### Steps

- [ ] **Step 1: Write failing test for working state**

Create `internal/trace/working_test.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorkingRefsInProgressTasks(t *testing.T) {
	tmpDir := t.TempDir()
	contextDir := tmpDir

	tasksContent := `# Tasks

- [ ] Implement auth handler
- [x] Write unit tests
- [ ] Add rate limiting
`
	if err := os.WriteFile(
		filepath.Join(contextDir, "TASKS.md"),
		[]byte(tasksContent), 0644,
	); err != nil {
		t.Fatal(err)
	}

	refs := WorkingRefs(contextDir)

	// Should find 2 in-progress tasks: task:1 and task:2
	found := map[string]bool{}
	for _, r := range refs {
		found[r] = true
	}

	if !found["task:1"] {
		t.Error("expected task:1 for 'Implement auth handler'")
	}
	if !found["task:2"] {
		t.Error("expected task:2 for 'Add rate limiting'")
	}
	if found["task:3"] {
		t.Error("should not find task:3 — completed tasks are excluded")
	}
}

func TestWorkingRefsSessionEnv(t *testing.T) {
	tmpDir := t.TempDir()
	contextDir := tmpDir

	// Write empty TASKS.md
	if err := os.WriteFile(
		filepath.Join(contextDir, "TASKS.md"),
		[]byte("# Tasks\n"), 0644,
	); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CTX_SESSION_ID", "test-session-42")

	refs := WorkingRefs(contextDir)

	found := false
	for _, r := range refs {
		if r == "session:test-session-42" {
			found = true
		}
	}
	if !found {
		t.Error("expected session:test-session-42 from env")
	}
}

func TestWorkingRefsNoTasksFile(t *testing.T) {
	tmpDir := t.TempDir()
	refs := WorkingRefs(tmpDir)

	// No TASKS.md should not panic, just return empty or session-only
	_ = refs
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -run TestWorkingRefs -v`
Expected: FAIL — function not defined

- [ ] **Step 3: Implement working.go**

Create `internal/trace/working.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/task"
)

const envSessionID = "CTX_SESSION_ID"

// WorkingRefs detects context references from the current working state.
// This includes in-progress tasks and the active AI session.
//
// Parameters:
//   - contextDir: Path to the .context/ directory
//
// Returns:
//   - []string: Detected references
func WorkingRefs(contextDir string) []string {
	var refs []string

	refs = append(refs, inProgressTaskRefs(contextDir)...)

	if sessionID := os.Getenv(envSessionID); sessionID != "" {
		refs = append(refs, "session:"+sessionID)
	}

	return refs
}

// inProgressTaskRefs reads TASKS.md and returns refs for in-progress
// (pending, non-subtask) tasks.
func inProgressTaskRefs(contextDir string) []string {
	tasksPath := filepath.Join(contextDir, ctxCfg.Task)
	content, err := os.ReadFile(filepath.Clean(tasksPath))
	if err != nil {
		return nil
	}

	var refs []string
	pendingCount := 0
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		match := regex.Task.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		if task.Sub(match) {
			continue // skip subtasks
		}
		if task.Pending(match) {
			pendingCount++
			refs = append(refs, fmt.Sprintf("task:%d", pendingCount))
		}
	}

	return refs
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -run TestWorkingRefs -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/trace/working.go internal/trace/working_test.go
git commit -m "feat(trace): add working state detection for in-progress tasks and sessions"
```

---

## Task 5: Collect — Merge and Deduplicate from All Sources

**Files:**
- Create: `internal/trace/collect.go`
- Create: `internal/trace/collect_test.go`

### Steps

- [ ] **Step 1: Write failing test for Collect**

Create `internal/trace/collect_test.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectDeduplicates(t *testing.T) {
	tmpDir := t.TempDir()
	contextDir := tmpDir
	stateDir := filepath.Join(contextDir, "state")
	if err := os.MkdirAll(stateDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Write TASKS.md with one in-progress task
	if err := os.WriteFile(
		filepath.Join(contextDir, "TASKS.md"),
		[]byte("# Tasks\n\n- [ ] Implement auth handler\n"), 0644,
	); err != nil {
		t.Fatal(err)
	}

	// Record the same task in pending context
	_ = Record("task:1", stateDir)
	// And a decision
	_ = Record("decision:5", stateDir)

	refs := Collect(contextDir)

	// task:1 appears in both pending and working state — should be deduplicated
	taskCount := 0
	decisionCount := 0
	for _, r := range refs {
		if r == "task:1" {
			taskCount++
		}
		if r == "decision:5" {
			decisionCount++
		}
	}

	if taskCount != 1 {
		t.Errorf("task:1 should appear exactly once, got %d", taskCount)
	}
	if decisionCount != 1 {
		t.Errorf("decision:5 should appear exactly once, got %d", decisionCount)
	}
}

func TestCollectEmptyReturnsNil(t *testing.T) {
	tmpDir := t.TempDir()
	stateDir := filepath.Join(tmpDir, "state")
	if err := os.MkdirAll(stateDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Empty TASKS.md, no pending context
	if err := os.WriteFile(
		filepath.Join(tmpDir, "TASKS.md"),
		[]byte("# Tasks\n"), 0644,
	); err != nil {
		t.Fatal(err)
	}

	refs := Collect(tmpDir)
	if len(refs) != 0 {
		t.Errorf("expected empty refs, got %v", refs)
	}
}

func TestFormatTrailer(t *testing.T) {
	refs := []string{"decision:12", "task:8", "session:abc123"}
	trailer := FormatTrailer(refs)
	want := "ctx-context: decision:12, task:8, session:abc123"
	if trailer != want {
		t.Errorf("got %q, want %q", trailer, want)
	}
}

func TestFormatTrailerEmpty(t *testing.T) {
	trailer := FormatTrailer(nil)
	if trailer != "" {
		t.Errorf("expected empty trailer, got %q", trailer)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -run TestCollect -v`
Expected: FAIL — functions not defined

- [ ] **Step 3: Implement collect.go**

Create `internal/trace/collect.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/dir"
)

const trailerKey = "ctx-context"

// Collect gathers context references from all three sources
// (pending, staged, working state), merges, and deduplicates them.
//
// Parameters:
//   - contextDir: Path to the .context/ directory
//
// Returns:
//   - []string: Deduplicated context references
func Collect(contextDir string) []string {
	stateDir := filepath.Join(contextDir, dir.State)

	var all []string

	// Source 1: Pending context
	pending, _ := ReadPending(stateDir)
	for _, p := range pending {
		all = append(all, p.Ref)
	}

	// Source 2: Staged file analysis
	all = append(all, StagedRefs(contextDir)...)

	// Source 3: Current working state
	all = append(all, WorkingRefs(contextDir)...)

	return deduplicate(all)
}

// FormatTrailer formats refs as a git commit trailer string.
// Returns empty string if refs is empty.
//
// Parameters:
//   - refs: Context references to include
//
// Returns:
//   - string: Formatted trailer line (e.g., "ctx-context: decision:12, task:8")
func FormatTrailer(refs []string) string {
	if len(refs) == 0 {
		return ""
	}
	return trailerKey + ": " + strings.Join(refs, ", ")
}

// deduplicate removes duplicate refs while preserving order.
func deduplicate(refs []string) []string {
	seen := make(map[string]bool, len(refs))
	var result []string
	for _, r := range refs {
		if !seen[r] {
			seen[r] = true
			result = append(result, r)
		}
	}
	return result
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/trace/collect.go internal/trace/collect_test.go
git commit -m "feat(trace): add three-source collection with deduplication"
```

---

## Task 6: Reference Resolution

**Files:**
- Create: `internal/trace/resolve.go`
- Create: `internal/trace/resolve_test.go`

### Steps

- [ ] **Step 1: Write failing test for resolution**

Create `internal/trace/resolve_test.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseRef(t *testing.T) {
	tests := []struct {
		input    string
		wantType string
		wantNum  int
		wantText string
	}{
		{"decision:12", "decision", 12, ""},
		{"learning:5", "learning", 5, ""},
		{"task:8", "task", 8, ""},
		{"convention:3", "convention", 3, ""},
		{"session:abc123", "session", 0, "abc123"},
		{`"Hotfix for prod outage"`, "note", 0, "Hotfix for prod outage"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			typ, num, text := ParseRef(tt.input)
			if typ != tt.wantType {
				t.Errorf("type: got %q, want %q", typ, tt.wantType)
			}
			if num != tt.wantNum {
				t.Errorf("number: got %d, want %d", num, tt.wantNum)
			}
			if text != tt.wantText {
				t.Errorf("text: got %q, want %q", text, tt.wantText)
			}
		})
	}
}

func TestResolveDecision(t *testing.T) {
	tmpDir := t.TempDir()
	contextDir := tmpDir

	decisionsContent := `# Decisions

## [2026-03-10-100000] Use short-lived tokens

**Context:** Security review needed a token strategy.

**Rationale:** Short-lived tokens reduce blast radius of token theft.

**Consequences:** Need refresh token handling.

## [2026-03-01-090000] Use PostgreSQL

**Context:** Database selection.

**Rationale:** Well-supported.

**Consequences:** Team needs training.
`
	if err := os.WriteFile(
		filepath.Join(contextDir, "DECISIONS.md"),
		[]byte(decisionsContent), 0644,
	); err != nil {
		t.Fatal(err)
	}

	resolved := Resolve("decision:1", contextDir)
	if !resolved.Found {
		t.Fatal("expected to resolve decision:1")
	}
	if resolved.Title != "Use short-lived tokens" {
		t.Errorf("title: got %q", resolved.Title)
	}
	if resolved.Type != "decision" {
		t.Errorf("type: got %q", resolved.Type)
	}
}

func TestResolveTask(t *testing.T) {
	tmpDir := t.TempDir()
	contextDir := tmpDir

	tasksContent := `# Tasks

- [ ] Implement auth handler
- [x] Write unit tests
- [ ] Add rate limiting
`
	if err := os.WriteFile(
		filepath.Join(contextDir, "TASKS.md"),
		[]byte(tasksContent), 0644,
	); err != nil {
		t.Fatal(err)
	}

	resolved := Resolve("task:1", contextDir)
	if !resolved.Found {
		t.Fatal("expected to resolve task:1")
	}
	if resolved.Title != "Implement auth handler" {
		t.Errorf("title: got %q", resolved.Title)
	}
}

func TestResolveNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	resolved := Resolve("decision:999", tmpDir)
	if resolved.Found {
		t.Error("should not resolve decision:999")
	}
}

func TestResolveNote(t *testing.T) {
	tmpDir := t.TempDir()
	resolved := Resolve(`"Hotfix for production outage"`, tmpDir)
	if !resolved.Found {
		t.Fatal("notes should always resolve")
	}
	if resolved.Title != "Hotfix for production outage" {
		t.Errorf("title: got %q", resolved.Title)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -run TestParseRef -v`
Expected: FAIL — functions not defined

- [ ] **Step 3: Implement resolve.go**

Create `internal/trace/resolve.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/task"
)

// ParseRef breaks a reference string into its type, number, and text.
//
// Examples:
//   - "decision:12" → ("decision", 12, "")
//   - "session:abc" → ("session", 0, "abc")
//   - "\"Some note\"" → ("note", 0, "Some note")
//
// Parameters:
//   - ref: Raw reference string
//
// Returns:
//   - refType: "decision", "learning", "task", "convention", "session", or "note"
//   - number: Entry number (0 for session/note)
//   - text: Session ID or note text (empty for numbered entries)
func ParseRef(ref string) (refType string, number int, text string) {
	// Check for quoted free-form note
	if strings.HasPrefix(ref, `"`) && strings.HasSuffix(ref, `"`) {
		return "note", 0, strings.Trim(ref, `"`)
	}

	parts := strings.SplitN(ref, ":", 2)
	if len(parts) != 2 {
		return "note", 0, ref
	}

	refType = parts[0]
	value := parts[1]

	if num, err := strconv.Atoi(value); err == nil {
		return refType, num, ""
	}

	return refType, 0, value
}

// Resolve looks up a reference and returns its resolved form.
//
// Parameters:
//   - ref: Raw reference string
//   - contextDir: Path to the .context/ directory
//
// Returns:
//   - ResolvedRef: Resolved reference with title and detail
func Resolve(ref, contextDir string) ResolvedRef {
	refType, number, text := ParseRef(ref)

	resolved := ResolvedRef{
		Raw:    ref,
		Type:   refType,
		Number: number,
	}

	switch refType {
	case "decision":
		return resolveEntry(resolved, contextDir, ctxCfg.Decision, number)
	case "learning":
		return resolveEntry(resolved, contextDir, ctxCfg.Learning, number)
	case "convention":
		return resolveEntry(resolved, contextDir, ctxCfg.Convention, number)
	case "task":
		return resolveTask(resolved, contextDir, number)
	case "session":
		resolved.Title = text
		resolved.Found = true
		return resolved
	case "note":
		resolved.Title = text
		resolved.Found = true
		return resolved
	default:
		resolved.Title = ref
		return resolved
	}
}

// resolveEntry resolves a numbered entry from a context file
// (DECISIONS.md, LEARNINGS.md, CONVENTIONS.md).
func resolveEntry(resolved ResolvedRef, contextDir, fileName string, number int) ResolvedRef {
	filePath := filepath.Join(contextDir, fileName)
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return resolved
	}

	entries := index.ParseHeaders(string(content))
	if number < 1 || number > len(entries) {
		return resolved
	}

	entry := entries[number-1]
	resolved.Title = entry.Title
	resolved.Detail = fmt.Sprintf("Date: %s", entry.Date)
	resolved.Found = true

	return resolved
}

// resolveTask resolves a task number from TASKS.md.
// Task numbers count only top-level pending tasks in file order.
func resolveTask(resolved ResolvedRef, contextDir string, number int) ResolvedRef {
	filePath := filepath.Join(contextDir, ctxCfg.Task)
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return resolved
	}

	lines := strings.Split(string(content), "\n")
	count := 0

	for _, line := range lines {
		match := regex.Task.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		// Count all top-level tasks (both pending and completed)
		if !task.Sub(match) {
			count++
			if count == number {
				resolved.Title = task.Content(match)
				if task.Completed(match) {
					resolved.Detail = "Status: completed"
				} else {
					resolved.Detail = "Status: pending"
				}
				resolved.Found = true
				return resolved
			}
		}
	}

	return resolved
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/parlakisik/projects/github/ctx && go test ./internal/trace/ -run "TestParseRef|TestResolve" -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/trace/resolve.go internal/trace/resolve_test.go
git commit -m "feat(trace): add reference resolution from context files"
```

---

## Task 7: Wire Recording into Existing Commands

**Files:**
- Modify: `internal/cli/add/cmd/root/run.go`
- Modify: `internal/cli/task/cmd/complete/run.go`

### Steps

- [ ] **Step 1: Understand entry numbering for add command**

The `ctx add` command does not return an entry number. Since entries are prepended (newest first), a newly added decision is always entry #1 in the file. We need to count entries after write to determine the new entry's number.

- [ ] **Step 2: Modify add command to record pending context**

In `internal/cli/add/cmd/root/run.go`, add the trace recording after the successful write. The entry number is determined by counting entries in the file after write, and the new entry is always #1 (prepended for decisions/learnings) or the last entry (appended for tasks/conventions).

Add import:

```go
"github.com/ActiveMemory/ctx/internal/trace"
"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
```

After `writeAdd.Added(cmd, fName)` and before `return nil`, add:

```go
	// Best-effort: record pending context for commit tracing.
	// Decisions and learnings are prepended (newest = #1).
	// Tasks and conventions are appended (newest = last).
	if fType == cfgEntry.Decision || fType == cfgEntry.Learning ||
		fType == cfgEntry.Convention {
		_ = trace.Record(fType+":1", state.Dir())
	}
```

Note: We record as entry #1 for prepended types because new entries are always inserted at the top. For tasks, recording happens in the `complete` command instead, since tasks are tracked by completion, not creation.

- [ ] **Step 3: Modify complete command to record pending context**

In `internal/cli/task/cmd/complete/run.go`, add trace recording after a successful completion.

Add import:

```go
"github.com/ActiveMemory/ctx/internal/trace"
"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
```

In the `Run` function, after `complete.Completed(cmd, matchedTask)` and before `return nil`, add:

```go
	// Best-effort: record pending context for commit tracing.
	_ = trace.Record("task:"+args[0], state.Dir())
```

- [ ] **Step 4: Run existing tests to verify no regressions**

Run: `cd /Users/parlakisik/projects/github/ctx && CTX_SKIP_PATH_CHECK=1 go test ./internal/cli/add/ ./internal/cli/task/... -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/cli/add/cmd/root/run.go internal/cli/task/cmd/complete/run.go
git commit -m "feat(trace): wire pending context recording into add and complete commands"
```

---

## Task 8: CLI — `ctx trace` Command Structure

**Files:**
- Create: `internal/cli/trace/trace.go`
- Create: `internal/cli/trace/cmd/show/cmd.go`
- Create: `internal/cli/trace/cmd/show/run.go`
- Create: `internal/config/embed/cmd/trace.go`
- Modify: `internal/config/embed/cmd/base.go`
- Modify: `internal/bootstrap/group.go`
- Modify: `internal/assets/commands/commands.yaml`

### Steps

- [ ] **Step 1: Add trace Use and DescKey constants**

Create `internal/config/embed/cmd/trace.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

const UseTrace = "trace [commit]"

const (
	DescKeyTrace        = "trace"
	DescKeyTraceFile    = "trace.file"
	DescKeyTraceTag     = "trace.tag"
	DescKeyTraceCollect = "trace.collect"
	DescKeyTraceHook    = "trace.hook"
)
```

- [ ] **Step 2: Add trace command descriptions to commands.yaml**

In `internal/assets/commands/commands.yaml`, add the trace command descriptions:

```yaml
trace:
  long: |-
    Show the context behind git commits.

    ctx trace links commits back to the decisions, tasks, learnings,
    and sessions that motivated them.

    Usage:
      ctx trace <commit>         Show context for a specific commit
      ctx trace --last 5         Show context for last N commits
      ctx trace file <path>      Show context trail for a file
      ctx trace tag <commit>     Manually tag a commit with context
      ctx trace collect          Collect context refs (used by hook)
      ctx trace hook enable      Install prepare-commit-msg hook

    Examples:
      ctx trace abc123
      ctx trace --last 10
      ctx trace file src/auth.go
      ctx trace tag HEAD --note "Hotfix for production outage"
  short: Show context behind git commits
trace.file:
  long: |-
    Show the context trail for a file.

    Combines git log with trailer resolution to show what decisions,
    tasks, and learnings motivated changes to a specific file.

    Supports optional line range with colon syntax:
      ctx trace file src/auth.go:42-60

    Examples:
      ctx trace file src/auth.go
      ctx trace file src/auth.go:42-60
  short: Show context trail for a file
trace.tag:
  long: |-
    Manually tag a commit with context.

    For commits made without the hook, or to add extra context
    after the fact. Tags are stored in .context/trace/overrides.jsonl
    since git trailers cannot be modified without rewriting history.

    Examples:
      ctx trace tag HEAD --note "Hotfix for production outage"
      ctx trace tag abc123 --note "Part of Q1 compliance initiative"
  short: Manually tag a commit with context
trace.collect:
  long: |-
    Collect context references from all sources.

    Gathers pending context, staged file analysis, and working state,
    then outputs a ctx-context trailer line. Used by the
    prepare-commit-msg hook.

    This command is not typically called directly.
  short: Collect context refs for hook
trace.hook:
  long: |-
    Enable or disable the prepare-commit-msg hook for automatic
    context tracing. The hook injects ctx-context trailers into
    commit messages.

    Usage:
      ctx trace hook enable     Install the hook
      ctx trace hook disable    Remove the hook

    Examples:
      ctx trace hook enable
      ctx trace hook disable
  short: Manage prepare-commit-msg hook
```

- [ ] **Step 3: Create show subcommand (ctx trace [commit] / ctx trace --last N)**

Create `internal/cli/trace/cmd/show/cmd.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the trace show command (the default action for ctx trace).
//
// Returns:
//   - *cobra.Command: Configured trace command
func Cmd() *cobra.Command {
	var (
		last       int
		jsonOutput bool
	)

	short, long := desc.Command(cmd.DescKeyTrace)

	c := &cobra.Command{
		Use:   cmd.UseTrace,
		Short: short,
		Long:  long,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, last, jsonOutput)
		},
	}

	c.Flags().IntVar(&last, cFlag.Last, 0, "Show context for last N commits")
	c.Flags().BoolVar(&jsonOutput, cFlag.JSON, false, "Output as JSON")

	return c
}
```

- [ ] **Step 4: Create show run logic**

Create `internal/cli/trace/cmd/show/run.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trace"
)

// Run executes the trace show command.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Optional commit hash as first argument
//   - last: Number of recent commits to show (0 = disabled)
//   - jsonOutput: Whether to output as JSON
//
// Returns:
//   - error: Non-nil on failure
func Run(cmd *cobra.Command, args []string, last int, jsonOutput bool) error {
	contextDir := rc.ContextDir()
	traceDir := filepath.Join(contextDir, dir.Trace)

	if last > 0 {
		return showLast(cmd, last, contextDir, traceDir, jsonOutput)
	}

	if len(args) == 0 {
		return showLast(cmd, 10, contextDir, traceDir, jsonOutput)
	}

	return showCommit(cmd, args[0], contextDir, traceDir, jsonOutput)
}

func showCommit(cmd *cobra.Command, commitHash, contextDir, traceDir string, jsonOutput bool) error {
	// Resolve full hash
	fullHash := resolveCommitHash(commitHash)
	if fullHash == "" {
		fullHash = commitHash
	}

	// Collect refs from all sources
	refs := collectRefsForCommit(fullHash, traceDir)

	if len(refs) == 0 {
		cmd.Printf("Commit: %s\n\nContext: (none)\n", shortHash(fullHash))
		return nil
	}

	if jsonOutput {
		return outputJSON(cmd, fullHash, refs, contextDir)
	}

	// Get commit message
	message := commitMessage(fullHash)
	date := commitDate(fullHash)

	cmd.Printf("Commit: %s %q\n", shortHash(fullHash), message)
	if date != "" {
		cmd.Printf("Date:   %s\n", date)
	}
	cmd.Println()
	cmd.Println("Context:")

	for _, ref := range refs {
		resolved := trace.Resolve(ref, contextDir)
		printResolved(cmd, resolved)
	}

	return nil
}

func showLast(cmd *cobra.Command, n int, contextDir, traceDir string, jsonOutput bool) error {
	// Get last N commit hashes
	out, err := exec.Command("git", "log", fmt.Sprintf("-%d", n), "--format=%H %s").Output()
	if err != nil {
		return fmt.Errorf("git log: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		cmd.Println("No commits found.")
		return nil
	}

	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		hash := parts[0]
		message := parts[1]

		refs := collectRefsForCommit(hash, traceDir)

		if len(refs) > 0 {
			cmd.Printf("%s  %-40s \u2192 %s\n", shortHash(hash), message, strings.Join(refs, ", "))
		} else {
			cmd.Printf("%s  %-40s   (no context)\n", shortHash(hash), message)
		}
	}

	return nil
}

func collectRefsForCommit(commitHash, traceDir string) []string {
	var allRefs []string

	// Source 1: history.jsonl (primary)
	entry, found := trace.ReadHistoryForCommit(commitHash, traceDir)
	if found {
		allRefs = append(allRefs, entry.Refs...)
	}

	// Source 2: git trailer
	allRefs = append(allRefs, readTrailerRefs(commitHash)...)

	// Source 3: overrides.jsonl
	allRefs = append(allRefs, trace.ReadOverridesForCommit(commitHash, traceDir)...)

	// Deduplicate
	seen := make(map[string]bool, len(allRefs))
	var result []string
	for _, r := range allRefs {
		if !seen[r] {
			seen[r] = true
			result = append(result, r)
		}
	}
	return result
}

func readTrailerRefs(commitHash string) []string {
	out, err := exec.Command("git", "log", "-1", "--format=%(trailers:key=ctx-context,valueonly)", commitHash).Output()
	if err != nil {
		return nil
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil
	}

	var refs []string
	for _, part := range strings.Split(raw, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			refs = append(refs, trimmed)
		}
	}
	return refs
}

func printResolved(cmd *cobra.Command, r trace.ResolvedRef) {
	prefix := strings.Title(r.Type)
	if r.Number > 0 {
		prefix = fmt.Sprintf("%s #%d", prefix, r.Number)
	}

	if r.Found {
		cmd.Printf("  %s: %s\n", prefix, r.Title)
		if r.Detail != "" {
			cmd.Printf("    %s\n", r.Detail)
		}
	} else {
		cmd.Printf("  %s: [not found \u2014 may have been archived]\n", prefix)
	}
	cmd.Println()
}

func outputJSON(cmd *cobra.Command, hash string, refs []string, contextDir string) error {
	type jsonRef struct {
		Raw      string `json:"raw"`
		Type     string `json:"type"`
		Number   int    `json:"number,omitempty"`
		Title    string `json:"title,omitempty"`
		Detail   string `json:"detail,omitempty"`
		Found    bool   `json:"found"`
	}

	type jsonOutput struct {
		Commit  string    `json:"commit"`
		Message string    `json:"message"`
		Refs    []jsonRef `json:"refs"`
	}

	var jRefs []jsonRef
	for _, ref := range refs {
		resolved := trace.Resolve(ref, contextDir)
		jRefs = append(jRefs, jsonRef{
			Raw:    resolved.Raw,
			Type:   resolved.Type,
			Number: resolved.Number,
			Title:  resolved.Title,
			Detail: resolved.Detail,
			Found:  resolved.Found,
		})
	}

	out := jsonOutput{
		Commit:  hash,
		Message: commitMessage(hash),
		Refs:    jRefs,
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	cmd.Println(string(data))
	return nil
}

func resolveCommitHash(short string) string {
	out, err := exec.Command("git", "rev-parse", short).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func commitMessage(hash string) string {
	out, err := exec.Command("git", "log", "-1", "--format=%s", hash).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func commitDate(hash string) string {
	out, err := exec.Command("git", "log", "-1", "--format=%ci", hash).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func shortHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}
```

- [ ] **Step 5: Create trace.go top-level command**

Create `internal/cli/trace/trace.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace provides the ctx trace CLI command for commit context tracing.
package trace

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/show"
)

// Cmd returns the trace command with all subcommands.
//
// Returns:
//   - *cobra.Command: The trace command
func Cmd() *cobra.Command {
	return show.Cmd()
}
```

- [ ] **Step 6: Register trace command in bootstrap**

In `internal/bootstrap/group.go`, add the import:

```go
"github.com/ActiveMemory/ctx/internal/cli/trace"
```

Add `{trace.Cmd, embedCmd.GroupDiagnostics}` to the `diagnostics()` function return slice.

- [ ] **Step 7: Run build to verify compilation**

Run: `cd /Users/parlakisik/projects/github/ctx && go build ./cmd/ctx/`
Expected: BUILD SUCCESS

- [ ] **Step 8: Commit**

```bash
git add internal/cli/trace/ internal/config/embed/cmd/trace.go internal/bootstrap/group.go internal/assets/commands/commands.yaml
git commit -m "feat(trace): add ctx trace command for querying commit context"
```

---

## Task 9: CLI — `ctx trace file` Subcommand

**Files:**
- Create: `internal/cli/trace/cmd/file/cmd.go`
- Create: `internal/cli/trace/cmd/file/run.go`
- Modify: `internal/cli/trace/trace.go`

### Steps

- [ ] **Step 1: Create file subcommand**

Create `internal/cli/trace/cmd/file/cmd.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the trace file subcommand.
//
// Returns:
//   - *cobra.Command: Configured trace file command
func Cmd() *cobra.Command {
	var last int

	short, long := desc.Command(cmd.DescKeyTraceFile)

	c := &cobra.Command{
		Use:   "file <path[:line-range]>",
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], last)
		},
	}

	c.Flags().IntVarP(&last, cFlag.Last, cFlag.ShortLast, 20, "Max commits to show")

	return c
}
```

- [ ] **Step 2: Create file run logic**

Create `internal/cli/trace/cmd/file/run.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trace"
)

// Run executes the trace file command.
//
// Parameters:
//   - cmd: Cobra command for output
//   - pathArg: File path, optionally with line range (e.g., "src/auth.go:42-60")
//   - last: Maximum number of commits to show
//
// Returns:
//   - error: Non-nil on failure
func Run(cmd *cobra.Command, pathArg string, last int) error {
	filePath, lineRange := parsePathArg(pathArg)
	contextDir := rc.ContextDir()
	traceDir := filepath.Join(contextDir, dir.Trace)

	// Get commit hashes for this file
	gitArgs := []string{"log", fmt.Sprintf("-%d", last), "--format=%H %ci %s"}
	if lineRange != "" {
		// Use -L for line ranges
		gitArgs = []string{"log", fmt.Sprintf("-%d", last), "--format=%H %ci %s", "-L", lineRange + ":" + filePath}
	} else {
		gitArgs = append(gitArgs, "--", filePath)
	}

	out, err := exec.Command("git", gitArgs...).Output()
	if err != nil {
		return fmt.Errorf("git log for %s: %w", filePath, err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		cmd.Printf("No commits found for %s\n", pathArg)
		return nil
	}

	for _, line := range lines {
		parts := strings.SplitN(line, " ", 4)
		if len(parts) < 4 {
			continue
		}
		hash := parts[0]
		date := parts[1]
		message := parts[3]

		refs := collectRefsForCommit(hash, traceDir)
		if len(refs) > 0 {
			cmd.Printf("%s  %s  %-35s \u2192 %s\n", shortHash(hash), date, message, strings.Join(refs, ", "))
		} else {
			cmd.Printf("%s  %s  %-35s   (no context)\n", shortHash(hash), date, message)
		}
	}

	return nil
}

func parsePathArg(arg string) (path, lineRange string) {
	// Check for path:line-range format (e.g., "src/auth.go:42-60")
	lastColon := strings.LastIndex(arg, ":")
	if lastColon == -1 {
		return arg, ""
	}

	potential := arg[lastColon+1:]
	if strings.Contains(potential, "-") || isNumeric(potential) {
		return arg[:lastColon], potential
	}

	return arg, ""
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func collectRefsForCommit(commitHash, traceDir string) []string {
	var allRefs []string

	entry, found := trace.ReadHistoryForCommit(commitHash, traceDir)
	if found {
		allRefs = append(allRefs, entry.Refs...)
	}

	allRefs = append(allRefs, trace.ReadOverridesForCommit(commitHash, traceDir)...)

	seen := make(map[string]bool, len(allRefs))
	var result []string
	for _, r := range allRefs {
		if !seen[r] {
			seen[r] = true
			result = append(result, r)
		}
	}
	return result
}

func shortHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}
```

- [ ] **Step 3: Register file subcommand**

Update `internal/cli/trace/trace.go` to add the file subcommand:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace provides the ctx trace CLI command for commit context tracing.
package trace

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/file"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/show"
)

// Cmd returns the trace command with all subcommands.
//
// Returns:
//   - *cobra.Command: The trace command
func Cmd() *cobra.Command {
	c := show.Cmd()
	c.AddCommand(file.Cmd())
	return c
}
```

- [ ] **Step 4: Run build to verify compilation**

Run: `cd /Users/parlakisik/projects/github/ctx && go build ./cmd/ctx/`
Expected: BUILD SUCCESS

- [ ] **Step 5: Commit**

```bash
git add internal/cli/trace/cmd/file/ internal/cli/trace/trace.go
git commit -m "feat(trace): add ctx trace file subcommand for file history"
```

---

## Task 10: CLI — `ctx trace tag` Subcommand

**Files:**
- Create: `internal/cli/trace/cmd/tag/cmd.go`
- Create: `internal/cli/trace/cmd/tag/run.go`
- Modify: `internal/cli/trace/trace.go`

### Steps

- [ ] **Step 1: Create tag subcommand**

Create `internal/cli/trace/cmd/tag/cmd.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tag

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the trace tag subcommand.
//
// Returns:
//   - *cobra.Command: Configured trace tag command
func Cmd() *cobra.Command {
	var note string

	short, long := desc.Command(cmd.DescKeyTraceTag)

	c := &cobra.Command{
		Use:   "tag <commit>",
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], note)
		},
	}

	c.Flags().StringVar(&note, "note", "", "Free-form context note")

	return c
}
```

- [ ] **Step 2: Create tag run logic**

Create `internal/cli/trace/cmd/tag/run.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tag

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trace"
)

// Run executes the trace tag command.
//
// Parameters:
//   - cmd: Cobra command for output
//   - commitRef: Git commit reference (hash or "HEAD")
//   - note: Free-form context note
//
// Returns:
//   - error: Non-nil on failure
func Run(cmd *cobra.Command, commitRef, note string) error {
	if note == "" {
		return fmt.Errorf("--note is required")
	}

	// Resolve commit hash
	hash, err := resolveHash(commitRef)
	if err != nil {
		return fmt.Errorf("cannot resolve %q: %w", commitRef, err)
	}

	contextDir := rc.ContextDir()
	traceDir := filepath.Join(contextDir, dir.Trace)

	entry := trace.OverrideEntry{
		Commit: hash,
		Refs:   []string{fmt.Sprintf("%q", note)},
	}

	if writeErr := trace.WriteOverride(entry, traceDir); writeErr != nil {
		return fmt.Errorf("write override: %w", writeErr)
	}

	cmd.Printf("Tagged %s with: %s\n", shortHash(hash), note)
	return nil
}

func resolveHash(ref string) (string, error) {
	out, err := exec.Command("git", "rev-parse", ref).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func shortHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}
```

- [ ] **Step 3: Register tag subcommand in trace.go**

Update `internal/cli/trace/trace.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace provides the ctx trace CLI command for commit context tracing.
package trace

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/file"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/show"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/tag"
)

// Cmd returns the trace command with all subcommands.
//
// Returns:
//   - *cobra.Command: The trace command
func Cmd() *cobra.Command {
	c := show.Cmd()
	c.AddCommand(file.Cmd())
	c.AddCommand(tag.Cmd())
	return c
}
```

- [ ] **Step 4: Run build to verify compilation**

Run: `cd /Users/parlakisik/projects/github/ctx && go build ./cmd/ctx/`
Expected: BUILD SUCCESS

- [ ] **Step 5: Commit**

```bash
git add internal/cli/trace/cmd/tag/ internal/cli/trace/trace.go
git commit -m "feat(trace): add ctx trace tag subcommand for manual commit tagging"
```

---

## Task 11: CLI — `ctx trace collect` and `ctx trace hook`

**Files:**
- Create: `internal/cli/trace/cmd/collect/cmd.go`
- Create: `internal/cli/trace/cmd/collect/run.go`
- Create: `internal/cli/trace/cmd/hook/cmd.go`
- Create: `internal/cli/trace/cmd/hook/run.go`
- Modify: `internal/cli/trace/trace.go`

### Steps

- [ ] **Step 1: Create collect subcommand**

Create `internal/cli/trace/cmd/collect/cmd.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package collect

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the trace collect subcommand.
//
// Returns:
//   - *cobra.Command: Configured trace collect command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyTraceCollect)

	c := &cobra.Command{
		Use:    "collect",
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd)
		},
	}

	return c
}
```

- [ ] **Step 2: Create collect run logic**

Create `internal/cli/trace/cmd/collect/run.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package collect

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trace"
)

// Run executes the trace collect command.
// Outputs the ctx-context trailer line to stdout for the hook to consume.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on failure
func Run(cmd *cobra.Command) error {
	contextDir := rc.ContextDir()
	refs := trace.Collect(contextDir)

	trailer := trace.FormatTrailer(refs)
	if trailer != "" {
		cmd.Println(trailer)
	}

	return nil
}
```

- [ ] **Step 3: Create hook subcommand**

Create `internal/cli/trace/cmd/hook/cmd.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the trace hook subcommand.
//
// Returns:
//   - *cobra.Command: Configured trace hook command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyTraceHook)

	c := &cobra.Command{
		Use:   "hook <enable|disable>",
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0])
		},
	}

	return c
}
```

- [ ] **Step 4: Create hook run logic**

Create `internal/cli/trace/cmd/hook/run.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
)

const hookScript = `#!/bin/sh
# ctx: prepare-commit-msg hook for commit context tracing.
# Installed by: ctx trace hook enable
# Remove with:  ctx trace hook disable

COMMIT_MSG_FILE="$1"
COMMIT_SOURCE="$2"

# Only inject on normal commits (not merges, squashes, or amends)
case "$COMMIT_SOURCE" in
  merge|squash) exit 0 ;;
esac

# Collect context refs
TRAILER=$(ctx trace collect 2>/dev/null)

if [ -n "$TRAILER" ]; then
  # Append trailer with a blank line separator
  echo "" >> "$COMMIT_MSG_FILE"
  echo "$TRAILER" >> "$COMMIT_MSG_FILE"
fi
`

// Run executes the trace hook command.
//
// Parameters:
//   - cmd: Cobra command for output
//   - action: "enable" or "disable"
//
// Returns:
//   - error: Non-nil on failure
func Run(cmd *cobra.Command, action string) error {
	switch strings.ToLower(action) {
	case "enable":
		return enable(cmd)
	case "disable":
		return disable(cmd)
	default:
		return fmt.Errorf("unknown action %q: use 'enable' or 'disable'", action)
	}
}

func enable(cmd *cobra.Command) error {
	hookPath, err := hookFilePath()
	if err != nil {
		return err
	}

	// Check if hook already exists
	if _, statErr := os.Stat(hookPath); statErr == nil {
		content, readErr := os.ReadFile(filepath.Clean(hookPath))
		if readErr == nil && strings.Contains(string(content), "ctx trace collect") {
			cmd.Println("Hook already installed.")
			return nil
		}
		return fmt.Errorf("a prepare-commit-msg hook already exists at %s; remove it first or add ctx integration manually", hookPath)
	}

	if writeErr := os.WriteFile(hookPath, []byte(hookScript), fs.PermExec); writeErr != nil {
		return fmt.Errorf("write hook: %w", writeErr)
	}

	cmd.Printf("Installed prepare-commit-msg hook at %s\n", hookPath)
	return nil
}

func disable(cmd *cobra.Command) error {
	hookPath, err := hookFilePath()
	if err != nil {
		return err
	}

	if _, statErr := os.Stat(hookPath); os.IsNotExist(statErr) {
		cmd.Println("No hook installed.")
		return nil
	}

	// Verify it's our hook before removing
	content, readErr := os.ReadFile(filepath.Clean(hookPath))
	if readErr != nil {
		return fmt.Errorf("read hook: %w", readErr)
	}

	if !strings.Contains(string(content), "ctx trace collect") {
		return fmt.Errorf("hook at %s is not a ctx trace hook; not removing", hookPath)
	}

	if removeErr := os.Remove(hookPath); removeErr != nil {
		return fmt.Errorf("remove hook: %w", removeErr)
	}

	cmd.Printf("Removed prepare-commit-msg hook from %s\n", hookPath)
	return nil
}

func hookFilePath() (string, error) {
	// Get git hooks directory
	out, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}

	gitDir := strings.TrimSpace(string(out))
	hooksDir := filepath.Join(gitDir, "hooks")

	if mkdirErr := os.MkdirAll(hooksDir, fs.PermExec); mkdirErr != nil {
		return "", fmt.Errorf("create hooks dir: %w", mkdirErr)
	}

	return filepath.Join(hooksDir, "prepare-commit-msg"), nil
}
```

- [ ] **Step 5: Register collect and hook subcommands in trace.go**

Update `internal/cli/trace/trace.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace provides the ctx trace CLI command for commit context tracing.
package trace

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/collect"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/file"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/hook"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/show"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/tag"
)

// Cmd returns the trace command with all subcommands.
//
// Returns:
//   - *cobra.Command: The trace command
func Cmd() *cobra.Command {
	c := show.Cmd()
	c.AddCommand(collect.Cmd())
	c.AddCommand(file.Cmd())
	c.AddCommand(hook.Cmd())
	c.AddCommand(tag.Cmd())
	return c
}
```

- [ ] **Step 6: Run build to verify compilation**

Run: `cd /Users/parlakisik/projects/github/ctx && go build ./cmd/ctx/`
Expected: BUILD SUCCESS

- [ ] **Step 7: Commit**

```bash
git add internal/cli/trace/cmd/collect/ internal/cli/trace/cmd/hook/ internal/cli/trace/trace.go
git commit -m "feat(trace): add collect and hook subcommands for prepare-commit-msg integration"
```

---

## Task 12: Error Package and Output Formatting

**Files:**
- Create: `internal/err/trace/doc.go`
- Create: `internal/err/trace/trace.go`

### Steps

- [ ] **Step 1: Create trace error package**

Create `internal/err/trace/doc.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace provides error constructors for trace operations.
package trace
```

Create `internal/err/trace/trace.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"errors"
	"fmt"
)

// CommitNotFound returns an error when a commit hash cannot be found.
//
// Parameters:
//   - hash: The commit hash that was not found
//
// Returns:
//   - error: Descriptive error
func CommitNotFound(hash string) error {
	return fmt.Errorf("commit not found: %s", hash)
}

// NotInGitRepo returns an error when the command is run outside a git repo.
//
// Returns:
//   - error: Descriptive error
func NotInGitRepo() error {
	return errors.New("not in a git repository")
}

// NoteRequired returns an error when --note flag is missing.
//
// Returns:
//   - error: Descriptive error
func NoteRequired() error {
	return errors.New("--note is required")
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/err/trace/
git commit -m "feat(trace): add error package for trace operations"
```

---

## Task 13: Integration Tests

**Files:**
- Create: `internal/cli/trace/trace_test.go`

### Steps

- [ ] **Step 1: Write integration test**

Create `internal/cli/trace/trace_test.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/trace"
)

func TestTraceTagAndShow(t *testing.T) {
	tmpDir := t.TempDir()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Init git repo
	run(t, "git", "init")
	run(t, "git", "config", "user.email", "test@test.com")
	run(t, "git", "config", "user.name", "Test")

	// Init ctx
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Create a file and commit
	if err := os.WriteFile("test.go", []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run(t, "git", "add", ".")
	run(t, "git", "commit", "-m", "Initial commit")

	// Record some pending context
	stateDir := filepath.Join(".context", "state")
	if err := os.MkdirAll(stateDir, 0750); err != nil {
		t.Fatal(err)
	}
	_ = trace.Record("decision:1", stateDir)

	// Write history for the commit
	traceDir := filepath.Join(".context", "trace")
	hash := strings.TrimSpace(runOutput(t, "git", "rev-parse", "HEAD"))

	err := trace.WriteHistory(trace.HistoryEntry{
		Commit:  hash,
		Refs:    []string{"decision:1"},
		Message: "Initial commit",
	}, traceDir)
	if err != nil {
		t.Fatalf("WriteHistory: %v", err)
	}

	// Test ctx trace <commit>
	traceCmd := Cmd()
	traceCmd.SetArgs([]string{hash[:7]})
	if err := traceCmd.Execute(); err != nil {
		t.Errorf("trace show failed: %v", err)
	}

	// Test ctx trace tag
	traceCmd = Cmd()
	traceCmd.SetArgs([]string{"tag", "HEAD", "--note", "Test tag"})
	if err := traceCmd.Execute(); err != nil {
		t.Errorf("trace tag failed: %v", err)
	}

	// Verify override was written
	overrides, err := trace.ReadOverrides(traceDir)
	if err != nil {
		t.Fatalf("ReadOverrides: %v", err)
	}
	if len(overrides) != 1 {
		t.Errorf("expected 1 override, got %d", len(overrides))
	}
}

func run(t *testing.T, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s %v failed: %v", name, args, err)
	}
}

func runOutput(t *testing.T, name string, args ...string) string {
	t.Helper()
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		t.Fatalf("%s %v failed: %v", name, args, err)
	}
	return string(out)
}
```

- [ ] **Step 2: Run integration test**

Run: `cd /Users/parlakisik/projects/github/ctx && CTX_SKIP_PATH_CHECK=1 go test ./internal/cli/trace/ -v -run TestTraceTagAndShow`
Expected: PASS

- [ ] **Step 3: Run all tests**

Run: `cd /Users/parlakisik/projects/github/ctx && make test`
Expected: All PASS

- [ ] **Step 4: Run lint**

Run: `cd /Users/parlakisik/projects/github/ctx && make lint`
Expected: No errors

- [ ] **Step 5: Commit**

```bash
git add internal/cli/trace/trace_test.go
git commit -m "test(trace): add integration tests for trace command"
```

---

## Task 14: Hook Post-Commit — Write History Entry

The prepare-commit-msg hook injects the trailer before the commit is finalized. But we also need to record the commit in `history.jsonl` after the commit succeeds. This is done by adding a post-commit behavior to the collect flow.

**Files:**
- Modify: `internal/cli/trace/cmd/collect/run.go`

### Steps

- [ ] **Step 1: Add commit-msg-file argument to collect**

The prepare-commit-msg hook passes the commit message file path. After outputting the trailer, we need to also record it. However, since the commit hasn't happened yet at prepare-commit-msg time, we need a separate mechanism.

Update the collect command to also accept a `--record` flag that writes to history after the fact. Or, more pragmatically, enhance the hook script to also call `ctx trace collect --record` in a post-commit hook.

Actually, the simpler approach: modify the hook script to write the history entry at prepare-commit-msg time (before commit), using a temporary marker. Then the trailer is the canonical source for the data. The `ctx trace` command already reads from trailers as a fallback.

Let's keep it simple: the hook injects the trailer, and `ctx trace` reads from the trailer at query time. The `history.jsonl` is a performance optimization we can add in a follow-up. For now, we'll only write history from the `ctx trace collect --record <hash>` subcommand, which can be called from a post-commit hook.

Update `internal/cli/trace/cmd/collect/cmd.go`:

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package collect

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the trace collect subcommand.
//
// Returns:
//   - *cobra.Command: Configured trace collect command
func Cmd() *cobra.Command {
	var record string

	short, long := desc.Command(cmd.DescKeyTraceCollect)

	c := &cobra.Command{
		Use:    "collect",
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if record != "" {
				return RecordCommit(cmd, record)
			}
			return Run(cmd)
		},
	}

	c.Flags().StringVar(&record, "record", "", "Record history entry for commit hash (called from post-commit)")

	return c
}
```

Update `internal/cli/trace/cmd/collect/run.go` — add `RecordCommit`:

```go
// RecordCommit writes a history entry for a completed commit.
// Called from the post-commit hook with the commit hash.
//
// Parameters:
//   - cmd: Cobra command for output
//   - commitHash: The commit hash to record
//
// Returns:
//   - error: Non-nil on failure
func RecordCommit(cmd *cobra.Command, commitHash string) error {
	contextDir := rc.ContextDir()
	stateDir := filepath.Join(contextDir, dir.State)
	traceDir := filepath.Join(contextDir, dir.Trace)

	// Read pending context before truncating
	refs := trace.Collect(contextDir)
	if len(refs) == 0 {
		return nil
	}

	// Get commit message
	message := commitMessage(commitHash)

	entry := trace.HistoryEntry{
		Commit:  commitHash,
		Refs:    refs,
		Message: message,
	}

	if err := trace.WriteHistory(entry, traceDir); err != nil {
		return err
	}

	// Truncate pending context
	_ = trace.TruncatePending(stateDir)

	return nil
}

func commitMessage(hash string) string {
	out, err := exec.Command("git", "log", "-1", "--format=%s", hash).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
```

Add necessary imports to collect/run.go:

```go
import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trace"
)
```

- [ ] **Step 2: Update hook script to include post-commit**

Update the hook script in `internal/cli/trace/cmd/hook/run.go` to also install a post-commit hook:

Add a `postCommitScript` constant:

```go
const postCommitScript = `#!/bin/sh
# ctx: post-commit hook for recording commit context history.
# Installed by: ctx trace hook enable
# Remove with:  ctx trace hook disable

COMMIT_HASH=$(git rev-parse HEAD)
ctx trace collect --record "$COMMIT_HASH" 2>/dev/null || true
`
```

Update `enable` to install both hooks:

```go
func enable(cmd *cobra.Command) error {
	prepareHook, err := hookFilePath("prepare-commit-msg")
	if err != nil {
		return err
	}
	postHook, err := hookFilePath("post-commit")
	if err != nil {
		return err
	}

	if err := installHook(prepareHook, hookScript, "prepare-commit-msg"); err != nil {
		return err
	}
	if err := installHook(postHook, postCommitScript, "post-commit"); err != nil {
		return err
	}

	cmd.Printf("Installed prepare-commit-msg and post-commit hooks\n")
	return nil
}

func installHook(path, script, name string) error {
	if _, statErr := os.Stat(path); statErr == nil {
		content, readErr := os.ReadFile(filepath.Clean(path))
		if readErr == nil && strings.Contains(string(content), "ctx trace") {
			return nil // already installed
		}
		return fmt.Errorf("a %s hook already exists at %s; remove it first or add ctx integration manually", name, path)
	}
	return os.WriteFile(path, []byte(script), fs.PermExec)
}
```

Update `disable` to remove both hooks:

```go
func disable(cmd *cobra.Command) error {
	prepareHook, err := hookFilePath("prepare-commit-msg")
	if err != nil {
		return err
	}
	postHook, err := hookFilePath("post-commit")
	if err != nil {
		return err
	}

	removeHook(prepareHook)
	removeHook(postHook)

	cmd.Println("Removed ctx trace hooks")
	return nil
}

func removeHook(path string) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return
	}
	if strings.Contains(string(content), "ctx trace") {
		_ = os.Remove(path)
	}
}
```

Update `hookFilePath` to accept the hook name:

```go
func hookFilePath(hookName string) (string, error) {
	out, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}

	gitDir := strings.TrimSpace(string(out))
	hooksDir := filepath.Join(gitDir, "hooks")

	if mkdirErr := os.MkdirAll(hooksDir, fs.PermExec); mkdirErr != nil {
		return "", fmt.Errorf("create hooks dir: %w", mkdirErr)
	}

	return filepath.Join(hooksDir, hookName), nil
}
```

- [ ] **Step 3: Run build to verify**

Run: `cd /Users/parlakisik/projects/github/ctx && go build ./cmd/ctx/`
Expected: BUILD SUCCESS

- [ ] **Step 4: Commit**

```bash
git add internal/cli/trace/cmd/collect/ internal/cli/trace/cmd/hook/
git commit -m "feat(trace): add post-commit history recording and dual hook management"
```

---

## Task 15: Final Verification

### Steps

- [ ] **Step 1: Run full test suite**

Run: `cd /Users/parlakisik/projects/github/ctx && make test`
Expected: All PASS

- [ ] **Step 2: Run linter**

Run: `cd /Users/parlakisik/projects/github/ctx && make lint`
Expected: No errors

- [ ] **Step 3: Run build**

Run: `cd /Users/parlakisik/projects/github/ctx && make build`
Expected: BUILD SUCCESS

- [ ] **Step 4: Manual smoke test**

Run these commands to verify the feature works end-to-end:

```bash
# Build and install
cd /Users/parlakisik/projects/github/ctx && go build -o /tmp/ctx ./cmd/ctx/

# Test trace --last (should show existing commits with no context)
/tmp/ctx trace --last 5

# Test trace tag
/tmp/ctx trace tag HEAD --note "Test: commit context tracing feature"

# Verify tag was written
cat .context/trace/overrides.jsonl

# Test trace on HEAD (should show the manual tag)
/tmp/ctx trace $(git rev-parse --short HEAD)

# Test hook enable (don't actually enable in this repo)
# /tmp/ctx trace hook enable
```

- [ ] **Step 5: Final commit (if any fixes needed)**

```bash
git add -A
git commit -m "fix(trace): final adjustments from smoke testing"
```
