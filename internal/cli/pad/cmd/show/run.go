//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
)

// runShow prints the raw text of entry at 1-based position n.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//   - outPath: File path for blob output (empty for stdout)
//
// Returns:
//   - error: Non-nil on invalid index, read failure, or write failure
func runShow(cmd *cobra.Command, n int, outPath string) error {
	entries, err := core.ReadEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return ctxerr.EntryRange(n, 0)
	}

	if err := core.ValidateIndex(n, entries); err != nil {
		return err
	}

	entry := entries[n-1]

	if label, data, ok := core.SplitBlob(entry); ok {
		_ = label
		if outPath != "" {
			if writeErr := os.WriteFile(outPath, data, 0600); writeErr != nil {
				return fmt.Errorf("write file: %w", writeErr)
			}
			cmd.Println(fmt.Sprintf("Wrote %d bytes to %s", len(data), outPath))
			return nil
		}
		cmd.Print(string(data))
		return nil
	}

	// Non-blob entry.
	if outPath != "" {
		return fmt.Errorf("--out can only be used with blob entries")
	}

	cmd.Println(entry)
	return nil
}
