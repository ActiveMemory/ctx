//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package loop defines configuration constants for the
// loop script generator, which produces shell scripts
// that drive AI tools in autonomous iteration cycles.
//
// # What Loop Scripts Do
//
// A loop script invokes an AI tool repeatedly until
// the tool emits a completion signal (by default
// DefaultCompletionSignal, "SYSTEM_CONVERGED"). Each
// iteration feeds a prompt from the prompt file
// (PromptMd, ".context/loop.md") to the selected
// tool, checks the output for convergence, and either
// loops again or exits.
//
// # Supported Tools
//
// Three tool backends are supported, tracked in the
// ValidTools map:
//
//   - DefaultTool ("claude") -- Claude Code CLI.
//   - ToolAider ("aider") -- the Aider coding agent.
//   - ToolGeneric ("generic") -- any tool that reads
//     stdin and writes stdout.
//
// The generated script is written to DefaultOutput
// ("loop.sh") by default.
//
// # Key Constants
//
//   - DefaultCompletionSignal -- string the loop
//     watches for to stop iterating.
//   - DefaultOutput -- default output filename.
//   - PromptMd -- path to the iteration prompt file.
//   - ValidTools -- set of recognized tool IDs.
package loop
