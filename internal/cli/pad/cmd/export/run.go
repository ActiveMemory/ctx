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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
	"github.com/ActiveMemory/ctx/internal/write"
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
		if mkErr := os.MkdirAll(dir, 0o750); mkErr != nil {
			return fmt.Errorf("mkdir %s: %w", dir, mkErr)
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
					write.InfoPathConversionExists(cmd, dir, label, newName)
					count++
					continue
				}
				outPath = filepath.Join(dir, newName)
				write.InfoExistsWritingAsAlternative(cmd, label, newName)
			}
		}

		if dryRun {
			cmd.Println(fmt.Sprintf("  %s → %s", label, outPath))
			count++
			continue
		}

		if writeErr := os.WriteFile(outPath, data, 0o600); writeErr != nil {
			cmd.PrintErrln(fmt.Sprintf("  ! failed to write %s: %v", label, writeErr))
			continue
		}

		cmd.Println(fmt.Sprintf("  + %s", label))
		count++
	}

	if count == 0 {
		cmd.Println("No blob entries to export.")
		return nil
	}

	verb := "Exported"
	if dryRun {
		verb = "Would export"
	}
	cmd.Println(fmt.Sprintf("%s %d blobs.", verb, count))
	return nil
}
