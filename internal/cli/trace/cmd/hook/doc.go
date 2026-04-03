//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hook provides the "ctx trace hook" CLI subcommand.
//
// It enables or disables the git hooks that power context tracing:
// prepare-commit-msg (injects context trailer) and post-commit
// (records refs to history). The action argument must be "enable"
// or "disable".
//
// Key exports: [Cmd], [Run].
// Used by the trace command tree to register the hook subcommand.
package hook
