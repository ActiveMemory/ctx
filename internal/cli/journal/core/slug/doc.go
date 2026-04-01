//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package slug generates URL-safe identifiers from session titles.
//
// [FromTitle] converts a title to a lowercase hyphenated slug.
// [CleanTitle] removes non-alphanumeric characters. [ForTitle]
// generates both a slug and cleaned title for a session, handling
// deduplication against existing slugs.
package slug
