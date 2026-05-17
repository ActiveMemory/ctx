//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sourcemap appends rows to the source-map artifact at
// `.context/kb/source-map.md` for the ctx knowledge-base
// editorial pipeline (Phase KB).
//
// One row per source cited from `evidence-index.md`. Each row
// records what a source *is* (short-name, kind, locator) and
// whether it was admitted against the kb's declared scope.
// The source map records identity-and-admission; the parallel
// source-coverage ledger records progress against admitted
// sources.
//
// The writer is append-only; when the artifact does not yet
// exist, [Append] initialises it with the schema's table
// header. The schema lives at
// `internal/assets/kb/templates/ingest/schemas/source-map.md`.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb] supplies
//     [cfgKB.SourceMap].
//   - [github.com/ActiveMemory/ctx/internal/write/kb/sourcecoverage]
//     is the parallel progress ledger for admitted sources.
package sourcemap
