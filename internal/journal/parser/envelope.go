//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import "encoding/json"

// parseOriginKind extracts the "kind" value from an origin
// JSON object. Returns empty string if origin is nil, not an
// object, or has no "kind" field.
//
// Parameters:
//   - raw: JSON-encoded origin value (e.g. {"kind":"task-notification"})
//
// Returns:
//   - string: the kind value, or empty on any failure
func parseOriginKind(raw json.RawMessage) string {
	if raw == nil {
		return ""
	}
	var origin struct {
		Kind string `json:"kind"`
	}
	if unmarshalErr := json.Unmarshal(raw, &origin); unmarshalErr != nil {
		return ""
	}
	return origin.Kind
}
