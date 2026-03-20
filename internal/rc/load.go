//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"gopkg.in/yaml.v3"
)

// loadRC loads configuration from the .ctxrc file and applies env
// overrides.
//
// Returns:
//   - *CtxRC: Configuration with file values and env overrides applied
func loadRC() *CtxRC {
	cfg := Default()

	// Try to load .ctxrc from the current directory
	data, err := os.ReadFile(file.CtxRC)
	if err == nil {
		if yamlErr := yaml.Unmarshal(data, cfg); yamlErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, desc.TextDesc(text.DescKeyRcParseWarning)+token.NewlineLF,
				file.CtxRC, yamlErr)
		}
	}

	// Apply environment variable overrides
	if envDir := os.Getenv(env.CtxDir); envDir != "" {
		cfg.ContextDir = envDir
	}
	if envBudget := os.Getenv(env.CtxTokenBudget); envBudget != "" {
		if budget, err := strconv.Atoi(envBudget); err == nil && budget > 0 {
			cfg.TokenBudget = budget
		}
	}

	return cfg
}
