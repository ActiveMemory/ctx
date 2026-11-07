//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/agy"
	cfgAgy "github.com/ActiveMemory/ctx/internal/config/agy"
	"github.com/ActiveMemory/ctx/internal/io"
)

// AgySessionDirs returns the directory where Antigravity sessions are
// stored: `$AGY_HOME/antigravity-cli/brain`, or
// `~/.gemini/antigravity-cli/brain` when the variable is unset.
// Transcripts nest as `<id>/.system_generated/logs/transcript.jsonl`
// beneath it; [ScanDirectory] walks recursively, so the root alone
// is returned.
//
// Returns:
//   - []string: the sessions directory, or nil when it does not exist
func AgySessionDirs() []string {
	agyHome := agy.Home()
	if agyHome == "" {
		return nil
	}

	dir := filepath.Join(
		agyHome, cfgAgy.DirAntigravityCLI, cfgAgy.DirBrain,
	)
	info, statErr := io.SafeStat(dir)
	if statErr != nil || !info.IsDir() {
		return nil
	}

	return []string{dir}
}
