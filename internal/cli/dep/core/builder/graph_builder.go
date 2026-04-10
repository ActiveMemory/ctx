//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package builder

// GraphBuilder produces a dependency adjacency list for a
// specific ecosystem.
type GraphBuilder interface {
	// Name returns the ecosystem label
	// (e.g. "go", "node", "python", "rust").
	Name() string

	// Detect returns true if the current directory contains
	// this ecosystem's manifest file.
	Detect() bool

	// Build produces an adjacency list of dependencies.
	// When external is false, only internal dependencies
	// are included.
	Build(external bool) (map[string][]string, error)
}
