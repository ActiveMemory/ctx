//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package decision implements the ctx decision command
// group for managing DECISIONS.md.
//
// DECISIONS.md records architectural choices with their
// rationale, trade-offs, and timestamps. The decision
// command group provides tooling to append entries; a
// quick-reference index is projected on demand by
// `ctx index DECISIONS.md`, not stored in the file.
//
// # Subcommands
//
//   - add: appends a new decision entry with structured
//     ADR-style fields (context, rationale, consequence)
//     plus required provenance metadata (session-id, branch,
//     commit)
//
// # Subpackages
//
//	cmd/add: cobra command for noun-first decision addition
package decision
