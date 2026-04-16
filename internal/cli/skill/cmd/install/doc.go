//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package install implements the "ctx skill install"
// subcommand for adding skills to the project.
//
// # Behavior
//
// The command accepts exactly one positional argument:
// the path to a source directory containing a
// SKILL.md manifest with YAML frontmatter. It parses
// the manifest, validates the skill metadata, and
// copies all files from the source directory into the
// project's skills directory under .context/skills/.
//
// The skill name is derived from the SKILL.md
// frontmatter. If a skill with the same name already
// exists, the install overwrites it. The installed
// skill becomes immediately available for use in
// agent sessions.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// On success, prints a confirmation showing the
// installed skill name and the destination directory
// path. On failure, returns an error for missing
// manifests, invalid frontmatter, or I/O problems.
//
// # Delegation
//
// Skill installation logic is handled by
// [skill.Install], which parses the manifest and
// copies files. The skills directory path is
// constructed from [rc.ContextDir] and [dir.Skills].
// Output is routed through [writeSkill.Installed].
package install
