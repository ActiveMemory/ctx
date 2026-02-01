//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config"
)

// loadRC loads configuration from the .contextrc file and applies env
// overrides.
//
// Returns:
//   - *CtxRC: Configuration with file values and env overrides applied
func loadRC() *CtxRC {
	cfg := Default()

	// Try to load .contextrc from the current directory
	data, err := os.ReadFile(config.FileContextRC)
	if err == nil {
		// Parse YAML, ignoring errors (use defaults for invalid config)
		_ = yaml.Unmarshal(data, cfg)
	}

	// Apply environment variable overrides
	if envDir := os.Getenv(config.EnvCtxDir); envDir != "" {
		cfg.ContextDir = envDir
	}
	if envBudget := os.Getenv(config.EnvCtxTokenBudget); envBudget != "" {
		if budget, err := strconv.Atoi(envBudget); err == nil && budget > 0 {
			cfg.TokenBudget = budget
		}
	}

	return cfg
}
