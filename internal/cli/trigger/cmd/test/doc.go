//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package test implements the "ctx trigger test" cobra
// subcommand.
//
// This command executes all enabled trigger scripts
// for a given hook type using mock input, allowing
// developers to verify their hooks work correctly
// without waiting for a real lifecycle event.
//
// # Usage
//
//	ctx trigger test <hook-type> [--tool <name>] [--path <file>]
//
// # Arguments
//
// Exactly one positional argument is required:
//
//   - hook-type: the trigger type to test (e.g.
//     "pre-tool-use"). Must be a valid type.
//
// # Flags
//
//	--tool   Optional tool name to include in the
//	         mock input JSON.
//	--path   Optional file path to include in the
//	         mock input parameters.
//
// # Behavior
//
// The command:
//
//   - Validates the hook type against valid trigger
//     types.
//   - Constructs a mock HookInput with the specified
//     tool, path, a mock session ID and model, the
//     current timestamp, and a mock version string.
//   - Prints the mock input JSON for inspection.
//   - Calls trigger.RunAll to execute every enabled
//     script for the hook type with a configured
//     timeout.
//   - Displays the aggregated result: context output
//     if any was produced, error lines if any scripts
//     failed, a cancellation message if a script
//     cancelled the operation, or a "no output"
//     message if scripts ran silently.
//
// # Output
//
// A test header, the mock input JSON, and the
// aggregated output from all executed scripts.
//
// # Delegation
//
// Hook execution uses trigger.RunAll. Type validation
// uses trigger.ValidTypes. Output formatting uses
// write/trigger.
package test
