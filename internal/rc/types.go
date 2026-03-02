//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

// CtxRC represents the configuration from the .ctxrc file.
//
// Fields:
//   - ContextDir: Name of the context directory (default ".context")
//   - TokenBudget: Default token budget for context assembly (default 8000)
//   - PriorityOrder: Custom file loading priority order
//   - AutoArchive: Whether to auto-archive completed tasks (default true)
//   - ArchiveAfterDays: Days before archiving completed tasks (default 7)
//   - ScratchpadEncrypt: Whether to encrypt the scratchpad (default true)
//   - AllowOutsideCwd: Skip boundary validation for external context dirs (default false)
//   - InjectionTokenWarn: Token threshold for oversize injection warning (default 15000, 0 = disabled)
//   - ContextWindow: Context window size in tokens for usage reporting (default 200000)
//   - EventLog: Whether to log hook events locally (default false)
//   - KeyRotationDays: Days before encryption key rotation nudge (default 90)
//   - KeyPathOverride: Explicit encryption key file path (default: auto-resolved)
type CtxRC struct {
	ContextDir          string        `yaml:"context_dir"`
	TokenBudget         int           `yaml:"token_budget"`
	PriorityOrder       []string      `yaml:"priority_order"`
	AutoArchive         bool          `yaml:"auto_archive"`
	ArchiveAfterDays    int           `yaml:"archive_after_days"`
	ScratchpadEncrypt   *bool         `yaml:"scratchpad_encrypt"`
	AllowOutsideCwd     bool          `yaml:"allow_outside_cwd"`
	EntryCountLearnings int           `yaml:"entry_count_learnings"`
	EntryCountDecisions int           `yaml:"entry_count_decisions"`
	ConventionLineCount int           `yaml:"convention_line_count"`
	InjectionTokenWarn  int           `yaml:"injection_token_warn"`
	ContextWindow       int           `yaml:"context_window"`
	EventLog            bool          `yaml:"event_log"`
	KeyRotationDays     int           `yaml:"key_rotation_days"`
	KeyPathOverride     string        `yaml:"key_path"`
	Notify              *NotifyConfig `yaml:"notify"`
}

// NotifyConfig holds webhook notification settings.
//
// KeyRotationDays is deprecated here; use the top-level CtxRC.KeyRotationDays
// instead. This field is retained for backwards compatibility with existing
// .ctxrc files that have key_rotation_days nested under notify.
type NotifyConfig struct {
	Events          []string `yaml:"events"`
	KeyRotationDays int      `yaml:"key_rotation_days"`
}
