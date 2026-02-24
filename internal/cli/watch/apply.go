//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/config"
)

// applyUpdate routes a context update to the appropriate handler.
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
func applyUpdate(update ContextUpdate) error {
	switch update.Type {
	case config.EntryTask:
		return runAddSilent(update)
	case config.EntryDecision:
		return runAddSilent(update)
	case config.EntryLearning:
		return runAddSilent(update)
	case config.EntryConvention:
		return runAddSilent(update)
	case config.EntryComplete:
		return runCompleteSilent([]string{update.Content})
	default:
		return fmt.Errorf("unknown update type: %s", update.Type)
	}
}
