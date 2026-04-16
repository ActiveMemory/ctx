//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_task_completion implements the
// **`ctx system check-task-completion`** hidden hook,
// which periodically nudges the agent to update
// TASKS.md after tool use.
//
// # What It Does
//
// The hook maintains a per-session prompt counter that
// increments on each post-tool-use invocation. When the
// counter reaches the configured interval (read from
// .ctxrc via [rc.TaskNudgeInterval]), the hook resets
// the counter and emits a context-level nudge asking
// the agent to review TASKS.md and mark any completed
// items as done.
//
// The nudge interval is configurable and can be
// disabled by setting it to zero or a negative value.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On interval reached: a context message reminding
// the agent to check task completion status. On below
// interval, disabled, or paused: no output.
//
// # Throttling
//
// The per-session counter resets after each nudge,
// so the nudge fires repeatedly at the configured
// interval throughout the session.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], manages the counter
// through [core/counter], loads the nudge message via
// [core/message.Load], and emits it as a context
// block through [write/setup.Context]. Relay
// notifications go through [core/nudge.Relay].
package check_task_completion
