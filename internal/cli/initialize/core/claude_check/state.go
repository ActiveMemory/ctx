//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude_check

import (
	"os/exec"

	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/plugin"
	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
)

// Detect returns the current combined state of Claude Code
// and the ctx plugin.
//
// Detect never returns an error: any failure to read plugin
// metadata is treated as a "not yet installed" signal,
// which matches how the init flow already handles missing
// files.
//
// Returns:
//   - State: the current combined state
func Detect() State {
	if _, lookErr := exec.LookPath(cfgClaude.Binary); lookErr != nil {
		return StateClaudeAbsent
	}
	if !plugin.Installed() {
		return StatePluginNotInstalled
	}
	if !plugin.EnabledGlobally() && !plugin.EnabledLocally() {
		return StatePluginInstalledNotEnabled
	}
	return StatePluginReady
}
