//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_version

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/claude"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	"github.com/ActiveMemory/ctx/internal/config/version"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Run executes the check-version hook logic.
//
// Compares the binary version against the embedded plugin version and
// emits a version mismatch warning if they differ. Also, piggybacks
// a key rotation age check. Throttled to once per day.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !core.Initialized() {
		return nil
	}

	input, _, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	tmpDir := core.StateDir()
	markerFile := filepath.Join(tmpDir, version.ThrottleID)

	if core.IsDailyThrottled(markerFile) {
		return nil
	}

	binaryVer := cmd.Root().Version

	// Skip check for dev builds
	if binaryVer == version.DevBuild {
		core.TouchFile(markerFile)
		return nil
	}

	pluginVer, pluginErr := claude.PluginVersion()
	if pluginErr != nil {
		return nil // embedded plugin.json missing — nothing to compare
	}

	bMajor, bMinor, bOK := core.ParseMajorMinor(binaryVer)
	pMajor, pMinor, pOK := core.ParseMajorMinor(pluginVer)

	if !bOK || !pOK {
		core.TouchFile(markerFile)
		return nil
	}

	if bMajor == pMajor && bMinor == pMinor {
		core.TouchFile(markerFile)
		return nil
	}

	// Version mismatch — emit warning
	fallback := fmt.Sprintf(desc.TextDesc(
		text.DescKeyCheckVersionFallback), binaryVer, pluginVer,
	)
	content := core.LoadMessage(hook.CheckVersion, hook.VariantMismatch,
		map[string]any{
			tpl.VarBinaryVersion: binaryVer,
			tpl.VarPluginVersion: pluginVer,
		}, fallback)
	if content == "" {
		core.TouchFile(markerFile)
		return nil
	}

	boxTitle := desc.TextDesc(text.DescKeyCheckVersionBoxTitle)
	relayPrefix := desc.TextDesc(text.DescKeyCheckVersionRelayPrefix)

	cmd.Println(core.NudgeBox(relayPrefix, boxTitle, content))

	ref := notify.NewTemplateRef(hook.CheckVersion, hook.VariantMismatch,
		map[string]any{
			tpl.VarBinaryVersion: binaryVer,
			tpl.VarPluginVersion: pluginVer,
		})
	versionMsg := hook.CheckVersion + ": " +
		fmt.Sprintf(
			desc.TextDesc(
				text.DescKeyCheckVersionMismatchRelayFormat,
			), binaryVer, pluginVer,
		)
	core.NudgeAndRelay(versionMsg, input.SessionID, ref)

	core.TouchFile(markerFile)

	// Key age check — piggyback on the daily version check
	core.CheckKeyAge(cmd, input.SessionID)

	return nil
}
