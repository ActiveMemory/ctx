//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package flag

// PrefixLong is a CLI flag prefix for display formatting.
const PrefixLong = "--"

// Activate / deactivate command flag names.
const (
	// Shell selects the shell dialect for `ctx activate` and
	// `ctx deactivate` emitters. When unset, the command
	// auto-detects from $SHELL, falling back to bash.
	Shell = "shell"
)

// Add command flag names: used for both flag registration and error display.
const (
	Application = "application"
	Branch      = "branch"
	Commit      = "commit"
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
	Skill    = "skill"
)

// Shared flag names used across commands.
const (
	After       = "after"
	All         = "all"
	AllProjects = "all-projects"
	Append      = "append"
	Archive     = "archive"
	BaseURL     = "base-url"
	Blob        = "blob"
	Build       = "build"
	Caller      = "caller"
	Check       = "check"
	Commands    = "commands"
	Completion  = "completion"
	Daemon      = "daemon"
	DataDir     = "data-dir"
	Days        = "days"
	Dir         = "dir"
	DryRun      = "dry-run"
	Event       = "event"

	IncludeHub      = "include-hub"
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
	Note            = "note"
	Message         = "message"
	Minimal         = "minimal"
	NoPluginEnable  = "no-plugin-enable"
	NoSteeringInit  = "no-steering-init"
	Out             = "out"
	Output          = "output"
	Path            = "path"
	Prepend         = "prepend"
	Project         = "project"
	Prompt          = "prompt"
	Quiet           = "quiet"
	Raw             = "raw"
	Record          = "record"
	Regenerate      = "regenerate"
	Scope           = "scope"
	Peers           = "peers"
	Port            = "port"
	Serve           = "serve"
	Share           = "share"
	Show            = "show"
	SessionID       = "session-id"
	Skills          = "skills"
	Tag             = "tag"
	Tool            = "tool"
	Token           = "token"
	Type            = "type"
	Variant         = "variant"
	Verbose         = "verbose"
	Width           = "width"
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
	ShortTag           = "t"
	ShortTool          = "t"
	ShortVerbose       = "v"
	ShortWrite         = "w"
	ShortYes           = "y"
)

// CLI flag names used in multiple commands.
const (
	Log   = "log"
	Since = "since"
	Until = "until"
)
