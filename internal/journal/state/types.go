//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package state

// State is the top-level state file structure.
type State struct {
	Version int             `json:"version"`
	Entries map[string]File `json:"entries"`
}

// File tracks processing stages for a single journal entry.
// Values are date strings (YYYY-MM-DD) indicating when the stage completed.
type File struct {
	Exported       string `json:"exported,omitempty"`
	Enriched       string `json:"enriched,omitempty"`
	Normalized     string `json:"normalized,omitempty"`
	FencesVerified string `json:"fences_verified,omitempty"`
	Locked         string `json:"locked,omitempty"`
}
