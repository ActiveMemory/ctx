//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package bootstrap implements the ctx system bootstrap hidden command.
//
// Prints the resolved context directory path for AI agents to
// anchor their session. Agent-only plumbing — no human types it
// interactively. Callable with --quiet for just the path,
// or --json for structured output.
//
// Key exports: [Cmd], [Run].
package bootstrap
