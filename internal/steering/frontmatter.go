//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
	errSteering "github.com/ActiveMemory/ctx/internal/err/steering"
)

// splitFrontmatter separates YAML frontmatter from the markdown body.
// Frontmatter must start with a --- line and end with a second --- line.
func splitFrontmatter(
	data []byte,
) (frontmatter []byte, body string, err error) {
	content := string(data)
	content = strings.TrimLeft(content, token.TrimCR)

	if !strings.HasPrefix(content, token.FrontmatterDelimiter) {
		return nil, "", errSteering.MissingOpeningDelimiter()
	}

	// Skip the opening delimiter line.
	rest := content[len(token.FrontmatterDelimiter):]
	rest = strings.TrimPrefix(rest, token.NewlineLF)

	needle := token.NewlineLF + token.FrontmatterDelimiter
	idx := strings.Index(rest, needle)
	if idx < 0 {
		return nil, "", errSteering.MissingClosingDelimiter()
	}

	fm := rest[:idx]

	// Skip past the closing delimiter line.
	after := rest[idx+1+len(token.FrontmatterDelimiter):]
	// Trim exactly one leading newline from the body if present.
	after = strings.TrimPrefix(after, token.NewlineLF)

	return []byte(fm), after, nil
}

// applyDefaults sets default values for fields not present in the
// parsed frontmatter.
func applyDefaults(sf *SteeringFile) {
	if sf.Inclusion == "" {
		sf.Inclusion = defaultInclusion
	}
	if sf.Priority == 0 {
		sf.Priority = defaultPriority
	}
	// Tools: nil means all tools — no default needed.
}
