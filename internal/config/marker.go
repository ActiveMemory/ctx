//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

// HTML comment markers for parsing and generation.
const (
	// CommentOpen is the HTML comment opening tag.
	CommentOpen = "<!--"
	// CommentClose is the HTML comment closing tag.
	CommentClose = "-->"
)

// Context block markers for embedding context in files.
const (
	// CtxMarkerStart marks the beginning of an embedded context block.
	CtxMarkerStart = "<!-- ctx:context -->"
	// CtxMarkerEnd marks the end of an embedded context block.
	CtxMarkerEnd = "<!-- ctx:end -->"
)

// Index markers for auto-generated table of contents sections.
const (
	// IndexStart marks the beginning of an auto-generated index.
	IndexStart = "<!-- INDEX:START -->"
	// IndexEnd marks the end of an auto-generated index.
	IndexEnd = "<!-- INDEX:END -->"
)

// Task checkbox prefixes for Markdown task lists.
const (
	// PrefixTaskUndone is the prefix for an unchecked task item.
	PrefixTaskUndone = "- [ ]"
	// PrefixTaskDone is the prefix for a checked (completed) task item.
	PrefixTaskDone = "- [x]"
)
