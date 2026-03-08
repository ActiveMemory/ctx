//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

// DocAliases maps user-facing names to embedded asset names.
var DocAliases = map[string]string{
	"manifesto":  "manifesto",
	"about":      "about",
	"invariants": "design-invariants",
}

// DocEntry pairs a document alias with its display label.
type DocEntry struct {
	Alias string
	Label string
}

// DocOrder defines the display order for the interactive menu.
var DocOrder = []DocEntry{
	{"manifesto", "The ctx Manifesto"},
	{"about", "About ctx"},
	{"invariants", "Design Invariants"},
}
