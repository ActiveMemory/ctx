//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package specs_nudge implements the hidden
// "ctx system specs-nudge" cobra subcommand.
//
// This hook fires on PreToolUse events and reminds
// the agent to save implementation plans to the
// specs/ directory when a new implementation is
// detected.
//
// # Behavior
//
// On each invocation the hook:
//
//   - Reads hook JSON from stdin to extract the
//     session ID.
//   - Loads a templated nudge message for the
//     specs-nudge variant, falling back to a
//     built-in default.
//   - Appends the context directory path to the
//     message and emits it as a PreToolUse context
//     block on stdout.
//   - Sends a relay notification through the webhook
//     channel with a template ref for the nudge
//     variant.
//
// The hook is skipped when the context directory is
// not initialized or the session is paused.
//
// # Flags
//
// None. The command reads hook JSON from stdin.
//
// # Output
//
// Emits a PreToolUse context block containing the
// specs nudge message. The agent sees this before
// the tool call executes.
//
// # Delegation
//
// Message loading uses system/core/message. Context
// directory resolution uses context/resolve. Relay
// notification uses system/core/nudge.
package specs_nudge
