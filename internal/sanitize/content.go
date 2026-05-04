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

// Content escapes Markdown structural patterns from s so untrusted
// input cannot inject entry headers, task checkboxes, or constitution
// rules into .context/ files.
//
// Parameters:
//   - s: raw content string from untrusted input
//
// Returns:
//   - string: s with Markdown structural patterns escaped and null
//     bytes removed
func Content(s string) string {
	s = regex.SanEntryHeader.ReplaceAllStringFunc(s, func(m string) string {
		return cfgSan.EscapePrefix + m
	})
	s = regex.SanTaskCheckbox.ReplaceAllStringFunc(s, func(m string) string {
		return cfgSan.EscapePrefix + m
	})
	s = regex.SanConstitutionRule.ReplaceAllStringFunc(s, func(m string) string {
		return cfgSan.EscapePrefix + m
	})
	s = strings.ReplaceAll(s, cfgSan.NullByte, "")
	return s
}

// StripControl removes ASCII control characters from s, preserving
// tab (\t), line feed (\n), and carriage return (\r).
//
// Parameters:
//   - s: input string potentially containing control characters
//
// Returns:
//   - string: s with control characters removed (except tab and newlines)
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
