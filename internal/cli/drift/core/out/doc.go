//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package out is the **rendering half** of `ctx drift`:
// given a [drift.Report], it formats the report for
// either humans (terminal text with icons and
// section grouping) or machines (JSON for tooling and
// CI pipelines).
//
// # Public Surface
//
//   - **[Text](report, w)** — writes the
//     human-readable report to `w`. Groups
//     violations and warnings by issue type
//     (path refs, staleness, missing files,
//     other) so similar issues cluster. Renders
//     each with status glyphs (`✗`, `⚠`, `✓`)
//     and a one-line message; passed checks are
//     listed at the bottom.
//   - **[JSON](report, w)** — writes a
//     structured JSON document with a UTC
//     timestamp, the per-issue detail (file,
//     line, type, message, path, rule), and the
//     passed-check list. Stable shape suitable
//     for `jq` parsing in CI scripts.
//
// # Why Two Renderers
//
// Humans want skimmable output with visual
// grouping; CI wants stable JSON with explicit
// types. Hoisting both into a single output
// package keeps the formatting choices in one
// place and the underlying data shape (the
// [drift.Report]) decoupled from how it's
// presented.
//
// # Concurrency
//
// Pure data → io.Writer transformation.
// Concurrent callers never race.
//
// # Related Packages
//
//   - [internal/cli/drift]   — chief consumer.
//   - [internal/drift]       — produces the
//     [Report] this package renders.
//   - [internal/cli/doctor]  — sister renderer
//     for the doctor's report (similar shape,
//     different roll-up).
package out
