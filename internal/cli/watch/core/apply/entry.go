//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package apply

import (
	"github.com/ActiveMemory/ctx/internal/cli/watch/core"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/entry"
)

// addEntry appends an entry to a context file without output.
//
// Used by [Update] to silently apply updates detected in the watch
// input stream. Uses shared validation and write logic from the
// entry package to ensure consistent behavior with ctx add.
//
// Parameters:
//   - update: The parsed ContextUpdate with type, content, and required
//     structured fields (context, lesson, application for learnings;
//     context, rationale, consequence for decisions)
//
// Returns:
//   - error: Non-nil if validation fails, type is unknown,
//     or file operations fail
func addEntry(update core.ContextUpdate) error {
	params := entity.EntryParams{
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
