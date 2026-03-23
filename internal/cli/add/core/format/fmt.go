//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import (
	"fmt"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
)

// Task formats a task entry as a Markdown checkbox item.
//
// The output includes a timestamp tag for session correlation and an optional
// priority tag. Format: "- [ ] content #priority:level #added:YYYY-MM-DD-HHMMSS"
//
// Parameters:
//   - content: Task description text
//   - priority: Priority level (high, medium, low); empty string omits the tag
//
// Returns:
//   - string: Formatted task line with trailing newline
func Task(content string, priority string) string {
	// Use YYYY-MM-DD-HHMMSS timestamp for session correlation
	timestamp := time.Now().Format(cfgTime.TimestampCompact)
	var priorityTag string
	if priority != "" {
		priorityTag = fmt.Sprintf(tpl.TaskPriority, priority)
	}
	return fmt.Sprintf(tpl.Task, content, priorityTag, timestamp)
}

// Learning formats a learning entry as a structured Markdown section.
//
// The output includes a timestamped heading and complete sections for context,
// lesson, and application.
//
// Parameters:
//   - title: Learning title/summary text
//   - context: What prompted this learning
//   - lesson: The key insight
//   - application: How to apply this going forward
//
// Returns:
//   - string: Formatted learning section with all fields
func Learning(title, context, lesson, application string) string {
	timestamp := time.Now().Format(cfgTime.TimestampCompact)
	return fmt.Sprintf(
		tpl.Learning, timestamp, title, context, lesson, application,
	)
}

// Convention formats a convention entry as a simple Markdown list item.
//
// Format: "- content"
//
// Parameters:
//   - content: Convention description text
//
// Returns:
//   - string: Formatted convention line with trailing newline
func Convention(content string) string {
	return fmt.Sprintf(tpl.Convention, content)
}

// Decision formats a decision entry as a structured Markdown section.
//
// The output includes a timestamped heading, status, and complete ADR sections
// for context, rationale, and consequence.
//
// Parameters:
//   - title: Decision title/summary text
//   - context: What prompted this decision
//   - rationale: Why this choice over alternatives
//   - consequence: What changes as a result
//
// Returns:
//   - string: Formatted decision section with all ADR fields
func Decision(title, context, rationale, consequence string) string {
	timestamp := time.Now().Format(cfgTime.TimestampCompact)
	return fmt.Sprintf(
		tpl.Decision,
		timestamp, title, context, title, rationale, consequence,
	)
}
