//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema provides the "ctx journal schema" parent
// command.
//
// # Overview
//
// This package groups schema-related subcommands under a
// single namespace. It does not contain business logic
// itself; it delegates to its two children:
//
//   - check: scans JSONL session files for format drift
//     and writes a Markdown report when violations are
//     found.
//   - dump: prints the full embedded schema definition
//     to stdout for human inspection.
//
// # Usage
//
//	ctx journal schema check [flags]
//	ctx journal schema dump
//
// # Behavior
//
// [Cmd] uses the parent.Cmd helper to build a
// cobra.Command with short and long descriptions loaded
// from embedded assets. It attaches the check and dump
// subcommands as children. Running the parent without a
// subcommand prints the help text.
//
// Both subcommands are designed for CI pipelines and
// nightly cron jobs as well as interactive use. The check
// subcommand returns exit code 1 on drift; the dump
// subcommand always succeeds.
package schema
