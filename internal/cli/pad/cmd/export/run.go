//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/export"
	"github.com/ActiveMemory/ctx/internal/write/pad"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
)

// runExport exports blob entries from the scratchpad to the given directory.
//
// Parameters:
//   - cmd: Cobra command for output routing
//   - dir: Target directory for exported files
//   - force: When true, overwrite existing files instead of timestamping
//   - dryRun: When true, report the plan without writing
//
// Returns:
//   - error: On directory creation or scratchpad read failure
func runExport(cmd *cobra.Command, dir string, force, dryRun bool) error {
	entries, readErr := core.ReadEntries()
	if readErr != nil {
		return readErr
	}

	if !dryRun {
		if mkErr := os.MkdirAll(dir, fs.PermExec); mkErr != nil {
			return ctxerr.Mkdir(dir, mkErr)
		}
	}

	var count int
	for _, entry := range entries {
		label, data, ok := core.SplitBlob(entry)
		if !ok {
			continue
		}

		outPath := filepath.Join(dir, label)

		if !force {
			if _, statErr := os.Stat(outPath); statErr == nil {
				ts := fmt.Sprintf("%d", time.Now().Unix())
				newName := ts + "-" + label
				if dryRun {
					pad.InfoPathConversionExists(cmd, dir, label, newName)
					count++
					continue
				}
				outPath = filepath.Join(dir, newName)
				export.InfoExistsWritingAsAlternative(cmd, label, newName)
			}
		}

		if dryRun {
			pad.PadExportPlan(cmd, label, outPath)
			count++
			continue
		}

		if writeErr := os.WriteFile(outPath, data, fs.PermSecret); writeErr != nil {
			pad.ErrPadExportWrite(cmd, label, writeErr)
			continue
		}

		pad.PadExportDone(cmd, label)
		count++
	}

	pad.PadExportSummary(cmd, count, dryRun)
	return nil
}
