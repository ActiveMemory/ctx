//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude_check

import (
	"github.com/spf13/cobra"

	writeInit "github.com/ActiveMemory/ctx/internal/write/initialize"
)

// InitHint prints stage-aware Claude Code setup guidance as
// a post-script at the end of `ctx init`. Never writes
// files, never errors, never fatal — it's a friendly nudge
// matching whichever step the user still needs to complete.
//
// State-to-output mapping:
//
//   - StateClaudeAbsent: print the two-step install path
//     (install Claude Code first, then the ctx plugin).
//   - StatePluginNotInstalled: print the dev-symlink
//     install flow with user-scope guidance.
//   - StatePluginInstalledNotEnabled: stay silent; the
//     EnableLocally step earlier in `ctx init` already
//     emitted its own confirmation for this state.
//   - StatePluginReady: print the detail block with
//     scope/version/source/clone path/enablement, or the
//     minimal one-liner if metadata parsing failed.
//
// Parameters:
//   - cmd: Cobra command for output
func InitHint(cmd *cobra.Command) {
	switch Detect() {
	case StateClaudeAbsent:
		writeInit.ClaudeAbsent(cmd)
	case StatePluginNotInstalled:
		writeInit.ClaudePluginMissing(cmd)
	case StatePluginReady:
		if d, ok := Details(); ok {
			scope, version, source, clonePath, enabled := renderDetails(d)
			writeInit.ClaudeReady(
				cmd, scope, version, source, clonePath, enabled,
			)
			return
		}
		writeInit.ClaudeReadyMinimal(cmd)
	case StatePluginInstalledNotEnabled:
		// Handled by the EnableLocally confirmation
		// printed earlier in the init flow.
	}
}

// SetupHint prints stage-aware Claude Code setup guidance
// as the primary output of `ctx setup claude-code`. Unlike
// the other `ctx setup <tool>` commands, Claude Code has no
// writable config file ctx can emit directly — the
// integration is delivered via the ctx plugin installed
// from the user's local clone.
//
// State-to-output mapping:
//
//   - StateClaudeAbsent: print the full two-step install
//     path.
//   - StatePluginNotInstalled: print the plugin install
//     flow.
//   - StatePluginInstalledNotEnabled: print the same
//     install flow, which ends with "re-run `ctx init` to
//     enable locally" — the action the user needs.
//   - StatePluginReady: print the detail block + setup
//     ready message, or the minimal variant on metadata
//     parse failure.
//
// Parameters:
//   - cmd: Cobra command for output
func SetupHint(cmd *cobra.Command) {
	switch Detect() {
	case StateClaudeAbsent:
		writeInit.ClaudeAbsent(cmd)
	case StatePluginNotInstalled:
		writeInit.ClaudePluginMissing(cmd)
	case StatePluginInstalledNotEnabled:
		writeInit.ClaudePluginMissing(cmd)
	case StatePluginReady:
		if d, ok := Details(); ok {
			scope, version, source, clonePath, enabled := renderDetails(d)
			writeInit.SetupClaudeReady(
				cmd, scope, version, source, clonePath, enabled,
			)
			return
		}
		writeInit.SetupClaudeReadyMinimal(cmd)
	}
}
