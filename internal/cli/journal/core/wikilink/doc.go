//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package wikilink converts standard Markdown links to
// **Obsidian-style `[[wikilinks]]`** during vault export so
// Obsidian's graph view, backlinks, and unlinked-mentions
// features pick up the journal's cross-references natively.
//
// The package is one of the per-renderer adapters in the
// site/vault pipeline; the site renderer keeps standard
// `[text](url.md)` markdown, the vault renderer routes
// links through here.
//
// # Public Surface
//
//   - **[ConvertMarkdownLinks](text)**: rewrites every
//     `[text](url.md)` in `text` to `[[target|text]]`
//     (Obsidian's display-text wikilink form).
//     Preserves URLs that are not journal entries (raw
//     `https://...` links, anchor-only refs).
//   - **[Format](target, display)**: builds a single
//     wikilink string. `display` is optional; pass
//     empty to get `[[target]]`.
//   - **[FormatEntry](entry)**: convenience that
//     produces the canonical wikilink for a journal
//     entry using its slug as target and its title
//     as display text.
//
// # The "Why Obsidian Form" Question
//
// Obsidian's wikilinks resolve **by note name**, not by
// path. A vault expects `[[my-note]]` regardless of where
// `my-note.md` lives in the folder hierarchy. Standard
// markdown links break the moment the vault is
// reorganized; wikilinks survive.
//
// # Concurrency
//
// All functions are pure. Concurrent callers never
// race.
package wikilink
