//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package glossary appends rows to the kb-scoped glossary
// artifact at `.context/kb/glossary.md` for the ctx knowledge-base
// editorial pipeline (Phase KB).
//
// Each row is one canonical term: a definition, a confidence
// band, and at least one backing `EV-###` row from
// `evidence-index.md`. The writer is append-only; idempotency
// at the call-site is the caller's responsibility. This package
// just opens O_CREATE|O_APPEND|O_WRONLY and writes one row.
//
// When the artifact does not exist, [Append] initialises it with
// the schema's table header so subsequent appends are
// table-shaped. The schema lives at
// `internal/assets/kb/templates/ingest/schemas/glossary.md`.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb] supplies
//     the artifact filename constant ([cfgKB.Glossary]).
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/core/path]
//     resolves the artifact path via [KBArtifactFile].
package glossary
