//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package plugin

import "encoding/json"

// installedPlugins represents the JSON structure for installed Claude
// Code plugins read from settings files.
type installedPlugins struct {
	Plugins map[string]json.RawMessage `json:"plugins"`
}

// globalSettings represents a Claude Code global settings file as a
// flat key-value map of raw JSON values.
type globalSettings map[string]json.RawMessage
