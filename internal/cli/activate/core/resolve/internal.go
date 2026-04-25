//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"os"

	errActivate "github.com/ActiveMemory/ctx/internal/err/activate"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// scan returns the innermost visible .context/ candidate from CWD
// alongside any additional candidates further up the path. The
// scan walks via [rc.ScanCandidates] (innermost-first); resolution
// itself never walks outside this function.
//
// Multi-candidate behavior is "innermost wins, the rest are
// reported." This matches what `git` and `make` do for nested
// project layouts (innermost project owns the working directory)
// and supports legitimate workspace-level shared `.context/` dirs
// next to per-project ones; the previous "refuse on multi" rule
// was overly conservative for that workflow. Callers receive the
// full list of additional candidates so they can surface them as
// informational comments in eval-able output without overriding
// the bind.
//
// Returns:
//   - string: absolute path of the innermost (selected) candidate.
//   - []string: zero-or-more additional candidates further up the
//     path, in the order [rc.ScanCandidates] returned them
//     (closest-first). Nil when only one candidate is visible.
//   - error: [errActivate.NoCandidates] when the upward walk finds
//     no `.context/` directory at all. Other errors are surfaced
//     for I/O failures (e.g., os.Getwd).
func scan() (string, []string, error) {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return "", nil, cwdErr
	}
	candidates := rc.ScanCandidates(cwd)
	if len(candidates) == 0 {
		return "", nil, errActivate.NoCandidates()
	}
	return candidates[0], candidates[1:], nil
}
