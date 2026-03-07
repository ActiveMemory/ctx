//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkresources

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

// Cmd returns the "ctx system check-resources" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-resources subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-resources",
		Short: "Resource pressure hook",
		Long: `Collects system resource metrics (memory, swap, disk, load) and outputs
a VERBATIM relay warning when any resource hits DANGER severity.
Silent at WARNING level and below.

  Memory DANGER: >= 90% used    Swap DANGER: >= 75% used
  Disk DANGER:   >= 95% full    Load DANGER: >= 1.5x CPUs

For full resource stats at any severity, use: ctx system resources

Hook event: UserPromptSubmit
Output: VERBATIM relay (DANGER only), silent otherwise
Silent when: all resources below DANGER thresholds`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckResources(cmd, os.Stdin)
		},
	}
}

func runCheckResources(cmd *cobra.Command, stdin *os.File) error {
	input := core.ReadInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}
	if core.Paused(sessionID) > 0 {
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
			alertMessages += "\u2716 " + a.Message + config.NewlineLF
		}
	}

	fallback := alertMessages +
		"\nSystem resources are critically low.\n" +
		"Persist unsaved context NOW with /ctx-wrap-up\n" +
		"and consider ending this session."
	content := core.LoadMessage("check-resources", "alert",
		map[string]any{"AlertMessages": alertMessages}, fallback)
	if content == "" {
		return nil
	}

	msg := "IMPORTANT: Relay this resource warning to the user VERBATIM.\n\n" +
		"\u250c\u2500 Resource Alert \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\n"
	msg += core.BoxLines(content)
	if line := core.ContextDirLine(); line != "" {
		msg += "\u2502 " + line + config.NewlineLF
	}
	msg += "\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500"
	cmd.Println(msg)

	ref := notify.NewTemplateRef("check-resources", "alert",
		map[string]any{"AlertMessages": alertMessages})
	_ = notify.Send("nudge", "check-resources: System resources critically low", input.SessionID, ref)
	_ = notify.Send("relay", "check-resources: System resources critically low", input.SessionID, ref)
	eventlog.Append("relay", "check-resources: System resources critically low", input.SessionID, ref)

	return nil
}
