//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package drift provides the **terminal-output helpers**
// the `ctx drift` and `ctx drift --fix` CLI surfaces use
// to render their per-issue progress and final summary.
//
// All exported functions take a `*cobra.Command` so
// they route through cobra's output stream (which tests
// can wire to a buffer for assertion).
//
// # Public Surface
//
// Output families:
//
//   - **Fix progress**: [FixHeader],
//     [FixRecheck], [FixedCount],
//     [SkippedCount], [FixError], [FixStaleness].
//     Used by `--fix` to narrate what is being
//     auto-remediated and what remained.
//   - **Per-issue lines**: formatters for the
//     individual issue rows (path refs,
//     staleness markers, missing files,
//     constitution violations) with the matching
//     status glyph.
//   - **Summaries**: final roll-up for `ctx
//     drift` (counts of warnings/violations/
//     passed) and for `--fix` (counts of fixed
//     vs skipped).
//
// # Why a Separate Output Package
//
// Same data, two surfaces (`ctx drift` and
// `ctx drift --fix`), each with its own preferred
// presentation. Hoisting both renderers keeps the
// drift detector
// ([internal/drift]) free of presentation
// concerns and the fix engine
// ([internal/cli/drift/core/fix]) free of UI
// strings.
//
// # Concurrency
//
// Pure data → io.Writer. Concurrent calls go
// through cobra's output stream which is
// serialized.
package drift
