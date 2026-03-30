//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package importer implements the ctx journal import subcommand.
//
// [Cmd] builds the cobra.Command with --all, --regenerate,
// --dry-run, and --keep-frontmatter flags. [Run] plans the import
// (which sessions to create, regenerate, or skip), confirms with
// the user, and executes the plan.
package importer
