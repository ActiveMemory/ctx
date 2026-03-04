//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

func schemaCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schema",
		Short: "Print JSON Schema for .ctxrc",
		Long: `Print the JSON Schema for .ctxrc to stdout.

Pipe-friendly — redirect to a file for IDE integration:

  ctx config schema > .ctxrc.schema.json

VS Code integration (requires redhat.vscode-yaml extension):

  // .vscode/settings.json
  {
    "yaml.schemas": {
      "./.ctxrc.schema.json": ".ctxrc"
    }
  }`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			data, readErr := assets.Schema()
			if readErr != nil {
				return fmt.Errorf("read embedded schema: %w", readErr)
			}
			cmd.Print(string(data))
			return nil
		},
	}
}
