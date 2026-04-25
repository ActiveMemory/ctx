//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

// Selected returns the innermost visible .context/ directory
// alongside any additional candidates further up the path.
//
// Single-source-anchor model
// (specs/single-source-context-anchor.md): activation is always
// project-local scan from CWD. The explicit-path mode that used
// to accept an argument was removed.
//
// Multi-candidate is no longer an error: workspace-level shared
// `.context/` dirs alongside per-project ones are a legitimate
// nested-project layout. Innermost wins (matching `git` / `make`
// behavior in nested layouts), and the additional candidates are
// surfaced so callers can include them as informational comments
// in eval-able output.
//
// Returns:
//   - string: absolute path of the resolved .context/ directory.
//   - []string: additional candidates further up the path, nil
//     when only one is visible.
//   - error: [errActivate.NoCandidates] when no `.context/` is
//     visible from CWD upward.
func Selected() (string, []string, error) {
	return scan()
}
