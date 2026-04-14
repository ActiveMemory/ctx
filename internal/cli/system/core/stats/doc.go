//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package stats reads, parses, and formats **per-session
// token-usage statistics** for the `ctx usage` command
// and the system hooks that rely on token-pressure
// signals.
//
// Statistics are stored as JSON files under
// `.context/state/` — one per session, named for the
// session ID. Each record carries a timestamp, the
// running token count, the budget, and the percentage
// used.
//
// # Public Surface
//
//   - **[ReadDir](contextDir)** — returns every
//     session-stats file as a typed [Stats] slice,
//     sorted by mtime descending.
//   - **[ExtractSessionID](path)** — pulls the
//     session ID out of a stats filename.
//   - **[ParseFile](path)** — reads one stats
//     file and returns the typed [Stats].
//   - **[FormatDump](stats)** — renders the
//     collection as the human-readable table the
//     `ctx usage` command displays.
//   - **[FormatJSON](stats)** — renders the same
//     collection as a JSON document for tooling.
//
// # Concurrency
//
// Filesystem-bound. Concurrent reads are safe;
// writers are the per-session tracking hook
// which is single-process per session.
//
// # Related Packages
//
//   - [internal/cli/usage]                  — the
//     `ctx usage` CLI surface.
//   - [internal/cli/system/cmd/check_context_size]
//     — writes the stats this package reads.
//   - [internal/cli/system/cmd/heartbeat]   — also
//     contributes mtime updates.
//   - [internal/format]                     — the
//     Number / TimeAgo formatters used here.
package stats
