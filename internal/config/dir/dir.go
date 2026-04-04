//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dir

// Directory path constants used throughout the application.
const (
	// Archive is the subdirectory for archived tasks within .context/.
	Archive = "archive"
	// Claude is the Claude Code configuration directory in the project root.
	Claude = ".claude"
	// Context is the default context directory name.
	Context = ".context"
	// Hooks is the subdirectory for lifecycle hook scripts within .context/.
	Hooks = "hooks"
	// HooksMessages is the subdirectory path for hook message
	// overrides within .context/.
	HooksMessages = "hooks/messages"
	// Ideas is the project-root directory for early-stage ideas and explorations.
	Ideas = "ideas"
	// Journal is the subdirectory for journal entries within .context/.
	Journal = "journal"
	// JournalObsidian is the Obsidian export of journal entries within .context/.
	JournalObsidian = "journal-obsidian"
	// JournalSite is the journal static site output directory within .context/.
	JournalSite = "journal-site"
	// Logs is the subdirectory name for log files within the context directory.
	Logs = "logs"
	// Memory is the subdirectory for memory bridge files within .context/.
	Memory = "memory"
	// MemoryArchive is the archive subdirectory within .context/memory/.
	MemoryArchive = "memory/archive"
	// Projects is the projects subdirectory within .claude/.
	Projects = "projects"
	// Sessions is the subdirectory for session summaries within .context/.
	Sessions = "sessions"
	// Skills is the subdirectory for skill definitions within .context/.
	Skills = "skills"
	// Steering is the subdirectory for steering files within .context/.
	Steering = "steering"
	// Specs is the project-root directory for formalized plans and feature specs.
	Specs = "specs"
	// State is the subdirectory for project-scoped runtime state within .context/.
	State = "state"
	// Trace is the subdirectory for commit context trace data within .context/.
	Trace = "trace"
	// Templates is the subdirectory for entry templates within .context/.
	Templates = "templates"
	// CtxData is the user-level ctx data directory (~/.ctx/).
	CtxData = ".ctx"
	// DefaultSteeringPath is the default steering directory
	// path relative to the project root.
	DefaultSteeringPath = ".context/steering"
	// DefaultHooksPath is the default hooks directory path
	// relative to the project root.
	DefaultHooksPath = ".context/hooks"
)

// Platform-specific home directory path components.
const (
	// HomeLinux is the home directory parent on Linux (/home/username).
	HomeLinux = "home"
	// HomeMacOS is the home directory parent on macOS (/Users/username).
	HomeMacOS = "Users"
)

// Journal site output directories.
const (
	// JournalDocs is the docs subdirectory in the generated site.
	JournalDocs = "docs"
	// JournTopics is the topics subdirectory in the generated site.
	JournTopics = "topics"
	// JournalFiles is the key files subdirectory in the generated site.
	JournalFiles = "files"
	// JournalTypes is the session types subdirectory in the generated site.
	JournalTypes = "types"
)
