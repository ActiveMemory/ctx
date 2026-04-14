//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package site implements **`ctx journal site`** — the
// subcommand that turns the project's enriched journal
// entries into a browsable static site, optionally
// invoking the zensical builder to produce the HTML.
//
// # Public Surface
//
//   - **[Cmd]** — cobra command with `--build` (also
//     run zensical) and `--output` (override the
//     destination directory).
//   - **[Run]** — orchestrates the full generation:
//     parse entries (parse), normalize each (normalize),
//     build month-grouped pages and topic indexes
//     (section + generate + moc), write the zensical
//     `README.md` (generate.SiteReadme), and — when
//     `--build` is set — shell out to `zensical build`.
//
// # Output Layout
//
//   - `<output>/README.md`        — zensical config
//   - `<output>/index.md`         — chronological index
//   - `<output>/topics/index.md`  — topic overview MOC
//   - `<output>/topics/<slug>.md` — per-topic pages
//   - `<output>/<YYYY>/<MM>/<slug>.md` — entries
//
// # Concurrency
//
// Single-process, sequential. The site build is
// `O(N)` over journal entries and typically
// completes in seconds.
//
// # Related Packages
//
//   - [internal/cli/journal/core/section] — topic
//     index builders.
//   - [internal/cli/journal/core/moc]     — Map of
//     Content pages.
//   - [internal/cli/journal/core/generate] — top-
//     level page templates.
//   - [internal/cli/journal/core/normalize] — runs
//     per-entry before rendering.
//   - [internal/cli/serve]                  — the
//     `ctx serve` command that hosts the built site.
package site
