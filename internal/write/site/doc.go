//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package site provides terminal output for the site
// feed generation command (ctx site feed).
//
// The site command generates RSS/Atom feeds from
// journal entries and blog posts in the context
// directory. Output renders the generation summary
// after the feed file is written.
//
// # Output
//
// [PrintFeedReport] outputs the complete feed
// generation summary. It prints the output path
// and number of entries included, then optionally
// lists skipped entries and warnings.
//
// Skipped entries appear when a journal entry
// lacks required frontmatter or fails validation.
// Warnings cover non-fatal issues like missing
// dates or truncated content.
//
// The function accepts a [scan.FeedReport] struct
// that carries pre-computed counts and message
// lists so the output function contains no
// business logic.
package site
