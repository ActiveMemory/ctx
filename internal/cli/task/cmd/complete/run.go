//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package complete

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	coreComplete "github.com/ActiveMemory/ctx/internal/cli/task/core/complete"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
	"github.com/ActiveMemory/ctx/internal/trace"
	writeComplete "github.com/ActiveMemory/ctx/internal/write/complete"
)

// Run executes the complete command logic.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Command arguments (first arg is the query)
//
// Returns:
//   - error: Non-nil on task match or write failure
func Run(cmd *cobra.Command, args []string) error {
	matchedTask, matchedNum, completeErr := coreComplete.Complete(args[0], "")
	if completeErr != nil {
		return completeErr
	}

	writeComplete.Completed(cmd, matchedTask)

	// Best-effort: record pending context for commit tracing.
	ref := fmt.Sprintf(
		cfgTrace.RefFormat, cfgTrace.RefTypeTask, matchedNum,
	)
	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		return dirErr
	}
	_ = trace.Record(ref, stateDir)

	return nil
}
