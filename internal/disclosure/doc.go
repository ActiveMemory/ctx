//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package disclosure implements progressive disclosure for ctx's
// canonical knowledge files — LEARNINGS.md, DECISIONS.md, and
// CONVENTIONS.md (see specs/progressive-disclosure.md).
//
// A canonical file grows without bound while its entries stay valid, so
// an agent that reads the whole file to understand the system eventually
// exhausts its context window. Progressive disclosure turns each file
// into a BOUNDED ROOT: a compact `## Themes` section (a "just enough"
// gist per theme plus a link to a theme file) alongside a staging zone
// of recent, not-yet-digested entries. Entry bodies roll out into
// per-theme files reachable only via the root's links; the agent reads
// the bounded root and drills into a theme file on demand.
//
// # What this package provides (milestone 1)
//
// [Parse] splits a root into its regions without normalizing a byte, so
// [Root.Reconstruct] returns the input exactly — the digesting pass must
// see verbatim bodies. [Validate] is the fail-loud precondition
// (zero-or-one `## Themes`, no entry below it, staging parses into
// discrete entries) that refuses a malformed root rather than regenerate
// from "what it recognized" — the failure mode behind the historical
// clobber bug. [CheckPairing], [CheckUniqueness], and [CheckLinks] are
// the cross-file invariants: they keep the root ↔ theme-file link graph
// 1:1 and every entry in exactly one place.
//
// The digesting pass that moves bodies and writes gists is NOT here; it
// is an agent-driven, human-gated skill in a later milestone, built on
// these guards. Nothing in this package writes a knowledge file.
//
// # Related packages
//
//   - internal/heading — parses the "## [ts] Title" entry blocks the
//     staging and uniqueness checks rely on.
//   - internal/config/disclosure — the structural heading vocabulary.
//   - internal/err/disclosure — the guard and invariant sentinels.
//
// # Concurrency
//
// Pure data plus read-only theme-file stat/read; callers own all writes.
package disclosure
