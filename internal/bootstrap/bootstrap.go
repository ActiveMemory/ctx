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
	"github.com/ActiveMemory/ctx/internal/cli/pad"
	"github.com/ActiveMemory/ctx/internal/cli/permissions"
	"github.com/ActiveMemory/ctx/internal/cli/recall"
	"github.com/ActiveMemory/ctx/internal/cli/serve"
	"github.com/ActiveMemory/ctx/internal/cli/status"
	"github.com/ActiveMemory/ctx/internal/cli/sync"
	"github.com/ActiveMemory/ctx/internal/cli/system"
	"github.com/ActiveMemory/ctx/internal/cli/task"
	"github.com/ActiveMemory/ctx/internal/cli/watch"
)

// Initialize registers all ctx subcommands with the root command.
//
// This function attaches all available subcommands to the provided root
// command, including init, status, load, add, complete, agent, drift,
// sync, compact, decision, watch, hook, learnings, tasks, loop, recall,
// journal, and serve.
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
		task.Cmd,
		loop.Cmd,
		pad.Cmd,
		permissions.Cmd,
		recall.Cmd,
		journal.Cmd,
		serve.Cmd,
		system.Cmd,
	} {
		cmd.AddCommand(c())
	}

	return cmd
}
