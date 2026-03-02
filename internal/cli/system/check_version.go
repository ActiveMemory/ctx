//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// checkVersionCmd returns the "ctx system check-version" command.
//
// Compares the binary version (set via ldflags) against the plugin version
// (embedded plugin.json). Warns when the binary is older than the plugin
// expects, which can happen when the marketplace plugin is updated but the
// user hasn't reinstalled the binary. Runs once per day (throttled).
func checkVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-version",
		Short: "Binary/plugin version drift detection hook",
		Long: `Compares the ctx binary version against the embedded plugin version.
Warns when the binary is older than the plugin expects, which happens
when the marketplace plugin updates but the binary hasn't been
reinstalled. Throttled to once per day. Skipped for dev builds.

Hook event: UserPromptSubmit
Output: VERBATIM relay with reinstall command, silent otherwise
Silent when: versions match, dev build, or already checked today`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckVersion(cmd, os.Stdin)
		},
	}
}

func runCheckVersion(cmd *cobra.Command, stdin *os.File) error {
	if !isInitialized() {
		return nil
	}

	input := readInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = sessionUnknown
	}
	if paused(sessionID) > 0 {
		return nil
	}

	tmpDir := secureTempDir()
	markerFile := filepath.Join(tmpDir, "version-checked")

	if isDailyThrottled(markerFile) {
		return nil
	}

	binaryVer := config.BinaryVersion

	// Skip check for dev builds
	if binaryVer == "dev" {
		touchFile(markerFile)
		return nil
	}

	pluginVer, err := assets.PluginVersion()
	if err != nil {
		return nil // embedded plugin.json missing — nothing to compare
	}

	bMajor, bMinor, bOK := parseMajorMinor(binaryVer)
	pMajor, pMinor, pOK := parseMajorMinor(pluginVer)

	if !bOK || !pOK {
		touchFile(markerFile)
		return nil
	}

	if bMajor == pMajor && bMinor == pMinor {
		touchFile(markerFile)
		return nil
	}

	// Version mismatch — emit warning
	fallback := fmt.Sprintf("Your ctx binary is v%s but the plugin expects v%s.\n", binaryVer, pluginVer) +
		"\nReinstall the binary to get the best out of ctx:\n" +
		"  go install github.com/ActiveMemory/ctx/cmd/ctx@latest"
	content := loadMessage("check-version", "mismatch",
		map[string]any{
			"BinaryVersion": binaryVer,
			"PluginVersion": pluginVer,
		}, fallback)
	if content == "" {
		touchFile(markerFile)
		return nil
	}

	msg := "IMPORTANT: Relay this version warning to the user VERBATIM before answering their question.\n\n" +
		"┌─ Version Mismatch ─────────────────────────────\n"
	msg += boxLines(content)
	if line := contextDirLine(); line != "" {
		msg += "│ " + line + "\n"
	}
	msg += "└────────────────────────────────────────────────"
	cmd.Println(msg)

	ref := notify.NewTemplateRef("check-version", "mismatch",
		map[string]any{"BinaryVersion": binaryVer, "PluginVersion": pluginVer})
	versionMsg := fmt.Sprintf("check-version: Binary v%s vs plugin v%s", binaryVer, pluginVer)
	_ = notify.Send("nudge", versionMsg, input.SessionID, ref)
	_ = notify.Send("relay", versionMsg, input.SessionID, ref)
	eventlog.Append("relay", versionMsg, input.SessionID, ref)

	touchFile(markerFile)

	// Key age check — piggyback on the daily version check
	checkKeyAge(cmd, input.SessionID)

	return nil
}

// checkKeyAge emits a nudge when the encryption key is older than the
// configured rotation threshold. Runs at most once per day (shares the
// daily throttle from the version check's marker file).
func checkKeyAge(cmd *cobra.Command, sessionID string) {
	cwd, _ := os.Getwd()
	config.MigrateKeyFile(rc.ContextDir(), cwd)
	kp := rc.KeyPath()
	info, err := os.Stat(kp)
	if err != nil {
		return // no key — nothing to check
	}

	ageDays := int(time.Since(info.ModTime()).Hours() / 24)
	threshold := rc.KeyRotationDays()

	if ageDays < threshold {
		return
	}

	keyFallback := fmt.Sprintf("Your encryption key is %d days old.\n", ageDays) +
		"Consider rotating: ctx pad rotate-key"
	keyContent := loadMessage("check-version", "key-rotation",
		map[string]any{"KeyAgeDays": ageDays}, keyFallback)
	if keyContent == "" {
		return
	}

	keyMsg := "\nIMPORTANT: Relay this security reminder to the user VERBATIM.\n\n" +
		"┌─ Key Rotation ──────────────────────────────────┐\n"
	keyMsg += boxLines(keyContent)
	if line := contextDirLine(); line != "" {
		keyMsg += "│ " + line + "\n"
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
	if _, err := fmt.Sscanf(parts[0], "%d", &m); err != nil {
		return 0, 0, false
	}
	if _, err := fmt.Sscanf(parts[1], "%d", &n); err != nil {
		return 0, 0, false
	}
	return m, n, true
}
