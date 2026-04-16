//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package post_commit implements the hidden
// "ctx system post-commit" cobra subcommand.
//
// This hook fires after a non-amend git commit and
// nudges the agent to capture context (a decision or
// learning) and to run lints and tests before pushing.
//
// # Behavior
//
// On each invocation the hook:
//
//   - Reads hook JSON from stdin to extract the tool
//     command and session ID.
//   - Ignores non-git-commit commands and amend
//     commits by matching against config/regex
//     patterns.
//   - Loads a templated nudge message for the
//     post-commit variant, falling back to a
//     built-in default.
//   - Appends the context directory path to the
//     message and emits it as a PostToolUse context
//     block on stdout.
//   - Sends a relay notification through the webhook
//     channel.
//   - Checks for version drift between the installed
//     binary and the latest release, appending a
//     warning block when a newer version exists.
//   - Scores the most recent commit against spec
//     enforcement rules, appending violation nudges
//     when any are detected.
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
// Emits a PostToolUse context block with the nudge
// message, optional version-drift warning, and
// optional spec-violation nudge.
//
// # Delegation
//
// Message loading uses system/core/message. Version
// drift checking uses system/core/drift. Spec scoring
// delegates to system/core/post_commit. Notification
// relay uses system/core/nudge.
package post_commit
