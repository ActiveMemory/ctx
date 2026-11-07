//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

import (
	"os"
	"path/filepath"

	cfgAgy "github.com/ActiveMemory/ctx/internal/config/agy"
)

// Home returns the Antigravity home directory: `$AGY_HOME` when set,
// otherwise `~/.gemini`.
//
// Returns:
//   - string: absolute Antigravity home path, or empty when neither
//     `$AGY_HOME` nor the user home directory can be resolved
func Home() string {
	if fromEnv := os.Getenv(cfgAgy.EnvHome); fromEnv != "" {
		return fromEnv
	}
	userHome, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return ""
	}
	return filepath.Join(userHome, cfgAgy.DirGemini)
}

// ConfigDir returns the Antigravity configuration directory:
// `$AGY_CONFIG_DIR` when set, otherwise `<Home()>/config`.
//
// Returns:
//   - string: absolute Antigravity config directory path, or empty
//     when Home() cannot be resolved
func ConfigDir() string {
	if fromEnv := os.Getenv(cfgAgy.EnvConfigDir); fromEnv != "" {
		return fromEnv
	}
	home := Home()
	if home == "" {
		return ""
	}
	return filepath.Join(home, cfgAgy.DirConfig)
}
