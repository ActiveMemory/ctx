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
	"github.com/ActiveMemory/ctx/internal/config/knowledge"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
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
	_, sessionID, ctxDir, stateDir, ok := check.FullPreamble(stdin)
	bailSilently := !ok
	if bailSilently {
		return nil
	}
	markerPath := filepath.Join(stateDir, knowledge.ThrottleID)
	if check.DailyThrottled(markerPath) {
		return nil
	}

	box, warned, checkErr := coreKnowledge.CheckHealth(sessionID, ctxDir)
	if checkErr != nil {
		return checkErr
	}
	if warned {
		writeSetup.Nudge(cmd, box)
		internalIo.TouchFile(markerPath)
	}

	return nil
}
