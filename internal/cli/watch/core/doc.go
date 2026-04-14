//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core holds the **shared helpers** behind
// `ctx watch`: the stdin streamer that scans an AI's
// output for `<ctx-update>` blocks, the parser that
// turns a block into a typed [PendingUpdate], and the
// dispatcher that routes the update to the right
// backend.
//
// `ctx watch` is the bridge between an AI's
// transcript and the project's `.context/` files: it
// reads stdin, picks out the structured update blocks
// the AI emits, and writes the corresponding entries.
//
// # Sub-Packages
//
//   - **[apply]** — dispatches a parsed update to
//     the right backend (entry, task, etc.).
//
// # Public Surface
//
// The shared helpers in this package include the
// stream scanner that reads stdin line by line, the
// XML-block extractor that finds `<ctx-update>` /
// `</ctx-update>` boundaries, and the rate limiter
// that throttles re-application when the same block
// reappears within a short window (Claude Code
// echoes the same update across multiple tool
// results sometimes).
//
// # Concurrency
//
// Single goroutine per `ctx watch` invocation. The
// scanner is sequential by design — order matters
// when adds and completes interleave.
//
// # Related Packages
//
//   - [internal/cli/watch]              — the
//     `ctx watch` CLI surface.
//   - [internal/cli/watch/core/apply]   — the
//     dispatcher.
//   - [internal/entry]                  — the
//     add-side backend.
//   - [internal/cli/task]               — the
//     complete-side backend.
package core
