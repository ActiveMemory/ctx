//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package unknown holds the RunE the `ctx system` group installs so
// that an unrecognised subcommand fails loud and legible instead of
// dumping help at exit 0.
//
// # Why this exists
//
// `ctx system` is a grouping command. Cobra raises an
// "unknown command" error only for the *root*; for a non-root group
// an unmatched subcommand falls through to the group's own Run/RunE,
// and a group with neither prints help and returns nil — exit 0. In
// a Claude Code UserPromptSubmit hook, exit 0 reads as "hook
// success", so the ~51-line help blob is injected into the agent's
// context every prompt. That is exactly how a stale `hooks.json`
// wiring `ctx system check-anchor-drift` (a command the binary later
// deleted) polluted sessions silently.
//
// A non-zero exit alone does not fix it: the harness swallows a
// failed hook's exit code, so the signal has to travel on hook
// stdout. [Handler] therefore emits a verbatim-relay box naming the
// unknown verb and hinting at the likely cause (plugin/binary
// version skew), best-effort records the event (event log + webhook)
// when a session is present on stdin, suppresses cobra's help dump,
// and returns a non-nil error.
//
// Scope is `ctx system` only — the single group wired into
// hooks.json. The shared [parent.Cmd] is untouched; other groups
// keep cobra's default behavior. See
// specs/system-unknown-subcommand-relay.md.
package unknown
