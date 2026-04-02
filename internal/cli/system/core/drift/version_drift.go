//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/claude"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	cfgVersion "github.com/ActiveMemory/ctx/internal/config/version"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// CheckVersion compares VERSION, plugin.json, and marketplace.json.
// If any differ, it emits a relay box listing the drift. Silent when all match.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - string: JSON hook response to print, or empty string if no drift
func CheckVersion(sessionID string) string {
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
		cfgVersion.VarFile:        fileVer,
		cfgVersion.VarPlugin:      pluginVer,
		cfgVersion.VarMarketplace: marketVer,
	}
	fallback := fmt.Sprintf(
		desc.Text(text.DescKeyWriteVersionDriftFallback),
		fileVer, pluginVer, marketVer,
	)
	msg := message.Load(hook.VersionDrift, hook.VariantNudge, vars, fallback)
	if msg == "" {
		return ""
	}
	response := coreSession.FormatContext(hook.EventPostToolUse, msg)

	ref := notify.NewTemplateRef(hook.VersionDrift, hook.VariantNudge, vars)
	nudge.Relay(fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.VersionDrift,
		desc.Text(text.DescKeyVersionDriftRelayMessage),
	), sessionID, ref)

	return response
}

// ReadVersionFile reads and trims the VERSION file from the project root.
//
// Returns:
//   - string: Version string or empty string
func ReadVersionFile() string {
	data, readErr := ctxIo.SafeReadUserFile("VERSION")
	if readErr != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// ReadMarketplaceVersion parses .claude-plugin/marketplace.json and returns
// plugins[0].version, or empty string if the file is missing or malformed.
//
// Returns:
//   - string: Version string or empty string
func ReadMarketplaceVersion() string {
	path := filepath.Clean(
		filepath.Join(cfgVersion.DirClaudePlugin, cfgVersion.FileMarketplace),
	)
	data, readErr := ctxIo.SafeReadUserFile(path)
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
