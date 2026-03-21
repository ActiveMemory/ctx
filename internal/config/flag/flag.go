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
	All             = "all"
	AllProjects     = "all-projects"
	Append          = "append"
	Archive         = "archive"
	Blob            = "blob"
	Build           = "build"
	Commands        = "commands"
	Completion      = "completion"
	DryRun          = "dry-run"
	Event           = "event"
	External        = "external"
	Fix             = "fix"
	Force           = "force"
	Full            = "full"
	Hook            = "hook"
	JSON            = "json"
	KeepFrontmatter = "keep-frontmatter"
	Key             = "key"
	Label           = "label"
	Latest          = "latest"
	Limit           = "limit"
	MaxIterations   = "max-iterations"
	Merge           = "merge"
	Message         = "message"
	Minimal         = "minimal"
	NoPluginEnable  = "no-plugin-enable"
	Out             = "out"
	Output          = "output"
	Prepend         = "prepend"
	Prompt          = "prompt"
	Ralph           = "ralph"
	Raw             = "raw"
	Regenerate      = "regenerate"
	Serve           = "serve"
	SessionID       = "session-id"
	Skills          = "skills"
	Stdin           = "stdin"
	Tool            = "tool"
	Type            = "type"
	Variant         = "variant"
	Write           = "write"
	Yes             = "yes"
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
	ShortYes           = "y"
)

// CLI flag names used in multiple commands.
const (
	Since = "since"
	Until = "until"
)
