//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package frontmatter handles YAML frontmatter transformation
// for journal entries and Obsidian vault generation.
//
// [Transform] converts raw frontmatter into a normalized format.
// [ExtractStringSlice] safely extracts []string from parsed YAML
// maps. The Obsidian struct provides the vault-specific frontmatter
// schema.
package frontmatter
