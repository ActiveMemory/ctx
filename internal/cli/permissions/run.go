//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package permissions

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/claude"
	"github.com/ActiveMemory/ctx/internal/config"
)

// runSnapshot saves settings.local.json as the golden image.
func runSnapshot(cmd *cobra.Command) error {
	content, err := os.ReadFile(config.FileSettings)
	if err != nil {
		if os.IsNotExist(err) {
			return errSettingsNotFound()
		}
		return errReadFile(config.FileSettings, err)
	}

	// Determine message based on whether golden already exists.
	verb := "Saved"
	if _, err := os.Stat(config.FileSettingsGolden); err == nil {
		verb = "Updated"
	}

	if err := os.WriteFile(config.FileSettingsGolden, content, config.PermFile); err != nil {
		return errWriteFile(config.FileSettingsGolden, err)
	}

	cmd.Printf("%s golden image: %s\n", verb, config.FileSettingsGolden)
	return nil
}

// runRestore resets settings.local.json from the golden image.
func runRestore(cmd *cobra.Command) error {
	goldenBytes, err := os.ReadFile(config.FileSettingsGolden)
	if err != nil {
		if os.IsNotExist(err) {
			return errGoldenNotFound()
		}
		return errReadFile(config.FileSettingsGolden, err)
	}

	localBytes, err := os.ReadFile(config.FileSettings)
	if err != nil {
		if os.IsNotExist(err) {
			// No local file — just copy golden.
			if writeErr := os.WriteFile(config.FileSettings, goldenBytes, config.PermFile); writeErr != nil {
				return errWriteFile(config.FileSettings, writeErr)
			}
			cmd.Println("Restored golden image (no local settings existed).")
			return nil
		}
		return errReadFile(config.FileSettings, err)
	}

	// Fast path: files are identical.
	if bytes.Equal(goldenBytes, localBytes) {
		cmd.Println("Settings already match golden image.")
		return nil
	}

	// Parse both to compute permission diff.
	var golden, local claude.Settings
	if err := json.Unmarshal(goldenBytes, &golden); err != nil {
		return errParseSettings(config.FileSettingsGolden, err)
	}
	if err := json.Unmarshal(localBytes, &local); err != nil {
		return errParseSettings(config.FileSettings, err)
	}

	restored, dropped := diffStringSlices(golden.Permissions.Allow, local.Permissions.Allow)
	denyRestored, denyDropped := diffStringSlices(golden.Permissions.Deny, local.Permissions.Deny)

	if len(dropped) > 0 {
		cmd.Printf("Dropped %d session allow permission(s):\n", len(dropped))
		for _, p := range dropped {
			cmd.Printf("  - %s\n", p)
		}
	}
	if len(restored) > 0 {
		cmd.Printf("Restored %d allow permission(s):\n", len(restored))
		for _, p := range restored {
			cmd.Printf("  + %s\n", p)
		}
	}
	if len(denyDropped) > 0 {
		cmd.Printf("Dropped %d session deny rule(s):\n", len(denyDropped))
		for _, p := range denyDropped {
			cmd.Printf("  - %s\n", p)
		}
	}
	if len(denyRestored) > 0 {
		cmd.Printf("Restored %d deny rule(s):\n", len(denyRestored))
		for _, p := range denyRestored {
			cmd.Printf("  + %s\n", p)
		}
	}
	allEmpty := len(dropped) == 0 && len(restored) == 0 && len(denyDropped) == 0 && len(denyRestored) == 0
	if allEmpty {
		cmd.Println("Permission lists match; other settings differ.")
	}

	// Write golden bytes (byte-for-byte copy).
	if err := os.WriteFile(config.FileSettings, goldenBytes, config.PermFile); err != nil {
		return errWriteFile(config.FileSettings, err)
	}

	cmd.Println("Restored from golden image.")
	return nil
}

// diffStringSlices computes the set difference between golden and local slices.
//
// Returns:
//   - restored: entries in golden but not in local
//   - dropped: entries in local but not in golden
//
// Both output slices preserve the source ordering of their respective inputs.
func diffStringSlices(golden, local []string) (restored, dropped []string) {
	goldenSet := make(map[string]struct{}, len(golden))
	for _, s := range golden {
		goldenSet[s] = struct{}{}
	}

	localSet := make(map[string]struct{}, len(local))
	for _, s := range local {
		localSet[s] = struct{}{}
	}

	for _, s := range golden {
		if _, ok := localSet[s]; !ok {
			restored = append(restored, s)
		}
	}

	for _, s := range local {
		if _, ok := goldenSet[s]; !ok {
			dropped = append(dropped, s)
		}
	}

	return restored, dropped
}
