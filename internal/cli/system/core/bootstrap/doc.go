//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package bootstrap provides helpers for the system
// bootstrap command, which initializes agent sessions
// by reporting context directory state.
//
// # Plugin Warning
//
// [PluginWarning] checks whether the ctx plugin is
// installed but not enabled in either global or local
// settings. When the plugin is installed but inactive,
// it returns a warning message. This helps agents
// detect misconfigured environments at session start.
//
// # Context File Listing
//
// [ListContextFiles] reads the context directory and
// returns a sorted list of markdown filenames. This
// gives agents an inventory of available context files
// without requiring them to scan the directory
// themselves.
//
// [WrapFileList] formats file names as a
// comma-separated list with line wrapping at a
// configured width. Continuation lines are indented.
// Returns a "none" label when the list is empty.
//
// # Numbered Line Parsing
//
// [ParseNumberedLines] splits a numbered multiline
// string into individual items, stripping the leading
// "N. " prefix from each line. This is used to parse
// structured output from agent responses.
//
// # Data Flow
//
// The system bootstrap command calls these helpers to
// build its output. PluginWarning checks installation
// state. ListContextFiles reads the directory. Both
// results are formatted through the write/ layer for
// agent consumption.
package bootstrap
