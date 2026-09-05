//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"strings"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// tomlAssign renders the ` = ` separator between a TOML key and
// its value.
//
// Returns:
//   - string: the assignment separator
func tomlAssign() string {
	return token.Space + token.KeyValueSep + token.Space
}

// tableHeaderPresent reports whether any line of the document,
// after trimming, equals the ctx MCP table header.
//
// Parameters:
//   - data: config.toml content
//
// Returns:
//   - bool: true when the header line exists
func tableHeaderPresent(data []byte) bool {
	for _, line := range splitLines(data) {
		if headerLine(line, cfgCodex.TOMLHeaderMCPCtx) {
			return true
		}
	}
	return false
}

// headerLine reports whether a config.toml line is the given table
// header, tolerating equivalent spellings: leading whitespace, a
// trailing comment or whitespace after the closing bracket, and a
// quoted last key segment ([t.k] vs [t."k"] vs [t.'k']). Inline
// tables ([parent] with k = {...}) remain undetected — an accepted
// limit of the never-parse-TOML design, documented in the spec.
//
// Parameters:
//   - line: raw config.toml line
//   - header: canonical header (e.g. `[mcp_servers.ctx]`)
//
// Returns:
//   - bool: true when the line spells this header
func headerLine(line, header string) bool {
	trimmed := strings.TrimSpace(line)
	inner := strings.TrimSuffix(
		strings.TrimPrefix(header, cfgCodex.TOMLBracketOpen),
		cfgCodex.TOMLBracketClose,
	)
	dot := strings.LastIndex(inner, cfgCodex.TOMLDot)
	variants := []string{header}
	if dot >= 0 {
		parent, key := inner[:dot], inner[dot+1:]
		bare := strings.Trim(
			key, cfgCodex.TOMLQuoteBasic+cfgCodex.TOMLQuoteLiteral,
		)
		for _, quoted := range []string{
			bare,
			cfgCodex.TOMLQuoteBasic + bare + cfgCodex.TOMLQuoteBasic,
			cfgCodex.TOMLQuoteLiteral + bare + cfgCodex.TOMLQuoteLiteral,
		} {
			variants = append(variants,
				cfgCodex.TOMLBracketOpen+parent+
					cfgCodex.TOMLDot+quoted+
					cfgCodex.TOMLBracketClose,
			)
		}
	}
	for _, v := range variants {
		if !strings.HasPrefix(trimmed, v) {
			continue
		}
		rest := strings.TrimSpace(trimmed[len(v):])
		if rest == "" || strings.HasPrefix(rest, cfgCodex.TOMLComment) {
			return true
		}
	}
	return false
}

// splitLines splits a document into lines without a line-length
// ceiling (bufio.Scanner caps lines at 64KB and stops scanning
// silently, which would make header detection miss everything
// after one very long line).
//
// Parameters:
//   - data: document content
//
// Returns:
//   - []string: the document's lines
func splitLines(data []byte) []string {
	return strings.Split(string(data), token.NewlineLF)
}
