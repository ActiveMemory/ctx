//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package specs_nudge

import (
	"fmt"
	"os"

	hook2 "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/resolve"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
)

// Run executes the specs-nudge hook logic.
//
// Emits a PreToolUse nudge reminding the agent to save plans to specs/
// when a new implementation is detected. Appends a context directory
// footer if available.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !state.Initialized() {
		return nil
	}
	input, _, paused := hook2.Preamble(stdin)
	if paused {
		return nil
	}
	fallback := desc.Text(text.DescKeySpecsNudgeFallback)
	msg := message.LoadMessage(
		hook.SpecsNudge, hook.VariantNudge, nil, fallback,
	)
	if msg == "" {
		return nil
	}
	msg = ctxContext.AppendDir(msg)
	writeHook.HookContext(cmd, coreSession.FormatContext(hook.EventPreToolUse, msg))
	nudgeMsg := desc.Text(text.DescKeySpecsNudgeNudgeMessage)
	ref := notify.NewTemplateRef(hook.SpecsNudge, hook.VariantNudge, nil)
	nudge.Relay(
		fmt.Sprintf(
			desc.Text(text.DescKeyRelayPrefixFormat),
			hook.SpecsNudge, nudgeMsg,
		),
		input.SessionID, ref,
	)
	return nil
}
