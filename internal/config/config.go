package config

const (
	ContextDirName         = ".context"
	ClaudeDirName          = ".claude"
	ClaudeHooksDirName     = ".claude/hooks"
	SettingsFileName       = ".claude/settings.local.json"
	AutoSaveScriptName     = "auto-save-session.sh"
	BlockNonPathScriptName = "block-non-path-ctx.sh"
	ClaudeMdFileName       = "CLAUDE.md"
	CtxMarkerStart         = "<!-- ctx:context -->"
	CtxMarkerEnd           = "<!-- ctx:end -->"
)

// FileType maps short names to actual file names
var FileType = map[string]string{
	"decision":    "DECISIONS.md",
	"decisions":   "DECISIONS.md",
	"task":        "TASKS.md",
	"tasks":       "TASKS.md",
	"learning":    "LEARNINGS.md",
	"learnings":   "LEARNINGS.md",
	"convention":  "CONVENTIONS.md",
	"conventions": "CONVENTIONS.md",
}

// FileReadOrder defines the priority order for reading context files.
var FileReadOrder = []string{
	"CONSTITUTION.md",
	"TASKS.md",
	"CONVENTIONS.md",
	"ARCHITECTURE.md",
	"DECISIONS.md",
	"LEARNINGS.md",
	"GLOSSARY.md",
	"DRIFT.md",
	"AGENT_PLAYBOOK.md",
}
