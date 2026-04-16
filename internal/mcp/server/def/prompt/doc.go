//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package prompt defines MCP prompt definitions and
// their argument schemas.
//
// The Defs variable holds all available MCP prompts
// that the server advertises through the prompts/list
// method. Each prompt has a name, description, and
// optional list of required arguments.
//
// # Available Prompts
//
// The following prompts are defined:
//
//   - ctx-session-start: loads full project context
//     at the beginning of a session. Takes no args.
//   - ctx-decision-add: formats an architectural
//     decision entry. Requires content, context,
//     rationale, and consequence arguments.
//   - ctx-learning-add: formats a learning entry.
//     Requires content, context, lesson, and
//     application arguments.
//   - ctx-reflect: guides end-of-session reflection.
//     Takes no arguments.
//   - ctx-checkpoint: reports session statistics.
//     Takes no arguments.
//
// # Argument Schema
//
// Prompts that accept arguments use PromptArgument
// structs with name, description, and required flag.
// The argument names match the CLI attribute names
// (e.g., "context", "rationale") so that values
// can flow between MCP clients and ctx commands.
//
// # Usage
//
// The Defs slice is consumed by the prompt list
// dispatcher, which returns it as-is in response to
// prompts/list requests.
//
//	for _, p := range prompt.Defs {
//	    fmt.Println(p.Name, p.Description)
//	}
package prompt
