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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/claude"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreVersion "github.com/ActiveMemory/ctx/internal/cli/system/core/version"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/version"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
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
	input, _, _, tmpDir, ok := coreCheck.FullPreamble(stdin)
	bailSilently := !ok
	if bailSilently {
		return nil
	}
	markerFile := filepath.Join(tmpDir, version.ThrottleID)
	if coreCheck.DailyThrottled(markerFile) {
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
		internalIo.TouchFile(markerFile)
		msg := fmt.Sprintf(
			desc.Text(text.DescKeyCheckVersionPluginReadError), pluginErr,
		)
		writeSetup.Nudge(cmd, msg)
		return nil
	}

	bMajor, bMinor, bOK := coreVersion.ParseMajorMinor(binaryVer)
	pMajor, pMinor, pOK := coreVersion.ParseMajorMinor(pluginVer)

	if !bOK || !pOK {
		internalIo.TouchFile(markerFile)
		return nil
	}

	if bMajor == pMajor && bMinor == pMinor {
		internalIo.TouchFile(markerFile)
		return nil
	}

	// Version mismatch - emit warning
	fallback := fmt.Sprintf(desc.Text(
		text.DescKeyCheckVersionFallback), binaryVer, pluginVer,
	)
	content := message.Load(hook.CheckVersion, hook.VariantMismatch,
		map[string]any{
			version.VarBinary: binaryVer,
			version.VarPlugin: pluginVer,
		}, fallback)
	if content == "" {
		internalIo.TouchFile(markerFile)
		return nil
	}

	boxTitle := desc.Text(text.DescKeyCheckVersionBoxTitle)
	relayPrefix := desc.Text(text.DescKeyCheckVersionRelayPrefix)

	writeSetup.Nudge(cmd, message.NudgeBox(relayPrefix, boxTitle, content))

	ref := notify.NewTemplateRef(hook.CheckVersion, hook.VariantMismatch,
		map[string]any{
			version.VarBinary: binaryVer,
			version.VarPlugin: pluginVer,
		})
	versionMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckVersion, fmt.Sprintf(
			desc.Text(text.DescKeyCheckVersionMismatchRelayFormat),
			binaryVer, pluginVer))
	if err := nudge.EmitAndRelay(versionMsg, input.SessionID, ref); err != nil {
		return err
	}

	internalIo.TouchFile(markerFile)

	// Key age check: piggyback on the daily version check.
	keyBox, keyErr := coreVersion.CheckKeyAge(input.SessionID)
	if keyErr != nil {
		return keyErr
	}
	writeSetup.Nudge(cmd, keyBox)

	return nil
}
