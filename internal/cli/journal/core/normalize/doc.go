//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package normalize sanitizes journal markdown for two
// downstream renderers: the **zensical site builder** and the
// **Obsidian vault exporter**. Raw enriched journal entries
// often carry constructs that one or both renderers cannot
// handle (or render confusingly): unbalanced code fences, H1
// headings that collide with the page title, raw HTML, etc.
// This package smoothes those out without losing meaning.
//
// # The Transformations
//
// The main `Content(text, opts)` entry point performs, in
// order:
//
//  1. **Fence stripping at boundaries**: orphan opening or
//     closing fences left over from incomplete code blocks
//     are removed so the renderer does not enter "code
//     mode" for the rest of the document.
//  2. **Heading demotion**: every H1 in the body is
//     demoted to H2 so it does not collide with the
//     frontmatter-derived page title rendered by both
//     zensical and Obsidian.
//  3. **HTML escaping**: bare `<tag>` patterns that are
//     not legitimate HTML are escaped so they do not get
//     swallowed silently.
//  4. **Turn-boundary normalization**: turn headers like
//     `## [12:34:56] User:` are recognized and given a
//     consistent shape via [MatchTurnHeader] /
//     [FindTurnBoundary] so the per-turn navigator on the
//     site can find them.
//  5. **Trim**: leading and trailing blank-line runs are
//     reduced to a single blank line via [TrimBlankLines].
//
// # The Public Helpers
//
//   - **[MatchTurnHeader](line)**: returns true plus the
//     parsed turn role + timestamp when the line matches
//     the canonical turn-header shape.
//   - **[FindTurnBoundary](lines, start)**: locates the
//     index of the next turn boundary at or after `start`,
//     used for slicing out a specific turn.
//   - **[TrimBlankLines](lines)**: strips leading and
//     trailing blank entries from a `[]string`.
//
// # Idempotency
//
// Every transformation is **idempotent**: running
// `Content` twice in a row produces no further changes.
// This is what makes the package safe to call from both
// the import pipeline (writes the normalized form to disk)
// and the renderers (re-normalize what they read).
//
// # Concurrency
//
// All exported functions are pure data transformations
// over `string` / `[]string`. Concurrent callers never
// race.
package normalize
