//   /    ctx:                         https://ctx.ist
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
	"github.com/ActiveMemory/ctx/internal/cli/changes"
	"github.com/ActiveMemory/ctx/internal/cli/compact"
	"github.com/ActiveMemory/ctx/internal/cli/config"
	"github.com/ActiveMemory/ctx/internal/cli/decision"
	"github.com/ActiveMemory/ctx/internal/cli/deps"
	"github.com/ActiveMemory/ctx/internal/cli/doctor"
	"github.com/ActiveMemory/ctx/internal/cli/drift"
	"github.com/ActiveMemory/ctx/internal/cli/guide"
	"github.com/ActiveMemory/ctx/internal/cli/hook"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/cli/journal"
	"github.com/ActiveMemory/ctx/internal/cli/learnings"
	"github.com/ActiveMemory/ctx/internal/cli/load"
	"github.com/ActiveMemory/ctx/internal/cli/loop"
	"github.com/ActiveMemory/ctx/internal/cli/mcp"
	"github.com/ActiveMemory/ctx/internal/cli/memory"
	"github.com/ActiveMemory/ctx/internal/cli/notify"
	"github.com/ActiveMemory/ctx/internal/cli/pad"
	"github.com/ActiveMemory/ctx/internal/cli/pause"
	"github.com/ActiveMemory/ctx/internal/cli/permissions"
	"github.com/ActiveMemory/ctx/internal/cli/prompt"
	"github.com/ActiveMemory/ctx/internal/cli/recall"
	"github.com/ActiveMemory/ctx/internal/cli/reindex"
	"github.com/ActiveMemory/ctx/internal/cli/remind"
	"github.com/ActiveMemory/ctx/internal/cli/resume"
	"github.com/ActiveMemory/ctx/internal/cli/serve"
	"github.com/ActiveMemory/ctx/internal/cli/site"
	"github.com/ActiveMemory/ctx/internal/cli/status"
	"github.com/ActiveMemory/ctx/internal/cli/sync"
	"github.com/ActiveMemory/ctx/internal/cli/system"
	"github.com/ActiveMemory/ctx/internal/cli/task"
	"github.com/ActiveMemory/ctx/internal/cli/watch"
	"github.com/ActiveMemory/ctx/internal/cli/why"
)

// Initialize registers all ctx subcommands with the root command.
//
// This function attaches all available subcommands to the provided root
// command, including init, status, load, add, agent, drift,
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
		add.Cmd,
		agent.Cmd,
		changes.Cmd,
		compact.Cmd,
		config.Cmd,
		decision.Cmd,
		deps.Cmd,
		doctor.Cmd,
		drift.Cmd,
		guide.Cmd,
		hook.Cmd,
		initialize.Cmd,
		journal.Cmd,
		learnings.Cmd,
		load.Cmd,
		loop.Cmd,
		mcp.Cmd,
		memory.Cmd,
		notify.Cmd,
		pad.Cmd,
		pause.Cmd,
		permissions.Cmd,
		prompt.Cmd,
		recall.Cmd,
		reindex.Cmd,
		remind.Cmd,
		resume.Cmd,
		serve.Cmd,
		site.Cmd,
		status.Cmd,
		sync.Cmd,
		system.Cmd,
		task.Cmd,
		watch.Cmd,
		why.Cmd,
	} {
		cmd.AddCommand(c())
	}

	return cmd
}
