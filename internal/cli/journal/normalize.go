//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"html"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ActiveMemory/ctx/internal/config"
)

// HTML tag constants for pre-formatted blocks.
const (
	tagPre      = "<pre>"
	tagPreClose = "</pre>"
)

// normalizeContent sanitizes journal Markdown for static site rendering:
//   - Strips code fence markers (eliminates nesting conflicts)
//   - Wraps Tool Output and User sections in <pre><code> with HTML-escaped content
//   - Sanitizes H1 headings (strips Claude tags, truncates to 75 chars)
//   - Demotes non-turn-header headings to bold (prevents broken page structure)
//   - Inserts blank lines before list items when missing (Python-Markdown requires them)
//   - Strips bold markers from tool-use lines (**Glob: *.md** -> Glob: *.md)
//   - Escapes glob-like * characters outside code blocks
//   - Replaces inline code spans containing angle brackets with quoted entities
//
// Heavy formatting (metadata tables, proper fence reconstruction) is left to
// the ctx-journal-normalize skill which uses AI for context-aware cleanup.
//
// Parameters:
//   - content: Raw Markdown content of a journal entry
//   - fencesVerified: Whether the file's fences have been verified via state
//
// Returns:
//   - string: Sanitized content ready for static site rendering
func normalizeContent(content string, fencesVerified bool) string {
	// Strip fences first — eliminates all nesting conflicts
	content = stripFences(content, fencesVerified)

	// Wrap Tool Output and User turn bodies in <pre><code> with
	// HTML-escaped content. Eliminates all markdown interpretation —
	// headings, separators, fence markers, HTML tags become inert.
	// Strips <details>/<pre> wrappers from the source pipeline and
	// re-wraps uniformly.
	content = wrapToolOutputs(content)
	content = wrapUserTurns(content)

	lines := strings.Split(content, config.NewlineLF)
	var out []string
	inFrontmatter := false
	inPreBlock := false // inside <pre>...</pre> from wrapToolOutputs/wrapUserTurns

	for i, line := range lines {
		// Skip frontmatter
		if i == 0 && strings.TrimSpace(line) == config.Separator {
			inFrontmatter = true
			out = append(out, line)
			continue
		}
		if inFrontmatter {
			out = append(out, line)
			if strings.TrimSpace(line) == config.Separator {
				inFrontmatter = false
			}
			continue
		}

		// Track <pre> blocks from wrapToolOutputs/wrapUserTurns.
		// Content inside is HTML-escaped — skip all transforms.
		trimmed := strings.TrimSpace(line)
		if trimmed == "<pre><code>" || trimmed == tagPre {
			inPreBlock = true
			out = append(out, line)
			continue
		}
		if inPreBlock {
			if trimmed == "</code></pre>" || trimmed == tagPreClose {
				inPreBlock = false
			}
			out = append(out, line)
			continue
		}

		// Sanitize H1 headings: strip Claude tags, truncate to max title len
		if strings.HasPrefix(line, config.HeadingLevelOneStart) {
			heading := strings.TrimPrefix(line, config.HeadingLevelOneStart)
			heading = strings.TrimSpace(
				config.RegExClaudeTag.ReplaceAllString(heading, ""),
			)
			if utf8.RuneCountInString(heading) > config.RecallMaxTitleLen {
				runes := []rune(heading)
				truncated := string(runes[:config.RecallMaxTitleLen])
				if idx := strings.LastIndex(truncated, " "); idx > 0 {
					truncated = truncated[:idx]
				}
				heading = truncated
			}
			line = config.HeadingLevelOneStart + heading
		}

		// Demote headings to bold: ## Foo → **Foo**
		// Preserves turn headers (### N. Role (HH:MM:SS)) and the H1 title.
		if hm := config.RegExMarkdownHeading.FindStringSubmatch(line); hm != nil {
			if hm[1] != "#" && !config.RegExTurnHeader.MatchString(strings.TrimSpace(line)) {
				line = "**" + hm[2] + "**"
			}
		}

		// Insert blank line before list items when previous line is non-empty.
		// Python-Markdown requires a blank line before the first list item.
		if config.RegExListStart.MatchString(line) &&
			len(out) > 0 && strings.TrimSpace(out[len(out)-1]) != "" {
			out = append(out, "")
		}

		// Strip bold from tool-use lines
		line = config.RegExToolBold.ReplaceAllString(line, `🔧 $1`)

		// Escape glob stars
		if !strings.HasPrefix(line, "    ") {
			line = config.RegExGlobStar.ReplaceAllString(line, `\*$1`)
		}

		// Replace inline code spans containing angle brackets:
		// `</com` → "&lt;/com" (quotes preserve visual signal,
		// entities prevent broken HTML in rendered output).
		line = config.RegExInlineCodeAngle.ReplaceAllStringFunc(line, func(m string) string {
			inner := m[1 : len(m)-1] // strip backticks
			inner = strings.ReplaceAll(inner, "<", "&lt;")
			inner = strings.ReplaceAll(inner, ">", "&gt;")
			return `"` + inner + `"`
		})

		out = append(out, line)
	}

	return strings.Join(out, config.NewlineLF)
}

// wrapToolOutputs finds Tool Output sections and wraps their body in
// <pre><code> with HTML-escaped content. This prevents all markdown
// interpretation — headings, separators, HTML tags, fence markers all
// become inert entities.
//
// Requires pymdownx.highlight with use_pygments=false in the zensical
// config (set in TplZensicalTheme) to prevent the highlight extension
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
func wrapToolOutputs(content string) string {
	lines := strings.Split(content, config.NewlineLF)
	mask := preBlockMask(lines)
	turnSeq := collectTurnNumbers(lines)
	var out []string
	i := 0

	for i < len(lines) {
		// Skip lines inside <pre> blocks — they are not real turn headers.
		if mask[i] {
			out = append(out, lines[i])
			i++
			continue
		}
		m := config.RegExTurnHeader.FindStringSubmatch(
			strings.TrimSpace(lines[i]),
		)
		if m == nil || m[2] != config.LabelToolOutput {
			out = append(out, lines[i])
			i++
			continue
		}

		// Tool Output header — emit it, then collect body
		turnNum, _ := strconv.Atoi(m[1])
		turnTime := m[3]
		out = append(out, lines[i])
		i++

		// The boundary target is the minimum turn number > turnNum.
		// If the same number appears multiple times (e.g., an embedded
		// ### 42. inside <pre> AND the real ### 42. after </details>),
		// use the LAST occurrence — the real turn is always positionally
		// after any embedded duplicates.
		expectedNext := nextInSequence(turnSeq, turnNum)

		// Scan ahead to find the last occurrence of expectedNext,
		// skipping lines inside <pre> blocks (embedded content).
		boundary := len(lines) // default: EOF
		for j := i; j < len(lines); j++ {
			if mask[j] {
				continue
			}
			nm := config.RegExTurnHeader.FindStringSubmatch(
				strings.TrimSpace(lines[j]),
			)
			if nm != nil {
				nextNum, _ := strconv.Atoi(nm[1])
				nextTime := nm[3]
				if nextNum == expectedNext && nextTime >= turnTime {
					boundary = j
				}
			}
		}

		body := lines[i:boundary]
		i = boundary

		// If we hit EOF, split off any trailing multipart navigation
		// footer (--- + **Part N of M**) so it's not swallowed.
		var footer []string
		if i >= len(lines) {
			body, footer = splitTrailingFooter(body)
		}

		// Extract raw content — strip existing <details>/<pre> wrappers
		// and unescape HTML entities from the export pipeline.
		raw := stripPreWrapper(body)

		// Drop empty or boilerplate tool outputs entirely (header + body).
		// The header was already appended to out — remove it.
		if isBoilerplateToolOutput(raw) {
			out = out[:len(out)-1]
			continue
		}

		// Trim leading/trailing blank lines.
		start, end := 0, len(raw)-1
		for start <= end && strings.TrimSpace(raw[start]) == "" {
			start++
		}
		for end >= start && strings.TrimSpace(raw[end]) == "" {
			end--
		}

		trimmed := raw[start : end+1]

		// HTML-escape and wrap in <pre><code>...</code></pre>.
		out = append(out, "")
		out = append(out, "<pre><code>")
		for _, line := range trimmed {
			out = append(out, html.EscapeString(line))
		}
		out = append(out, "</code></pre>")
		out = append(out, "")

		// Emit footer after the block if present.
		if len(footer) > 0 {
			out = append(out, footer...)
		}
	}

	return strings.Join(out, config.NewlineLF)
}

// wrapUserTurns finds User turn bodies and wraps them in <pre><code>
// with HTML-escaped content. This is the "defencify" strategy: user input
// is treated as plain preformatted text, which eliminates an entire class
// of rendering bugs caused by stray/unclosed fence markers in user messages.
//
// Requires pymdownx.highlight with use_pygments=false in the zensical
// config (set in TplZensicalTheme). With Pygments enabled, the highlight
// extension hijacks <pre><code> and transforms block boundaries.
//
// Type 1 HTML block (<pre>) survives blank lines (ends at </pre>, not at a
// blank line). HTML escaping prevents ALL inner content conflicts — fence
// markers, headings, HTML tags, etc. all become inert entities.
//
// Trade-off: markdown formatting in user messages (bold, links, lists) is
// flattened to plain text. This is acceptable — preserving user input
// verbatim is more valuable than rendering decorative formatting.
//
// Boundary detection reuses the same pre-scan + last-match-wins approach
// as wrapToolOutputs.
func wrapUserTurns(content string) string {
	lines := strings.Split(content, config.NewlineLF)
	mask := preBlockMask(lines)
	turnSeq := collectTurnNumbers(lines)
	var out []string
	i := 0

	for i < len(lines) {
		// Skip lines inside <pre> blocks — they are not real turn headers.
		if mask[i] {
			out = append(out, lines[i])
			i++
			continue
		}
		m := config.RegExTurnHeader.FindStringSubmatch(
			strings.TrimSpace(lines[i]),
		)
		if m == nil || m[2] != config.LabelRoleUser {
			out = append(out, lines[i])
			i++
			continue
		}

		// User turn header — emit it, then collect body
		turnNum, _ := strconv.Atoi(m[1])
		turnTime := m[3]
		out = append(out, lines[i])
		i++

		expectedNext := nextInSequence(turnSeq, turnNum)

		// Scan ahead to find the last occurrence of expectedNext,
		// skipping lines inside <pre> blocks (embedded content).
		boundary := len(lines) // default: EOF
		for j := i; j < len(lines); j++ {
			if mask[j] {
				continue
			}
			nm := config.RegExTurnHeader.FindStringSubmatch(
				strings.TrimSpace(lines[j]),
			)
			if nm != nil {
				nextNum, _ := strconv.Atoi(nm[1])
				nextTime := nm[3]
				if nextNum == expectedNext && nextTime >= turnTime {
					boundary = j
				}
			}
		}

		body := lines[i:boundary]
		i = boundary

		// Trim leading/trailing blank lines from user body.
		start, end := 0, len(body)-1
		for start <= end && strings.TrimSpace(body[start]) == "" {
			start++
		}
		for end >= start && strings.TrimSpace(body[end]) == "" {
			end--
		}

		if start > end {
			// Empty user turn — emit blank lines as-is
			out = append(out, body...)
			continue
		}

		trimmed := body[start : end+1]

		// HTML-escape the content and wrap in <pre><code>...</code></pre>.
		out = append(out, "")
		out = append(out, "<pre><code>")
		for _, line := range trimmed {
			out = append(out, html.EscapeString(line))
		}
		out = append(out, "</code></pre>")
		out = append(out, "")
	}

	return strings.Join(out, config.NewlineLF)
}

// stripPreWrapper removes <details>, <summary>, <pre>, </pre>, </details>
// wrapper lines from tool output body. When <pre> tags are found (the old
// export format that HTML-escapes content), entities are unescaped. When
// only <details>/<summary> are found (collapseToolOutputs format), inner
// content is returned as-is since it was never HTML-escaped.
//
// Returns raw content lines ready for wrapping.
func stripPreWrapper(body []string) []string {
	var inner []string
	hadPre := false

	for _, line := range body {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "<details>" || trimmed == "</details>":
			continue
		case trimmed == tagPre || trimmed == tagPreClose:
			hadPre = true
			continue
		case strings.HasPrefix(trimmed, "<summary>") &&
			strings.HasSuffix(trimmed, "</summary>"):
			continue
		default:
			inner = append(inner, line)
		}
	}

	// Only unescape when <pre> was found — the old export format
	// HTML-escapes content inside <pre> blocks. The collapseToolOutputs
	// format (just <details>/<summary>) does not escape content.
	if hadPre {
		for j, line := range inner {
			inner[j] = html.UnescapeString(line)
		}
	}

	return inner
}

// isBoilerplateToolOutput returns true if the tool output body contains only
// empty lines or low-value confirmation messages that add no information to
// the rendered journal page. Both the Tool Output header and body are dropped.
//
// Detected patterns:
//   - Empty body (no non-blank lines)
//   - "No matches found" (grep/glob with zero results)
//   - Edit confirmations ("The file ... has been updated successfully.")
//   - Hook denials ("Hook PreToolUse:... denied this tool")
func isBoilerplateToolOutput(raw []string) bool {
	// Collect non-blank lines.
	var nonBlank []string
	for _, line := range raw {
		if strings.TrimSpace(line) != "" {
			nonBlank = append(nonBlank, strings.TrimSpace(line))
		}
	}

	// Empty body — no content at all.
	if len(nonBlank) == 0 {
		return true
	}

	// Join all non-blank lines for multi-line pattern matching.
	// softWrapContent can split single messages across lines.
	joined := strings.Join(nonBlank, " ")

	switch {
	case joined == "No matches found":
		return true
	case strings.HasPrefix(joined, "The file ") &&
		strings.HasSuffix(joined, "has been updated successfully."):
		return true
	case strings.Contains(joined, "denied this tool"):
		return true
	}

	return false
}

// preBlockMask returns a boolean slice where mask[i] is true if line i is
// inside a <pre> block (between <pre>/<pre><code> and </pre>/</code></pre>).
// This allows turn-header scanning to skip embedded headers from tool outputs
// that quote other journal files.
func preBlockMask(lines []string) []bool {
	mask := make([]bool, len(lines))
	inPre := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inPre && (trimmed == tagPre || trimmed == "<pre><code>") {
			inPre = true
			continue
		}
		if inPre {
			if trimmed == tagPreClose || trimmed == "</code></pre>" ||
				trimmed == "</details>" {
				inPre = false
				continue
			}
			mask[i] = true
		}
	}
	return mask
}

// collectTurnNumbers extracts all turn numbers from turn headers in the
// document, returning them sorted and deduplicated. Headers inside <pre>
// blocks are skipped — they are embedded content from tool outputs that
// read other journal files.
func collectTurnNumbers(lines []string) []int {
	mask := preBlockMask(lines)
	seen := make(map[int]bool)
	for i, line := range lines {
		if mask[i] {
			continue
		}
		if m := config.RegExTurnHeader.FindStringSubmatch(
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

// nextInSequence returns the smallest number in the sorted slice that is
// strictly greater than n. Returns -1 if no such number exists.
func nextInSequence(sorted []int, n int) int {
	idx := sort.SearchInts(sorted, n+1)
	if idx < len(sorted) {
		return sorted[idx]
	}
	return -1
}

// splitTrailingFooter splits a multipart navigation footer from the end of
// tool output body lines. The footer pattern is: a "---" separator followed
// (possibly across multiple lines) by a "**Part N of M**" label with
// navigation links. Returns (body without footer, footer lines). If no
// footer is found, returns the original body and nil.
func splitTrailingFooter(body []string) ([]string, []string) {
	// Find the last "---" separator and check if a "**Part " line follows.
	sepIdx := -1
	for j := len(body) - 1; j >= 0; j-- {
		if strings.TrimSpace(body[j]) == config.Separator {
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
		if strings.HasPrefix(strings.TrimSpace(body[j]), "**Part ") {
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
