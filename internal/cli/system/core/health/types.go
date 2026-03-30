//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package health

// MapTrackingInfo holds the minimal fields needed from map-tracking.json.
//
// Fields:
//   - OptedOut: User opted out of architecture mapping
//   - LastRun: ISO 8601 timestamp of last mapping run
type MapTrackingInfo struct {
	OptedOut bool   `json:"opted_out"`
	LastRun  string `json:"last_run"`
}
