//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"os"

	"github.com/spf13/cobra"

	coreExport "github.com/ActiveMemory/ctx/internal/cli/pad/core/export"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	writeExport "github.com/ActiveMemory/ctx/internal/write/export"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run exports blob entries from the scratchpad to the given directory.
//
// Parameters:
//   - cmd: Cobra command for output routing
//   - dir: Target directory for exported files
//   - force: When true, overwrite existing files instead of timestamping
//   - dryRun: When true, report the plan without writing
//
// Returns:
//   - error: On directory creation or scratchpad read failure
func Run(cmd *cobra.Command, dir string, force, dryRun bool) error {
	if !dryRun {
		if mkErr := os.MkdirAll(dir, fs.PermExec); mkErr != nil {
			return errFs.Mkdir(dir, mkErr)
		}
	}

	items, planErr := coreExport.Plan(dir, force)
	if planErr != nil {
		return planErr
	}

	var count int
	for _, item := range items {
		if dryRun {
			if item.AltName != "" {
				writePad.InfoPathConversionExists(cmd, dir, item.Label, item.AltName)
			} else {
				writePad.ExportPlan(cmd, item.Label, item.OutPath)
			}
			count++
			continue
		}

		if item.Exists {
			writeExport.InfoExistsWritingAsAlternative(cmd, item.Label, item.AltName)
		}

		if writeErr := os.WriteFile(
			item.OutPath, item.Data, fs.PermSecret,
		); writeErr != nil {
			writePad.ErrExportWrite(cmd, item.Label, writeErr)
			continue
		}

		writePad.ExportDone(cmd, item.Label)
		count++
	}

	writePad.ExportSummary(cmd, count, dryRun)
	return nil
}
