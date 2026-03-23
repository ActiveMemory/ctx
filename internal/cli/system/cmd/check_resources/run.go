//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_resources

import (
	"fmt"
	"os"

	hook2 "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
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
	input, _, paused := hook2.Preamble(stdin)
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
			alertMessages += stats.IconError + token.Space +
				a.Message + token.NewlineLF
		}
	}

	fallback := alertMessages +
		token.NewlineLF + desc.Text(
		text.DescKeyCheckResourcesFallbackLow) + token.NewlineLF +
		desc.Text(
			text.DescKeyCheckResourcesFallbackPersist) + token.NewlineLF +
		desc.Text(
			text.DescKeyCheckResourcesFallbackEnd)
	vars := map[string]any{stats.VarAlertMessages: alertMessages}
	content := core.LoadMessage(
		hook.CheckResources, hook.VariantAlert, vars, fallback,
	)
	if content == "" {
		return nil
	}

	writeHook.Nudge(cmd, core.NudgeBox(
		desc.Text(text.DescKeyCheckResourcesRelayPrefix),
		desc.Text(text.DescKeyCheckResourcesBoxTitle),
		content))

	ref := notify.NewTemplateRef(
		hook.CheckResources, hook.VariantAlert, vars,
	)
	core.NudgeAndRelay(fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckResources, desc.Text(text.DescKeyCheckResourcesRelayMessage)),
		input.SessionID, ref,
	)

	return nil
}
