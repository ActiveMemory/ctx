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
	Notify              *NotifyConfig `yaml:"notify"`
}

// NotifyConfig holds webhook notification settings.
type NotifyConfig struct {
	Events          []string `yaml:"events"`
	KeyRotationDays int      `yaml:"key_rotation_days"`
}
