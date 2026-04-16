//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hook implements the "ctx trace hook" cobra
// subcommand.
//
// This command enables or disables the git hooks that
// power context tracing. When enabled, the hooks
// automatically inject context trailers into commit
// messages and record refs to history after each
// commit.
//
// # Usage
//
//	ctx trace hook <enable|disable>
//
// # Arguments
//
// Exactly one positional argument is required:
//
//   - action: must be "enable" or "disable". Any
//     other value returns an error.
//
// # Behavior
//
// When the action is "enable":
//
//   - Installs a prepare-commit-msg hook that calls
//     "ctx trace collect" to inject a context trailer
//     into the commit message.
//   - Installs a post-commit hook that calls
//     "ctx trace collect --record" to persist the
//     trailer refs into trace history.
//
// When the action is "disable":
//
//   - Removes both git hooks, stopping automatic
//     context tracing.
//
// # Output
//
// Prints a confirmation message indicating whether
// the hooks were enabled or disabled.
//
// # Delegation
//
// Hook installation and removal are handled by
// trace/core/hook. Action constants are defined in
// config/trace.
package hook
