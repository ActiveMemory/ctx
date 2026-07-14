//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// Heading is a single Markdown heading projected from a file by the
// generic heading matcher (see internal/heading.Headings).
//
// Unlike IndexEntry, a Heading carries no timestamp semantics: it is any
// ATX heading (`##`, `###`, …), which is what lets one projector serve
// timestamped entry files (DECISIONS/LEARNINGS) and untimestamped ones
// (TASKS `## Phase …`, CONVENTIONS) alike.
//
// Fields:
//   - Level: Number of leading `#` characters (2 for `##`, 3 for `###`).
//   - Text: Heading text with the `#` markers and surrounding space removed.
type Heading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}
