//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import cfgMemory "github.com/ActiveMemory/ctx/internal/config/memory"

// CtxRC represents the configuration from the .ctxrc file.
//
// Fields:
//   - ContextDir: Name of the context directory (default ".context")
//   - TokenBudget: Default token budget for context assembly (default 8000)
//   - PriorityOrder: Custom file loading priority order
//   - AutoArchive: Whether to auto-archive completed tasks (default true)
//   - ArchiveAfterDays: Days before archiving completed tasks (default 7)
//   - ScratchpadEncrypt: Whether to encrypt the scratchpad (default true)
//   - AllowOutsideCwd: Skip boundary validation for
//     external context dirs (default false)
//   - InjectionTokenWarn: Token threshold for oversize
//     injection warning (default 15000, 0 = disabled)
//   - ContextWindow: Context window size in tokens for
//     usage reporting (default 200000). No-op for Claude
//     Code users: auto-detected from settings.json.
//     Only needed for non-Claude AI tools.
//   - BillingTokenWarn: Absolute token threshold for
//     billing nudge (default 0 = disabled). When set,
//     a one-shot VERBATIM warning fires the first time
//     session tokens exceed this value. Useful for Claude
//     Pro users with 1M context where tokens beyond the
//     included allowance incur extra cost.
//   - EventLog: Whether to log hook events locally
//     (default false)
//   - KeyRotationDays: Days before encryption key
//     rotation nudge (default 90)
//   - TaskNudgeInterval: Edit/Write calls between task
//     completion nudges (default 5, 0 = disabled)
//   - KeyPathOverride: Explicit encryption key file
//     path (default: auto-resolved)
//   - SessionPrefixes: Recognized session header
//     prefixes for Markdown parser (default: Session:)
//   - StaleAgeDays: Days before a context file is
//     flagged as stale by drift detection
//     (default 30, 0 = disabled)
//   - FreshnessFiles: Files to track for
//     technology-dependent constant staleness (opt-in)
//   - CompanionCheck: Check companion tool availability
//     during /ctx-remember (default true)
//   - ClassifyRules: Custom keyword rules for memory
//     entry classification (overrides defaults when set)
//   - SpecSignalWords: Terms that trigger a spec nudge
//     when adding tasks (overrides defaults when set)
//   - SpecNudgeMinLen: Task content length threshold
//     for spec nudge (default 150)
//   - Tool: Active AI tool identifier (e.g., claude,
//     cursor, cline, kiro, codex)
//   - Steering: Steering layer configuration overrides
//   - Hooks: Hook system configuration overrides
type CtxRC struct {
	Profile             string                   `yaml:"profile"`
	Tool                string                   `yaml:"tool"`
	ContextDir          string                   `yaml:"context_dir"`
	TokenBudget         int                      `yaml:"token_budget"`
	PriorityOrder       []string                 `yaml:"priority_order"`
	AutoArchive         bool                     `yaml:"auto_archive"`
	ArchiveAfterDays    int                      `yaml:"archive_after_days"`
	ScratchpadEncrypt   *bool                    `yaml:"scratchpad_encrypt"`
	AllowOutsideCwd     bool                     `yaml:"allow_outside_cwd"`
	EntryCountLearnings int                      `yaml:"entry_count_learnings"`
	EntryCountDecisions int                      `yaml:"entry_count_decisions"`
	ConventionLineCount int                      `yaml:"convention_line_count"`
	InjectionTokenWarn  int                      `yaml:"injection_token_warn"`
	ContextWindow       int                      `yaml:"context_window"`
	BillingTokenWarn    int                      `yaml:"billing_token_warn"`
	EventLog            bool                     `yaml:"event_log"`
	KeyRotationDays     int                      `yaml:"key_rotation_days"`
	TaskNudgeInterval   int                      `yaml:"task_nudge_interval"`
	KeyPathOverride     string                   `yaml:"key_path"`
	StaleAgeDays        int                      `yaml:"stale_age_days"`
	SessionPrefixes     []string                 `yaml:"session_prefixes"`
	FreshnessFiles      []FreshnessFile          `yaml:"freshness_files"`
	CompanionCheck      *bool                    `yaml:"companion_check"`
	ClassifyRules       []cfgMemory.ClassifyRule `yaml:"classify_rules"`
	SpecSignalWords     []string                 `yaml:"spec_signal_words"`
	SpecNudgeMinLen     int                      `yaml:"spec_nudge_min_len"`
	Notify              *NotifyConfig            `yaml:"notify"`
	Steering            *SteeringRC              `yaml:"steering"`
	Hooks               *HooksRC                 `yaml:"hooks"`
}

// FreshnessFile describes a source file containing technology-dependent
// constants that should be periodically reviewed.
//
// Fields:
//   - Path: File path relative to the project root
//   - Desc: Summary of what constants live in the file
//   - ReviewURL: Optional URL to check against when
//     reviewing (e.g., vendor docs)
type FreshnessFile struct {
	Path      string `yaml:"path"`
	Desc      string `yaml:"desc"`
	ReviewURL string `yaml:"review_url"`
}

// NotifyConfig holds webhook notification settings.
//
// KeyRotationDays is deprecated here; use the top-level CtxRC.KeyRotationDays
// instead. This field is retained for backwards compatibility with existing
// .ctxrc files that have key_rotation_days nested under notify.
// Fields:
//   - Events: Event filter list (loop, nudge, relay, heartbeat)
//   - KeyRotationDays: Deprecated; use top-level CtxRC.KeyRotationDays
type NotifyConfig struct {
	Events          []string `yaml:"events"`
	KeyRotationDays int      `yaml:"key_rotation_days"`
}

// SteeringRC holds steering layer configuration from .ctxrc.
//
// Fields:
//   - Dir: Path override for the steering directory
//     (default ".context/steering")
//   - DefaultInclusion: Default inclusion mode for new
//     steering files (default "manual")
//   - DefaultTools: Default tool identifier list for new
//     steering files (default: all tools)
type SteeringRC struct {
	Dir              string   `yaml:"dir"`
	DefaultInclusion string   `yaml:"default_inclusion"`
	DefaultTools     []string `yaml:"default_tools"`
}

// HooksRC holds hook system configuration from .ctxrc.
//
// Fields:
//   - Dir: Path override for the hooks directory
//     (default ".context/hooks")
//   - Timeout: Per-hook execution timeout in seconds
//     (default 10)
//   - Enabled: Whether hook execution is enabled
//     (default true). Pointer type distinguishes unset
//     (nil → true) from explicitly set to false.
type HooksRC struct {
	Dir     string `yaml:"dir"`
	Timeout int    `yaml:"timeout"`
	Enabled *bool  `yaml:"enabled"`
}
