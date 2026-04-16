//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package obsidian builds an **Obsidian vault** from the
// project's enriched journal entries, the engine behind
// the user-facing `ctx journal obsidian` command.
//
// The vault is a complete Obsidian-friendly directory
// tree: per-entry notes with vault-specific frontmatter,
// `[[wikilinks]]` instead of markdown links, MOC pages
// for navigation, and a `Home.md` landing page that
// surfaces recent entries and top topics.
//
// # Public Surface
//
//   - **[BuildVault](journalDir, vaultDir, opts)**:
//     end-to-end pipeline: scan entries (parse),
//     create directory structure, transform
//     frontmatter (frontmatter), convert links
//     (wikilink), build MOC pages (moc), write
//     `Home.md`. Idempotent: re-running with the
//     same inputs produces byte-identical output.
//
// # Layout Produced
//
//   - `<vault>/Home.md`: landing MOC
//   - `<vault>/MOC.md`: topics overview
//   - `<vault>/topics/<slug>.md`: per-topic pages
//   - `<vault>/<YYYY>/<MM>/<slug>.md`: entries
//
// # Concurrency
//
// Single-process, sequential. `O(N)` over journal
// entries.
package obsidian
