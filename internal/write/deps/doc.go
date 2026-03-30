//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package deps provides terminal output for the dependency graph
// command (ctx dep).
//
// Output functions handle multiple rendering formats:
// [Table] for human-readable tabular output, [Mermaid] for
// diagram-ready graph notation, and [JSON] for machine-readable
// output. [InfoNoProject] and [NoDeps] handle the empty cases.
//
// Example:
//
//	switch format {
//	case "table":
//	    write.Table(cmd, content)
//	case "mermaid":
//	    write.Mermaid(cmd, content)
//	}
package deps
