//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package block_non_path_ctx

import (
	"encoding/json"
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/spf13/cobra"

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
		fallback = desc.Text(text.DescKeyBlockDotSlash)
	}

	if regex.CtxGoRun.MatchString(command) {
		variant = hook.VariantGoRun
		fallback = desc.Text(text.DescKeyBlockGoRun)
	}

	if variant == "" && (regex.CtxAbsoluteStart.MatchString(command) ||
		regex.AbsoluteSep.MatchString(command)) {
		if !regex.CtxTestException.MatchString(command) {
			variant = hook.VariantAbsolutePath
			fallback = desc.Text(text.DescKeyBlockAbsolutePath)
		}
	}

	var reason string
	if variant != "" {
		reason = core.LoadMessage(hook.BlockNonPathCtx, variant, nil, fallback)
	}

	if reason != "" {
		resp := core.BlockResponse{
			Decision: hook.DecisionBlock,
			Reason: reason + token.NewlineLF + token.NewlineLF +
				desc.Text(text.DescKeyBlockConstitutionSuffix),
		}
		data, _ := json.Marshal(resp)
		cmd.Println(string(data))
		blockRef := notify.NewTemplateRef(hook.BlockNonPathCtx, variant, nil)
		core.Relay(hook.BlockNonPathCtx+": "+
			desc.Text(text.DescKeyBlockNonPathRelayMessage),
			input.SessionID, blockRef,
		)
	}

	return nil
}
