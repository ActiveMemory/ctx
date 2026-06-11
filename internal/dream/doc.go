//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dream is the executor-agnostic engine for the ctx-dream
// memory-consolidation feature. It owns the data contract — proposals,
// per-source state, the append-only ledger — and the structural safety
// guards (write-scope, don't-leak) that any executor must enforce.
//
// The dream only ever PROPOSES (Option B): nothing here writes canonical
// memory. Cognition (classify, ground, propose) runs in the ctx-dream
// skill under an executor (cron `claude -p` is the reference); this
// package provides the types, persistence, and guards those passes build
// on, so a Claude Code hook and a raw-API loop can enforce the same
// invariants. See specs/ctx-dream.md.
package dream
