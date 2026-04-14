//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package obsidian implements **`ctx journal obsidian`** —
// the subcommand that exports the project's enriched
// journal entries as a complete **Obsidian vault** (with
// MOC pages, wikilinks, and graph-friendly frontmatter)
// for users who consume the journal in Obsidian rather
// than the zensical site.
//
// # Public Surface
//
//   - **[Cmd]** — cobra command with `--output` to
//     control the destination directory (default
//     `vault/`).
//   - **[Run]** — delegates to
//     [internal/cli/journal/core/obsidian.BuildVault]
//     which handles the full file generation pipeline
//     (scan, transform frontmatter, convert links to
//     `[[wikilinks]]`, build MOC pages, write
//     `Home.md`).
//
// # Why a Separate Vault
//
// Obsidian and the zensical site both consume the same
// raw entries but render them very differently
// (wikilinks vs markdown links, MOC vs topic index,
// graph view vs sidebar nav). Producing two output
// trees from one input set keeps each rendering
// idiomatic for its environment.
//
// # Concurrency
//
// Single-process, sequential. `O(N)` over journal
// entries.
//
// # Related Packages
//
//   - [internal/cli/journal/core/obsidian] — the
//     vault-building engine.
//   - [internal/cli/journal/core/wikilink] — markdown
//     → wikilink conversion.
//   - [internal/cli/journal/core/frontmatter] —
//     Obsidian-flavored frontmatter assembly.
//   - [internal/cli/journal/core/moc]      — MOC
//     pages (Obsidian flavor).
//   - [internal/cli/journal/cmd/site]      — sister
//     command for the zensical-flavored output.
package obsidian
