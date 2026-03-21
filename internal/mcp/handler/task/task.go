//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/parse"
	"github.com/ActiveMemory/ctx/internal/task"
)

// ForEachPending iterates pending top-level tasks in TASKS.md,
// skipping the Completed section and subtasks. It calls fn for each
// match; if fn returns true, iteration stops early.
//
// Parameters:
//   - lines: TASKS.md split by newline
//   - fn: visitor called with each pending task; return true to stop
func ForEachPending(lines []string, fn func(Pending) bool) {
	inCompletedSection := false
	idx := 0

	for _, line := range lines {
		if strings.HasPrefix(line, desc.Text(text.DescKeyHeadingCompleted)) {
			inCompletedSection = true
			continue
		}
		if strings.HasPrefix(
			line, token.HeadingLevelTwoStart,
		) && inCompletedSection {
			inCompletedSection = false
		}
		if inCompletedSection {
			continue
		}

		match := regex.Task.FindStringSubmatch(line)
		if match == nil || !task.Pending(match) {
			continue
		}
		if task.SubTask(match) {
			continue
		}

		idx++
		if fn(Pending{Index: idx, Content: task.Content(match)}) {
			return
		}
	}
}

// ContainsOverlap checks if two strings share meaningful words.
//
// Uses word-set intersection rather than substring matching to avoid
// false positives (e.g., "test" matching inside "contestant").
//
// Parameters:
//   - action: the recent action description
//   - taskText: the task text to compare against
//
// Returns:
//   - bool: true if at least 2 significant words overlap
func ContainsOverlap(action, taskText string) bool {
	actionWords := parse.WordSet(strings.ToLower(action))
	taskWords := strings.Fields(strings.ToLower(taskText))

	matchCount := 0
	for _, w := range taskWords {
		if len(w) < cfg.MinWordLen {
			continue // Skip short common words.
		}
		if actionWords[w] {
			matchCount++
		}
	}

	return matchCount >= cfg.MinWordOverlap
}
