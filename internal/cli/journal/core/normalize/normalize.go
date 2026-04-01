//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package normalize

import (
	"html"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/reduce"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Content sanitizes journal Markdown for static site rendering:
//   - Strips code fence markers (eliminates nesting conflicts)
//   - Wraps Tool Output and User sections in
//     <pre><code> with HTML-escaped content
//   - Sanitizes H1 headings (strips Claude tags, truncates to 75 chars)
//   - Demotes non-turn-header headings to bold (prevents broken page structure)
//   - Inserts blank lines before list items when missing
//     (Python-Markdown requires them)
//   - Strips bold markers from tool-use lines (**Glob: *.md** -> Glob: *.md)
//   - Escapes glob-like * characters outside code blocks
//   - Replaces inline code spans containing angle brackets with quoted entities
//
// Heavy formatting (metadata tables, proper fence reconstruction) is handled
// programmatically. Edge cases may require parser-level fixes.
//
// Parameters:
//   - content: Raw Markdown content of a journal entry
//   - fencesVerified: Whether the file's fences have been verified via state
//
// Returns:
//   - string: Sanitized content ready for static site rendering
func Content(content string, fencesVerified bool) string {
	// Strip fences first - eliminates all nesting conflicts
	content = reduce.StripFences(content, fencesVerified)

	// Wrap Tool Output and User turn bodies in <pre><code> with
	// HTML-escaped content. Eliminates all markdown interpretation:
	// headings, separators, fence markers, HTML tags become inert.
	// Strips <details>/<pre> wrappers from the source pipeline and
	// re-wraps uniformly.
	content = WrapToolOutputs(content)
	content = WrapUserTurns(content)

	lines := strings.Split(content, token.NewlineLF)
	var out []string
	inFrontmatter := false
	inPreBlock := false // inside <pre>...</pre> from WrapToolOutputs/WrapUserTurns

	for i, line := range lines {
		// Skip frontmatter
		if i == 0 && strings.TrimSpace(line) == token.Separator {
			inFrontmatter = true
			out = append(out, line)
			continue
		}
		if inFrontmatter {
			out = append(out, line)
			if strings.TrimSpace(line) == token.Separator {
				inFrontmatter = false
			}
			continue
		}

		// Track <pre> blocks from WrapToolOutputs/WrapUserTurns.
		// Content inside is HTML-escaped - skip all transforms.
		trimmed := strings.TrimSpace(line)
		if trimmed == marker.TagPreCode || trimmed == marker.TagPre {
			inPreBlock = true
			out = append(out, line)
			continue
		}
		if inPreBlock {
			if trimmed == marker.TagPreCodeClose || trimmed == marker.TagPreClose {
				inPreBlock = false
			}
			out = append(out, line)
			continue
		}

		// Sanitize H1 headings: strip Claude tags, truncate to max title len
		if strings.HasPrefix(line, token.HeadingLevelOneStart) {
			heading := strings.TrimPrefix(line, token.HeadingLevelOneStart)
			heading = strings.TrimSpace(
				regex.SystemClaudeTag.ReplaceAllString(heading, ""),
			)
			if utf8.RuneCountInString(heading) > journal.MaxTitleLen {
				runes := []rune(heading)
				truncated := string(runes[:journal.MaxTitleLen])
				if idx := strings.LastIndex(truncated, " "); idx > 0 {
					truncated = truncated[:idx]
				}
				heading = truncated
			}
			line = token.HeadingLevelOneStart + heading
		}

		// Demote headings to bold: ## Foo → **Foo**
		// Preserves turn headers (### N. Role (HH:MM:SS)) and the H1 title.
		if hm := regex.MarkdownHeading.FindStringSubmatch(line); hm != nil {
			if hm[1] != token.PrefixHeading && !regex.TurnHeader.MatchString(strings.TrimSpace(line)) {
				line = marker.BoldWrap + hm[2] + marker.BoldWrap
			}
		}

		// Insert a blank line before list items when previous line is non-empty.
		// Python-Markdown requires a blank line before the first list item.
		if regex.ListStart.MatchString(line) &&
			len(out) > 0 && strings.TrimSpace(out[len(out)-1]) != "" {
			out = append(out, "")
		}

		// Strip bold from tool-use lines
		line = regex.ToolBold.ReplaceAllString(line, `🔧 $1`)

		// Escape glob stars
		if !strings.HasPrefix(line, "    ") {
			line = regex.GlobStar.ReplaceAllString(line, `\*$1`)
		}

		// Replace inline code spans containing angle brackets:
		// `</com` → "&lt;/com" (quotes preserve visual signal,
		// entities prevent broken HTML in rendered output).
		replacer := func(m string) string {
			inner := m[1 : len(m)-1] // strip backticks
			inner = strings.ReplaceAll(inner, marker.AngleLT, marker.EntityLT)
			inner = strings.ReplaceAll(inner, marker.AngleGT, marker.EntityGT)
			return `"` + inner + `"`
		}
		line = regex.InlineCodeAngle.ReplaceAllStringFunc(
			line, replacer,
		)

		out = append(out, line)
	}

	return strings.Join(out, token.NewlineLF)
}

// WrapToolOutputs finds Tool Output sections and wraps their body in
// <pre><code> with HTML-escaped content. This prevents all markdown
// interpretation: headings, separators, HTML tags, fence markers all
// become inert entities.
//
// Requires pymdownx.highlight with use_pygments=false in the zensical
// config (set in ZensicalTheme) to prevent the highlight extension
// from hijacking <pre><code> blocks.
//
// Tool outputs already wrapped in <details><pre> by the export pipeline
// are unwrapped and unescaped before re-escaping uniformly.
//
// Boundary detection: all turn numbers are pre-scanned and sorted. For
// turn N, the boundary target is the minimum turn number > N across the
// entire document. This correctly skips embedded turn headers from other
// journal files (e.g., ### 802. Assistant inside a tool output that read
// another session's file) because the real next turn (### 42.) is always
// the smallest number > N.
//
// Parameters:
//   - content: Raw markdown content containing tool output sections
//
// Returns:
//   - string: Transformed markdown with tool outputs wrapped in pre/code blocks
func WrapToolOutputs(content string) string {
	return ProcessTurns(content, desc.Text(text.DescKeyLabelToolOutput),
		func(out, body []string, atEOF bool) []string {
			// If we hit EOF, split off any trailing multipart navigation
			// footer (--- + **Part N of M**) so it's not swallowed.
			var footer []string
			if atEOF {
				body, footer = SplitTrailingFooter(body)
			}

			// Extract raw content: strip existing <details>/<pre> wrappers
			// and unescape HTML entities from the export pipeline.
			raw := StripPreWrapper(body)

			// Drop empty or boilerplate tool outputs entirely (header + body).
			// The header was already appended to out - remove it.
			if IsBoilerplateToolOutput(raw) {
				return out[:len(out)-1]
			}

			trimmed := TrimBlankLines(raw)
			if len(trimmed) == 0 {
				return out[:len(out)-1]
			}

			// HTML-escape and wrap in <pre><code>...</code></pre>.
			out = append(out, "", marker.TagPreCode)
			for _, line := range trimmed {
				out = append(out, html.EscapeString(line))
			}
			out = append(out, marker.TagPreCodeClose, "")

			if len(footer) > 0 {
				out = append(out, footer...)
			}
			return out
		})
}

// WrapUserTurns finds User turn bodies and wraps them in <pre><code>
// with HTML-escaped content. This is the "defencify" strategy: user input
// is treated as plain preformatted text, which eliminates an entire class
// of rendering bugs caused by stray/unclosed fence markers in user messages.
//
// Requires pymdownx.highlight with use_pygments=false in the zensical
// config (set in ZensicalTheme). With Pygments enabled, the highlight
// extension hijacks <pre><code> and transforms block boundaries.
//
// Type 1 HTML block (<pre>) survives blank lines (ends at </pre>, not at a
// blank line). HTML escaping prevents ALL inner content conflicts: fence
// markers, headings, HTML tags, etc. all become inert entities.
//
// Trade-off: markdown formatting in user messages (bold, links, lists) is
// flattened to plain text. This is acceptable: preserving user input
// verbatim is more valuable than rendering decorative formatting.
//
// Boundary detection reuses the same pre-scan + last-match-wins approach
// as WrapToolOutputs.
//
// Parameters:
//   - content: Raw markdown content containing user turn sections
//
// Returns:
//   - string: Transformed markdown with user turns wrapped in pre/code blocks
func WrapUserTurns(content string) string {
	return ProcessTurns(content, desc.Text(text.DescKeyLabelRoleUser),
		func(out, body []string, _ bool) []string {
			trimmed := TrimBlankLines(body)
			if len(trimmed) == 0 {
				out = append(out, body...)
				return out
			}

			// HTML-escape and wrap in <pre><code>...</code></pre>.
			out = append(out, "", marker.TagPreCode)
			for _, line := range trimmed {
				out = append(out, html.EscapeString(line))
			}
			out = append(out, marker.TagPreCodeClose, "")
			return out
		})
}

// StripPreWrapper removes <details>, <summary>, <pre>, </pre>, </details>
// wrapper lines from tool output body. When <pre> tags are found (the old
// export format that HTML-escapes content), entities are unescaped. When
// only <details>/<summary> are found (ToolOutputs format), inner
// content is returned as-is since it was never HTML-escaped.
//
// Returns raw content lines ready for wrapping.
//
// Parameters:
//   - body: Lines of tool output body to unwrap
//
// Returns:
//   - []string: Raw content lines with wrapper tags
//     removed and entities unescaped
func StripPreWrapper(body []string) []string {
	var inner []string
	hadPre := false

	for _, line := range body {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == marker.TagDetails || trimmed == marker.TagDetailsClose:
			continue
		case trimmed == marker.TagPre || trimmed == marker.TagPreClose:
			hadPre = true
			continue
		case strings.HasPrefix(trimmed, marker.TagSummaryOpen) &&
			strings.HasSuffix(trimmed, marker.TagSummaryClose):
			continue
		default:
			inner = append(inner, line)
		}
	}

	// Only unescape when <pre> was found: the old export format
	// HTML-escapes content inside <pre> blocks. The ToolOutputs
	// format (just <details>/<summary>) does not escape content.
	if hadPre {
		for j, line := range inner {
			inner[j] = html.UnescapeString(line)
		}
	}

	return inner
}

// IsBoilerplateToolOutput returns true if the tool output body contains only
// empty lines or low-value confirmation messages that add no information to
// the rendered journal page. Both the Tool Output header and body are dropped.
//
// Detected patterns:
//   - Empty body (no non-blank lines)
//   - "No matches found" (grep/glob with zero results)
//   - Edit confirmations ("The file ... has been updated successfully.")
//   - Hook denials ("Hook PreToolUse:... denied this tool")
//
// Parameters:
//   - raw: Lines of raw tool output content to check
//
// Returns:
//   - bool: True if the output is empty or matches a boilerplate pattern
func IsBoilerplateToolOutput(raw []string) bool {
	// Collect non-blank lines.
	var nonBlank []string
	for _, line := range raw {
		if strings.TrimSpace(line) != "" {
			nonBlank = append(nonBlank, strings.TrimSpace(line))
		}
	}

	// Empty body - no content at all.
	if len(nonBlank) == 0 {
		return true
	}

	// Join all non-blank lines for multi-line pattern matching.
	// Content can split single messages across lines.
	joined := strings.Join(nonBlank, token.Space)

	switch {
	case joined == journal.BoilerplateNoMatch:
		return true
	case strings.HasPrefix(joined, journal.BoilerplateFilePrefix) &&
		strings.HasSuffix(joined, journal.BoilerplateFileSuffix):
		return true
	case strings.Contains(joined, journal.BoilerplateDenied):
		return true
	}

	return false
}

// PreBlockMask returns a boolean slice where mask[i] is true if line i is
// inside a <pre> block (between <pre>/<pre><code> and </pre>/</code></pre>).
// This allows turn-header scanning to skip embedded headers from tool outputs
// that quote other journal files.
//
// Parameters:
//   - lines: Document lines to scan for pre block regions
//
// Returns:
//   - []bool: Boolean slice where true indicates the line is inside a pre block
func PreBlockMask(lines []string) []bool {
	mask := make([]bool, len(lines))
	inPre := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inPre && (trimmed == marker.TagPre || trimmed == marker.TagPreCode) {
			inPre = true
			continue
		}
		if inPre {
			if trimmed == marker.TagPreClose || trimmed == marker.TagPreCodeClose ||
				trimmed == marker.TagDetailsClose {
				inPre = false
				continue
			}
			mask[i] = true
		}
	}
	return mask
}

// CollectTurnNumbers extracts all turn numbers from turn headers in the
// document, returning them sorted and deduplicated. Headers inside <pre>
// blocks are skipped: they are embedded content from tool outputs that
// read other journal files.
//
// Parameters:
//   - lines: Document lines to extract turn numbers from
//
// Returns:
//   - []int: Sorted, deduplicated turn numbers found in the document
func CollectTurnNumbers(lines []string) []int {
	mask := PreBlockMask(lines)
	seen := make(map[int]bool)
	for i, line := range lines {
		if mask[i] {
			continue
		}
		if m := regex.TurnHeader.FindStringSubmatch(
			strings.TrimSpace(line),
		); m != nil {
			num, _ := strconv.Atoi(m[1])
			seen[num] = true
		}
	}
	nums := make([]int, 0, len(seen))
	for n := range seen {
		nums = append(nums, n)
	}
	sort.Ints(nums)
	return nums
}

// NextInSequence returns the smallest number in the sorted slice that is
// strictly greater than n. Returns -1 if no such number exists.
//
// Parameters:
//   - sorted: Sorted slice of turn numbers
//   - n: Reference number to search after
//
// Returns:
//   - int: Smallest number greater than n, or -1 if none exists
func NextInSequence(sorted []int, n int) int {
	idx := sort.SearchInts(sorted, n+1)
	if idx < len(sorted) {
		return sorted[idx]
	}
	return -1
}

// SplitTrailingFooter splits a multipart navigation footer from the end of
// tool output body lines. The footer pattern is: a "---" separator followed
// (possibly across multiple lines) by a "**Part N of M**" label with
// navigation links. Returns (body without footer, footer lines). If no
// footer is found, returns the original body and nil.
//
// Parameters:
//   - body: Lines of tool output body to split
//
// Returns:
//   - []string: Body lines without footer
//   - []string: Footer lines, or nil if none found
func SplitTrailingFooter(body []string) ([]string, []string) {
	// Find the last "---" separator and check if a "**Part " line follows.
	sepIdx := -1
	for j := len(body) - 1; j >= 0; j-- {
		if strings.TrimSpace(body[j]) == token.Separator {
			sepIdx = j
			break
		}
	}
	if sepIdx < 0 {
		return body, nil
	}

	// Verify a "**Part " line exists after the separator.
	hasPartLabel := false
	for j := sepIdx + 1; j < len(body); j++ {
		if strings.HasPrefix(strings.TrimSpace(body[j]), journal.PartPrefix) {
			hasPartLabel = true
			break
		}
	}
	if !hasPartLabel {
		return body, nil
	}

	// Strip trailing blank lines before the separator.
	cutIdx := sepIdx
	for cutIdx > 0 && strings.TrimSpace(body[cutIdx-1]) == "" {
		cutIdx--
	}

	return body[:cutIdx], body[sepIdx:]
}
