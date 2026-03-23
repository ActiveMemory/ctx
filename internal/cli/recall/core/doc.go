//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides shared operations used by all recall subcommand
// packages. It owns formatting, querying, validation, frontmatter handling,
// index building, slug generation, and export planning; subcommands import
// core, never each other.
package core
