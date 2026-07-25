//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package theme

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// splitSpec divides "<name> — <gist>" on the theme metadata separator.
//
// Parameters:
//   - content: the caller-supplied spec
//
// Returns:
//   - string: the theme name
//   - string: the gist
//   - bool: false when either half is missing
func splitSpec(content string) (string, string, bool) {
	idx := strings.Index(content, token.MetaSeparator)
	if idx == -1 {
		return "", "", false
	}
	name := strings.TrimSpace(content[:idx])
	gist := strings.TrimSpace(content[idx+len(token.MetaSeparator):])
	if name == "" || gist == "" {
		return "", "", false
	}
	return name, gist, true
}
