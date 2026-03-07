//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkversion

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func runCheckVersion(cmd *cobra.Command, stdin *os.File) error {
	if !core.IsInitialized() {
		return nil
	}

	input := core.ReadInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}
	if core.Paused(sessionID) > 0 {
		return nil
	}

	tmpDir := core.StateDir()
	markerFile := filepath.Join(tmpDir, "version-checked")

	if core.IsDailyThrottled(markerFile) {
		return nil
	}

	binaryVer := config.BinaryVersion

	// Skip check for dev builds
	if binaryVer == "dev" {
		core.TouchFile(markerFile)
		return nil
	}

	pluginVer, pluginErr := assets.PluginVersion()
	if pluginErr != nil {
		return nil // embedded plugin.json missing — nothing to compare
	}

	bMajor, bMinor, bOK := parseMajorMinor(binaryVer)
	pMajor, pMinor, pOK := parseMajorMinor(pluginVer)

	if !bOK || !pOK {
		core.TouchFile(markerFile)
		return nil
	}

	if bMajor == pMajor && bMinor == pMinor {
		core.TouchFile(markerFile)
		return nil
	}

	// Version mismatch — emit warning
	fallback := fmt.Sprintf("Your ctx binary is v%s but the plugin expects v%s.\n", binaryVer, pluginVer) +
		"\nReinstall the binary to get the best out of ctx:\n" +
		"  go install github.com/ActiveMemory/ctx/cmd/ctx@latest"
	content := core.LoadMessage("check-version", "mismatch",
		map[string]any{
			"BinaryVersion": binaryVer,
			"PluginVersion": pluginVer,
		}, fallback)
	if content == "" {
		core.TouchFile(markerFile)
		return nil
	}

	msg := "IMPORTANT: Relay this version warning to the user VERBATIM before answering their question.\n\n" +
		"┌─ Version Mismatch ─────────────────────────────\n"
	msg += core.BoxLines(content)
	if line := core.ContextDirLine(); line != "" {
		msg += "│ " + line + config.NewlineLF
	}
	msg += "└────────────────────────────────────────────────"
	cmd.Println(msg)

	ref := notify.NewTemplateRef("check-version", "mismatch",
		map[string]any{"BinaryVersion": binaryVer, "PluginVersion": pluginVer})
	versionMsg := fmt.Sprintf("check-version: Binary v%s vs plugin v%s", binaryVer, pluginVer)
	_ = notify.Send("nudge", versionMsg, input.SessionID, ref)
	_ = notify.Send("relay", versionMsg, input.SessionID, ref)
	eventlog.Append("relay", versionMsg, input.SessionID, ref)

	core.TouchFile(markerFile)

	// Key age check — piggyback on the daily version check
	checkKeyAge(cmd, input.SessionID)

	return nil
}

// checkKeyAge emits a nudge when the encryption key is older than the
// configured rotation threshold.
func checkKeyAge(cmd *cobra.Command, sessionID string) {
	config.MigrateKeyFile(rc.ContextDir())
	kp := rc.KeyPath()
	info, statErr := os.Stat(kp)
	if statErr != nil {
		return // no key — nothing to check
	}

	ageDays := int(time.Since(info.ModTime()).Hours() / 24)
	threshold := rc.KeyRotationDays()

	if ageDays < threshold {
		return
	}

	keyFallback := fmt.Sprintf("Your encryption key is %d days old.\n", ageDays) +
		"Consider rotating: ctx pad rotate-key"
	keyContent := core.LoadMessage("check-version", "key-rotation",
		map[string]any{"KeyAgeDays": ageDays}, keyFallback)
	if keyContent == "" {
		return
	}

	keyMsg := "\nIMPORTANT: Relay this security reminder to the user VERBATIM.\n\n" +
		"┌─ Key Rotation ──────────────────────────────────┐\n"
	keyMsg += core.BoxLines(keyContent)
	if line := core.ContextDirLine(); line != "" {
		keyMsg += "│ " + line + config.NewlineLF
	}
	keyMsg += "└──────────────────────────────────────────────────┘"
	cmd.Println(keyMsg)

	keyRef := notify.NewTemplateRef("check-version", "key-rotation",
		map[string]any{"KeyAgeDays": ageDays})
	keyNotifyMsg := fmt.Sprintf("check-version: Encryption key is %d days old", ageDays)
	_ = notify.Send("nudge", keyNotifyMsg, sessionID, keyRef)
	_ = notify.Send("relay", keyNotifyMsg, sessionID, keyRef)
	eventlog.Append("relay", keyNotifyMsg, sessionID, keyRef)
}

// parseMajorMinor extracts major and minor version numbers from a semver
// string like "1.2.3". Returns ok=false for unparseable versions.
func parseMajorMinor(ver string) (major, minor int, ok bool) {
	parts := strings.SplitN(ver, ".", 3)
	if len(parts) < 2 {
		return 0, 0, false
	}
	var m, n int
	if _, scanErr := fmt.Sscanf(parts[0], "%d", &m); scanErr != nil {
		return 0, 0, false
	}
	if _, scanErr := fmt.Sscanf(parts[1], "%d", &n); scanErr != nil {
		return 0, 0, false
	}
	return m, n, true
}
