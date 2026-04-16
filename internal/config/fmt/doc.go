//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package fmt defines output format identifiers used
// by CLI commands that support multiple rendering
// modes.
//
// Several ctx commands (status, agent, doctor, drift,
// and others) accept a --format flag that selects how
// output is rendered. This package provides the
// canonical string constants for those format values
// so that flag registration and output branching use
// the same identifiers.
//
// # Format Constants
//
//   - FormatJSON ("json") -- selects JSON output,
//     useful for piping into jq or programmatic
//     consumption
//   - FormatMarkdown ("md") -- selects Markdown
//     output, the default human-readable format for
//     most commands
//
// # Usage Pattern
//
// A typical command registers the format flag using
// config/flag.Format and compares the user's value
// against these constants:
//
//	switch outputFmt {
//	case fmt.FormatJSON:
//	    renderJSON(result)
//	default:
//	    renderMarkdown(result)
//	}
//
// # Why Centralized
//
// Format identifiers are referenced by flag
// registration code, output rendering branches, and
// shell completion generators. Defining them here
// prevents typos (e.g. "Json" vs "json") and ensures
// that adding a new format is a single-point change
// with compile-time verification across all commands.
package fmt
