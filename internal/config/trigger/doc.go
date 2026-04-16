//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trigger defines lifecycle event types and
// hook testing constants for the ctx trigger system.
//
// Triggers fire at specific points in an AI session's
// lifecycle: before and after tool invocations, at
// session start and end, on file saves, and when
// context is added. This package declares the event
// type identifiers the dispatcher and hook configs
// use to wire up behavior.
//
// # Lifecycle Event Types
//
// The [TriggerType] alias and its constants define
// the supported events:
//
//   - [PreToolUse]    — fires before an AI tool
//     invocation (used for permission checks and
//     context injection).
//   - [PostToolUse]   — fires after an AI tool
//     invocation (used for context updates).
//   - [SessionStart]  — fires when a session begins
//     (used for bootstrap and status).
//   - [SessionEnd]    — fires when a session ends
//     (used for wrap-up nudges).
//   - [FileSave]      — fires when a file is saved.
//   - [ContextAdd]    — fires when context is added.
//
// # Hook Status Labels
//
//   - [StatusEnabled], [StatusDisabled] — display
//     labels for hook list output.
//
// # Mock Constants
//
//   - [MockSessionID], [MockModel], [MockVersion] —
//     fixed values used in hook test input to enable
//     deterministic testing.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package trigger
