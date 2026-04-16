//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core is the umbrella for the guide command's
// listing and discovery logic.
//
// # Overview
//
// The guide command helps users discover available ctx
// commands and skills. This package groups sub-packages
// that enumerate CLI commands and parse skill metadata
// for display.
//
// # Sub-packages
//
//   - command: lists all non-hidden cobra subcommands
//     from the root command tree. Exports [command.List].
//   - skill: lists available SKILL.md files, parses
//     their YAML frontmatter, and truncates descriptions
//     for display. Exports [skill.List],
//     [skill.ParseFrontmatter], [skill.TruncateDescription],
//     and the [skill.Meta] type.
//
// # Data Flow
//
// The cmd layer calls into the appropriate sub-package
// based on the subcommand:
//
//  1. ctx guide commands: calls command.List, which
//     walks the cobra command tree and prints each
//     non-hidden command with its short description.
//  2. ctx guide skills: calls skill.List, which
//     reads SKILL.md files from the claude skill
//     directory, extracts frontmatter metadata, and
//     prints each skill name with a truncated
//     description.
package core
