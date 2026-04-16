//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package preview implements the **`ctx steering preview`**
// subcommand, which shows which steering files would be
// included for a given prompt.
//
// # What It Does
//
// Given a single positional argument (the prompt text),
// the command loads all steering files, applies the
// inclusion filter, and displays the subset that would
// be injected into the agent context. The filter
// includes files with inclusion mode "always" and those
// whose "auto" keywords match the prompt. Files with
// inclusion mode "manual" are excluded from the preview
// since they require explicit naming.
//
// For each matching file, the output shows:
//
//   - name
//   - inclusion mode
//   - priority
//   - target tools (or "all tools")
//
// # Arguments
//
//   - PROMPT (required): the prompt text to match
//     against auto-inclusion keywords.
//
// # Output
//
// A header line quoting the prompt, followed by one
// line per matching file and a count footer. If no
// files match, a "no files match" message is shown.
//
// # Delegation
//
// [Cmd] builds the cobra command and enforces
// ExactArgs(1). [Run] calls [steering.LoadAll] and
// [steering.Filter] to select matching files, then
// emits results through the [write/steering] preview
// formatters.
package preview
