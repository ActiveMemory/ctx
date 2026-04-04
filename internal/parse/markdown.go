//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
)

// StripLineNumbers removes Claude Code's line number prefixes from content.
//
// Parameters:
//   - content: Text potentially containing "    1→" style prefixes
//
// Returns:
//   - string: Content with line number prefixes removed
func StripLineNumbers(content string) string {
	return regex.LineNumber.ReplaceAllString(content, "")
}

// ExtractSystemReminders separates system-reminder content from tool output.
//
// Claude Code injects <system-reminder> tags into tool results. This function
// extracts them so they can be rendered as Markdown outside code fences.
//
// Parameters:
//   - content: Tool result content potentially containing system-reminder tags
//
// Returns:
//   - string: Content with system-reminder tags removed
//   - []string: Extracted reminder texts (may be empty)
func ExtractSystemReminders(content string) (string, []string) {
	matches := regex.SystemReminder.FindAllStringSubmatch(content, -1)
	var reminders []string
	for _, m := range matches {
		if len(m) > 1 && m[1] != "" {
			reminders = append(reminders, m[1])
		}
	}
	cleaned := regex.SystemReminder.ReplaceAllString(content, "")
	return cleaned, reminders
}

// NormalizeCodeFences ensures code fences are on their own lines with
// proper spacing.
//
// Users often type "text: ```code" without proper line
// breaks. Markdown requires
// code fences to be on their own lines with blank lines separating them from
// surrounding content.
//
// Parameters:
//   - content: Text that may contain inline code fences
//
// Returns:
//   - string: Content with code fences properly separated by blank lines
func NormalizeCodeFences(content string) string {
	doubleNL := token.NewlineLF + token.NewlineLF
	result := regex.CodeFenceInline.ReplaceAllString(content, "$1"+doubleNL+"$2")
	result = regex.CodeFenceClose.ReplaceAllString(result, "$1"+doubleNL+"$2")
	return result
}

// FenceForContent returns the appropriate code fence for content.
//
// Uses longer fences when content contains backticks to avoid
// nested Markdown rendering issues. Starts with ``` and adds
// more backticks as needed.
//
// Parameters:
//   - content: The content to be fenced
//
// Returns:
//   - string: A fence string (e.g., "`"x3, "`"x4)
func FenceForContent(content string) string {
	fence := token.CodeFence
	for strings.Contains(content, fence) {
		fence += token.Backtick
	}
	return fence
}

// SplitFrontmatter separates YAML frontmatter from a
// markdown body. Frontmatter must start with a ---
// line and end with a second --- line.
//
// Parameters:
//   - data: Raw file bytes
//
// Returns:
//   - []byte: YAML frontmatter (between delimiters)
//   - string: Body after the closing delimiter
//   - error: Non-nil if delimiters are missing
func SplitFrontmatter(
	data []byte,
) ([]byte, string, error) {
	content := strings.TrimLeft(
		string(data), token.TrimCR,
	)

	if !strings.HasPrefix(
		content, token.FrontmatterDelimiter,
	) {
		return nil, "", errParser.MissingOpenDelim()
	}

	rest := content[len(token.FrontmatterDelimiter):]
	rest = strings.TrimPrefix(rest, token.NewlineLF)

	needle := token.NewlineLF +
		token.FrontmatterDelimiter
	idx := strings.Index(rest, needle)
	if idx < 0 {
		return nil, "", errParser.MissingCloseDelim()
	}

	fm := rest[:idx]
	after := rest[idx+1+len(
		token.FrontmatterDelimiter,
	):]
	after = strings.TrimPrefix(
		after, token.NewlineLF,
	)

	return []byte(fm), after, nil
}
