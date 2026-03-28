//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/config/file"
	writeRC "github.com/ActiveMemory/ctx/internal/write/rc"
)

// loadRC loads configuration from the .ctxrc file and applies env
// overrides.
//
// Returns:
//   - *CtxRC: Configuration with file values and env overrides applied
func loadRC() *CtxRC {
	cfg := Default()

	// Try to load .ctxrc from the current directory
	data, readErr := os.ReadFile(file.CtxRC)
	if readErr == nil {
		if yamlErr := yaml.Unmarshal(data, cfg); yamlErr != nil {
			writeRC.ParseWarning(file.CtxRC, yamlErr)
		}
	}

	// Apply environment variable overrides
	if envDir := os.Getenv(env.CtxDir); envDir != "" {
		cfg.ContextDir = envDir
	}
	if envBudget := os.Getenv(env.CtxTokenBudget); envBudget != "" {
		if budget, parseErr := strconv.Atoi(envBudget); parseErr == nil && budget > 0 {
			cfg.TokenBudget = budget
		}
	}

	return cfg
}
