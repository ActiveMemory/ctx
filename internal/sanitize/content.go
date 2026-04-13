//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

import (
	"strings"
	"unicode"

	"github.com/ActiveMemory/ctx/internal/config/regex"
	cfgSan "github.com/ActiveMemory/ctx/internal/config/sanitize"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Content neutralizes Markdown structure characters in entry content
// that could corrupt .context/ file parsing.
//
// Escapes entry headers, task checkboxes, and constitution rule
// patterns so they render as literal text instead of structural
// elements.
//
// Parameters:
//   - s: raw content string from MCP client
//
// Returns:
//   - string: content safe for appending to .context/ Markdown files
func Content(s string) string {
	// Escape entry headers: "## [2026-" → "\\## [2026-"
	s = regex.SanEntryHeader.ReplaceAllStringFunc(
		s, func(m string) string {
			return cfgSan.EscapePrefix + m
		},
	)

	// Escape task checkboxes: "- [ ]" → "\\- [ ]"
	s = regex.SanTaskCheckbox.ReplaceAllStringFunc(
		s, func(m string) string {
			return cfgSan.EscapePrefix + m
		},
	)

	// Escape constitution rules.
	s = regex.SanConstitutionRule.ReplaceAllStringFunc(
		s, func(m string) string {
			return cfgSan.EscapePrefix + m
		},
	)

	// Strip null bytes.
	s = strings.ReplaceAll(s, cfgSan.NullByte, "")

	return s
}

// StripControl removes ASCII control characters (except tab and
// newline) from a string.
//
// Parameters:
//   - s: input string potentially containing control characters
//
// Returns:
//   - string: input with control characters removed
func StripControl(s string) string {
	return strings.Map(func(r rune) rune {
		if r == rune(token.Tab[0]) ||
			r == rune(token.NewlineLF[0]) ||
			r == rune(token.NewlineCRLF[0]) {
			return r
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}
