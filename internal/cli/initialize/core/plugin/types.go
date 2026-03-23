//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package plugin

import "encoding/json"

type installedPlugins struct {
	Plugins map[string]json.RawMessage `json:"plugins"`
}

type globalSettings map[string]json.RawMessage
