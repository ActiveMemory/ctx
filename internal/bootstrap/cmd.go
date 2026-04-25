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
	cfgBootstrap "github.com/ActiveMemory/ctx/internal/config/bootstrap"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/flag"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/validate"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeBootstrap "github.com/ActiveMemory/ctx/internal/write/bootstrap"
)

// version is set at build time via ldflags:
//
//	-X github.com/ActiveMemory/ctx/internal/bootstrap.version=$(cat VERSION)
var version = cfgBootstrap.DefaultVersion

// RootCmd creates and returns the root cobra command for the ctx CLI.
//
// The root command provides the entry point for all ctx subcommands and
// displays help information when invoked without arguments.
//
// Returns:
//   - *cobra.Command: The configured root command with usage and version info
func RootCmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyCtx)

	c := &cobra.Command{
		Use:     cmd.DescKeyCtx,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyCtx),
		Version: version,
		// Cobra auto-prints returned errors to stderr by default;
		// main.go also prints them via writeErr.With, producing a
		// double-printed error. Silence cobra's path so writeErr is
		// the sole printer. (SilenceUsage stays per-return so
		// genuine cobra parse errors keep their help dump.)
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip every downstream check for administrative commands
			// that must run without a declared or initialized context:
			//   - Hidden commands (e.g. ctx system bootstrap; hooks
			//     supply their own guards).
			//   - Cobra's built-in shell-completion subcommands.
			//   - Commands annotated with AnnotationSkipInit (init,
			//     activate, deactivate, guide, why, doctor, config
			//     switch/status, hub *).
			//   - Grouping commands without a Run / RunE of their own
			//     (they just print help for their subtree).
			if cmd.Hidden {
				return nil
			}
			if p := cmd.Parent(); p != nil && p.Name() == cli.CmdCompletion {
				return nil
			}
			if _, ok := cmd.Annotations[cli.AnnotationSkipInit]; ok {
				return nil
			}
			if cmd.RunE == nil && cmd.Run == nil {
				return nil
			}

			// Under the single-source-anchor model, every non-exempt
			// command requires CTX_DIR to be declared and to point at
			// an existing .context/ directory. RequireContextDir
			// returns a tailored error (with a next-step hint based on
			// how many .context/ candidates are visible from CWD) when
			// the declaration is missing or broken. The parent of the
			// declared directory is the project root by contract; CWD
			// has no say in project identity.
			ctxDir, reqErr := rc.RequireContextDir()
			if reqErr != nil {
				// Actionable error, not a usage problem. Suppress
				// cobra's help dump so the call-to-action stays
				// the only thing on stderr. Genuine cobra errors
				// (unknown subcommand, bad flag) still print usage
				// because they happen before PreRunE runs.
				cmd.SilenceUsage = true
				return reqErr
			}

			// Require initialization: the declared directory must
			// have been initialized before other commands operate.
			if !ctxContext.Initialized(ctxDir) {
				cmd.SilenceUsage = true
				return errInit.NotInitialized()
			}

			return nil
		},
	}

	// Cobra's c.Print() defaults to stderr (OutOrStderr). Set stdout
	// explicitly so all subcommands inherit the correct output, and shell
	// redirection (>) works as expected.
	c.SetOut(os.Stdout)

	// Append a community footer to the root help output only.
	defaultHelp := c.HelpFunc()
	c.SetHelpFunc(func(helpCmd *cobra.Command, args []string) {
		defaultHelp(helpCmd, args)
		if helpCmd == c {
			writeBootstrap.CommunityFooter(helpCmd)
		}
	})

	c.PersistentFlags().String(
		flag.Tool,
		"",
		desc.Flag(embedFlag.DescKeyTool),
	)

	return c
}
