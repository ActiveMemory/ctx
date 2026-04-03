//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// normalizeInternalPkg truncates a package path to its top-level
// directory (e.g. "internal/cli/pad" → "internal/cli").
//
// Parameters:
//   - path: slash-separated package path
//
// Returns:
//   - string: first two segments joined, or the original path if
//     fewer than two segments exist
func normalizeInternalPkg(path string) string {
	parts := strings.SplitN(path, token.Slash, 3)
	if len(parts) < 2 {
		return path
	}
	return parts[0] + token.Slash + parts[1]
}
