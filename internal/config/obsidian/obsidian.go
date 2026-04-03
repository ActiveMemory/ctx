//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package obsidian

// Obsidian vault output directory constants.
const (
	// DirName is the default output directory for the Obsidian vault
	// within .context/.
	DirName = "journal-obsidian"
	// DirEntries is the subdirectory for journal entry files.
	DirEntries = "entries"
	// DirConfig is the Obsidian configuration directory name.
	DirConfig = ".obsidian"
)

// Obsidian file constants.
const (
	// AppConfigFile is the Obsidian app configuration filename.
	AppConfigFile = "app.json"
)

// Obsidian vault configuration content.
const (
	// AppConfig is the minimal app.json content that enforces
	// wikilink mode and shows frontmatter properties.
	AppConfig = `{
  "useMarkdownLinks": false,
  "showFrontmatter": true,
  "strictLineBreaks": false
}
`
)

// Obsidian wikilink format templates.
const (
	// WikilinkFmt formats a wikilink with display text.
	// Args: target, display.
	WikilinkFmt = "[[%s|%s]]"
	// WikilinkPlain formats a plain wikilink (no display text).
	// Args: target.
	WikilinkPlain = "[[%s]]"
)

// MaxRelated is the maximum number of "see also" entries in the
// related sessions footer.
const MaxRelated = 5

// Obsidian MOC (Map of Content) page filenames.
const (
	// MOCHome is the root navigation hub filename.
	MOCHome = "Home.md"
	// MOCTopics is the topics index MOC filename.
	MOCTopics = "_Topics.md"
	// MOCFiles is the key files index MOC filename.
	MOCFiles = "_Key Files.md"
	// MOCTypes is the session types index MOC filename.
	MOCTypes = "_Session Types.md"
)
