//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/claude"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"

	"github.com/ActiveMemory/ctx/internal/notify"
)

// CheckVersionDrift compares VERSION, plugin.json, and marketplace.json.
// If any differ, it emits a relay box listing the drift. Silent when all match.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - string: JSON hook response to print, or empty string if no drift
func CheckVersionDrift(sessionID string) string {
	fileVer := ReadVersionFile()
	if fileVer == "" {
		return ""
	}

	pluginVer, pluginErr := claude.PluginVersion()
	if pluginErr != nil || pluginVer == "" {
		return ""
	}

	marketVer := ReadMarketplaceVersion()
	if marketVer == "" {
		return ""
	}

	if fileVer == pluginVer && pluginVer == marketVer {
		return ""
	}

	vars := map[string]any{
		"FileVersion":        fileVer,
		"PluginVersion":      pluginVer,
		"MarketplaceVersion": marketVer,
	}
	fallback := "VERSION (" + fileVer + "), plugin.json (" + pluginVer +
		"), marketplace.json (" + marketVer + ") are out of sync. Update all three before releasing."
	msg := message.LoadMessage(hook.VersionDrift, hook.VariantNudge, vars, fallback)
	if msg == "" {
		return ""
	}
	response := coreSession.FormatContext(hook.EventPostToolUse, msg)

	ref := notify.NewTemplateRef(hook.VersionDrift, hook.VariantNudge, vars)
	nudge.Relay(fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.VersionDrift, desc.Text(text.DescKeyVersionDriftRelayMessage)), sessionID, ref)

	return response
}

// ReadVersionFile reads and trims the VERSION file from the project root.
//
// Returns:
//   - string: Version string or empty string
func ReadVersionFile() string {
	data, readErr := os.ReadFile("VERSION")
	if readErr != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// MarketplaceManifest is the structure of .claude-plugin/marketplace.json.
type MarketplaceManifest struct {
	Plugins []struct {
		Version string `json:"version"`
	} `json:"plugins"`
}

// ReadMarketplaceVersion parses .claude-plugin/marketplace.json and returns
// plugins[0].version, or empty string if the file is missing or malformed.
//
// Returns:
//   - string: Version string or empty string
func ReadMarketplaceVersion() string {
	path := filepath.Clean(filepath.Join(".claude-plugin", "marketplace.json"))
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return ""
	}
	var manifest MarketplaceManifest
	if parseErr := json.Unmarshal(data, &manifest); parseErr != nil {
		return ""
	}
	if len(manifest.Plugins) == 0 {
		return ""
	}
	return manifest.Plugins[0].Version
}
