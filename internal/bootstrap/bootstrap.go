//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package bootstrap initializes the ctx CLI application.
//
// It provides functions to create the root command and register all
// subcommands. The typical usage pattern is:
//
//	cmd := bootstrap.Initialize(bootstrap.RootCmd())
//	if err := cmd.Execute(); err != nil {
//	    // handle error
//	}
package bootstrap

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/agent"
	"github.com/ActiveMemory/ctx/internal/cli/compact"
	"github.com/ActiveMemory/ctx/internal/cli/complete"
	"github.com/ActiveMemory/ctx/internal/cli/decision"
	"github.com/ActiveMemory/ctx/internal/cli/drift"
	"github.com/ActiveMemory/ctx/internal/cli/hook"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/cli/journal"
	"github.com/ActiveMemory/ctx/internal/cli/learnings"
	"github.com/ActiveMemory/ctx/internal/cli/load"
	"github.com/ActiveMemory/ctx/internal/cli/loop"
	"github.com/ActiveMemory/ctx/internal/cli/recall"
	"github.com/ActiveMemory/ctx/internal/cli/serve"
	"github.com/ActiveMemory/ctx/internal/cli/session"
	"github.com/ActiveMemory/ctx/internal/cli/status"
	"github.com/ActiveMemory/ctx/internal/cli/sync"
	"github.com/ActiveMemory/ctx/internal/cli/task"
	"github.com/ActiveMemory/ctx/internal/cli/watch"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// version is set at build time via ldflags
const version = "dev"

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
	var contextDir string
	var noColor bool

	cmd := &cobra.Command{
		Use:   "ctx",
		Short: "Context - persistent context for AI coding assistants",
		Long: `Context (ctx) maintains persistent context files that help
  AI coding assistants understand your project's architecture, conventions,
  decisions, and current tasks.

  Use 'ctx init' to create a .context/ directory in your project,
  then use 'ctx status', 'ctx load', and 'ctx agent' to work with context.`,
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Apply global flag values
			if contextDir != "" {
				rc.OverrideContextDir(contextDir)
			}
			if noColor {
				color.NoColor = true
			}
		},
	}

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

	return cmd
}

// Initialize registers all ctx subcommands with the root command.
//
// This function attaches all available subcommands (init, status, load, add,
// complete, agent, drift, sync, compact, watch, hook, session, tasks, loop)
// to the provided root command.
//
// Parameters:
//   - cmd: The root cobra command to attach subcommands to
//
// Returns:
//   - *cobra.Command: The same command with all subcommands registered
func Initialize(cmd *cobra.Command) *cobra.Command {
	for _, c := range []func() *cobra.Command{
		initialize.Cmd,
		status.Cmd,
		load.Cmd,
		add.Cmd,
		complete.Cmd,
		agent.Cmd,
		drift.Cmd,
		sync.Cmd,
		compact.Cmd,
		decision.Cmd,
		watch.Cmd,
		hook.Cmd,
		learnings.Cmd,
		session.Cmd,
		task.Cmd,
		loop.Cmd,
		recall.Cmd,
		journal.Cmd,
		serve.Cmd,
	} {
		cmd.AddCommand(c())
	}

	return cmd
}
