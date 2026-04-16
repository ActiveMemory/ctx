//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package lookup is the **eager-init in-memory cache** for
// every embedded YAML asset ctx ships: command-help text,
// flag-help text, prompt templates, examples, stop-word
// lists, glob patterns, default permission lists.
//
// The package is what sits between
// [internal/assets/read/desc] (the typed lookup helpers
// every CLI command calls) and the embedded YAML files
// (slow to parse on the hot path). Loading once at process
// start trades a few milliseconds at boot for fast lookups
// every time a hook fires.
//
// # Public Surface
//
//   - **[Init]** — loads every embedded YAML map into
//     memory. Called exactly once from `main()` before
//     the CLI starts dispatching. Idempotent: repeat
//     calls are noops (the [sync.Once] guard short-
//     circuits).
//   - **[TextDesc](key)** — resolves a
//     [internal/config/embed/text].DescKey to its
//     rendered string.
//   - **[StopWords]** — returns the embedded English
//     stop-word set used by
//     [internal/cli/agent/core/score] for relevance
//     scoring.
//   - **[ConfigPatterns]** — returns the embedded glob
//     pattern list used to detect "this file is a
//     config" in drift checks and skill heuristics.
//   - **[PermAllowListDefault]** /
//     **[PermDenyListDefault]** — return the default
//     allow/deny entries for Claude Code permissions
//     (used by `ctx init` and the
//     `_ctx-permission-sanitize` skill).
//
// # Why Eager Loading
//
// Lazy parsing per call would dominate the time budget
// for fast-fire hooks (some run on every tool call).
// One up-front parse means the per-call cost is just
// a map lookup. The maps are read-only after [Init];
// concurrent readers never race.
//
// # Concurrency
//
// All readers are safe for concurrent use after [Init]
// returns. The single-init guard ensures no race
// between concurrent first-callers.
package lookup
