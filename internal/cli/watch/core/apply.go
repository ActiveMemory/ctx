//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	entry2 "github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entry"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/task"
)

// ApplyUpdate routes a context update to the appropriate handler.
//
// Dispatches based on update type to add entries to context files
// or mark tasks complete. For learnings and decisions, uses structured
// fields (context, lesson, application, rationale, consequences) if
// provided in the XML attributes.
//
// Parameters:
//   - update: ContextUpdate containing type, content, and optional metadata
//
// Returns:
//   - error: Non-nil if type is unknown or the handler fails
func ApplyUpdate(update ContextUpdate) error {
	switch update.Type {
	case entry2.Task:
		return RunAddSilent(update)
	case entry2.Decision:
		return RunAddSilent(update)
	case entry2.Learning:
		return RunAddSilent(update)
	case entry2.Convention:
		return RunAddSilent(update)
	case entry2.Complete:
		return RunCompleteSilent([]string{update.Content})
	default:
		return ctxerr.UnknownUpdateType(update.Type)
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
//     context, rationale, consequences for decisions)
//
// Returns:
//   - error: Non-nil if validation fails, type is unknown,
//     or file operations fail
func RunAddSilent(update ContextUpdate) error {
	params := entry.Params{
		Type:         update.Type,
		Content:      update.Content,
		Context:      update.Context,
		Rationale:    update.Rationale,
		Consequences: update.Consequences,
		Lesson:       update.Lesson,
		Application:  update.Application,
	}

	// Validate required fields (same as ctx add)
	if err := entry.Validate(params, nil); err != nil {
		return err
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
		return ctxerr.NoTaskSpecified()
	}

	query := args[0]
	filePath := filepath.Join(rc.ContextDir(), ctx.Task)
	nl := token.NewlineLF

	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return err
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
		return ctxerr.NoTaskMatch(query)
	}

	lines[matchedLine] = regex.Task.ReplaceAllString(
		lines[matchedLine], regex.TaskCompleteReplace,
	)
	return os.WriteFile(filePath, []byte(strings.Join(lines, nl)), fs.PermFile)
}
