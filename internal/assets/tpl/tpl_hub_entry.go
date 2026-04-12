//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// Hub entry markdown rendering template.
const (
	// TplEntryMarkdown formats a single hub entry as markdown
	// with a date header, origin tag, and horizontal rule.
	//
	// Args (in order):
	//   - date: formatted date string
	//   - title: first line of content (used as heading)
	//   - origin: entry origin identifier
	//   - content: full entry content
	TplEntryMarkdown = "## [%s] %s\n\n**Origin**: %s\n\n%s\n\n---\n\n"
)
