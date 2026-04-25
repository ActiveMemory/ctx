//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resolve picks the `.context/` directory that `ctx
// activate` should emit a shell export for. It is the ONE place
// in the CLI that walks the filesystem during context resolution;
// all other commands honor `CTX_DIR` or error via
// [rc.RequireContextDir].
//
// [Selected] is the single entry point. It walks upward from CWD
// via [rc.ScanCandidates] and returns:
//
//   - the **innermost** visible `.context/` (selected),
//   - any **additional** candidates further up the path,
//   - or [errActivate.NoCandidates] when the walk finds none.
//
// Multi-candidate is not an error: workspace-level shared
// `.context/` dirs alongside per-project ones are a legitimate
// nested-project layout. Innermost wins (matching git / make
// behavior in nested layouts), and the additional candidates are
// surfaced so callers can include them as informational comments
// in eval-able output.
package resolve
