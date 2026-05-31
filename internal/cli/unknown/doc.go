//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package unknown provides a per-group, opt-in RunE that makes an
// unrecognised subcommand fail loud and legible instead of dumping help
// at exit 0. A group opts in with `c.RunE = unknown.HandlerFor(cfg)`.
//
// # Why this exists
//
// A cobra grouping command raises an "unknown command" error only for
// the *root*; for a non-root group an unmatched subcommand falls
// through to the group's own Run/RunE, and a group with neither prints
// help and returns nil — exit 0. That silent exit-0 is the failure mode
// this package kills, in two shapes:
//
//   - `ctx system` is wired into Claude Code's hooks.json. There, exit 0
//     reads as "hook success", so a stale wiring (e.g.
//     `ctx system check-anchor-drift`, a verb the binary later deleted)
//     injects the ~51-line help blob into the agent's context every
//     prompt. See specs/system-unknown-subcommand-relay.md.
//   - `ctx hook` is consumed by name from skills and loop scripts
//     (`ctx hook event|message|pause|resume|notify`). If such a verb
//     drifts out of the binary, the caller silently receives help at
//     exit 0 — the agent misreads or ignores it, and for
//     `ctx hook notify` the human is never told. See
//     specs/unknown-subcommand-relay-generalization.md.
//
// # Behavior
//
// A non-zero exit alone does not fix the hooks.json case: the harness
// swallows a failed hook's exit code, so the signal must travel on hook
// stdout. [handle] therefore emits a verbatim-relay box naming the
// unknown verb and its likely cause, best-effort records the event
// (event log + webhook) when a session is present on stdin, suppresses
// cobra's help dump, and returns a non-nil error. A bare group
// invocation (no leftover args) still prints help and exits 0.
//
// # Parameterization
//
// [Config] supplies the per-group copy (relay prefix, box title, body,
// relay message) and relay ref. [SystemConfig] and [HookConfig] are the
// two opt-ins today; the shared [parent.Cmd] stays untouched, so groups
// that want cobra's default keep it.
package unknown
