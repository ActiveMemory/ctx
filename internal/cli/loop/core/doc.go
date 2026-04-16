//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides loop script generation logic.
//
// The "ctx loop" command creates a bash script that
// runs an AI tool repeatedly with the same prompt file
// until a completion signal appears in the output.
// This core package holds the script generation
// business logic.
//
// # Script Generation
//
// The script sub-package exports [script.Generate],
// which builds a complete bash script string. The
// function accepts four parameters: the prompt file
// path, the AI tool name, a maximum iteration count,
// and a completion message string.
//
// The tool parameter selects the AI command template:
//
//   - "claude": runs Claude Code with the prompt
//     file via the LoopCmdClaude template.
//   - "aider": runs Aider with the prompt file
//     via the LoopCmdAider template.
//   - "generic": runs a generic command via the
//     LoopCmdGeneric template.
//
// When maxIterations is greater than zero, the script
// includes an iteration-limit guard that stops after
// the specified number of runs and sends a
// notification. The script monitors output for the
// completion message and exits cleanly when detected.
//
// # Data Flow
//
// The cmd/loop layer validates user inputs and calls
// script.Generate. The resulting script string is
// written to stdout or a file for the user to execute
// in their shell.
package core
