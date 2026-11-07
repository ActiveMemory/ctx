//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// State represents the combined install state of the Antigravity CLI
// and the ctx Antigravity plugin.
type State int

const (
	// StateAbsent means the `agy` binary is not on PATH.
	StateAbsent State = iota
	// StatePluginNotInstalled means `agy` is present but the
	// ctx plugin is not in <configDir>/plugins/ctx/plugin.json.
	StatePluginNotInstalled
	// StatePluginInstalledNotEnabled means the plugin is installed
	// but plugin.json sets "disabled": true.
	StatePluginInstalledNotEnabled
	// StatePluginReady means `agy` is present and the ctx
	// plugin is installed and enabled.
	StatePluginReady
)

// stateKeys maps each State to the text key of its label.
var stateKeys = map[State]string{
	StateAbsent:                    text.DescKeyWriteHookAgyStateAbsent,
	StatePluginNotInstalled:        text.DescKeyWriteHookAgyStateNotInstalled,
	StatePluginInstalledNotEnabled: text.DescKeyWriteHookAgyStateNotEnabled,
	StatePluginReady:               text.DescKeyWriteHookAgyStateReady,
}

// String returns the user-facing label of the state.
//
// Returns:
//   - string: label resolved from the text assets
func (s State) String() string {
	return desc.Text(stateKeys[s])
}

// Outcome reports what a merge helper decided about a file.
type Outcome int

const (
	// OutcomeCreated means there was no existing content and the
	// result is the embedded asset.
	OutcomeCreated Outcome = iota
	// OutcomeMerged means existing content was combined with the
	// embedded asset and the result differs from the input.
	OutcomeMerged
	// OutcomeSkipped means the existing content already matches
	// the desired state; nothing should be written.
	OutcomeSkipped
)
