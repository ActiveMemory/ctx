//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/flag"
	errCli "github.com/ActiveMemory/ctx/internal/err/cli"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Tool returns the active tool identifier from the --tool flag
// or the .ctxrc tool field.
//
// Resolution order:
//  1. --tool flag (if explicitly set on the command)
//  2. rc.Tool() (from .ctxrc)
//
// Returns an error if neither source provides a value.
//
// Parameters:
//   - cmd: The cobra command to read the --tool flag from
//
// Returns:
//   - string: The resolved tool identifier
//   - error: Non-nil when no tool is configured
func Tool(cmd *cobra.Command) (string, error) {
	if cmd.Flags().Changed(flag.Tool) {
		v, _ := cmd.Flags().GetString(flag.Tool)
		if v != "" {
			return v, nil
		}
	}

	if t := rc.Tool(); t != "" {
		return t, nil
	}

	return "", errCli.NoToolSpecified()
}
