//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package guide provides terminal output for the help/guide command
// that lists available skills and commands.
//
// Functions render two sections: skills ([InfoSkillsHeader],
// [InfoSkillLine]) and CLI commands ([CommandsHeader], [CommandLine]).
// [Default] outputs the combined guide when no subcommand is given.
//
// Example:
//
//	write.InfoSkillsHeader(cmd)
//	for _, s := range skills {
//	    write.InfoSkillLine(cmd, s.Name, s.Description)
//	}
package guide
