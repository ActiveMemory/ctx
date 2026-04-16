//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resource provides the **terminal-output
// helpers** the resource-related CLI surfaces use to
// render their results in either human-readable text or
// JSON for tooling.
//
// All functions take a `*cobra.Command` so they route
// through cobra's output stream (which tests can wire
// to a buffer for assertion).
//
// # Public Surface
//
//   - **[Text](cmd, payload)** — renders the
//     resource snapshot in human-readable form —
//     section headers, glyph-prefixed counts,
//     summary line.
//   - **[JSON](cmd, payload)** — emits the same
//     payload as a structured JSON document with
//     a UTC timestamp wrapper, suitable for
//     `jq` consumption in CI.
//
// # Why a Separate Output Package
//
// The same data shape needs two different
// renderings (human vs machine). Hoisting both
// into a write-side package keeps the producer
// (the CLI command) free of presentation choices
// and the renderer free of business logic.
//
// # Concurrency
//
// Pure data → io.Writer transformation.
// Concurrent calls each write to the cobra
// command's output stream; cobra serializes them.
package resource
