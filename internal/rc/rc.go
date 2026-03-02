//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rc provides runtime configuration loading from .ctxrc files.
package rc

import (
	"os"
	"sync"

	"github.com/ActiveMemory/ctx/internal/config"
)

// Default returns a new CtxRC with hardcoded default values.
//
// Returns:
//   - *CtxRC: Configuration with defaults
//     (8000 token budget, 7-day archive, etc.)
func Default() *CtxRC {
	return &CtxRC{
		ContextDir:          config.DirContext,
		TokenBudget:         DefaultTokenBudget,
		PriorityOrder:       nil, // nil means use config.FileReadOrder
		AutoArchive:         true,
		ArchiveAfterDays:    DefaultArchiveAfterDays,
		EntryCountLearnings: DefaultEntryCountLearnings,
		EntryCountDecisions: DefaultEntryCountDecisions,
		ConventionLineCount: DefaultConventionLineCount,
		InjectionTokenWarn:  DefaultInjectionTokenWarn,
		ContextWindow:       DefaultContextWindow,
	}
}

// RC returns the loaded configuration, initializing it on the first call.
//
// It loads from .ctxrc if present, then applies environment overrides.
// The result is cached for subsequent calls.
//
// Returns:
//   - *CtxRC: The loaded and cached configuration
func RC() *CtxRC {
	rcOnce.Do(func() {
		rc = loadRC()
	})
	return rc
}

// ContextDir returns the configured context directory.
//
// Priority: CLI override > env var > .ctxrc > default.
//
// Returns:
//   - string: The context directory path (e.g., ".context")
func ContextDir() string {
	rcMu.RLock()
	defer rcMu.RUnlock()
	if rcOverrideDir != "" {
		return rcOverrideDir
	}
	return RC().ContextDir
}

// TokenBudget returns the configured default token budget.
//
// Priority: env var > .ctxrc > default (8000).
//
// Returns:
//   - int: The token budget for context assembly
func TokenBudget() int {
	return RC().TokenBudget
}

// PriorityOrder returns the configured file priority order.
//
// Returns:
//   - []string: File names in priority order, or nil if not configured
//     (callers should fall back to config.FileReadOrder)
func PriorityOrder() []string {
	return RC().PriorityOrder
}

// AutoArchive returns whether auto-archiving is enabled.
//
// Returns:
//   - bool: True if completed tasks should be auto-archived
func AutoArchive() bool {
	return RC().AutoArchive
}

// ArchiveAfterDays returns the configured days before archiving.
//
// Returns:
//   - int: Number of days after which completed tasks are archived (default 7)
func ArchiveAfterDays() int {
	return RC().ArchiveAfterDays
}

// ScratchpadEncrypt returns whether the scratchpad should be encrypted.
//
// Returns true (default) when the field is not set in .ctxrc.
//
// Returns:
//   - bool: True if scratchpad encryption is enabled (default true)
func ScratchpadEncrypt() bool {
	v := RC().ScratchpadEncrypt
	if v == nil {
		return true
	}
	return *v
}

// EntryCountLearnings returns the entry count threshold for LEARNINGS.md.
//
// Returns 0 if the check is disabled. Default: 30.
//
// Returns:
//   - int: Threshold above which a drift warning is emitted
func EntryCountLearnings() int {
	return RC().EntryCountLearnings
}

// EntryCountDecisions returns the entry count threshold for DECISIONS.md.
//
// Returns 0 if the check is disabled. Default: 20.
//
// Returns:
//   - int: Threshold above which a drift warning is emitted
func EntryCountDecisions() int {
	return RC().EntryCountDecisions
}

// ConventionLineCount returns the line count threshold for CONVENTIONS.md.
//
// Returns 0 if the check is disabled. Default: 200.
//
// Returns:
//   - int: Threshold above which a drift warning is emitted
func ConventionLineCount() int {
	return RC().ConventionLineCount
}

// InjectionTokenWarn returns the token threshold for oversize injection warning.
//
// Returns 0 if the check is disabled. Default: 15000.
//
// Returns:
//   - int: Threshold above which an oversize flag is written
func InjectionTokenWarn() int {
	return RC().InjectionTokenWarn
}

// ContextWindow returns the configured context window size in tokens.
//
// Returns 200000 (default) for Claude Opus/Sonnet. Override via
// `context_window` in .ctxrc for different models.
//
// Returns:
//   - int: Context window size in tokens
func ContextWindow() int {
	v := RC().ContextWindow
	if v <= 0 {
		return DefaultContextWindow
	}
	return v
}

// NotifyEvents returns the configured event filter list for notifications.
//
// Returns nil if Notify is nil (no filtering — all events pass).
//
// Returns:
//   - []string: Event names to allow, or nil for all
func NotifyEvents() []string {
	n := RC().Notify
	if n == nil {
		return nil
	}
	return n.Events
}

// KeyPath returns the resolved encryption key file path.
//
// Priority: key_path in .ctxrc (explicit) > user-level (~/.local/ctx/keys/) > legacy project-local.
//
// Returns:
//   - string: Resolved path to the encryption key file
func KeyPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return config.ResolveKeyPath(ContextDir(), cwd, RC().KeyPathOverride)
}

// KeyRotationDays returns the configured key rotation threshold in days.
//
// The encryption key is shared by both ctx pad and ctx notify, so the
// rotation threshold is a project-wide setting.
//
// Priority: top-level key_rotation_days > notify.key_rotation_days (legacy) > default (90).
//
// Returns:
//   - int: Number of days before a key rotation nudge
func KeyRotationDays() int {
	cfg := RC()
	if cfg.KeyRotationDays > 0 {
		return cfg.KeyRotationDays
	}
	if cfg.Notify != nil && cfg.Notify.KeyRotationDays > 0 {
		return cfg.Notify.KeyRotationDays
	}
	return DefaultKeyRotationDays
}

// EventLog returns whether local hook event logging is enabled.
//
// Returns false (default) when the field is not set in .ctxrc.
//
// Returns:
//   - bool: True if hook events should be logged to .context/state/events.jsonl
func EventLog() bool {
	return RC().EventLog
}

// AllowOutsideCwd returns whether boundary validation should be skipped.
//
// Returns false (default) when the field is not set in .ctxrc.
//
// Returns:
//   - bool: True if context directory is allowed outside the project root
func AllowOutsideCwd() bool {
	return RC().AllowOutsideCwd
}

// OverrideContextDir sets a CLI-provided override for the context directory.
//
// This takes precedence over all other configuration sources.
//
// Parameters:
//   - dir: Directory path to use as an override
func OverrideContextDir(dir string) {
	rcMu.Lock()
	defer rcMu.Unlock()
	rcOverrideDir = dir
}

// Reset clears the cached configuration, forcing reload on the next access.
// This is primarily useful for testing.
func Reset() {
	rcMu.Lock()
	defer rcMu.Unlock()
	rcOnce = sync.Once{}
	rc = nil
	rcOverrideDir = ""
}

// FilePriority returns the priority of a context file.
//
// If a priority_order is configured in .ctxrc, that order is used.
// Otherwise, the default config.FileReadOrder is used.
//
// Lower numbers indicate higher priority (1 = highest).
// Unknown files return 100.
//
// Parameters:
//   - name: Filename to look up (e.g., "TASKS.md")
//
// Returns:
//   - int: Priority value (1-9 for known files, 100 for unknown)
func FilePriority(name string) int {
	// Check for .ctxrc override first
	if order := PriorityOrder(); order != nil {
		for i, fName := range order {
			if fName == name {
				return i + 1
			}
		}
		// File not in custom order gets the lowest priority
		return 100
	}

	// Use the default priority from config.FileReadOrder
	for i, fName := range config.FileReadOrder {
		if fName == name {
			return i + 1
		}
	}
	return 100
}
