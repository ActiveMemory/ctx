//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package frontmatter handles the **YAML frontmatter
// transformations** that journal entries undergo as they
// pass through the pipeline: importer → normalizer →
// renderer (site or vault).
//
// The package owns the per-renderer adapters that map the
// canonical [entity.JournalFrontmatter] into the slightly
// different shapes each downstream renderer expects.
//
// # Public Surface
//
//   - **[Transform](raw)** — converts a raw frontmatter
//     map (untyped, just-parsed YAML) into the
//     normalized journal frontmatter shape: enforces
//     field types, fills in defaults, drops fields the
//     schema does not recognize. Used by the importer
//     when ingesting hand-edited entries.
//   - **[ExtractStringSlice](m, key)** — safely pulls a
//     `[]string` from a `map[string]any`, tolerating
//     both `[]string` and `[]any` source types (YAML
//     decoders produce one or the other depending on
//     content). Returns nil when the key is missing.
//   - **[Obsidian]** — the Obsidian-vault frontmatter
//     struct: subset/extension of the canonical shape
//     with additional `aliases:`, `tags:`, and graph
//     metadata Obsidian renders.
//
// # Why a Separate Package
//
// Frontmatter handling looks trivial on the surface but
// is one of the most bug-prone surfaces in any markdown
// pipeline because YAML's loose typing produces
// `[]string` in some cases and `[]any` in others for
// "the same" structure. Hoisting the conversions here
// means every renderer benefits from the same
// safe-decode helpers.
//
// # Concurrency
//
// All functions are pure. Concurrent callers never
// race.
package frontmatter
