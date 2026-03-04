//	/    ctx:                         https://ctx.ist
//
// ,'`./    do you remember?
//
//	`.,'\
//	  \    Copyright 2026-present Context contributors.
//	                SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// version is set at build time via ldflags:
//
//	-X github.com/ActiveMemory/ctx/internal/bootstrap.version=0.2.0
var version = "dev"

// RootCmd creates and returns the root cobra command for the ctx CLI.
//
// The root command provides the entry point for all ctx subcommands and
// displays help information when invoked without arguments.
//
// Global flags:
//   - --context-dir: Override the context directory path (default: .context)
//   - --no-color: Disable colored output
//
// Returns:
//   - *cobra.Command: The configured root command with usage and version info
func RootCmd() *cobra.Command {
	config.BinaryVersion = version

	var contextDir string
	var noColor bool
	var allowOutsideCwd bool

	cmd := &cobra.Command{
		Use:   "ctx",
		Short: "Context - persistent context for AI coding assistants",
		Long: `ctx (Context) maintains persistent context files that help
  AI coding assistants understand your project's architecture, conventions,
  decisions, and current tasks.

  Use 'ctx init' to create a .context/ directory in your project,
  then use 'ctx status', 'ctx load', and 'ctx agent' to work with context.`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Apply global flag values
			if contextDir != "" {
				rc.OverrideContextDir(contextDir)
			}
			if noColor {
				color.NoColor = true
			}

			// Validate that the context directory stays within the project root.
			// Skip if CLI flag is set or .ctxrc has allow_outside_cwd: true.
			if !allowOutsideCwd && !rc.AllowOutsideCwd() {
				if validateErr := validation.ValidateBoundary(rc.ContextDir()); validateErr != nil {
					return fmt.Errorf("%w\nUse --allow-outside-cwd to override this check", validateErr)
				}
			}

			// Skip init check for hidden commands (hooks have their own guards)
			// and cobra's built-in completion subcommands (bash, zsh, fish,
			// powershell) which must work in any directory.
			if cmd.Hidden {
				return nil
			}
			if p := cmd.Parent(); p != nil && p.Name() == "completion" {
				return nil
			}

			// Skip init check for annotated commands.
			if _, ok := cmd.Annotations[config.AnnotationSkipInit]; ok {
				return nil
			}

			// Skip init check for grouping commands (no Run/RunE = just shows help).
			if cmd.RunE == nil && cmd.Run == nil {
				return nil
			}

			// Require initialization.
			if !config.Initialized(rc.ContextDir()) {
				return fmt.Errorf("ctx: not initialized — run \"ctx init\" first")
			}

			return nil
		},
	}

	// Cobra's cmd.Print() defaults to stderr (OutOrStderr). Set stdout
	// explicitly so all subcommands inherit correct output and shell
	// redirection (>) works as expected.
	cmd.SetOut(os.Stdout)

	// Global flags available to all subcommands
	cmd.PersistentFlags().StringVar(
		&contextDir,
		"context-dir",
		"",
		"Override context directory path (default: .context)",
	)
	cmd.PersistentFlags().BoolVar(
		&noColor,
		"no-color",
		false,
		"Disable colored output",
	)
	cmd.PersistentFlags().BoolVar(
		&allowOutsideCwd,
		"allow-outside-cwd",
		false,
		"Allow context directory outside current working directory",
	)

	return cmd
}
