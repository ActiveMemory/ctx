//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parent provides a helper for creating pure
// grouping commands that have no RunE, no flags; only
// desc-loaded descriptions and subcommand registration.
//
// Many ctx commands are namespace parents that exist
// solely to group related subcommands (e.g. config,
// hub, journal, steering). Rather than duplicating the
// same cobra.Command boilerplate in each package, those
// packages call [Cmd] with a description key, a use
// string, and a list of subcommands.
//
// # How It Works
//
// [Cmd] loads Short and Long descriptions from embedded
// assets via [internal/assets/read/desc.Command], creates
// a cobra.Command with no RunE, and registers all
// provided subcommands via AddCommand. The result is a
// ready-to-use parent command.
//
// # Example
//
//	parent.Cmd(cmd.DescKeyConfig, cmd.UseConfig,
//	    schema.Cmd(),
//	    status.Cmd(),
//	    switchcmd.Cmd(),
//	)
//
// [Cmd] builds a cobra.Command that has no RunE of its
// own, loads its Short and Long text from embedded desc
// assets, and registers every supplied subcommand via
// AddCommand.
package parent
