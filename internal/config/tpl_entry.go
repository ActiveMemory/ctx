//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

// Markdown format templates for context entries.
//
// These templates define the structure of entries written to .context/ files
// by the add command. Each uses fmt.Sprintf verbs for interpolation.
const (
	// TplTask formats a task checkbox line.
	// Args: content, priorityTag, timestamp.
	TplTask = "- [ ] %s%s #added:%s\n"

	// TplTaskPriority formats the inline priority tag.
	// Args: priority level.
	TplTaskPriority = " #priority:%s"

	// TplLearning formats a learning section with all ADR-style fields.
	// Args: timestamp, title, context, lesson, application.
	TplLearning = `## [%s] %s

**Context**: %s

**Lesson**: %s

**Application**: %s
`

	// TplConvention formats a convention list item.
	// Args: content.
	TplConvention = "- %s\n"

	// TplDecision formats a decision section with all ADR fields.
	// Args: timestamp, title, context, title (repeated), rationale, consequences.
	TplDecision = `## [%s] %s

**Status**: Accepted

**Context**: %s

**Decision**: %s

**Rationale**: %s

**Consequences**: %s
`
)
