//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd returns the deps command.
func Cmd() *cobra.Command {
	var (
		format   string
		external bool
	)

	cmd := &cobra.Command{
		Use:   "deps",
		Short: "Show package dependency graph",
		Long: `Generate a dependency graph from source code.

Outputs a Mermaid graph of internal package dependencies by default.
Use --external to include external module dependencies.

Supported project types: Go (detected via go.mod).

Output formats:
  mermaid   Mermaid graph definition (default)
  table     Package | Imports table
  json      Machine-readable adjacency list`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDeps(cmd, format, external)
		},
	}

	cmd.Flags().StringVar(&format, "format", "mermaid", "Output format: mermaid, table, json")
	cmd.Flags().BoolVar(&external, "external", false, "Include external module dependencies")

	return cmd
}

func runDeps(cmd *cobra.Command, format string, external bool) error {
	ptype := detectProjectType()
	if ptype == "" {
		cmd.Println("No supported project detected (looking for go.mod).")
		cmd.Println("Currently supported: Go projects.")
		return nil
	}

	switch format {
	case "mermaid", "table", "json":
	default:
		return fmt.Errorf("unknown format %q (supported: mermaid, table, json)", format)
	}

	var graph map[string][]string
	var err error

	if external {
		graph, err = buildFullGraph()
	} else {
		graph, err = buildInternalGraph()
	}
	if err != nil {
		return err
	}

	if len(graph) == 0 {
		cmd.Println("No dependencies found.")
		return nil
	}

	switch format {
	case "mermaid":
		cmd.Print(renderMermaid(graph))
	case "table":
		cmd.Print(renderTable(graph))
	case "json":
		cmd.Print(renderJSON(graph))
	}

	return nil
}
