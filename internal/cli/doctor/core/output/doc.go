//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package output formats doctor reports for display in two
// modes: JSON and human-readable text.
//
// # JSON Mode
//
// The [JSON] function serializes a [check.Report] as
// indented JSON and writes it to the command's output
// stream. This is intended for machine consumption and
// piping into other tools. Indentation uses two-space
// indent for readability.
//
// # Human Mode
//
// The [Human] function renders the report as a categorized
// list with status indicators. Each check result is
// converted to a [writeDoctor.ResultItem] carrying its
// category, status, and message. Warnings and errors from
// the report are appended as summary counts at the bottom.
//
// Both functions accept a cobra command for output routing
// and return an error for interface consistency, though
// Human always returns nil.
package output
