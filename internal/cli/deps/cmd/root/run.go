//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/deps/core"
)

// Run executes the deps command logic.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - format: Output format ("mermaid", "table", or "json")
//   - external: If true, include external module dependencies
//   - projType: Force project type override; empty for auto-detect
//
// Returns:
//   - error: Non-nil if format is invalid, project type unknown,
//     or graph building fails
func Run(cmd *cobra.Command, format string, external bool, projType string) error {
	switch format {
	case "mermaid", "table", "json":
	default:
		return fmt.Errorf("unknown format %q (supported: mermaid, table, json)", format)
	}

	var builder core.GraphBuilder
	if projType != "" {
		builder = core.FindBuilder(projType)
		if builder == nil {
			return fmt.Errorf("unknown project type %q (supported: %s)", projType, strings.Join(core.BuilderNames(), ", "))
		}
	} else {
		builder = core.DetectBuilder()
		if builder == nil {
			cmd.Println("No supported project detected.")
			cmd.Println("Looking for: go.mod, package.json, requirements.txt, pyproject.toml, Cargo.toml")
			cmd.Println("Use --type to force: " + strings.Join(core.BuilderNames(), ", "))
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
		cmd.Print(core.RenderMermaid(graph))
	case "table":
		cmd.Print(core.RenderTable(graph))
	case "json":
		cmd.Print(core.RenderJSON(graph))
	}

	return nil
}
