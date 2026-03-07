//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"github.com/spf13/cobra"

	depsroot "github.com/ActiveMemory/ctx/internal/cli/deps/cmd/root"
)

// Cmd returns the deps command.
//
// Flags:
//   - --format: Output format (mermaid, table, json)
//   - --external: Include external module dependencies
//   - --type: Force project type override
//
// Returns:
//   - *cobra.Command: Configured deps command with flags registered
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
			return depsroot.Run(cmd, format, external, projType)
		},
	}

	cmd.Flags().StringVar(&format, "format", "mermaid", "Output format: mermaid, table, json")
	cmd.Flags().BoolVar(&external, "external", false, "Include external module dependencies")
	cmd.Flags().StringVar(&projType, "type", "", "Force project type: go, node, python, rust")

	return cmd
}
