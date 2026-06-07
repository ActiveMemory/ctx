//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"path/filepath"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
)

// statePath returns the absolute path to the state file within the
// dreams/ notebook under dreamsDir.
//
// Parameters:
//   - dreamsDir: absolute path to the dreams/ notebook directory
//
// Returns:
//   - string: <dreamsDir>/state.json
func statePath(dreamsDir string) string {
	return filepath.Join(dreamsDir, cfgDream.FileState)
}

// ledgerPath returns the absolute path to the ledger file within the
// dreams/ notebook under dreamsDir.
//
// Parameters:
//   - dreamsDir: absolute path to the dreams/ notebook directory
//
// Returns:
//   - string: <dreamsDir>/ledger.md
func ledgerPath(dreamsDir string) string {
	return filepath.Join(dreamsDir, cfgDream.FileLedger)
}
