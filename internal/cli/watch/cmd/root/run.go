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

	"github.com/ActiveMemory/ctx/internal/cli/watch/core/stream"
	"github.com/ActiveMemory/ctx/internal/context/validate"
	"github.com/ActiveMemory/ctx/internal/err/initialize"
	errJournal "github.com/ActiveMemory/ctx/internal/err/journal"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/watch"
)

// Run executes the watch command logic.
//
// Sets up a reader from either a log file (logPath) or stdin, then
// processes the stream for context update commands. Displays status
// messages and respects the dryRun flag.
//
// Parameters:
//   - cmd: Cobra command for output
//   - logPath: Path to the log file, or empty for stdin
//   - dryRun: If true, show what would be updated without making changes
//
// Returns:
//   - error: Non-nil if the context directory is missing, the log file cannot
//     be opened, or stream processing fails
func Run(cmd *cobra.Command, logPath string, dryRun bool) error {
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}
	exists, existsErr := validate.Exists("")
	if existsErr != nil {
		return existsErr
	}
	if !exists {
		return initialize.ContextNotInitialized()
	}

	watch.Started(cmd)
	if dryRun {
		watch.DryRun(cmd)
	}
	watch.StopHint(cmd)
	watch.Separator(cmd)

	var reader io.Reader
	if logPath != "" {
		file, err := internalIo.SafeOpenUserFile(logPath)
		if err != nil {
			return errJournal.OpenLogFile(err)
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

	return stream.Process(cmd, reader, dryRun)
}
