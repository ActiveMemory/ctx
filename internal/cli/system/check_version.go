//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
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
			return runCheckVersion(cmd)
		},
	}
}

func runCheckVersion(cmd *cobra.Command) error {
	if !isInitialized() {
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
	cmd.Println("IMPORTANT: Relay this version warning to the user VERBATIM before answering their question.")
	cmd.Println()
	cmd.Println("┌─ Version Mismatch ─────────────────────────────")
	cmd.Println(fmt.Sprintf("│ Your ctx binary is v%s but the plugin expects v%s.", binaryVer, pluginVer))
	cmd.Println("│")
	cmd.Println("│ Reinstall the binary to get the best out of ctx:")
	cmd.Println("│   go install github.com/ActiveMemory/ctx/cmd/ctx@latest")
	cmd.Println("└────────────────────────────────────────────────")

	touchFile(markerFile)
	return nil
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
