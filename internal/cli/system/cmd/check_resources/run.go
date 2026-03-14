//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_resources

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run executes the check-resources hook logic.
//
// Collects system resource snapshots, evaluates alert thresholds, and
// emits a relay warning box when any resource is at danger level.
// Throttled by session pause state.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	input, _, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	snap := sysinfo.Collect(".")
	alerts := sysinfo.Evaluate(snap)

	if sysinfo.MaxSeverity(alerts) < sysinfo.SeverityDanger {
		return nil
	}

	// Build pre-formatted alert messages for the template variable
	var alertMessages string
	for _, a := range alerts {
		if a.Severity == sysinfo.SeverityDanger {
			alertMessages += "✖ " +
				a.Message + token.NewlineLF
		}
	}

	fallback := alertMessages +
		token.NewlineLF + assets.TextDesc(
		assets.TextDescKeyCheckResourcesFallbackLow) + token.NewlineLF +
		assets.TextDesc(
			assets.TextDescKeyCheckResourcesFallbackPersist) + token.NewlineLF +
		assets.TextDesc(
			assets.TextDescKeyCheckResourcesFallbackEnd)
	vars := map[string]any{tpl.VarAlertMessages: alertMessages}
	content := core.LoadMessage(
		hook.CheckResources, hook.VariantAlert, vars, fallback,
	)
	if content == "" {
		return nil
	}

	write.HookNudge(cmd, core.NudgeBox(
		assets.TextDesc(assets.TextDescKeyCheckResourcesRelayPrefix),
		assets.TextDesc(assets.TextDescKeyCheckResourcesBoxTitle),
		content))

	ref := notify.NewTemplateRef(
		hook.CheckResources, hook.VariantAlert, vars,
	)
	core.NudgeAndRelay(hook.CheckResources+": "+
		assets.TextDesc(assets.TextDescKeyCheckResourcesRelayMessage),
		input.SessionID, ref,
	)

	return nil
}
