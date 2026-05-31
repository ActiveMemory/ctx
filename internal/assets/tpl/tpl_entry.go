//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// Markdown format templates for context entries.
//
// These templates define the structure of entries written to .context/ files
// by the add command. Each uses fmt.Sprintf verbs for interpolation.
const (
	// Task formats a task checkbox line.
	// Args: content, priorityTag, sessionTag, branchTag, commitTag, timestamp.
	Task = "- [ ] %s%s%s%s%s #added:%s\n"

	// TaskPriority formats the inline priority tag.
	// Args: priority level.
	TaskPriority = " #priority:%s"

	// TaskSession formats the inline session provenance tag.
	// Args: session ID (8-char short ID or "unknown").
	TaskSession = " #session:%s"

	// TaskBranch formats the inline branch provenance tag.
	// Args: branch name.
	TaskBranch = " #branch:%s"

	// TaskCommit formats the inline commit provenance tag.
	// Args: short commit hash.
	TaskCommit = " #commit:%s"

	// Convention formats a convention list item.
	// Args: content.
	Convention = "- %s\n"
)
