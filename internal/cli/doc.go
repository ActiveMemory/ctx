//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cli contains the implementation of all ctx
// subcommands.
//
// Each command lives in its own subpackage following a
// consistent taxonomy:
//
//   - parent.go: exports a Cmd() function that wires the
//     cobra command and registers subcommands
//   - cmd/root/ or cmd/<sub>/: cobra command definitions,
//     flag binding, and RunE entry points
//   - core/: shared helpers, business logic, and
//     formatting used by one or more cmd/ packages
//
// The [internal/bootstrap] package registers all command
// packages into the root cobra.Command tree at startup.
// Commands that are pure namespace groupings (no RunE)
// use [internal/cli/parent.Cmd] to create the parent
// with desc-loaded descriptions and subcommand wiring.
//
// # Package Categories
//
// Context file commands: add, compact, decision, drift,
// fmt, load, agent, reindex, status, watch.
//
// Session lifecycle: pause, resume, event, system.
//
// Publishing and export: journal, serve, site, memory.
//
// Infrastructure: backup, config, connection, hub, mcp,
// prune, setup, steering, trigger, usage, sysinfo.
//
// Utilities: guide, loop, pad, resolve, skill, trace,
// why, parent, hook, message, notify.
package cli
