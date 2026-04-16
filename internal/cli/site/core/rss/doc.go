//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rss generates Atom 1.0 XML feeds from blog
// post metadata collected by the scan subpackage.
//
// # Feed Generation
//
// [Atom] is the sole entry point. It takes a slice of
// parsed blog posts, an output file path, and a base
// URL. The function:
//
//  1. Constructs an AtomFeed with the configured title,
//     namespace, and self/alternate links derived from
//     the base URL.
//  2. Sets the feed updated timestamp to the most
//     recent post date.
//  3. Builds an AtomEntry for each post with title,
//     link, ID, updated timestamp, optional summary,
//     author (defaulting to a configured fallback),
//     and category tags from topics.
//  4. Marshals the feed as indented XML with an XML
//     declaration header.
//  5. Creates the output directory if needed and writes
//     the file.
//
// Summaries containing a skip sentinel string are
// omitted from the entry to avoid including placeholder
// text in the published feed.
//
// # Types
//
// The package defines five types that map to Atom XML
// elements: [AtomFeed], [AtomEntry], [AtomLink],
// [AtomAuthor], and [AtomCategory]. Each uses xml
// struct tags for serialization.
package rss
