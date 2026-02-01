//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rc provides runtime configuration loading from .contextrc files.
package rc

import (
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
		ContextDir:       config.DirContext,
		TokenBudget:      DefaultTokenBudget,
		PriorityOrder:    nil, // nil means use config.FileReadOrder
		AutoArchive:      true,
		ArchiveAfterDays: DefaultArchiveAfterDays,
	}
}

// RC returns the loaded configuration, initializing it on the first call.
//
// It loads from .contextrc if present, then applies environment overrides.
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
// Priority: CLI override > env var > .contextrc > default.
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
// Priority: env var > .contextrc > default (8000).
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
// If a priority_order is configured in .contextrc, that order is used.
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
	// Check for .contextrc override first
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
