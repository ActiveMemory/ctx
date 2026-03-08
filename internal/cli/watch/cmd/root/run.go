//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/watch/core"
	"github.com/ActiveMemory/ctx/internal/context"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
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
	if !context.Exists("") {
		return ctxerr.ContextNotInitialized()
	}

	write.WatchWatching(cmd)
	if dryRun {
		write.WatchDryRun(cmd)
	}
	write.WatchStopHint(cmd)
	cmd.Println()

	var reader io.Reader
	if logPath != "" {
		file, err := os.Open(logPath) //nolint:gosec // user-provided path via --log flag
		if err != nil {
			return ctxerr.OpenLogFile(err)
		}
		defer func(file *os.File) {
			if closeErr := file.Close(); closeErr != nil {
				write.WatchCloseLogError(cmd, closeErr)
			}
		}(file)
		reader = file
	} else {
		reader = os.Stdin
	}

	return core.ProcessStream(cmd, reader, dryRun)
}
