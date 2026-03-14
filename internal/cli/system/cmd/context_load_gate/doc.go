//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package context_load_gate implements the ctx system context-load-gate
// subcommand.
//
// It auto-injects project context into the agent's context window on the
// first tool use per session, with subsequent calls silently skipped.
package context_load_gate
