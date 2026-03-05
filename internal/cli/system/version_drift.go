//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// checkVersionDrift compares VERSION, plugin.json, and marketplace.json.
// If any differ, it emits a relay box listing the drift. Silent when all match.
func checkVersionDrift(cmd *cobra.Command, sessionID string) {
	fileVer := readVersionFile()
	if fileVer == "" {
		return
	}

	pluginVer, pluginErr := assets.PluginVersion()
	if pluginErr != nil || pluginVer == "" {
		return
	}

	marketVer := readMarketplaceVersion()
	if marketVer == "" {
		return
	}

	if fileVer == pluginVer && pluginVer == marketVer {
		return
	}

	vars := map[string]any{
		"FileVersion":        fileVer,
		"PluginVersion":      pluginVer,
		"MarketplaceVersion": marketVer,
	}
	fallback := "VERSION (" + fileVer + "), plugin.json (" + pluginVer +
		"), marketplace.json (" + marketVer + ") are out of sync. Update all three before releasing."
	msg := loadMessage("version-drift", "nudge", vars, fallback)
	if msg == "" {
		return
	}
	printHookContext(cmd, "PostToolUse", msg)

	ref := notify.NewTemplateRef("version-drift", "nudge", vars)
	_ = notify.Send("relay", "version-drift: versions out of sync", sessionID, ref)
	eventlog.Append("relay", "version-drift: versions out of sync", sessionID, ref)
}

// readVersionFile reads and trims the VERSION file from the project root.
func readVersionFile() string {
	data, readErr := os.ReadFile("VERSION")
	if readErr != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// marketplaceManifest is the structure of .claude-plugin/marketplace.json.
type marketplaceManifest struct {
	Plugins []struct {
		Version string `json:"version"`
	} `json:"plugins"`
}

// readMarketplaceVersion parses .claude-plugin/marketplace.json and returns
// plugins[0].version, or empty string if the file is missing or malformed.
func readMarketplaceVersion() string {
	path := filepath.Clean(filepath.Join(".claude-plugin", "marketplace.json"))
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return ""
	}
	var manifest marketplaceManifest
	if parseErr := json.Unmarshal(data, &manifest); parseErr != nil {
		return ""
	}
	if len(manifest.Plugins) == 0 {
		return ""
	}
	return manifest.Plugins[0].Version
}
