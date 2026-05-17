//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package timeline appends rows to the kb-scoped timeline
// artifact at `.context/kb/timeline.md` for the ctx
// knowledge-base editorial pipeline (Phase KB).
//
// Each row records a dated event in the domain the kb is
// studying, pinned to one or more `EV-###` rows that ground
// it. The timeline is not a changelog of kb activity; that
// lives in closeouts and SESSION_LOG.md.
//
// The writer is append-only; when the artifact does not yet
// exist, [Append] initialises it with the schema's table
// header. The schema lives at
// `internal/assets/kb/templates/ingest/schemas/timeline.md`.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb] supplies
//     [cfgKB.Timeline].
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/core/path]
//     resolves the artifact path via [KBArtifactFile].
package timeline
