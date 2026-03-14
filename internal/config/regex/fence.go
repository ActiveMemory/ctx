//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// CodeFenceInline matches code fences that appear inline after text.
// E.g., "some text: ```code" where fence should be on its own line.
// Groups:
//   - 1: preceding non-whitespace character
//   - 2: the code fence (3+ backticks)
var CodeFenceInline = regexp.MustCompile("(\\S) *(```+)")

// CodeFenceClose matches code fences immediately followed by text.
// E.g., "```text" where text should be on its own line after the fence.
// Groups:
//   - 1: the code fence (3+ backticks)
//   - 2: following non-whitespace character
var CodeFenceClose = regexp.MustCompile("(```+) *(\\S)")

// CodeFenceLine matches lines that are code fence markers (3+ backticks or
// tildes, optionally followed by a language tag).
var CodeFenceLine = regexp.MustCompile("^\\s*(`{3,}|~{3,})(.*)$")

// CodeFencePath matches file paths in Markdown backticks.
//
// Groups:
//   - 1: file path
var CodeFencePath = regexp.MustCompile("`([^`]+\\.[a-zA-Z]{1,5})`")
