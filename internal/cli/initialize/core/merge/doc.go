//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package merge implements the **create-or-merge** file
// operations that make `ctx init` safely idempotent: each
// foundation file is either created from a template or, if
// already present, has only its **ctx-managed marker
// section** updated — never the user-edited surrounding
// content.
//
// The package solves the "I edited my CONSTITUTION.md and
// re-ran `ctx init` — did I lose my edits?" problem by
// making "yes, you keep them" the only possible answer.
//
// # Public Surface
//
//   - **[OrCreate](path, template, vars)** — file does
//     not exist → write the template (with `vars`
//     interpolated). File exists → run
//     [UpdateMarkedSection] on it. Always creates a
//     timestamped `.bak` before writing. Returns a
//     report indicating which path was taken.
//   - **[UpdateMarkedSection](existing, newSection,
//     start, end)** — finds the `start` and `end` marker
//     lines in `existing` and replaces only the content
//     between them. If the markers are missing, the
//     section is inserted at the bottom of the file with
//     the markers added so the next run becomes a true
//     in-place update.
//   - **[SettingsPermissions](path, allow, deny)** —
//     specialized merger for Claude Code permission
//     lists; preserves user-added entries while ensuring
//     the ctx-required entries are present.
//   - **[Permissions](existing, additions)** —
//     deduplicating list union used by the settings
//     merger and by `_ctx-permission-sanitize`.
//
// # Marker Convention
//
// ctx-managed sections are bracketed by HTML-comment
// markers:
//
//	<!-- ctx:section-name -->
//	... ctx-managed content ...
//	<!-- ctx:section-name:end -->
//
// The markers are invisible in rendered Markdown but
// trivially greppable. Constants for the well-known
// pairs live in [internal/config/marker].
//
// # Backup Policy
//
// Every write goes through a timestamped backup
// (`<path>.bak.YYYY-MM-DD-HHMMSS`). Backups accumulate;
// `ctx prune` cleans them on schedule. The trade-off is
// disk space for accident recovery, which the user
// always wants.
//
// # Concurrency
//
// Filesystem-bound and stateless; serialized through
// process-level execution.
package merge
