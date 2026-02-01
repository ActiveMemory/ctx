//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/task"
)

// runWatch executes the watch command logic.
//
// Sets up a reader from either a log file (--log) or stdin, then
// processes the stream for context update commands. Displays status
// messages and respects the --dry-run flag.
//
// Parameters:
//   - cmd: Cobra command for output
//   - _: Unused positional arguments
//
// Returns:
//   - error: Non-nil if the context directory is missing, the log file cannot
//     be opened, or stream processing fails
func runWatch(cmd *cobra.Command, _ []string) error {
	if !context.Exists("") {
		return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	cmd.Println(cyan("Watching for context updates..."))
	if watchDryRun {
		yellow := color.New(color.FgYellow).SprintFunc()
		cmd.Println(yellow("DRY RUN — No changes will be made"))
	}
	cmd.Println("Press Ctrl+C to stop")
	cmd.Println()

	var reader io.Reader
	if watchLog != "" {
		file, err := os.Open(watchLog)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				cmd.Printf("failed to close log file: %v\n", err)
			}
		}(file)
		reader = file
	} else {
		reader = os.Stdin
	}

	return processStream(cmd, reader)
}

// runAddSilent appends an entry to a context file without output.
//
// Used by the watch command to silently apply updates detected in
// the input stream. Uses shared validation and write logic from the
// add package to ensure consistent behavior with `ctx add`.
//
// Parameters:
//   - update: The parsed ContextUpdate with type, content, and required
//     structured fields (context, lesson, application for learnings;
//     context, rationale, consequences for decisions)
//
// Returns:
//   - error: Non-nil if validation fails, type is unknown,
//     or file operations fail
func runAddSilent(update ContextUpdate) error {
	params := add.EntryParams{
		Type:         update.Type,
		Content:      update.Content,
		Context:      update.Context,
		Rationale:    update.Rationale,
		Consequences: update.Consequences,
		Lesson:       update.Lesson,
		Application:  update.Application,
	}

	// Validate required fields (same as ctx add)
	if err := add.ValidateEntry(params); err != nil {
		return err
	}

	// Write using the shared function
	// (handles formatting, append, and index update)
	return add.WriteEntry(params)
}

// runCompleteSilent marks a task as complete without output.
//
// Used by the watch command to silently complete tasks detected in
// the input stream. Searches for an unchecked task matching the query
// and marks it as done by changing [ ] to [x].
//
// Parameters:
//   - args: Slice where args[0] is the search query to match against
//     task descriptions (case-insensitive substring match)
//
// Returns:
//   - error: Non-nil if args is empty, no matching task is found,
//     or file operations fail
func runCompleteSilent(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no task specified")
	}

	query := args[0]
	filePath := filepath.Join(rc.ContextDir(), config.FileTask)
	nl := config.NewlineLF

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), nl)

	matchedLine := -1
	for i, line := range lines {
		match := config.RegExTask.FindStringSubmatch(line)
		if match != nil && task.IsPending(match) {
			if strings.Contains(
				strings.ToLower(task.Content(match)),
				strings.ToLower(query),
			) {
				matchedLine = i
				break
			}
		}
	}

	if matchedLine == -1 {
		return fmt.Errorf("no task matching %q found", query)
	}

	lines[matchedLine] = config.RegExTask.ReplaceAllString(
		lines[matchedLine], "$1- [x] $3",
	)
	return os.WriteFile(filePath, []byte(strings.Join(lines, nl)), 0644)
}
