//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/drift/core/fix"
	"github.com/ActiveMemory/ctx/internal/cli/drift/core/out"
	"github.com/ActiveMemory/ctx/internal/context/load"
	"github.com/ActiveMemory/ctx/internal/drift"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeDrift "github.com/ActiveMemory/ctx/internal/write/drift"
)

// Run executes the drift command logic.
//
// Loads context, runs drift detection, and outputs results
// in the specified format. When `doFix` is true, attempts
// to auto-fix supported issue types.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - jsonOutput: If true, output as JSON
//   - doFix: If true, attempt to auto-fix supported issues
//
// Returns:
//   - error: Non-nil if context loading fails
func Run(
	cmd *cobra.Command, jsonOutput, doFix bool,
) error {
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}
	ctx, err := load.Do("")
	if err != nil {
		if _, ok := errors.AsType[*errCtx.NotFoundError](err); ok {
			return errInit.NotInitialized()
		}
		return err
	}

	report := drift.Detect(ctx)

	// Apply fixes if requested
	if doFix && (len(report.Warnings) > 0 ||
		len(report.Violations) > 0) {
		writeDrift.FixHeader(cmd)

		result := fix.Apply(cmd, ctx, report)

		writeDrift.BlankLine(cmd)
		if result.Fixed > 0 {
			writeDrift.FixedCount(cmd, result.Fixed)
		}
		if result.Skipped > 0 {
			writeDrift.SkippedCount(cmd, result.Skipped)
		}
		for _, errMsg := range result.Errors {
			writeDrift.FixError(cmd, errMsg)
		}

		// Re-run detection to show the updated status
		if result.Fixed > 0 {
			writeDrift.FixRecheck(cmd)
			ctx, _ = load.Do("")
			report = drift.Detect(ctx)
		}
	}

	if jsonOutput {
		return out.DriftJSON(cmd, report)
	}

	return out.DriftText(cmd, report)
}
