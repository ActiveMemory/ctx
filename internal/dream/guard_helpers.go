//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"path/filepath"
	"strings"

	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
)

// relUnderRoot resolves target relative to projectRoot and returns the
// cleaned relative path. It reports an error when either path cannot be
// made absolute or the relative path cannot be computed.
//
// Parameters:
//   - projectRoot: absolute path to the project root
//   - target: the write target (absolute or relative)
//
// Returns:
//   - string: cleaned relative path from projectRoot to target
//   - error: non-nil when root resolution or rel computation fails
func relUnderRoot(projectRoot, target string) (string, error) {
	absRoot, rootErr := filepath.Abs(projectRoot)
	if rootErr != nil {
		return "", errDream.ResolveRoot(rootErr)
	}
	absTarget := target
	if !filepath.IsAbs(absTarget) {
		absTarget = filepath.Join(absRoot, absTarget)
	}
	rel, relErr := filepath.Rel(absRoot, filepath.Clean(absTarget))
	if relErr != nil {
		return "", errDream.RelPath(relErr)
	}
	return rel, nil
}

// underDir reports whether rel (a cleaned relative path) resolves at or
// below dir.
//
// Parameters:
//   - rel: cleaned relative path from the project root
//   - dir: the top-level directory to test containment against
//
// Returns:
//   - bool: true when rel equals dir or is nested under dir
func underDir(rel, dir string) bool {
	if rel == dir {
		return true
	}
	return strings.HasPrefix(rel, dir+string(filepath.Separator))
}
