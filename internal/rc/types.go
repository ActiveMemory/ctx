//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

// CtxRC represents the configuration from the .contextrc file.
//
// Fields:
//   - ContextDir: Name of the context directory (default ".context")
//   - TokenBudget: Default token budget for context assembly (default 8000)
//   - PriorityOrder: Custom file loading priority order
//   - AutoArchive: Whether to auto-archive completed tasks (default true)
//   - ArchiveAfterDays: Days before archiving completed tasks (default 7)
type CtxRC struct {
	ContextDir       string   `yaml:"context_dir"`
	TokenBudget      int      `yaml:"token_budget"`
	PriorityOrder    []string `yaml:"priority_order"`
	AutoArchive      bool     `yaml:"auto_archive"`
	ArchiveAfterDays int      `yaml:"archive_after_days"`
}
