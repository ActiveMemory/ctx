//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package block_non_path_ctx

import (
	"encoding/json"
	"os"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
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

	if config.RegExRelativeCtxStart.MatchString(command) ||
		config.RegExRelativeCtxSep.MatchString(command) {
		variant = file.VariantDotSlash
		fallback = assets.TextDesc(assets.TextDescKeyBlockDotSlash)
	}

	if config.RegExGoRunCtx.MatchString(command) {
		variant = file.VariantGoRun
		fallback = assets.TextDesc(assets.TextDescKeyBlockGoRun)
	}

	if variant == "" && (config.RegExAbsoluteCtxStart.MatchString(command) ||
		config.RegExAbsoluteCtxSep.MatchString(command)) {
		if !config.RegExCtxTestException.MatchString(command) {
			variant = file.VariantAbsolutePath
			fallback = assets.TextDesc(assets.TextDescKeyBlockAbsolutePath)
		}
	}

	var reason string
	if variant != "" {
		reason = core.LoadMessage(file.HookBlockNonPathCtx, variant, nil, fallback)
	}

	if reason != "" {
		resp := core.BlockResponse{
			Decision: file.HookDecisionBlock,
			Reason: reason + config.NewlineLF + config.NewlineLF +
				assets.TextDesc(assets.TextDescKeyBlockConstitutionSuffix),
		}
		data, _ := json.Marshal(resp)
		cmd.Println(string(data))
		blockRef := notify.NewTemplateRef(file.HookBlockNonPathCtx, variant, nil)
		core.Relay(file.HookBlockNonPathCtx+": "+
			assets.TextDesc(assets.TextDescKeyBlockNonPathRelayMessage),
			input.SessionID, blockRef,
		)
	}

	return nil
}
