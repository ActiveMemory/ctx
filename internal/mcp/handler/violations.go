//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	ctxio "github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// readAndClearViolations reads violations from
// .context/state/violations.json and removes the file to prevent
// repeated escalation.
//
// Parameters:
//   - contextDir: path to the project context directory
//
// Returns:
//   - []violation: parsed violations, or nil if contextDir is empty,
//     no file exists, or on read/parse error
func readAndClearViolations(contextDir string) []violation {
	if contextDir == "" {
		return nil
	}
	stateDir := filepath.Join(contextDir, dir.State)
	data, readErr := ctxio.SafeReadFile(stateDir, file.Violations)
	if readErr != nil {
		return nil
	}
	// Remove the file immediately to prevent duplicate alerts. A
	// failed remove means the next read re-reports these violations,
	// so surface it rather than swallowing.
	violationsPath := filepath.Join(stateDir, file.Violations)
	if rmErr := os.Remove(violationsPath); rmErr != nil {
		logWarn.Warn(cfgWarn.Remove, violationsPath, rmErr)
	}

	var vd violationsData
	if unmarshalErr := json.Unmarshal(data, &vd); unmarshalErr != nil {
		return nil
	}
	return vd.Entries
}
