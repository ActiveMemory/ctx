//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package extract

import (
	"strings"

	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/task"
)

// BulletItems extracts Markdown bullet items up to a limit.
//
// Skips empty items and lines starting with "#" (headers).
//
// Parameters:
//   - content: Markdown content to parse
//   - limit: Maximum number of items to return
//
// Returns:
//   - []string: Bullet item text without the "- " prefix
func BulletItems(content string, limit int) []string {
	matches := regex.BulletItem.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, limit)
	for i, m := range matches {
		if i >= limit {
			break
		}
		text := strings.TrimSpace(m[1])
		// Skip empty or header-only items
		if text != "" && !strings.HasPrefix(text, token.PrefixHeading) {
			items = append(items, text)
		}
	}
	return items
}

// CheckboxItems extracts text from Markdown checkbox items.
//
// Matches both checked "- [x]" and unchecked "- [ ]" items.
//
// Parameters:
//   - content: Markdown content to parse
//
// Returns:
//   - []string: Text content of each checkbox item
func CheckboxItems(content string) []string {
	matches := regex.Task.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, len(matches))
	for _, m := range matches {
		items = append(items, strings.TrimSpace(task.Content(m)))
	}
	return items
}

// ConstitutionRules extracts checkbox items from CONSTITUTION.md.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: List of constitution rules; nil if the file is not found
func ConstitutionRules(ctx *entity.Context) []string {
	if f := ctx.File(cfgCtx.Constitution); f != nil {
		return CheckboxItems(string(f.Content))
	}
	return nil
}

// UncheckedTasks extracts unchecked Markdown checkbox items.
//
// Only matches "- [ ]" items (not checked). Returns items with the
// "- [ ]" prefix preserved for display.
//
// Parameters:
//   - content: Markdown content to parse
//
// Returns:
//   - []string: Unchecked task items with "- [ ]" prefix
func UncheckedTasks(content string) []string {
	matches := regex.TaskMultiline.FindAllStringSubmatch(content, -1)
	items := make([]string, 0, len(matches))
	for _, m := range matches {
		if task.Pending(m) {
			text := strings.TrimSpace(task.Content(m))
			items = append(items, marker.PrefixTaskUndone+token.Space+text)
		}
	}
	return items
}

// ActiveTasks extracts unchecked task items from TASKS.md.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: List of active tasks with "- [ ]" prefix; nil if
//     the file is not found
func ActiveTasks(ctx *entity.Context) []string {
	if f := ctx.File(cfgCtx.Task); f != nil {
		return UncheckedTasks(string(f.Content))
	}
	return nil
}
