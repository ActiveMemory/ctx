//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package list implements the "ctx skill list"
// subcommand for displaying all installed skills.
//
// # Behavior
//
// The command scans the project's skills directory
// under .context/skills/, loads the SKILL.md manifest
// from each subdirectory, and displays a summary of
// all installed skills.
//
// Each skill is printed with its name and description
// (when available). Skills without a description show
// only the name. After the list, a total count line
// is printed.
//
// When no skills are installed, the command prints a
// notice and exits.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// Prints one line per skill showing name and
// description, followed by a count summary. When
// the skills directory is empty or missing, prints
// a "no skills found" notice.
//
// # Delegation
//
// Skill loading is handled by [skill.LoadAll], which
// reads and parses all SKILL.md manifests. The skills
// directory path is constructed from [rc.ContextDir]
// and [dir.Skills]. Output is routed through the
// [writeSkill] package.
package list
