//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tag

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trace"
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
		return errors.New("--note is required")
	}

	hash, err := trace.ResolveCommitHash(commitRef)
	if err != nil {
		return fmt.Errorf("resolve commit %q: %w", commitRef, err)
	}

	traceDir := filepath.Join(rc.ContextDir(), dir.Trace)

	entry := trace.OverrideEntry{
		Commit: hash,
		Refs:   []string{fmt.Sprintf("%q", note)},
	}

	if err := trace.WriteOverride(entry, traceDir); err != nil {
		return fmt.Errorf("write override: %w", err)
	}

	cmd.Println(fmt.Sprintf("Tagged %s with: %s", trace.ShortHash(hash), note))
	return nil
}
