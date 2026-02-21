//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/sysinfo"
	"github.com/spf13/cobra"
)

// checkResourcesCmd returns the "ctx system check-resources" hook command.
//
// Collects system resource metrics and outputs a VERBATIM relay message when
// any resource is at DANGER severity. Silent otherwise.
func checkResourcesCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "check-resources",
		Short:  "Resource pressure hook",
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckResources(cmd, os.Stdin)
		},
	}
}

func runCheckResources(cmd *cobra.Command, stdin *os.File) error {
	_ = readInput(stdin) // consume stdin per hook contract

	snap := sysinfo.Collect(".")
	alerts := sysinfo.Evaluate(snap)

	if sysinfo.MaxSeverity(alerts) < sysinfo.SeverityDanger {
		return nil
	}

	cmd.Println("IMPORTANT: Relay this resource warning to the user VERBATIM.")
	cmd.Println()
	cmd.Println("\u250c\u2500 Resource Alert \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500")
	for _, a := range alerts {
		if a.Severity == sysinfo.SeverityDanger {
			cmd.Println("\u2502 \u2716 " + a.Message)
		}
	}
	cmd.Println("\u2502")
	cmd.Println("\u2502 System resources are critically low.")
	cmd.Println("\u2502 Persist unsaved context NOW with /ctx-wrap-up")
	cmd.Println("\u2502 and consider ending this session.")
	cmd.Println("\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500")

	return nil
}
