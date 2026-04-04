//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

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
