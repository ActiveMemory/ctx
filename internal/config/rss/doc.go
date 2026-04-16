//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rss defines constants for generating Atom XML
// feeds from the ctx blog.
//
// The ctx documentation site publishes an Atom feed at
// /feed.xml so users can subscribe to blog updates. This
// package centralizes every string the feed generator
// needs: source and output paths, the XML namespace URI,
// default author metadata, and URL-construction helpers.
//
// # Key Constants
//
//   - [DefaultFeedInputDir]: where blog Markdown lives
//     (docs/blog).
//   - [DefaultFeedOutPath]: where the compiled feed is
//     written (site/feed.xml).
//   - [DefaultFeedBaseURL]: root URL prepended to
//     entry permalinks (https://ctx.ist).
//   - [FeedAtomNS]: the Atom XML namespace.
//   - [FeedXMLHeader]: the XML declaration line
//     prepended to every feed file.
//   - [FeedTitle], [FeedDefaultAuthor]: metadata
//     defaults.
//
// URL path constants ([FeedPath], [BlogPath],
// [LinkRelSelf]) drive link generation inside the
// feed entries.
//
// # Why Centralized
//
// The feed generator, the blog build pipeline, and
// the site configuration all need the same paths and
// URLs. Keeping them in one place prevents silent
// mismatches when a path changes.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package rss
