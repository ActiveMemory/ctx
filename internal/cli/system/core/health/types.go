//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package health

// MapTrackingInfo holds the minimal fields needed from map-tracking.json.
type MapTrackingInfo struct {
	OptedOut bool   `json:"opted_out"`
	LastRun  string `json:"last_run"`
}
