//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_knowledge

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	coreKnowledge "github.com/ActiveMemory/ctx/internal/cli/system/core/knowledge"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/knowledge"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
)

// Run executes the check-knowledge hook logic.
//
// Reads hook input from stdin, checks knowledge file sizes against
// configured thresholds (entry counts for DECISIONS.md and LEARNINGS.md,
// line count for CONVENTIONS.md), and emits a relay warning if any
// file exceeds its limit. Throttled to once per day.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !state.Initialized() {
		return nil
	}

	_, sessionID, paused := check.Preamble(stdin)
	if paused {
		return nil
	}

	markerPath := filepath.Join(state.Dir(), knowledge.KnowledgeThrottleID)
	if check.DailyThrottled(markerPath) {
		return nil
	}

	if box, warned := coreKnowledge.CheckKnowledgeHealth(sessionID); warned {
		writeHook.Nudge(cmd, box)
		internalIo.TouchFile(markerPath)
	}

	return nil
}
