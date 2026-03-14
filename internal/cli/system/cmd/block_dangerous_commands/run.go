//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package block_dangerous_commands

import (
	"encoding/json"
	"os"

	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Run executes the block-dangerous-commands hook logic.
//
// Reads a hook input from stdin, checks the command against dangerous
// patterns (mid-command sudo, git push, cp/mv to bin), and emits a
// block response if matched.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	input := core.ReadInput(stdin)
	command := input.ToolInput.Command

	if command == "" {
		return nil
	}

	var variant, fallback string

	if regex.MidSudo.MatchString(command) {
		variant = hook.VariantMidSudo
		fallback = assets.TextDesc(assets.TextDescKeyBlockMidSudo)
	}

	if variant == "" && regex.MidGitPush.MatchString(command) {
		variant = hook.VariantMidGitPush
		fallback = assets.TextDesc(assets.TextDescKeyBlockMidGitPush)
	}

	if variant == "" && regex.CpMvToBin.MatchString(command) {
		variant = hook.VariantCpToBin
		fallback = assets.TextDesc(assets.TextDescKeyBlockCpToBin)
	}

	if variant == "" && regex.InstallToLocalBin.MatchString(command) {
		variant = hook.VariantInstallToLocalBin
		fallback = assets.TextDesc(assets.TextDescKeyBlockInstallToLocalBin)
	}

	var reason string
	if variant != "" {
		reason = core.LoadMessage(
			hook.BlockDangerousCommands, variant, nil, fallback,
		)
	}

	if reason != "" {
		resp := core.BlockResponse{
			Decision: hook.HookDecisionBlock,
			Reason:   reason,
		}
		data, _ := json.Marshal(resp)
		cmd.Println(string(data))
		ref := notify.NewTemplateRef(hook.BlockDangerousCommands, variant, nil)
		core.Relay(hook.BlockDangerousCommands+": "+reason, input.SessionID, ref)
	}

	return nil
}
