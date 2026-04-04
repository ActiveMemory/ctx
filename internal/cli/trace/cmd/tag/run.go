//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tag

import (
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	errTrace "github.com/ActiveMemory/ctx/internal/err/trace"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trace"
	writeTrace "github.com/ActiveMemory/ctx/internal/write/trace"
)

// Run executes the trace tag command logic.
//
// Resolves commitRef to a full hash, attaches the note as an override entry,
// and writes it to the trace directory via trace.WriteOverride.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - commitRef: commit ref or hash to tag (e.g. "HEAD", "abc1234")
//   - note: context note to attach to the commit
//
// Returns:
//   - error: non-nil on execution failure or empty note
func Run(cmd *cobra.Command, commitRef, note string) error {
	if note == "" {
		return errTrace.NoteRequired()
	}

	hash, resolveErr := trace.ResolveCommitHash(commitRef)
	if resolveErr != nil {
		return errTrace.ResolveCommit(commitRef, resolveErr)
	}

	traceDir := filepath.Join(rc.ContextDir(), dir.Trace)

	entry := trace.OverrideEntry{
		Commit: hash,
		Refs:   []string{strconv.Quote(note)},
	}

	if writeErr := trace.WriteOverride(entry, traceDir); writeErr != nil {
		return errTrace.WriteOverride(writeErr)
	}

	writeTrace.Tagged(cmd, trace.ShortHash(hash), note)
	return nil
}
