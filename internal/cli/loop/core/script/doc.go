//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package script generates bash scripts for running AI tool
// iteration loops. The generated script repeatedly invokes
// an AI tool with a prompt file until a completion signal is
// detected in the output or a maximum iteration count is
// reached.
//
// # Supported Tools
//
// [Generate] accepts a tool identifier that selects the
// invocation command template:
//
//   - "claude" (default): runs Claude Code in headless
//     mode with the prompt file piped as input.
//   - "aider": runs the Aider CLI with the prompt file
//     as the message argument.
//   - "generic": runs a shell command that reads the
//     prompt file, suitable for custom tool wrappers.
//
// # Iteration Control
//
// The maxIterations parameter caps the loop. When set to
// zero, the loop runs indefinitely until the completion
// signal appears. Each iteration checks the tool's output
// for the completion message string; a match exits the
// loop cleanly with a desktop notification.
//
// # Path Handling
//
// The prompt file path is resolved to an absolute path
// via filepath.Abs before being embedded in the script,
// so the generated script works regardless of the working
// directory at execution time.
package script
