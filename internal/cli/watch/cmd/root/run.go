//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"io"
	"os"

	"github.com/ActiveMemory/ctx/internal/context/validate"
	"github.com/ActiveMemory/ctx/internal/err/initialize"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/recall"
	io2 "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/watch"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/watch/core"
)

// Run executes the watch command logic.
//
// Sets up a reader from either a log file (logPath) or stdin, then
// processes the stream for context update commands. Displays status
// messages and respects the dryRun flag.
//
// Parameters:
//   - cmd: Cobra command for output
//   - logPath: Path to log file, or empty for stdin
//   - dryRun: If true, show what would be updated without making changes
//
// Returns:
//   - error: Non-nil if the context directory is missing, the log file cannot
//     be opened, or stream processing fails
func Run(cmd *cobra.Command, logPath string, dryRun bool) error {
	if !validate.Exists("") {
		return initialize.ContextNotInitialized()
	}

	watch.Watching(cmd)
	if dryRun {
		watch.DryRun(cmd)
	}
	watch.StopHint(cmd)
	cmd.Println()

	var reader io.Reader
	if logPath != "" {
		file, err := io2.SafeOpenUserFile(logPath)
		if err != nil {
			return ctxerr.OpenLogFile(err)
		}
		defer func(file *os.File) {
			if closeErr := file.Close(); closeErr != nil {
				watch.CloseLogError(cmd, closeErr)
			}
		}(file)
		reader = file
	} else {
		reader = os.Stdin
	}

	return core.ProcessStream(cmd, reader, dryRun)
}
