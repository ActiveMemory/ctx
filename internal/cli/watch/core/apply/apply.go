//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package apply

import (
	"github.com/ActiveMemory/ctx/internal/cli/watch/core"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/err/config"
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
		return addEntry(update)
	case cfgEntry.Decision:
		return addEntry(update)
	case cfgEntry.Learning:
		return addEntry(update)
	case cfgEntry.Convention:
		return addEntry(update)
	case cfgEntry.Complete:
		return completeTask(update.Content)
	default:
		return config.UnknownUpdateType(update.Type)
	}
}
