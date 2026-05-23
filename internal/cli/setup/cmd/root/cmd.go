//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreBackend "github.com/ActiveMemory/ctx/internal/cli/setup/core/backend"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	errSetup "github.com/ActiveMemory/ctx/internal/err/setup"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx setup" command for generating AI tool integrations.
//
// The command outputs configuration snippets and instructions for integrating
// Context with various AI coding tools like Claude Code, Cursor, Aider, etc.
//
// Flags:
//   - --write, -w: Write the configuration file instead of printing
//
// Returns:
//   - *cobra.Command: Configured setup command that
//     accepts a tool name argument
func Cmd() *cobra.Command {
	var (
		write     bool
		backend   string
		endpoint  string
		apiKeyEnv string
	)

	short, long := desc.Command(cmd.DescKeySetup)
	c := &cobra.Command{
		Use:         cmd.UseSetup,
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		Long:        long,
		Example:     desc.Example(cmd.DescKeySetup),
		Args:        cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if backend != "" {
				if len(args) > 0 {
					return errSetup.ErrBackendAndToolConflict
				}
				return coreBackend.Setup(
					cmd, backend, endpoint, apiKeyEnv,
				)
			}
			if len(args) == 0 {
				return errSetup.ErrMissingToolOrBackend
			}
			return Run(cmd, args, write)
		},
	}

	flagbind.BoolFlagP(c, &write,
		cFlag.Write, cFlag.ShortWrite,
		flag.DescKeySetupWrite,
	)
	flagbind.StringFlag(c, &backend,
		cFlag.Backend, flag.DescKeySetupBackend,
	)
	flagbind.StringFlag(c, &endpoint,
		cFlag.Endpoint, flag.DescKeySetupEndpoint,
	)
	flagbind.StringFlag(c, &apiKeyEnv,
		cFlag.APIKeyEnv, flag.DescKeySetupAPIKeyEnv,
	)

	return c
}
