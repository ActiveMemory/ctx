//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/codex"
	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	"github.com/ActiveMemory/ctx/internal/io"
)

// CodexSessionDirs returns the directory where Codex rollouts are
// stored: `$CODEX_HOME/sessions`, or `~/.codex/sessions` when the
// variable is unset. Rollouts nest as `YYYY/MM/DD/rollout-*.jsonl`
// beneath it; [ScanDirectory] walks recursively, so the root alone
// is returned.
//
// Returns:
//   - []string: the sessions directory, or nil when it does not exist
func CodexSessionDirs() []string {
	codexHome := codex.Home()
	if codexHome == "" {
		return nil
	}

	dir := filepath.Join(codexHome, cfgCodex.DirSessions)
	info, statErr := io.SafeStat(dir)
	if statErr != nil || !info.IsDir() {
		return nil
	}

	return []string{dir}
}
