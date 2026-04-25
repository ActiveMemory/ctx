//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errFmt "github.com/ActiveMemory/ctx/internal/err/fmt"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/wrap"
	writeFmt "github.com/ActiveMemory/ctx/internal/write/fmt"
)

// contextFiles lists the context files to format.
var contextFiles = []string{
	cfgCtx.Task,
	cfgCtx.Decision,
	cfgCtx.Learning,
	cfgCtx.Convention,
}

// Run formats all context files to the target line width.
//
// Parameters:
//   - cmd: Cobra command for output
//   - width: Target line width in characters
//   - check: If true, only check without modifying files
//
// Returns:
//   - error: Non-nil if context directory is missing or file
//     operations fail; exits 1 in check mode if files would change
func Run(cmd *cobra.Command, width int, check bool) error {
	ctxDir, ctxErr := rc.RequireContextDir()
	if ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}
	if _, statErr := os.Stat(ctxDir); os.IsNotExist(statErr) {
		return errFmt.NoContextDir()
	}

	formatted := 0
	total := 0
	wouldChange := false

	for _, name := range contextFiles {
		fPath := filepath.Join(ctxDir, name)

		if _, statErr := os.Stat(fPath); os.IsNotExist(statErr) {
			continue
		}
		total++

		content, readErr := io.SafeReadUserFile(
			filepath.Clean(fPath),
		)
		if readErr != nil {
			return errFmt.FileRead(name, readErr)
		}

		wrapped := wrap.ContextFile(string(content), width)

		if wrapped == string(content) {
			continue
		}

		if check {
			wouldChange = true
			writeFmt.NeedsFormatting(cmd, name)
			continue
		}

		if writeErr := io.SafeWriteFile(
			fPath, []byte(wrapped), fs.PermFile,
		); writeErr != nil {
			return errFmt.FileWrite(name, writeErr)
		}
		formatted++
	}

	if total == 0 {
		return errFmt.NoFiles(ctxDir)
	}

	if check && wouldChange {
		return errFmt.NeedsFormatting()
	}

	if !check {
		writeFmt.Summary(cmd, formatted, total)
	}

	return nil
}
