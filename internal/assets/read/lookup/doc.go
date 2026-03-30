//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package lookup owns the YAML description maps and eager
// initialization for all embedded text lookups.
//
// [Init] loads all YAML files (commands, flags, text, examples)
// into in-memory maps. [TextDesc] resolves a text DescKey.
// [StopWords] returns the stop word set for relevance scoring.
// [ConfigPatterns] returns glob patterns for config file detection.
// [PermAllowListDefault] and [PermDenyListDefault] return the
// default permission lists for Claude Code settings.
package lookup
