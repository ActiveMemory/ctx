//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package source contains the **rendering helpers** that
// turn a parsed AI session ([entity.Session]) into the
// markdown that `ctx journal source --show <id>` and
// `ctx journal import` write to disk.
//
// The package is split into two focused subpackages:
//
//   - **[format]** — small format primitives:
//     part-navigation links for multipart sessions
//     ([PartNavigation]), [Duration] for human-readable
//     time spans, and [ToolUse] for one-line tool-call
//     summaries.
//   - **[frontmatter]** — YAML frontmatter assembly:
//     heading resolution from session content, field
//     writing, and ordering so re-import produces
//     byte-identical output.
//
// The top-level `source.go` here defines the [Opts]
// flag-bag the `ctx journal source` subcommand fills in
// (`--show`, `--latest`, `--limit`, `--full`, `--project`,
// `--since`, etc.) and used by callers that need to ask
// "which session(s) does the user mean?".
//
// # Public Surface
//
//   - **[Opts]**             — flag-bag for source
//     selection.
//   - **[format]**           — see subpackage docs.
//   - **[frontmatter]**      — see subpackage docs.
//
// # Concurrency
//
// All helpers are pure data transformations over
// [entity.Session]. Concurrent callers never race.
//
// # Related Packages
//
//   - [internal/cli/journal/cmd/source]    — the
//     `ctx journal source` CLI surface.
//   - [internal/cli/journal/cmd/importer]  — the
//     importer that consumes the same renderer to
//     write enriched journal entries to disk.
//   - [internal/journal/parser]            — produces
//     the [entity.Session] this package renders.
package source
