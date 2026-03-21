//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/drift/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/context/load"
	"github.com/ActiveMemory/ctx/internal/drift"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/initialize"
)

// Run executes the drift command logic.
//
// Loads context, runs drift detection, and outputs results in the
// specified format. When `fix` is true, attempts to auto-fix supported
// issue types (staleness, missing_file).
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - jsonOutput: If true, output as JSON; otherwise output as text
//   - fix: If true, attempt to auto-fix supported issues
//
// Returns:
//   - error: Non-nil if context loading fails or .context/ is not found
func Run(cmd *cobra.Command, jsonOutput, fix bool) error {
	ctx, err := load.Do("")
	if err != nil {
		var notFoundError *errCtx.NotFoundError
		if errors.As(err, &notFoundError) {
			return ctxErr.NotInitialized()
		}
		return err
	}

	report := drift.Detect(ctx)

	// Apply fixes if requested
	if fix && (len(report.Warnings) > 0 || len(report.Violations) > 0) {
		cmd.Println(desc.TextDesc(text.DescKeyDriftApplying))
		cmd.Println()

		result := core.ApplyFixes(cmd, ctx, report)

		cmd.Println()
		if result.Fixed > 0 {
			cmd.Println(fmt.Sprintf(
				desc.TextDesc(text.DescKeyDriftFixedCount), result.Fixed))
		}
		if result.Skipped > 0 {
			cmd.Println(fmt.Sprintf(
				desc.TextDesc(text.DescKeyDriftSkippedCount), result.Skipped))
		}
		for _, errMsg := range result.Errors {
			cmd.Println(fmt.Sprintf(
				desc.TextDesc(text.DescKeyDriftFixError), errMsg))
		}

		// Re-run detection to show the updated status
		if result.Fixed > 0 {
			cmd.Println()
			cmd.Println(desc.TextDesc(text.DescKeyDriftRechecking))
			ctx, _ = load.Do("")
			report = drift.Detect(ctx)
		}
	}

	if jsonOutput {
		return core.OutputDriftJSON(cmd, report)
	}

	return core.OutputDriftText(cmd, report)
}
