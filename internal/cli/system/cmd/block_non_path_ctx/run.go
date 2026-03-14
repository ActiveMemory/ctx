//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package block_non_path_ctx

import (
	"encoding/json"
	"os"

	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Run executes the block-non-path-ctx hook logic.
//
// Reads a hook input from stdin, checks the command against patterns
// that invoke ctx via relative paths, go run, or absolute paths
// instead of the PATH-installed binary, and emits a block response
// if matched.
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

	if regex.CtxRelativeStart.MatchString(command) ||
		regex.CtxRelativeSep.MatchString(command) {
		variant = hook.VariantDotSlash
		fallback = assets.TextDesc(assets.TextDescKeyBlockDotSlash)
	}

	if regex.CtxGoRun.MatchString(command) {
		variant = hook.VariantGoRun
		fallback = assets.TextDesc(assets.TextDescKeyBlockGoRun)
	}

	if variant == "" && (regex.CtxAbsoluteStart.MatchString(command) ||
		regex.AbsoluteSep.MatchString(command)) {
		if !regex.CtxTestException.MatchString(command) {
			variant = hook.VariantAbsolutePath
			fallback = assets.TextDesc(assets.TextDescKeyBlockAbsolutePath)
		}
	}

	var reason string
	if variant != "" {
		reason = core.LoadMessage(hook.BlockNonPathCtx, variant, nil, fallback)
	}

	if reason != "" {
		resp := core.BlockResponse{
			Decision: hook.HookDecisionBlock,
			Reason: reason + token.NewlineLF + token.NewlineLF +
				assets.TextDesc(assets.TextDescKeyBlockConstitutionSuffix),
		}
		data, _ := json.Marshal(resp)
		cmd.Println(string(data))
		blockRef := notify.NewTemplateRef(hook.BlockNonPathCtx, variant, nil)
		core.Relay(hook.BlockNonPathCtx+": "+
			assets.TextDesc(assets.TextDescKeyBlockNonPathRelayMessage),
			input.SessionID, blockRef,
		)
	}

	return nil
}
