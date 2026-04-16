//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx loop" command.
//
// # Overview
//
// The loop command generates a shell script that runs an
// AI assistant in a repeated loop until a completion
// signal is detected. This enables iterative development
// where the AI builds on its previous work across
// multiple invocations.
//
// The generated script is written to a file (default
// "loop.sh") with executable permissions and includes
// usage instructions printed to stdout.
//
// # Flags
//
//	-p, --prompt <file>      Prompt file for the AI
//	                          (default ".context/loop.md").
//	-t, --tool <name>        AI tool: claude, aider,
//	                          or generic (default "claude").
//	-n, --max-iterations <n> Maximum loop iterations;
//	                          0 means unlimited (default 0).
//	-c, --completion <sig>   Completion signal string
//	                          (default "SYSTEM_CONVERGED").
//	-o, --output <file>      Output script filename
//	                          (default "loop.sh").
//
// # Behavior
//
// [Cmd] builds the cobra.Command and registers all five
// flags with their defaults. [Run] validates the tool
// selection against a known-good set, generates the
// script via core/script.Generate, writes it with
// executable permissions, and prints usage instructions.
//
// If the tool name is not recognized, the command returns
// an "invalid tool" error listing valid options.
//
// # Output
//
// Prints the output filename, the tool being used, the
// prompt file path, iteration limit, and completion
// signal. The generated script itself is written to
// the output file, not to stdout.
package root
