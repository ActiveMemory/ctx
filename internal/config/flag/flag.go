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
	Follow   = "follow"
	Format   = "format"
	Session  = "session"
)

// Shared flag names used across commands.
const (
	After           = "after"
	All             = "all"
	AllProjects     = "all-projects"
	Append          = "append"
	Archive         = "archive"
	BaseURL         = "base-url"
	Blob            = "blob"
	Build           = "build"
	Check           = "check"
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
	Last            = "last"
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
	Project         = "project"
	Prompt          = "prompt"
	Quiet           = "quiet"
	Raw             = "raw"
	Regenerate      = "regenerate"
	Scope           = "scope"
	Serve           = "serve"
	Show            = "show"
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
	ShortAfter         = "a"
	ShortAll           = "A"
	ShortCompletion    = "c"
	ShortEvent         = "e"
	ShortForce         = "f"
	ShortHook          = "K"
	ShortFollow        = "f"
	ShortJSON          = "j"
	ShortKey           = "k"
	ShortLast          = "n"
	ShortMaxIterations = "M"
	ShortMinimal       = "m"
	ShortOutput        = "o"
	ShortQuiet         = "q"
	ShortProject       = "p"
	ShortPrompt        = "p"
	ShortSessionID     = "s"
	ShortShow          = "S"
	ShortTool          = "t"
	ShortWrite         = "w"
	ShortYes           = "y"
)

// CLI flag names used in multiple commands.
const (
	Since = "since"
	Until = "until"
)
