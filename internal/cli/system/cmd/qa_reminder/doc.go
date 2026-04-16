//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package qa_reminder implements the hidden
// "ctx system qa-reminder" cobra subcommand.
//
// This hook fires before any git command and injects
// a hard gate reminding the agent to lint, test, and
// verify a clean working tree before committing.
//
// # Behavior
//
// On each invocation the hook:
//
//   - Reads hook JSON from stdin to extract the tool
//     command and session ID.
//   - Checks whether the command string contains the
//     git binary name; non-git commands are ignored.
//   - Loads a templated gate message for the
//     qa-reminder variant, falling back to a built-in
//     default.
//   - Appends the context directory path to the
//     message and emits it as a PreToolUse context
//     block on stdout.
//   - Sends a relay notification through the webhook
//     channel.
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
// QA gate message. The agent sees this before the
// git command executes.
//
// # Delegation
//
// Message loading uses system/core/message. Context
// directory resolution uses context/resolve. Relay
// notification uses system/core/nudge.
package qa_reminder
