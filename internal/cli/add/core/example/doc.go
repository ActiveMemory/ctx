//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package example provides type-specific usage examples
// for the add command help text.
//
// # ForType
//
// [ForType] returns a formatted example string for a given
// entry type such as "decision", "task", "learning", or
// "convention". The examples are loaded from the embedded
// commands.yaml asset via the desc package.
//
// The lookup key is formed by prefixing the entry type with
// the add-command example key prefix defined in
// config/embed/cmd. When the type is unrecognized or the
// key is missing, ForType falls back to a generic default
// example keyed by the default suffix.
//
// # Usage
//
// The cmd/ layer calls ForType during cobra command setup
// to populate the Example field of each add subcommand.
// This keeps example text centralized in YAML rather than
// scattered across Go source files.
package example
