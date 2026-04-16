//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package remove implements the "ctx skill remove"
// subcommand for uninstalling skills from the
// project.
//
// # Behavior
//
// The command accepts exactly one positional argument:
// the name of the skill to remove. It locates the
// skill's directory under .context/skills/ and
// deletes the entire directory along with all
// contained files.
//
// If the named skill does not exist, the command
// returns an error. Removal is immediate and
// permanent with no confirmation prompt.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// On success, prints a confirmation showing the
// removed skill name. On failure, returns an error
// identifying the missing skill or any I/O problem
// encountered during deletion.
//
// # Delegation
//
// Skill removal logic is handled by [skill.Remove],
// which validates existence and deletes the
// directory tree. The skills directory path is
// constructed from [rc.ContextDir] and [dir.Skills].
// Output is routed through [writeSkill.Removed].
package remove
