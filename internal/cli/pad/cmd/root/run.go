//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/load"
)

// Run imports entries into the scratchpad from a file, stdin, or directory.
//
// When blobs is true, imports directory contents as blob entries.
// Otherwise reads lines from a file (or stdin when path is "-").
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: File path, "-" for stdin, or directory path (when blobs is true)
//   - blobs: When true, import directory contents as blob entries
//
// Returns:
//   - error: Non-nil on read/write failure
func Run(cmd *cobra.Command, path string, blobs bool) error {
	if blobs {
		return load.Blobs(cmd, path)
	}
	return load.Lines(cmd, path)
}
