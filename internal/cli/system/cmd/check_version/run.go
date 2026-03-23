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

	hook2 "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	version2 "github.com/ActiveMemory/ctx/internal/cli/system/core/version"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/claude"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/version"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
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
	if !state.Initialized() {
		return nil
	}

	input, _, paused := hook2.Preamble(stdin)
	if paused {
		return nil
	}

	tmpDir := state.StateDir()
	markerFile := filepath.Join(tmpDir, version.ThrottleID)

	if hook2.DailyThrottled(markerFile) {
		return nil
	}

	binaryVer := cmd.Root().Version

	// Skip check for dev builds
	if binaryVer == version.DevBuild {
		internalIo.TouchFile(markerFile)
		return nil
	}

	pluginVer, pluginErr := claude.PluginVersion()
	if pluginErr != nil {
		return nil // embedded plugin.json missing — nothing to compare
	}

	bMajor, bMinor, bOK := version2.ParseMajorMinor(binaryVer)
	pMajor, pMinor, pOK := version2.ParseMajorMinor(pluginVer)

	if !bOK || !pOK {
		internalIo.TouchFile(markerFile)
		return nil
	}

	if bMajor == pMajor && bMinor == pMinor {
		internalIo.TouchFile(markerFile)
		return nil
	}

	// Version mismatch — emit warning
	fallback := fmt.Sprintf(desc.Text(
		text.DescKeyCheckVersionFallback), binaryVer, pluginVer,
	)
	content := message.LoadMessage(hook.CheckVersion, hook.VariantMismatch,
		map[string]any{
			version.VarBinaryVersion: binaryVer,
			version.VarPluginVersion: pluginVer,
		}, fallback)
	if content == "" {
		internalIo.TouchFile(markerFile)
		return nil
	}

	boxTitle := desc.Text(text.DescKeyCheckVersionBoxTitle)
	relayPrefix := desc.Text(text.DescKeyCheckVersionRelayPrefix)

	writeHook.Nudge(cmd, message.NudgeBox(relayPrefix, boxTitle, content))

	ref := notify.NewTemplateRef(hook.CheckVersion, hook.VariantMismatch,
		map[string]any{
			version.VarBinaryVersion: binaryVer,
			version.VarPluginVersion: pluginVer,
		})
	versionMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckVersion, fmt.Sprintf(
			desc.Text(text.DescKeyCheckVersionMismatchRelayFormat),
			binaryVer, pluginVer))
	nudge.NudgeAndRelay(versionMsg, input.SessionID, ref)

	internalIo.TouchFile(markerFile)

	// Key age check: piggyback on the daily version check
	writeHook.Nudge(cmd, version2.CheckKeyAge(input.SessionID))

	return nil
}
