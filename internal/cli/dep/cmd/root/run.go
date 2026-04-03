//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/cli/dep/core/builder"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/dep/core/render"
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

	var b builder.GraphBuilder
	if projType != "" {
		b = builder.Find(projType)
		if b == nil {
			names := strings.Join(builder.Names(), token.CommaSpace)
			return errConfig.UnknownProjectType(projType, names)
		}
	} else {
		b = builder.Detect()
		if b == nil {
			names := strings.Join(
				builder.Names(), token.CommaSpace,
			)
			deps.InfoNoProject(cmd, names)
			return nil
		}
	}

	graph, buildErr := b.Build(external)
	if buildErr != nil {
		return buildErr
	}

	if len(graph) == 0 {
		deps.None(cmd)
		return nil
	}

	switch format {
	case fmt.FormatMermaid:
		deps.Mermaid(cmd, render.Mermaid(graph))
	case fmt.FormatTable:
		deps.Table(cmd, render.Table(graph))
	default:
		deps.JSON(cmd, render.JSON(graph))
	}

	return nil
}
