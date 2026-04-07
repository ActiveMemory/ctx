//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema provides the journal schema parent command.
//
// It groups the check and dump subcommands under
// "ctx journal schema". The check subcommand scans JSONL
// session files for format drift and writes a report. The
// dump subcommand prints the embedded schema definition for
// inspection. Both are designed for use in CI pipelines and
// nightly cron jobs as well as interactive use.
package schema
