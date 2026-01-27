//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"strconv"
	"sync"

	"gopkg.in/yaml.v3"
)

// RC represents the configuration from .contextrc file.
type RC struct {
	ContextDir       string   `yaml:"context_dir"`
	TokenBudget      int      `yaml:"token_budget"`
	PriorityOrder    []string `yaml:"priority_order"`
	AutoArchive      bool     `yaml:"auto_archive"`
	ArchiveAfterDays int      `yaml:"archive_after_days"`
}

// DefaultTokenBudget is the default token budget when not configured.
const DefaultTokenBudget = 8000

// DefaultArchiveAfterDays is the default days before archiving.
const DefaultArchiveAfterDays = 7

var (
	rc           *RC
	rcOnce       sync.Once
	rcOverrideDir string
)

// DefaultRC returns a new RC with hardcoded default values.
func DefaultRC() *RC {
	return &RC{
		ContextDir:       DirContext,
		TokenBudget:      DefaultTokenBudget,
		PriorityOrder:    nil, // nil means use FileReadOrder
		AutoArchive:      true,
		ArchiveAfterDays: DefaultArchiveAfterDays,
	}
}

// GetRC returns the loaded configuration, initializing it on first call.
// It loads from .contextrc if present, then applies environment overrides.
func GetRC() *RC {
	rcOnce.Do(func() {
		rc = loadRC()
	})
	return rc
}

// loadRC loads configuration from .contextrc file and applies env overrides.
func loadRC() *RC {
	cfg := DefaultRC()

	// Try to load .contextrc from current directory
	data, err := os.ReadFile(".contextrc")
	if err == nil {
		// Parse YAML, ignoring errors (use defaults for invalid config)
		_ = yaml.Unmarshal(data, cfg)
	}

	// Apply environment variable overrides
	if envDir := os.Getenv("CTX_DIR"); envDir != "" {
		cfg.ContextDir = envDir
	}
	if envBudget := os.Getenv("CTX_TOKEN_BUDGET"); envBudget != "" {
		if budget, err := strconv.Atoi(envBudget); err == nil && budget > 0 {
			cfg.TokenBudget = budget
		}
	}

	return cfg
}

// GetContextDir returns the configured context directory.
// Priority: CLI override > env var > .contextrc > default.
func GetContextDir() string {
	if rcOverrideDir != "" {
		return rcOverrideDir
	}
	return GetRC().ContextDir
}

// GetTokenBudget returns the configured default token budget.
// Priority: env var > .contextrc > default.
func GetTokenBudget() int {
	return GetRC().TokenBudget
}

// GetPriorityOrder returns the configured file priority order.
// Returns nil if not configured (callers should fall back to FileReadOrder).
func GetPriorityOrder() []string {
	return GetRC().PriorityOrder
}

// GetAutoArchive returns whether auto-archiving is enabled.
func GetAutoArchive() bool {
	return GetRC().AutoArchive
}

// GetArchiveAfterDays returns the configured days before archiving.
func GetArchiveAfterDays() int {
	return GetRC().ArchiveAfterDays
}

// OverrideContextDir sets a CLI-provided override for the context directory.
// This takes precedence over all other configuration sources.
func OverrideContextDir(dir string) {
	rcOverrideDir = dir
}

// ResetRC clears the cached configuration, forcing reload on next access.
// This is primarily useful for testing.
func ResetRC() {
	rcOnce = sync.Once{}
	rc = nil
	rcOverrideDir = ""
}
