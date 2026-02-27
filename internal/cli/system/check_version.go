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

	_ = notify.Send("nudge", fmt.Sprintf("check-version: Binary v%s vs plugin v%s", binaryVer, pluginVer), input.SessionID, msg)
	_ = notify.Send("relay", fmt.Sprintf("check-version: Binary v%s vs plugin v%s", binaryVer, pluginVer), input.SessionID, msg)

	touchFile(markerFile)

	// Key age check — piggyback on the daily version check
	checkKeyAge(cmd, input.SessionID)

	return nil
}

// checkKeyAge emits a nudge when the encryption key is older than the
// configured rotation threshold. Runs at most once per day (shares the
// daily throttle from the version check's marker file).
func checkKeyAge(cmd *cobra.Command, sessionID string) {
	config.MigrateKeyFile(rc.ContextDir())
	keyPath := filepath.Join(rc.ContextDir(), config.FileContextKey)
	info, err := os.Stat(keyPath)
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

	_ = notify.Send("nudge", fmt.Sprintf("check-version: Encryption key is %d days old", ageDays), sessionID, keyMsg)
	_ = notify.Send("relay", fmt.Sprintf("check-version: Encryption key is %d days old", ageDays), sessionID, keyMsg)
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
