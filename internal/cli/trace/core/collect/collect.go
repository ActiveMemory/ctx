//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package collect

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	errTrace "github.com/ActiveMemory/ctx/internal/err/trace"
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
		// No trailer found — the commit was made without the
		// prepare-commit-msg hook (e.g. --no-verify, external tool,
		// or hook not installed). Pending refs are still truncated
		// because they were accumulated for *this* commit window;
		// keeping them would attach stale context to the next commit.
		stateDir := filepath.Join(contextDir, dir.State)
		_ = trace.TruncatePending(stateDir)
		return nil
	}

	message, msgErr := trace.CommitMessage(commitHash)
	if msgErr != nil {
		return errTrace.GitLog(msgErr)
	}

	traceDir := filepath.Join(contextDir, dir.Trace)
	entry := trace.HistoryEntry{
		Commit:  commitHash,
		Refs:    refs,
		Message: message,
	}
	if histErr := trace.WriteHistory(entry, traceDir); histErr != nil {
		return errTrace.WriteHistory(histErr)
	}

	stateDir := filepath.Join(contextDir, dir.State)
	_ = trace.TruncatePending(stateDir)

	return nil
}
