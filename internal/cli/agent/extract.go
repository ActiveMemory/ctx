//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/task"
)

// extractBulletItems extracts Markdown bullet items up to a limit.
//
// Skips empty items and lines starting with "#" (headers).
//
// Parameters:
//   - content: Markdown content to parse
//   - limit: Maximum number of items to return
//
// Returns:
//   - []string: Bullet item text without the "- " prefix
func extractBulletItems(content string, limit int) []string {
	matches := config.RegExBulletItem.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, limit)
	for i, m := range matches {
		if i >= limit {
			break
		}
		text := strings.TrimSpace(m[1])
		// Skip empty or header-only items
		if text != "" && !strings.HasPrefix(text, "#") {
			items = append(items, text)
		}
	}
	return items
}

// extractCheckboxItems extracts text from Markdown checkbox items.
//
// Matches both checked "- [x]" and unchecked "- [ ]" items.
//
// Parameters:
//   - content: Markdown content to parse
//
// Returns:
//   - []string: Text content of each checkbox item
func extractCheckboxItems(content string) []string {
	matches := config.RegExTask.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, len(matches))
	for _, m := range matches {
		items = append(items, strings.TrimSpace(task.Content(m)))
	}
	return items
}

// extractConstitutionRules extracts checkbox items from CONSTITUTION.md.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: List of constitution rules; nil if the file is not found
func extractConstitutionRules(ctx *context.Context) []string {
	if f := ctx.File(config.FileConstitution); f != nil {
		return extractCheckboxItems(string(f.Content))
	}
	return nil
}

// extractUncheckedTasks extracts unchecked Markdown checkbox items.
//
// Only matches "- [ ]" items (not checked). Returns items with the
// "- [ ]" prefix preserved for display.
//
// Parameters:
//   - content: Markdown content to parse
//
// Returns:
//   - []string: Unchecked task items with "- [ ]" prefix
func extractUncheckedTasks(content string) []string {
	matches := config.RegExTaskMultiline.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, len(matches))
	for _, m := range matches {
		if task.Pending(m) {
			items = append(items, "- [ ] "+strings.TrimSpace(task.Content(m)))
		}
	}
	return items
}

// extractActiveTasks extracts unchecked task items from TASKS.md.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: List of active tasks with "- [ ]" prefix; nil if
//     the file is not found
func extractActiveTasks(ctx *context.Context) []string {
	if f := ctx.File(config.FileTask); f != nil {
		return extractUncheckedTasks(string(f.Content))
	}
	return nil
}
