//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package fix is the **auto-remediation half** of
// `ctx drift`: given a [drift.Report], it applies the
// fixes the package knows how to apply safely (archiving
// completed tasks, creating missing required files from
// templates) and skips the issues that need human
// judgment (dead paths, leaked secrets, constitution
// violations).
//
// The package is the conservative side of the drift
// loop. Anything that could be wrong if applied
// blindly stays in the report and the user fixes it
// by hand.
//
// # What Gets Auto-Fixed
//
//   - **Stale-completed tasks** — tasks marked `[x]`
//     in the body of TASKS.md (not in a Completed
//     section) are archived via [internal/tidy].
//   - **Missing required files** — empty placeholders
//     for the foundation files (CONSTITUTION,
//     CONVENTIONS, etc.) are deployed from the
//     embedded templates.
//
// # What Stays Manual
//
//   - **Dead path references** — the package cannot
//     know whether a path is genuinely gone or just
//     temporarily missing.
//   - **Leaked secrets** — the user must redact and
//     rotate; auto-removal could corrupt history.
//   - **Constitution violations** — the user agreed
//     to the rule and must un-violate it.
//   - **File-age warnings** — staleness is
//     informational, not fixable.
//
// # Public Surface
//
//   - **[Apply](report, contextDir)** — walks the
//     report, applies fixable issues, returns a
//     summary of what was changed and what was
//     skipped.
//
// # Concurrency
//
// Filesystem-bound. Single-process, sequential.
//
// # Related Packages
//
//   - [internal/cli/drift]    — the `ctx drift
//     --fix` CLI surface.
//   - [internal/drift]        — produces the
//     report this package consumes.
//   - [internal/tidy]         — supplies the
//     archive primitives.
package fix
