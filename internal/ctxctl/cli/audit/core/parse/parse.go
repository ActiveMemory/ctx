//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgAudit "github.com/ActiveMemory/ctx/internal/ctxctl/config/audit"
	errAudit "github.com/ActiveMemory/ctx/internal/ctxctl/err/audit"
)

// Re-export sentinel errors so callers can `errors.Is`
// against this package without importing internal/err
// directly.
var (
	// ErrNoFrontmatter is returned when the report does
	// not open with the YAML frontmatter delimiter.
	ErrNoFrontmatter = errAudit.ErrNoFrontmatter

	// ErrUnterminatedFrontmatter is returned when the
	// closing frontmatter delimiter is missing.
	ErrUnterminatedFrontmatter = errAudit.ErrUnterminatedFrontmatter
)

// Frontmatter splits a report's bytes into typed header
// plus the verbatim body that follows.
//
// Accepted shape (CR-stripped, trailing-whitespace-tolerant):
//
//	---
//	kind: surface
//	status: findings
//	...
//	---
//	<body>
//
// Parameters:
//   - data: full report file contents
//
// Returns:
//   - Header: parsed frontmatter struct
//   - string: report body (everything after the closing
//     delimiter, with a single leading newline stripped)
//   - error: [ErrNoFrontmatter], [ErrUnterminatedFrontmatter],
//     or a YAML unmarshal error
func Frontmatter(data []byte) (Header, string, error) {
	s := strings.ReplaceAll(
		string(data), token.NewlineCRLF, token.NewlineLF,
	)
	lines := strings.Split(s, token.NewlineLF)

	if len(lines) == 0 ||
		strings.TrimSpace(lines[0]) != cfgAudit.FrontmatterDelimiter {
		return Header{}, "", ErrNoFrontmatter
	}

	closeIdx := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == cfgAudit.FrontmatterDelimiter {
			closeIdx = i
			break
		}
	}
	if closeIdx == -1 {
		return Header{}, "", ErrUnterminatedFrontmatter
	}

	yamlBody := strings.Join(lines[1:closeIdx], token.NewlineLF)
	var hdr Header
	if err := yaml.Unmarshal([]byte(yamlBody), &hdr); err != nil {
		return Header{}, "", err
	}

	bodyLines := lines[closeIdx+1:]
	// Drop one leading blank line (common shape) so the
	// body starts at content, not whitespace.
	if len(bodyLines) > 0 && strings.TrimSpace(bodyLines[0]) == "" {
		bodyLines = bodyLines[1:]
	}
	return hdr, strings.Join(bodyLines, token.NewlineLF), nil
}
