//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package journal implements the "ctx journal" command for analyzing and
// publishing exported AI session files.
//
// The journal system provides two output formats from .context/journal/ entries:
//
//   - ctx journal site: generates a zensical-compatible static site with
//     browsable session history, topic/file/type indices, and search.
//   - ctx journal obsidian: generates an Obsidian vault with wikilinks,
//     MOC (Map of Content) pages, and graph-optimized cross-linking.
//
// Both formats reuse the same scan/parse/index infrastructure and consume
// the same enriched journal entries (YAML frontmatter with topics, type,
// outcome, technologies, key_files).
package journal
