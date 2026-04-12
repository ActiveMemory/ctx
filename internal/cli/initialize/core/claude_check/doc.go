//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude_check detects the state of Claude Code and
// the ctx plugin so `ctx init` and `ctx setup claude-code`
// can print stage-aware setup guidance.
//
// The detector answers four questions, ordered:
//
//  1. Is the `claude` binary on PATH?
//  2. Is the ctx plugin registered in
//     ~/.claude/plugins/installed_plugins.json?
//  3. Is the plugin enabled globally or locally?
//  4. (derived) Is the setup ready to use?
//
// Key exports: [State], [Detect].
package claude_check
