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
	"github.com/ActiveMemory/ctx/internal/cli/permissions/core"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run resets settings.local.json from the golden image.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on read/write/parse failure or missing golden file
func Run(cmd *cobra.Command) error {
	goldenBytes, goldenReadErr := os.ReadFile(config.FileSettingsGolden)
	if goldenReadErr != nil {
		if os.IsNotExist(goldenReadErr) {
			return ctxerr.GoldenNotFound()
		}
		return ctxerr.FileRead(config.FileSettingsGolden, goldenReadErr)
	}

	localBytes, localReadErr := os.ReadFile(config.FileSettings)
	if localReadErr != nil {
		if os.IsNotExist(localReadErr) {
			if writeErr := os.WriteFile(
				config.FileSettings, goldenBytes, config.PermFile,
			); writeErr != nil {
				return ctxerr.FileWrite(config.FileSettings, writeErr)
			}
			write.RestoreNoLocal(cmd)
			return nil
		}
		return ctxerr.FileRead(config.FileSettings, localReadErr)
	}

	if bytes.Equal(goldenBytes, localBytes) {
		write.RestoreMatch(cmd)
		return nil
	}

	var golden, local claude.Settings
	if goldenParseErr := json.Unmarshal(goldenBytes, &golden); goldenParseErr != nil {
		return ctxerr.ParseFile(config.FileSettingsGolden, goldenParseErr)
	}
	if localParseErr := json.Unmarshal(localBytes, &local); localParseErr != nil {
		return ctxerr.ParseFile(config.FileSettings, localParseErr)
	}

	restored, dropped := core.DiffStringSlices(
		golden.Permissions.Allow, local.Permissions.Allow,
	)
	denyRestored, denyDropped := core.DiffStringSlices(
		golden.Permissions.Deny, local.Permissions.Deny,
	)

	write.RestoreDiff(cmd, dropped, restored, denyDropped, denyRestored)

	if writeErr := os.WriteFile(
		config.FileSettings, goldenBytes, config.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(config.FileSettings, writeErr)
	}

	write.RestoreDone(cmd)
	return nil
}
