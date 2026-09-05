//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"os"
	"path/filepath"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
)

// Home returns the Codex home directory: `$CODEX_HOME` when set,
// otherwise `~/.codex`.
//
// Returns:
//   - string: absolute Codex home path, or empty when neither
//     `$CODEX_HOME` nor the user home directory can be resolved
func Home() string {
	if fromEnv := os.Getenv(cfgCodex.EnvHome); fromEnv != "" {
		return fromEnv
	}
	userHome, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return ""
	}
	return filepath.Join(userHome, cfgCodex.DirHome)
}
