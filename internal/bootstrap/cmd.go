//	/    ctx:                         https://ctx.ist
//
// ,'`./    do you remember?
//
//	`.,'\
//	  \    Copyright 2026-present Context contributors.
//	                SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/flag"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/validate"
	"github.com/ActiveMemory/ctx/internal/err/fs"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// version is set at build time via ldflags:
//
//	-X github.com/ActiveMemory/ctx/internal/bootstrap.version=$(cat VERSION)
var version = "dev"

// RootCmd creates and returns the root cobra command for the ctx CLI.
//
// The root command provides the entry point for all ctx subcommands and
// displays help information when invoked without arguments.
//
// Global flags:
//   - --context-dir: Override the context directory path (default: .context)
//   - --allow-outside-cwd: Allow context directory outside project root
//
// Returns:
//   - *cobra.Command: The configured root command with usage and version info
func RootCmd() *cobra.Command {
	var contextDir string
	var allowOutsideCwd bool

	short, long := desc.Command(cmd.DescKeyCtx)

	c := &cobra.Command{
		Use:     cmd.DescKeyCtx,
		Short:   short,
		Long:    long,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Apply global flag values
			if contextDir != "" {
				rc.OverrideContextDir(contextDir)
			}
			// Validate that the context directory stays within the project root.
			// Skip if the CLI flag is set or .ctxrc has allow_outside_cwd: true.
			if !allowOutsideCwd && !rc.AllowOutsideCwd() {
				if validateErr := validation.ValidateBoundary(
					rc.ContextDir(),
				); validateErr != nil {
					return fs.BoundaryViolation(validateErr)
				}
			}

			// Skip init check for hidden commands (hooks have their own guards)
			// and cobra's built-in completion subcommands (bash, zsh, fish,
			// PowerShell) which must work in any directory.
			if cmd.Hidden {
				return nil
			}
			if p := cmd.Parent(); p != nil && p.Name() == cli.CmdCompletion {
				return nil
			}

			// Skip init check for annotated commands.
			if _, ok := cmd.Annotations[cli.AnnotationSkipInit]; ok {
				return nil
			}

			// Skip init check for grouping commands (no Run/RunE = just shows help).
			if cmd.RunE == nil && cmd.Run == nil {
				return nil
			}

			// Require initialization.
			if !ctxContext.Initialized(rc.ContextDir()) {
				return errInit.NotInitialized()
			}

			return nil
		},
	}

	// Cobra's c.Print() defaults to stderr (OutOrStderr). Set stdout
	// explicitly so all subcommands inherit the correct output, and shell
	// redirection (>) works as expected.
	c.SetOut(os.Stdout)

	// Global flags available to all subcommands
	c.PersistentFlags().StringVar(
		&contextDir,
		flag.ContextDir,
		"",
		desc.Flag(embedFlag.DescKeyContextDir),
	)
	c.PersistentFlags().BoolVar(
		&allowOutsideCwd,
		flag.AllowOutsideCwd,
		false,
		desc.Flag(embedFlag.DescKeyAllowOutsideCwd),
	)

	return c
}
