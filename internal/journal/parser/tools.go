//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

// find returns a parser for the specified tool.
//
// Parameters:
//   - tool: Tool identifier (e.g., "claude-code")
//
// Returns:
//   - Session: The parser, or nil if not found
func find(tool string) Session {
	for _, p := range registeredParsers {
		if p.Tool() == tool {
			return p
		}
	}
	return nil
}

// registeredTools returns the list of supported tools.
//
// Returns:
//   - []string: Tool identifiers for all registered
//     parsers
func registeredTools() []string {
	tools := make([]string, len(registeredParsers))
	for i, p := range registeredParsers {
		tools[i] = p.Tool()
	}
	return tools
}
