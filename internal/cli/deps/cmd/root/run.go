//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/fmt"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/deps/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run executes the deps command logic.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - format: Output format (config.FormatMermaid, config.FormatTable, or config.FormatJSON)
//   - external: If true, include external module dependencies
//   - projType: Force project type override; empty for auto-detect
//
// Returns:
//   - error: Non-nil if format is invalid, project type unknown,
//     or graph building fails
func Run(cmd *cobra.Command, format string, external bool, projType string) error {
	supportedFormats := strings.Join([]string{
		fmt.FormatMermaid, fmt.FormatTable, fmt.FormatJSON,
	}, ", ")

	switch format {
	case fmt.FormatMermaid, fmt.FormatTable, fmt.FormatJSON:
	default:
		return ctxerr.UnknownFormat(format, supportedFormats)
	}

	var builder core.GraphBuilder
	if projType != "" {
		builder = core.FindBuilder(projType)
		if builder == nil {
			return ctxerr.UnknownProjectType(projType, strings.Join(core.BuilderNames(), ", "))
		}
	} else {
		builder = core.DetectBuilder()
		if builder == nil {
			write.InfoDepsNoProject(cmd, strings.Join(core.BuilderNames(), ", "))
			return nil
		}
	}

	graph, buildErr := builder.Build(external)
	if buildErr != nil {
		return buildErr
	}

	if len(graph) == 0 {
		write.InfoDepsNoDeps(cmd)
		return nil
	}

	switch format {
	case fmt.FormatMermaid:
		cmd.Print(core.RenderMermaid(graph))
	case fmt.FormatTable:
		cmd.Print(core.RenderTable(graph))
	case fmt.FormatJSON:
		cmd.Print(core.RenderJSON(graph))
	}

	return nil
}
