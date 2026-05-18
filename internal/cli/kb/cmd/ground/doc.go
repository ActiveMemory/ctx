//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ground implements `ctx kb ground`.
//
// Refreshes sources listed in .context/ingest/grounding-sources.md.
// The pass itself is performed by the /ctx-kb-ground skill; this
// CLI surface verifies that the grounding sources file exists and
// is non-empty before printing the canonical skill invocation.
//
// # Refusal contract
//
//   - Missing grounding-sources.md: [GroundingMissing] wraps an
//     errors.Is-friendly refusal carrying the resolved path.
//   - Empty grounding-sources.md: [GroundingEmpty] wraps the same.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb/cli] supplies
//     the hint and refusal strings.
//   - [github.com/ActiveMemory/ctx/internal/err/kb/cli] supplies the
//     refusal constructors.
package ground
