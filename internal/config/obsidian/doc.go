//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package obsidian defines directory, file, and format
// constants for exporting journal entries into an
// Obsidian vault.
//
// # Vault Structure
//
// The exporter writes into DirName ("journal-obsidian")
// under .context/. Inside that directory:
//
//   - DirEntries ("entries") holds individual journal
//     entry Markdown files.
//   - DirConfig (".obsidian") holds Obsidian app
//     configuration.
//
// AppConfigFile ("app.json") is written with the
// AppConfig content, which enforces wikilink mode,
// shows frontmatter properties, and disables strict
// line breaks.
//
// # Map of Content (MOC) Pages
//
// Four index pages provide navigation hubs:
//
//   - MOCHome ("Home.md"): root navigation hub.
//   - MOCTopics ("_Topics.md"): topics index.
//   - MOCFiles ("_Key Files.md"): key files index.
//   - MOCTypes ("_Session Types.md"): session type
//     index.
//
// # Wikilink Formatting
//
// WikilinkFmt and WikilinkPlain provide printf format
// strings for wikilinks with and without display text.
// These produce the [[target|display]] and [[target]]
// forms that Obsidian uses for internal links.
//
// # Related Entries
//
// MaxRelated (5) caps the number of "see also" links
// in the related sessions footer of each entry page.
package obsidian
