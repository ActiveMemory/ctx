//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dream implements the "ctx dream" command: a gated,
// proposing memory-consolidation pass over the gitignored ideas/
// directory.
//
// Invoked with no subcommand, it runs one executor-agnostic pass —
// scan ideas/ by the discipline clock, lock, invoke the configured
// executor to classify and ground each idea, validate the proposals
// written into the gitignored dreams/ notebook, and print a digest.
// It never autonomously mutates canonical memory.
//
// # Subcommands
//
//   - review: list pending proposals for a serendipity round
//   - accept: accept a proposal's recommended disposition
//   - reject: reject a proposal (recorded; not re-surfaced)
//   - amend: accept a proposal with a different action
//
// Mechanical dispositions (archive, mark-blog, keep, reject) apply
// instantly and pass both structural guards before any write;
// generative ones (promote, merge) are recorded as intent and
// completed by /ctx-serendipity from the full source.
//
// # Subpackages
//
//	cmd/review: pending-proposal listing
//	cmd/accept, cmd/reject, cmd/amend: disposition primitives
//	core/paths: project-root and notebook resolution
//	core/pass: one executor-agnostic run
//	core/dispose: load-by-id and apply a decision
//	core/review: pending-proposal filtering and rendering
package dream
