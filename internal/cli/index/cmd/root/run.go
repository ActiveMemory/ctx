//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"path/filepath"

	"github.com/spf13/cobra"

	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/heading"
	"github.com/ActiveMemory/ctx/internal/io"
	writeIndex "github.com/ActiveMemory/ctx/internal/write/index"
)

// Run executes the index command logic: read the file, project its headings,
// and render them.
//
// Parameters:
//   - cmd: Cobra command for the output stream.
//   - path: Path to the Markdown file to project.
//   - depth: Deepest heading level to include (2 = ## only; 3 adds ###).
//   - jsonOutput: If true, emit a JSON array instead of lines.
//
// Returns:
//   - error: Non-nil if the file cannot be read (path-bearing) or JSON
//     marshaling fails.
func Run(cmd *cobra.Command, path string, depth int, jsonOutput bool) error {
	content, readErr := io.SafeReadUserFile(filepath.Clean(path))
	if readErr != nil {
		cmd.SilenceUsage = true
		return errFs.FileRead(path, readErr)
	}

	headings := heading.Headings(string(content), depth)

	if jsonOutput {
		return writeIndex.JSON(cmd, headings)
	}
	writeIndex.Lines(cmd, headings)
	return nil
}
