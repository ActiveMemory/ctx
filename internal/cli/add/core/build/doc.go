//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package build assembles a cobra add subcommand bound to a
// specific noun (task / decision / learning / convention).
//
// Each noun-first parent (ctx task, ctx decision, ctx
// learning, ctx convention) registers an add subcommand by
// calling [Cmd] with its noun, description key, and Use
// string. The returned command shares all flag wiring and
// completion logic; the noun is prepended to the args slice
// before delegating to [run.Run].
//
// # Why a Builder?
//
// Without this helper each noun's cmd/add package would have
// to duplicate ~100 lines of cobra flag wiring (priority,
// section, file, context, rationale, consequence, lesson,
// application, session-id, branch, commit, share, plus
// completion). Centralizing here keeps the per-noun adapter
// down to a single Cmd() that selects the right description
// key.
package build
