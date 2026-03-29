//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/dep/core"
	"github.com/ActiveMemory/ctx/internal/config/fmt"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errConfig "github.com/ActiveMemory/ctx/internal/err/config"
	"github.com/ActiveMemory/ctx/internal/write/deps"
)

// Run executes the deps command logic.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - format: Output format (config.FormatMermaid,
//     config.FormatTable, or config.FormatJSON)
//   - external: If true, include external module dependencies
//   - projType: Force project type override; empty for auto-detect
//
// Returns:
//   - error: Non-nil if format is invalid, project type unknown,
//     or graph building fails
func Run(
	cmd *cobra.Command, format string,
	external bool, projType string,
) error {
	supportedFormats := strings.Join([]string{
		fmt.FormatMermaid, fmt.FormatTable, fmt.FormatJSON,
	}, token.CommaSpace)

	switch format {
	case fmt.FormatMermaid, fmt.FormatTable, fmt.FormatJSON:
	default:
		return errConfig.UnknownFormat(format, supportedFormats)
	}

	var builder core.GraphBuilder
	if projType != "" {
		builder = core.FindBuilder(projType)
		if builder == nil {
			names := strings.Join(core.BuilderNames(), token.CommaSpace)
			return errConfig.UnknownProjectType(projType, names)
		}
	} else {
		builder = core.DetectBuilder()
		if builder == nil {
			deps.InfoNoProject(cmd, strings.Join(core.BuilderNames(), token.CommaSpace))
			return nil
		}
	}

	graph, buildErr := builder.Build(external)
	if buildErr != nil {
		return buildErr
	}

	if len(graph) == 0 {
		deps.NoDeps(cmd)
		return nil
	}

	switch format {
	case fmt.FormatMermaid:
		deps.Mermaid(cmd, core.RenderMermaid(graph))
	case fmt.FormatTable:
		deps.Table(cmd, core.RenderTable(graph))
	default:
		deps.JSON(cmd, core.RenderJSON(graph))
	}

	return nil
}
