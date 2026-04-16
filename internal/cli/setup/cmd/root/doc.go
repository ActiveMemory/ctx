//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx setup" command
// for generating AI tool integration configurations.
//
// # Behavior
//
// The command accepts exactly one positional argument:
// the name of an AI coding tool (e.g., claude,
// cursor, copilot, kiro, cline, aider, windsurf,
// copilot-cli, agents). It outputs configuration
// snippets and setup instructions specific to that
// tool.
//
// Without --write, the command prints the integration
// instructions and configuration content to stdout
// for manual review. With --write, it generates and
// writes the configuration file directly to the
// expected location for the target tool.
//
// Unsupported tool names produce an error listing
// the available options.
//
// This command skips context initialization (via the
// SkipInit annotation), so it can run before ctx is
// fully set up in a project.
//
// # Flags
//
//	--write, -w    Write the configuration file
//	               directly instead of printing to
//	               stdout.
//
// # Output
//
// In print mode, outputs tool-specific integration
// instructions followed by the configuration content.
// In write mode, writes the file and prints a
// confirmation. For unsupported tools, prints a
// notice and returns an error.
//
// # Delegation
//
// Each supported tool has a dedicated core package
// (e.g., core/cursor, core/copilot) that handles
// deployment logic. Output formatting is routed
// through the [writeSetup] package.
package root
