//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"html"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ActiveMemory/ctx/internal/config"
)

// codeFence wraps tool output in a fenced code block. Safe because
// stripFences runs first and removes all fence lines from content.
const codeFence = "```"

// normalizeContent sanitizes journal Markdown for static site rendering:
//   - Strips code fence markers (eliminates nesting conflicts)
//   - Wraps Tool Output sections in <pre> with HTML-escaped content
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

	// Wrap tool output sections in <pre> to prevent ---/# from
	// breaking markdown rendering.
	content = wrapToolOutputs(content)

	lines := strings.Split(content, config.NewlineLF)
	var out []string
	inFrontmatter := false
	inFence := false

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

		// Track fenced code blocks (from wrapToolOutputs).
		// Skip all transforms inside fences.
		if strings.TrimSpace(line) == codeFence {
			inFence = !inFence
			out = append(out, line)
			continue
		}
		if inFence {
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
// fenced code blocks. This is safe because stripFences runs first and
// removes all fence lines from the content. Fenced code blocks
// correctly survive blank lines and prevent all markdown interpretation.
//
// Tool outputs already wrapped in <details><pre> by the export pipeline
// are unwrapped, unescaped, and re-wrapped uniformly.
//
// Boundary detection: a Tool Output section starts at a turn header
// matching "### N. Tool Output (HH:MM:SS)" and ends at the next turn
// header "### M. Role (HH:MM:SS)" where M > N and the timestamp is
// >= the tool output's timestamp. This prevents false matches from tool
// output content that happens to contain turn-header-like text, while
// tolerating gaps in turn numbering.
func wrapToolOutputs(content string) string {
	lines := strings.Split(content, config.NewlineLF)
	var out []string
	i := 0

	for i < len(lines) {
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

		// Collect body lines until the next valid turn header
		var body []string
		for i < len(lines) {
			nm := config.RegExTurnHeader.FindStringSubmatch(
				strings.TrimSpace(lines[i]),
			)
			if nm != nil {
				nextNum, _ := strconv.Atoi(nm[1])
				nextTime := nm[3]
				if nextNum > turnNum && nextTime >= turnTime {
					break
				}
			}
			body = append(body, lines[i])
			i++
		}

		// If we hit EOF, split off any trailing multipart navigation
		// footer (--- + **Part N of M**) so it's not swallowed by the fence.
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

		// Wrap in a fenced code block. Content is emitted verbatim —
		// fenced blocks prevent all markdown/HTML interpretation.
		out = append(out, "")
		out = append(out, codeFence)
		out = append(out, raw...)
		out = append(out, codeFence)
		out = append(out, "")

		// Emit footer after the fence if present.
		if len(footer) > 0 {
			out = append(out, footer...)
		}
	}

	return strings.Join(out, config.NewlineLF)
}

// stripPreWrapper removes <details>, <summary>, <pre>, </pre>, </details>
// wrapper lines from tool output body and unescapes HTML entities in the
// inner content. Returns raw content lines ready for wrapping.
func stripPreWrapper(body []string) []string {
	var inner []string
	wasWrapped := false

	for _, line := range body {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "<details>" || trimmed == "</details>" ||
			trimmed == "<pre>" || trimmed == "</pre>":
			wasWrapped = true
			continue
		case strings.HasPrefix(trimmed, "<summary>") &&
			strings.HasSuffix(trimmed, "</summary>"):
			wasWrapped = true
			continue
		default:
			inner = append(inner, line)
		}
	}

	// If the body had export-pipeline wrapping, the content has
	// HTML entities from html.EscapeString — decode them so the
	// fenced block shows the original text.
	if wasWrapped {
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

	// Single-line boilerplate patterns.
	if len(nonBlank) == 1 {
		line := nonBlank[0]
		switch {
		case line == "No matches found":
			return true
		case strings.HasPrefix(line, "The file ") &&
			strings.HasSuffix(line, "has been updated successfully."):
			return true
		case strings.Contains(line, "denied this tool"):
			return true
		}
	}

	return false
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
