//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides business logic for the site
// command, which generates static site assets from
// project blog content.
//
// The core package itself is a namespace that groups
// two subpackages: scan and rss. Together they form a
// pipeline that reads blog post markdown files, parses
// their frontmatter, and generates an Atom XML feed.
//
// # Subpackages
//
// The scan subpackage reads a directory of markdown
// blog posts, parses YAML frontmatter, validates
// required fields (title, date, finalized flag), and
// extracts a summary paragraph. It returns posts
// sorted by date descending along with a report of
// skipped and warned entries.
//
// The rss subpackage takes the parsed blog posts and
// generates an Atom 1.0 XML feed. It builds feed and
// entry elements with links, authors, categories, and
// summaries, then writes the XML to a file.
//
// # Data Flow
//
// The cmd/ layer calls scan.BlogPosts to collect post
// metadata, then passes the result to rss.Atom to
// generate the feed file. Warnings and skip reports
// are surfaced to the user through the write/ layer.
package core
