//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package permission implements **`ctx permission`** — the
// CLI for capturing and restoring **golden-image
// snapshots** of `.claude/settings.local.json` so a team
// can maintain a curated permission baseline and
// automatically drop session-accumulated permissions at
// session start.
//
// # The Problem It Solves
//
// During a Claude Code session, the user often grants
// one-off permissions (a `Bash(git commit*)`, a
// `Read(/tmp/...)`) that they did not mean to keep
// permanently. By session end the `allow:` list has
// drifted from the team's intended baseline. Manually
// pruning the file every few days is tedious and
// error-prone.
//
// # The Workflow
//
//  1. **Snapshot** — once, after a careful curation
//     pass, the user runs `ctx permission snapshot`.
//     The current `settings.local.json` is copied to
//     `.context/permissions.golden.json` and committed
//     to git as the team's baseline.
//  2. **Restore** — at session start (often via the
//     `_ctx-permission-sanitize` skill or a simple
//     `make` target), `ctx permission restore` resets
//     `settings.local.json` to the golden image.
//     Today's session starts clean.
//  3. **Iterate** — when the user finds a
//     permission they actually want to keep, they
//     re-snapshot to lock it in.
//
// # Subcommands
//
//   - **snapshot** — copy `settings.local.json` →
//     `permissions.golden.json` (overwrites previous
//     golden image; the git history is the safety
//     net).
//   - **restore** — copy `permissions.golden.json` →
//     `settings.local.json` (creates a `.bak` first).
//
// # Concurrency
//
// Filesystem-bound and stateless. Concurrent
// invocations would race on the destination file;
// single-process is the assumed model.
package permission
