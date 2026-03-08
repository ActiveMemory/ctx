//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package post_commit implements the ctx system post-commit subcommand.
//
// It detects git commit commands and nudges the agent to capture context
// (decisions or learnings) and run lints/tests after committing.
package post_commit
