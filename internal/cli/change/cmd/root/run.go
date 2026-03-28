//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/change/core/detect"
	"github.com/ActiveMemory/ctx/internal/cli/change/core/render"
	"github.com/ActiveMemory/ctx/internal/cli/change/core/scan"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	writeChange "github.com/ActiveMemory/ctx/internal/write/change"
)

// Run executes the changes command logic.
//
// Detects a reference time from the --since flag or session markers,
// finds context and code changes since that time, and renders a summary.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - since: Time reference string (duration like "24h" or date like "2026-03-01")
//
// Returns:
//   - error: Non-nil if reference time detection fails
func Run(cmd *cobra.Command, since string) error {
	refTime, refLabel, err := detect.ReferenceTime(since)
	if err != nil {
		return errInit.DetectReferenceTime(err)
	}

	ctxChanges, _ := scan.FindContextChanges(refTime)
	codeChanges, _ := scan.SummarizeCodeChanges(refTime)

	writeChange.Changes(cmd, render.Changes(refLabel, ctxChanges, codeChanges))
	return nil
}
