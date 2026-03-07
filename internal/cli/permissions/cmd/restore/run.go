//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package restore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
	"github.com/ActiveMemory/ctx/internal/cli/permissions/core"
	"github.com/ActiveMemory/ctx/internal/config"
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
			return core.ErrGoldenNotFound()
		}
		return core.ErrReadFile(config.FileSettingsGolden, goldenReadErr)
	}

	localBytes, localReadErr := os.ReadFile(config.FileSettings)
	if localReadErr != nil {
		if os.IsNotExist(localReadErr) {
			// No local file — just copy golden.
			if writeErr := os.WriteFile(config.FileSettings, goldenBytes, config.PermFile); writeErr != nil {
				return core.ErrWriteFile(config.FileSettings, writeErr)
			}
			cmd.Println("Restored golden image (no local settings existed).")
			return nil
		}
		return core.ErrReadFile(config.FileSettings, localReadErr)
	}

	// Fast path: files are identical.
	if bytes.Equal(goldenBytes, localBytes) {
		cmd.Println("Settings already match golden image.")
		return nil
	}

	// Parse both to compute permission diff.
	var golden, local claude.Settings
	if goldenParseErr := json.Unmarshal(goldenBytes, &golden); goldenParseErr != nil {
		return core.ErrParseSettings(config.FileSettingsGolden, goldenParseErr)
	}
	if localParseErr := json.Unmarshal(localBytes, &local); localParseErr != nil {
		return core.ErrParseSettings(config.FileSettings, localParseErr)
	}

	restored, dropped := core.DiffStringSlices(golden.Permissions.Allow, local.Permissions.Allow)
	denyRestored, denyDropped := core.DiffStringSlices(golden.Permissions.Deny, local.Permissions.Deny)

	if len(dropped) > 0 {
		cmd.Println(fmt.Sprintf("Dropped %d session allow permission(s):", len(dropped)))
		for _, p := range dropped {
			cmd.Println(fmt.Sprintf("  - %s", p))
		}
	}
	if len(restored) > 0 {
		cmd.Println(fmt.Sprintf("Restored %d allow permission(s):", len(restored)))
		for _, p := range restored {
			cmd.Println(fmt.Sprintf("  + %s", p))
		}
	}
	if len(denyDropped) > 0 {
		cmd.Println(fmt.Sprintf("Dropped %d session deny rule(s):", len(denyDropped)))
		for _, p := range denyDropped {
			cmd.Println(fmt.Sprintf("  - %s", p))
		}
	}
	if len(denyRestored) > 0 {
		cmd.Println(fmt.Sprintf("Restored %d deny rule(s):", len(denyRestored)))
		for _, p := range denyRestored {
			cmd.Println(fmt.Sprintf("  + %s", p))
		}
	}
	allEmpty := len(dropped) == 0 && len(restored) == 0 && len(denyDropped) == 0 && len(denyRestored) == 0
	if allEmpty {
		cmd.Println("Permission lists match; other settings differ.")
	}

	// Write golden bytes (byte-for-byte copy).
	if writeErr := os.WriteFile(config.FileSettings, goldenBytes, config.PermFile); writeErr != nil {
		return core.ErrWriteFile(config.FileSettings, writeErr)
	}

	cmd.Println("Restored from golden image.")
	return nil
}
