//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package relationship appends rows to the relationship-map
// artifact at `.context/kb/relationship-map.md` for the ctx
// knowledge-base editorial pipeline (Phase KB).
//
// The relationship map is the kb's edge list: nodes are topic
// slugs and `EV-###` IDs, edges are typed (depends-on,
// refines, contradicts, supersedes, relates-to). Edges key on
// slug + EV id, never on file path, so folder reorganisations
// do not invalidate the graph.
//
// The writer is append-only; when the artifact does not yet
// exist, [Append] initialises it with the schema's table
// header. The schema lives at
// `internal/assets/kb/templates/ingest/schemas/relationship-map.md`.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb] supplies
//     [cfgKB.RelationshipMap].
package relationship
