//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"encoding/json"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// PluginVersion returns the version string from the embedded plugin.json.
//
// Returns:
//   - string: Semver version (e.g. "0.8.1")
//   - error: If the embedded file cannot be read or parsed
func PluginVersion() (string, error) {
	data, readErr := assets.FS.ReadFile(asset.PathPluginJSON)
	if readErr != nil {
		return "", readErr
	}
	var manifest map[string]json.RawMessage
	if unmarshalErr := json.Unmarshal(data, &manifest); unmarshalErr != nil {
		return "", unmarshalErr
	}
	raw, ok := manifest[asset.JSONKeyVersion]
	if !ok {
		return "", nil
	}
	var version string
	if parseErr := json.Unmarshal(raw, &version); parseErr != nil {
		return "", parseErr
	}
	return version, nil
}
