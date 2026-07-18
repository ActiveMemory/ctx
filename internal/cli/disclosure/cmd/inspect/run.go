//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package inspect

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/disclosure"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/io"
	writeDisc "github.com/ActiveMemory/ctx/internal/write/disclosure"
)

// Run executes the inspect command: resolve the file's kind, read it,
// and report its staged entries and current themes. It never writes.
//
// Parameters:
//   - cmd: Cobra command for the output stream
//   - path: path to the canonical knowledge file
//   - jsonOutput: if true, emit JSON instead of a human summary
//
// Returns:
//   - error: NotAKnowledgeFile when path is not a canonical file, a
//     path-bearing read error, or a JSON-marshal error
func Run(cmd *cobra.Command, path string, jsonOutput bool) error {
	kind, ok := disclosure.KindFor(filepath.Base(path))
	if !ok {
		cmd.SilenceUsage = true
		return errDisc.NotAKnowledgeFile(path)
	}

	content, readErr := io.SafeReadUserFile(filepath.Clean(path))
	if readErr != nil {
		cmd.SilenceUsage = true
		return errFs.FileRead(path, readErr)
	}

	insp := disclosure.Inspect(string(content), kind)
	if jsonOutput {
		return writeDisc.JSON(cmd, insp)
	}
	writeDisc.Human(cmd, insp)
	return nil
}
