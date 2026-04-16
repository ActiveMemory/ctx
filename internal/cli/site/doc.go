//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package site provides the "ctx site" command for
// managing the ctx.ist static site.
//
// The site command group hosts subcommands for building
// and maintaining the project's public website and its
// syndication feeds.
//
// # Subcommands
//
//   - feed: generates an Atom 1.0 feed from finalized
//     blog posts in docs/blog/. Parses YAML frontmatter
//     for title, date, author, and topics; extracts
//     summary from the first paragraph after the
//     heading. Drafts (reviewed_and_finalized != true)
//     are skipped.
//
// # Subpackages
//
//	cmd/feed -- Atom feed generation from blog posts
//	core -- blog post scanning, frontmatter parsing,
//	  and feed XML assembly
package site
