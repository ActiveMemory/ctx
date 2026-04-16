//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package guide provides terminal output for the help
// and guide commands (ctx guide, ctx help).
//
// # Skills Section
//
// [InfoSkillsHeader] prints the skills list heading
// followed by a blank line. [InfoSkillLine] prints a
// single skill entry with its name and a truncated
// description.
//
// # Commands Section
//
// [CommandsHeader] prints the CLI commands list heading
// followed by a blank line. [CommandLine] prints a
// single command entry with its name and short
// description.
//
// # Default Guide
//
// [Default] outputs the combined default guide text
// when no subcommand is specified. The content is
// loaded from the embedded descriptor system and
// printed verbatim.
//
// # Message Categories
//
//   - Info: skill and command listings, default guide
//
// # Nil Safety
//
// All functions treat a nil *cobra.Command as a no-op.
//
// # Usage
//
//	guide.InfoSkillsHeader(cmd)
//	for _, s := range skills {
//	    guide.InfoSkillLine(cmd, s.Name, s.Desc)
//	}
package guide
