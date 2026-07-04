//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/err/config"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// SettingsStatusLine merges the ctx status line into
// settings.local.json.
//
// Deploys a statusLine entry running "ctx system statusline"
// (spec: specs/statusline.md). A pre-existing statusLine that is not
// ours is backed up to .context/state/ before being replaced, and is
// restored from that backup when statusline.enabled is set to false
// in .ctxrc. A foreign statusLine is never removed on disable: it is
// not ours to delete.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if file operations fail
func SettingsStatusLine(cmd *cobra.Command) error {
	if !rc.StatuslineEnabled() {
		return removeStatusLine(cmd)
	}
	raw, _, readErr := readSettingsRaw()
	if readErr != nil {
		return readErr
	}
	desired := claude.StatusLineConfig{
		Type:    cfgClaude.StatusLineType,
		Command: cfgClaude.StatusLineCommand,
	}
	if rawEntry, exists := raw[cfgClaude.FieldStatusLine]; exists {
		var existing claude.StatusLineConfig
		unmarshalErr := json.Unmarshal(rawEntry, &existing)
		if unmarshalErr == nil && existing.Command == desired.Command {
			initialize.NoChanges(cmd, cfgClaude.Settings)
			return nil
		}
		backupPath, backupErr := backupStatusLine(rawEntry)
		if backupErr != nil {
			return backupErr
		}
		initialize.StatuslineBackedUp(cmd, backupPath)
	}
	section, marshalErr := marshalSettingsSection(desired)
	if marshalErr != nil {
		return config.MarshalSettings(marshalErr)
	}
	raw[cfgClaude.FieldStatusLine] = section
	if writeErr := writeSettingsRaw(raw); writeErr != nil {
		return writeErr
	}
	initialize.StatuslineDeployed(cmd, cfgClaude.Settings)
	return nil
}
