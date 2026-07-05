//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"bytes"
	"encoding/json"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// removeStatusLine handles statusline.enabled: false. When the
// deployed statusLine is ours, it is replaced by the backed-up
// previous entry when one exists, or dropped. Foreign entries are
// left untouched.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if file operations fail
func removeStatusLine(cmd *cobra.Command) error {
	raw, fileExists, readErr := readSettingsRaw()
	if readErr != nil {
		return readErr
	}
	rawEntry, exists := raw[cfgClaude.FieldStatusLine]
	if !fileExists || !exists {
		initialize.NoChanges(cmd, cfgClaude.Settings)
		return nil
	}
	var existing claude.StatusLineConfig
	unmarshalErr := json.Unmarshal(rawEntry, &existing)
	if unmarshalErr != nil ||
		existing.Command != cfgClaude.StatusLineCommand {
		initialize.NoChanges(cmd, cfgClaude.Settings)
		return nil
	}
	backupPath, previous := readStatusLineBackup()
	if len(previous) > 0 {
		raw[cfgClaude.FieldStatusLine] = previous
		if writeErr := writeSettingsRaw(raw); writeErr != nil {
			return writeErr
		}
		initialize.StatuslineRestored(cmd, backupPath)
		return nil
	}
	delete(raw, cfgClaude.FieldStatusLine)
	if writeErr := writeSettingsRaw(raw); writeErr != nil {
		return writeErr
	}
	initialize.StatuslineRemoved(cmd, cfgClaude.Settings)
	return nil
}

// statusLineBackupPath resolves the backup file location under the
// project's .context/state/ directory.
//
// Returns:
//   - string: absolute backup file path
//   - error: Non-nil when the context directory cannot be resolved
func statusLineBackupPath() (string, error) {
	ctxDir, dirErr := rc.ContextDir()
	if dirErr != nil {
		return "", dirErr
	}
	return filepath.Join(ctxDir, dir.State, cfgClaude.PreviousStatusLine), nil
}

// backupStatusLine persists a displaced statusLine entry so a later
// statusline.enabled: false can restore it.
//
// Parameters:
//   - rawEntry: statusLine entry being displaced, verbatim
//
// Returns:
//   - string: backup file path the entry was written to
//   - error: Non-nil if path resolution or file operations fail
func backupStatusLine(rawEntry json.RawMessage) (string, error) {
	path, pathErr := statusLineBackupPath()
	if pathErr != nil {
		return "", pathErr
	}
	stateDir := filepath.Dir(path)
	if mkdirErr := ctxIo.SafeMkdirAll(stateDir, fs.PermExec); mkdirErr != nil {
		return "", errFs.Mkdir(stateDir, mkdirErr)
	}
	content := append(bytes.TrimSpace(rawEntry), token.NewlineLF...)
	if writeErr := ctxIo.SafeWriteFile(
		path, content, fs.PermFile,
	); writeErr != nil {
		return "", errFs.FileWrite(path, writeErr)
	}
	return path, nil
}

// readStatusLineBackup loads a previously backed-up statusLine entry.
// Restore is best-effort, never fatal.
//
// Returns:
//   - string: backup file path (may be "" when unresolvable)
//   - json.RawMessage: backed-up entry, or nil when no valid backup
//     exists
func readStatusLineBackup() (string, json.RawMessage) {
	path, pathErr := statusLineBackupPath()
	if pathErr != nil {
		return "", nil
	}
	content, readErr := ctxIo.SafeReadUserFile(path)
	if readErr != nil {
		return path, nil
	}
	trimmed := bytes.TrimSpace(content)
	if !json.Valid(trimmed) {
		return path, nil
	}
	return path, json.RawMessage(trimmed)
}
