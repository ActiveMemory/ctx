//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

// Directory path constants used throughout the application.
const (
	// DirArchive is the subdirectory for archived tasks within .context/.
	DirArchive = "archive"
	// DirClaude is the Claude Code configuration directory in the project root.
	DirClaude = ".claude"
	// DirClaudeHooks is the hooks subdirectory within .claude/.
	DirClaudeHooks = ".claude/hooks"
	// DirContext is the default context directory name.
	DirContext = ".context"
	// DirJournal is the subdirectory for journal entries within .context/.
	DirJournal = "journal"
	// DirTools is the subdirectory for tool scripts within .context/.
	DirTools = "tools"
	// DirJournalSite is the journal static site output directory within .context/.
	DirJournalSite = "journal-site"
	// DirSessions is the subdirectory for session summaries within .context/.
	DirSessions = "sessions"
)

// GitignoreEntries lists the recommended .gitignore entries added by ctx init.
var GitignoreEntries = []string{
	".context/journal/",
	".context/journal-site/",
	".context/journal-obsidian/",
	".context/logs/",
	".context/.context.key",
	".context/.scratchpad.key",
	".claude/settings.local.json",
}

// Journal site output directories.
const (
	// JournalDirDocs is the docs subdirectory in the generated site.
	JournalDirDocs = "docs"
	// JournalDirTopics is the topics subdirectory in the generated site.
	JournalDirTopics = "topics"
	// JournalDirFiles is the key files subdirectory in the generated site.
	JournalDirFiles = "files"
	// JournalDirTypes is the session types subdirectory in the generated site.
	JournalDirTypes = "types"
)
