//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package render is the **client-side renderer** that
// turns hub entries received from `ctx connection sync` /
// `ctx connection listen` into markdown files under
// `.context/hub/` so the local agent can read them.
//
// Each entry becomes a markdown block with a date
// header, an origin tag (which project published it),
// and the entry body, all separated by horizontal
// rules. The format is the same one defined by
// [internal/assets/tpl.HubEntryMarkdown] so what
// ships through the gRPC pipe is what lands on disk.
//
// # Public Surface
//
//   - **[WriteEntries](entries)**: appends each
//     entry to the matching per-type file
//     (`decisions.md`, `learnings.md`,
//     `conventions.md`, `tasks.md`) under
//     `.context/hub/`, formatting via
//     [HubEntryMarkdown]. It appends
//     unconditionally: it is not idempotent and does
//     not deduplicate, so handing it the same entry
//     twice writes it twice. Skipping what has
//     already landed is the caller's job —
//     `ctx connection sync` does it by passing the
//     hub only the sequences above its last-seen
//     watermark, while `ctx connection listen` asks
//     for sequence 0 on every run and therefore
//     re-appends its backlog.
//
// # File Layout
//
//   - `.context/hub/decisions.md`
//   - `.context/hub/learnings.md`
//   - `.context/hub/conventions.md`
//   - `.context/hub/tasks.md`
//   - `.context/hub/.sync-state.json`: the single
//     last-seen hub sequence, written by
//     `ctx connection sync` so its resume is exact.
//     `ctx connection listen` neither reads nor
//     writes it.
//
// # Concurrency
//
// Filesystem-bound. Concurrent renderers against
// the same hub directory would race; the
// `ctx connection listen` daemon is single-instance
// per project by convention.
package render
