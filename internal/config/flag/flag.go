//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package flag

// Global CLI flag names.
const (
	ContextDir      = "context-dir"
	AllowOutsideCwd = "allow-outside-cwd"
)

// PrefixLong is a CLI flag prefix for display formatting.
const PrefixLong = "--"

// Add command flag names: used for both flag registration and error display.
const (
	Application = "application"
	Consequence = "consequence"
	Context     = "context"
	File        = "file"
	Lesson      = "lesson"
	Priority    = "priority"
	Rationale   = "rationale"
	Section     = "section"
)

// Flag shorthand letters for the add command.
const (
	ShortApplication = "a"
	ShortContext     = "c"
	ShortFile        = "f"
	ShortLesson      = "l"
	ShortPriority    = "p"
	ShortRationale   = "r"
	ShortSection     = "s"
)

// Agent command flag names.
const (
	Budget   = "budget"
	Cooldown = "cooldown"
	Format   = "format"
	Session  = "session"
)

// Shared flag names used across commands.
const (
	AllProjects    = "all-projects"
	Append         = "append"
	Archive        = "archive"
	Blob           = "blob"
	Build          = "build"
	Commands       = "commands"
	Completion     = "completion"
	DryRun         = "dry-run"
	Event          = "event"
	External       = "external"
	Fix            = "fix"
	Force          = "force"
	Hook           = "hook"
	JSON           = "json"
	Key            = "key"
	Label          = "label"
	MaxIterations  = "max-iterations"
	Merge          = "merge"
	Message        = "message"
	Minimal        = "minimal"
	NoPluginEnable = "no-plugin-enable"
	Out            = "out"
	Output         = "output"
	Prepend        = "prepend"
	Prompt         = "prompt"
	Ralph          = "ralph"
	Raw            = "raw"
	Serve          = "serve"
	SessionID      = "session-id"
	Skills         = "skills"
	Tool           = "tool"
	Type           = "type"
	Variant        = "variant"
	Write          = "write"
)

// Shorthand letters for shared flags.
const (
	ShortCompletion    = "c"
	ShortEvent         = "e"
	ShortForce         = "f"
	ShortJSON          = "j"
	ShortKey           = "k"
	ShortMaxIterations = "n"
	ShortMinimal       = "m"
	ShortOutput        = "o"
	ShortPrompt        = "p"
	ShortSessionID     = "s"
	ShortTool          = "t"
	ShortWrite         = "w"
)

// CLI flag names used in multiple commands.
const (
	Since = "since"
	Until = "until"
)
