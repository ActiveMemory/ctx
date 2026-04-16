//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package remind implements **`ctx remind`**, the
// session-scoped reminders that surface automatically at
// session start and repeat until dismissed.
//
// Reminders are how a user (or skill) leaves a note for
// **the next session**: "remember to update CHANGELOG
// before merging", "the failing test in `auth_test.go`
// needs review", "Bob owes us the production credentials
// by Friday".
//
// # Subcommands
//
//   - **`ctx remind add <text>`**: appends a reminder
//     with optional `--after <YYYY-MM-DD>` date gate
//     and `--once` (auto-dismiss after first surface).
//   - **`ctx remind list`**: prints all open
//     reminders (with date-gating respected).
//   - **`ctx remind dismiss <id>`**: marks one
//     reminder dismissed (or `--all`).
//
// # The Surface Path
//
// At session start, the `check_reminder` system hook
// (`internal/cli/system/cmd/check_reminder`) reads the
// reminder store, filters by date and dismissal
// status, and emits the un-dismissed entries through
// the VERBATIM relay so the user (and the agent) both
// see them as the first interaction of the session.
//
// # Storage
//
// `.context/state/reminders.jsonl`: append-only,
// one [Reminder] per line. Dismissals are recorded as
// new lines (state-machine-over-log style); the
// reader collapses to the latest state per ID.
//
// # Concurrency
//
// Filesystem-bound and stateless. Concurrent writes
// would race on the JSONL append; ctx is single-process.
package remind
