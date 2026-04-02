//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package restore

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
	"github.com/ActiveMemory/ctx/internal/cli/permission/core"
	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/err/config"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/restore"
)

// Run resets settings.local.json from the golden image.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read/write/parse failure or missing golden file
func Run(cmd *cobra.Command) error {
	goldenBytes, goldenReadErr := io.SafeReadUserFile(cfgClaude.SettingsGolden)
	if goldenReadErr != nil {
		if os.IsNotExist(goldenReadErr) {
			return config.GoldenNotFound()
		}
		return errFs.FileRead(cfgClaude.SettingsGolden, goldenReadErr)
	}

	localBytes, localReadErr := io.SafeReadUserFile(cfgClaude.Settings)
	if localReadErr != nil {
		if os.IsNotExist(localReadErr) {
			if writeErr := io.SafeWriteFile(
				cfgClaude.Settings, goldenBytes, fs.PermFile,
			); writeErr != nil {
				return errFs.FileWrite(cfgClaude.Settings, writeErr)
			}
			restore.NoLocal(cmd)
			return nil
		}
		return errFs.FileRead(cfgClaude.Settings, localReadErr)
	}

	if bytes.Equal(goldenBytes, localBytes) {
		restore.Match(cmd)
		return nil
	}

	var golden, local claude.Settings
	goldenParseErr := json.Unmarshal(goldenBytes, &golden)
	if goldenParseErr != nil {
		return errParser.ParseFile(cfgClaude.SettingsGolden, goldenParseErr)
	}
	if localParseErr := json.Unmarshal(localBytes, &local); localParseErr != nil {
		return errParser.ParseFile(cfgClaude.Settings, localParseErr)
	}

	restored, dropped := core.DiffStringSlices(
		golden.Permissions.Allow, local.Permissions.Allow,
	)
	denyRestored, denyDropped := core.DiffStringSlices(
		golden.Permissions.Deny, local.Permissions.Deny,
	)

	restore.Diff(cmd, dropped, restored, denyDropped, denyRestored)

	if writeErr := io.SafeWriteFile(
		cfgClaude.Settings, goldenBytes, fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(cfgClaude.Settings, writeErr)
	}

	restore.Done(cmd)
	return nil
}
