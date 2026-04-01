//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package collect

import (
	"fmt"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trace"
)

// RecordCommit records context refs for a specific commit hash to history.
//
// Called from the post-commit hook after a commit is made. Reads refs from
// the commit trailer (not re-collected — the trailer is the single source
// of truth), writes a history entry, and truncates pending state.
//
// Pending context is always consumed (truncated) per commit, even when no
// hook ran and the trailer is empty. This prevents stale refs from leaking
// into future commits.
//
// Parameters:
//   - commitHash: full commit hash to record context for
//
// Returns:
//   - error: non-nil on execution failure
func RecordCommit(commitHash string) error {
	contextDir := rc.ContextDir()

	// Read refs from the commit trailer — single source of truth.
	// This matches exactly what was injected by the prepare-commit-msg hook.
	refs := trace.ReadTrailerRefs(commitHash)
	if len(refs) == 0 {
		// No trailer injected — truncate pending and exit.
		stateDir := filepath.Join(contextDir, dir.State)
		_ = trace.TruncatePending(stateDir)
		return nil
	}

	message, err := trace.CommitMessage(commitHash)
	if err != nil {
		return fmt.Errorf("git log: %w", err)
	}

	traceDir := filepath.Join(contextDir, dir.Trace)
	entry := trace.HistoryEntry{
		Commit:  commitHash,
		Refs:    refs,
		Message: message,
	}
	if err := trace.WriteHistory(entry, traceDir); err != nil {
		return fmt.Errorf("write history: %w", err)
	}

	stateDir := filepath.Join(contextDir, dir.State)
	_ = trace.TruncatePending(stateDir)

	return nil
}
