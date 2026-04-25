//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

import (
	"path/filepath"

	"github.com/spf13/cobra"

	coreFile "github.com/ActiveMemory/ctx/internal/cli/trace/core/file"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the trace file command logic.
//
// Parses the pathArg into a file path (stripping any line-range suffix),
// then runs git log to retrieve commits touching that file. For each
// commit, context refs are collected from history and overrides and
// printed as a table.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - pathArg: file path with optional line range suffix
//     (e.g. "src/auth.go:42-60")
//   - last: maximum number of commits to show
//
// Returns:
//   - error: non-nil on execution failure
func Run(cmd *cobra.Command, pathArg string, last int) error {
	contextDir, err := rc.RequireContextDir()
	if err != nil {
		cmd.SilenceUsage = true
		return err
	}
	traceDir := filepath.Join(contextDir, dir.Trace)

	filePath := coreFile.ParsePathArg(pathArg)

	return coreFile.Trace(cmd, filePath, last, traceDir)
}
