//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check

import (
	"github.com/spf13/cobra"

	coreSchema "github.com/ActiveMemory/ctx/internal/cli/journal/core/schema"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	errSchema "github.com/ActiveMemory/ctx/internal/err/schema"
	"github.com/ActiveMemory/ctx/internal/journal/schema"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
	writeSchema "github.com/ActiveMemory/ctx/internal/write/schema"
)

// Run executes the schema check command.
//
// Parameters:
//   - cmd: Cobra command for output
//   - opts: command flags
//
// Returns:
//   - error: non-nil when drift is detected or scan fails
func Run(cmd *cobra.Command, opts coreSchema.CheckOpts) error {
	c, checkErr := coreSchema.Check(opts)
	if checkErr != nil {
		return checkErr
	}

	if c.Meta.FilesScanned == 0 {
		if len(coreSchema.ScanDirs(opts)) == 0 {
			writeSchema.NoDirs(cmd)
		} else {
			writeSchema.NoFiles(cmd)
		}
		return nil
	}

	reportErr := coreSchema.WriteReport(c)
	if reportErr != nil {
		ctxLog.Warn(
			warn.Write, file.SchemaDrift, reportErr,
		)
	}

	if opts.Quiet {
		if c.Drift() {
			return errSchema.Drift()
		}
		return nil
	}

	if !c.Drift() {
		writeSchema.Clean(
			cmd,
			c.Meta.FilesScanned,
			c.Meta.LinesScanned,
		)
		return nil
	}

	writeSchema.DriftSummary(cmd, schema.Summary(c))
	return errSchema.Drift()
}
