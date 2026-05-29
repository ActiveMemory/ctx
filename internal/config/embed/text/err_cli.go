//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for CLI errors.
const (
	// DescKeyErrCliNoToolSpecified is the text key for err cli no tool specified
	// messages.
	DescKeyErrCliNoToolSpecified = "err.cli.no-tool-specified"
	// DescKeyErrCliUnknownSubcommand is the text key for the terse
	// error returned when `ctx system` receives an unrecognised
	// subcommand (the rich guidance goes in the stdout relay box).
	DescKeyErrCliUnknownSubcommand = "err.cli.unknown-subcommand"
)
