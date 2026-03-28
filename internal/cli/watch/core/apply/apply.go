//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package apply

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/cli/watch/core"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entry"
	"github.com/ActiveMemory/ctx/internal/err/config"
	errTask "github.com/ActiveMemory/ctx/internal/err/task"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/task"
)

// Update routes a context update to the appropriate handler.
//
// Dispatches based on the update type to add entries to context files
// or mark tasks complete. For learnings and decisions, uses structured
// fields (context, lesson, application, rationale, consequence) if
// provided in the XML attributes.
//
// Parameters:
//   - update: ContextUpdate containing type, content, and optional metadata
//
// Returns:
//   - error: Non-nil if type is unknown or the handler fails
func Update(update core.ContextUpdate) error {
	switch update.Type {
	case cfgEntry.Task:
		return RunAddSilent(update)
	case cfgEntry.Decision:
		return RunAddSilent(update)
	case cfgEntry.Learning:
		return RunAddSilent(update)
	case cfgEntry.Convention:
		return RunAddSilent(update)
	case cfgEntry.Complete:
		return RunCompleteSilent([]string{update.Content})
	default:
		return config.UnknownUpdateType(update.Type)
	}
}

// RunAddSilent appends an entry to a context file without output.
//
// Used by the watch command to silently apply updates detected in
// the input stream. Uses shared validation and write logic from the
// add package to ensure consistent behavior with `ctx add`.
//
// Parameters:
//   - update: The parsed ContextUpdate with type, content, and required
//     structured fields (context, lesson, application for learnings;
//     context, rationale, consequence for decisions)
//
// Returns:
//   - error: Non-nil if validation fails, type is unknown,
//     or file operations fail
func RunAddSilent(update core.ContextUpdate) error {
	params := entry.Params{
		Type:        update.Type,
		Content:     update.Content,
		Context:     update.Context,
		Rationale:   update.Rationale,
		Consequence: update.Consequence,
		Lesson:      update.Lesson,
		Application: update.Application,
	}

	// Validate required fields (same as ctx add)
	if validateErr := entry.Validate(params, nil); validateErr != nil {
		return validateErr
	}

	// Write using the shared function
	// (handles formatting, append, and index update)
	return entry.Write(params)
}

// RunCompleteSilent marks a task as complete without output.
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
func RunCompleteSilent(args []string) error {
	if len(args) < 1 {
		return errTask.NoneSpecified()
	}

	query := args[0]
	filePath := filepath.Join(rc.ContextDir(), ctx.Task)
	nl := token.NewlineLF

	content, readErr := os.ReadFile(filepath.Clean(filePath))
	if readErr != nil {
		return readErr
	}

	lines := strings.Split(string(content), nl)

	matchedLine := -1
	for i, line := range lines {
		match := regex.Task.FindStringSubmatch(line)
		if match != nil && task.Pending(match) {
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
		return errTask.NoMatch(query)
	}

	lines[matchedLine] = regex.Task.ReplaceAllString(
		lines[matchedLine], regex.TaskCompleteReplace,
	)
	return os.WriteFile(filePath, []byte(strings.Join(lines, nl)), fs.PermFile)
}
