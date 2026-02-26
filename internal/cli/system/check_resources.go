//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

// checkResourcesCmd returns the "ctx system check-resources" hook command.
//
// Collects system resource metrics and outputs a VERBATIM relay message when
// any resource is at DANGER severity. Silent otherwise.
//
// Unlike other hook commands, this does not read stdin — the hook input
// (session_id, tool_input) is not needed for resource checks, and blocking
// on io.ReadAll would hang if stdin is a terminal.
func checkResourcesCmd() *cobra.Command {
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
	input := readInput(stdin)

	snap := sysinfo.Collect(".")
	alerts := sysinfo.Evaluate(snap)

	if sysinfo.MaxSeverity(alerts) < sysinfo.SeverityDanger {
		return nil
	}

	msg := "IMPORTANT: Relay this resource warning to the user VERBATIM.\n\n" +
		"\u250c\u2500 Resource Alert \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\n"
	for _, a := range alerts {
		if a.Severity == sysinfo.SeverityDanger {
			msg += "\u2502 \u2716 " + a.Message + "\n"
		}
	}
	msg += "\u2502\n" +
		"\u2502 System resources are critically low.\n" +
		"\u2502 Persist unsaved context NOW with /ctx-wrap-up\n" +
		"\u2502 and consider ending this session.\n"
	if line := contextDirLine(); line != "" {
		msg += "\u2502 " + line + "\n"
	}
	msg += "\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500"
	cmd.Println(msg)

	_ = notify.Send("nudge", "check-resources: System resources critically low", input.SessionID, msg)
	_ = notify.Send("relay", "check-resources: System resources critically low", input.SessionID, msg)

	return nil
}
