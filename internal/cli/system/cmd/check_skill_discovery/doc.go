//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_skill_discovery implements the
// **`ctx system check-skill-discovery`** hidden hook,
// which fires a one-shot nudge mid-session to remind
// the agent about available skills.
//
// # What It Does
//
// The hook monitors the per-session prompt counter
// (shared with check-context-size). When the counter
// reaches the configured threshold (typically around
// 25 prompts), it fires a single nudge surfacing
// mid-session skills that are easy to forget:
//
//   - /ctx-reflect
//   - /ctx-learning-add
//   - /ctx-decision-add
//   - /ctx-prompt-audit
//
// The nudge fires exactly once per session. A guard
// file is written after the first fire to prevent
// repeat nudges.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On threshold reached (first time): a nudge block
// listing the discoverable skills. On already fired,
// below threshold, or paused: no output.
//
// # Throttling
//
// One-shot per session using a per-session guard file
// in the state directory. Once fired, subsequent
// invocations are no-ops for that session.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], reads the prompt
// counter from [core/counter.Read], loads the skill
// list message via [core/message.Load], and emits
// the nudge through [write/setup.NudgeBlock].
package check_skill_discovery
