//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// State represents the combined install state of the Codex CLI
// and the ctx Codex plugin.
type State int

const (
	// StateAbsent means the `codex` binary is not on PATH.
	StateAbsent State = iota
	// StatePluginNotInstalled means `codex` is present but the
	// ctx plugin is not in the Codex plugin cache.
	StatePluginNotInstalled
	// StatePluginInstalledNotEnabled means the plugin is cached
	// but `~/.codex/config.toml` does not enable it.
	StatePluginInstalledNotEnabled
	// StatePluginReady means `codex` is present and the ctx
	// plugin is installed and enabled.
	StatePluginReady
)

// stateKeys maps each State to the text key of its label.
var stateKeys = map[State]string{
	StateAbsent:                    text.DescKeyWriteHookCodexStateAbsent,
	StatePluginNotInstalled:        text.DescKeyWriteHookCodexStateNotInstalled,
	StatePluginInstalledNotEnabled: text.DescKeyWriteHookCodexStateNotEnabled,
	StatePluginReady:               text.DescKeyWriteHookCodexStateReady,
}

// String returns the user-facing label of the state.
//
// Returns:
//   - string: label resolved from the text assets (empty for an
//     unknown state)
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
