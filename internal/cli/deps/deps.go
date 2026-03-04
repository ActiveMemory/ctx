//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Cmd returns the deps command.
func Cmd() *cobra.Command {
	var (
		format   string
		external bool
		projType string
	)

	cmd := &cobra.Command{
		Use:   "deps",
		Short: "Show package dependency graph",
		Long: `Generate a dependency graph from source code.

Outputs a Mermaid graph of internal package dependencies by default.
Use --external to include external module dependencies.

Supported project types: Go, Node.js, Python, Rust.
Auto-detected from manifest files (go.mod, package.json,
requirements.txt/pyproject.toml, Cargo.toml). Use --type to override.

Output formats:
  mermaid   Mermaid graph definition (default)
  table     Package | Imports table
  json      Machine-readable adjacency list`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDeps(cmd, format, external, projType)
		},
	}

	cmd.Flags().StringVar(&format, "format", "mermaid", "Output format: mermaid, table, json")
	cmd.Flags().BoolVar(&external, "external", false, "Include external module dependencies")
	cmd.Flags().StringVar(&projType, "type", "", "Force project type: "+strings.Join(builderNames(), ", "))

	return cmd
}

func runDeps(cmd *cobra.Command, format string, external bool, projType string) error {
	switch format {
	case "mermaid", "table", "json":
	default:
		return fmt.Errorf("unknown format %q (supported: mermaid, table, json)", format)
	}

	var builder GraphBuilder
	if projType != "" {
		builder = findBuilder(projType)
		if builder == nil {
			return fmt.Errorf("unknown project type %q (supported: %s)", projType, strings.Join(builderNames(), ", "))
		}
	} else {
		builder = detectBuilder()
		if builder == nil {
			cmd.Println("No supported project detected.")
			cmd.Println("Looking for: go.mod, package.json, requirements.txt, pyproject.toml, Cargo.toml")
			cmd.Println("Use --type to force: " + strings.Join(builderNames(), ", "))
			return nil
		}
	}

	graph, buildErr := builder.Build(external)
	if buildErr != nil {
		return buildErr
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
