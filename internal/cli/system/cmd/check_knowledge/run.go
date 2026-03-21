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

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/knowledge"
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
	if !core.Initialized() {
		return nil
	}

	_, sessionID, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	markerPath := filepath.Join(core.StateDir(), knowledge.KnowledgeThrottleID)
	if core.IsDailyThrottled(markerPath) {
		return nil
	}

	if core.CheckKnowledgeHealth(cmd, sessionID) {
		core.TouchFile(markerPath)
	}

	return nil
}
