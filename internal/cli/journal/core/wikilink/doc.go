//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package wikilink converts standard Markdown links to Obsidian
// wikilink format for vault generation.
//
// [ConvertMarkdownLinks] rewrites [text](url) links to [[target|text]]
// syntax. [Format] builds a single wikilink string. [FormatEntry]
// builds a wikilink for a journal entry using its filename and title.
package wikilink
