//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/watch/core"
	"github.com/ActiveMemory/ctx/internal/context"
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
		return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
	}

	cmd.Println("Watching for context updates...")
	if dryRun {
		cmd.Println("DRY RUN — No changes will be made")
	}
	cmd.Println("Press Ctrl+C to stop")
	cmd.Println()

	var reader io.Reader
	if logPath != "" {
		file, err := os.Open(logPath) //nolint:gosec // user-provided path via --log flag
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				cmd.Println(fmt.Sprintf("failed to close log file: %v", err))
			}
		}(file)
		reader = file
	} else {
		reader = os.Stdin
	}

	return core.ProcessStream(cmd, reader, dryRun)
}
