//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package extract pulls YAML frontmatter from journal
// markdown content.
//
// # Overview
//
// Journal entries begin with a YAML frontmatter block
// delimited by --- lines. This package provides
// functions to extract that block or strip it from the
// content, which is needed during import regeneration
// to preserve enriched metadata.
//
// # Public Surface
//
//   - [Frontmatter] -- returns the raw frontmatter
//     block including the --- delimiters and trailing
//     newline.
//   - [StripFrontmatter] -- removes the frontmatter
//     block and returns only the body content.
//
// # Algorithm
//
// Frontmatter works by:
//
//  1. Checking whether the content starts with the
//     opening delimiter (--- followed by newline).
//  2. Searching for the closing delimiter (newline,
//     ---, newline) after the opening.
//  3. Returning the substring from the start through
//     the closing delimiter, or empty string if no
//     valid frontmatter block is found.
//
// StripFrontmatter calls Frontmatter, then returns
// everything after the block with leading newlines
// trimmed. If no frontmatter exists, the original
// content is returned unchanged.
package extract
