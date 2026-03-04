//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

// GraphBuilder produces a dependency adjacency list for a specific ecosystem.
type GraphBuilder interface {
	// Name returns the ecosystem label (e.g. "go", "node", "python", "rust").
	Name() string

	// Detect returns true if the current directory contains this ecosystem's
	// manifest file (e.g. go.mod, package.json).
	Detect() bool

	// Build produces an adjacency list of dependencies.
	// When external is false, only internal/project dependencies are included.
	// When external is true, third-party dependencies are included too.
	Build(external bool) (map[string][]string, error)
}

// builders is the ordered registry of graph builders.
// Detection walks this list; first match wins.
var builders = []GraphBuilder{
	&goBuilder{},
	&nodeBuilder{},
	&pythonBuilder{},
	&rustBuilder{},
}

// detectBuilder returns the first builder whose Detect() returns true,
// or nil if no ecosystem is detected.
func detectBuilder() GraphBuilder {
	for _, b := range builders {
		if b.Detect() {
			return b
		}
	}
	return nil
}

// findBuilder returns the builder matching the given name, or nil.
func findBuilder(name string) GraphBuilder {
	for _, b := range builders {
		if b.Name() == name {
			return b
		}
	}
	return nil
}

// builderNames returns all registered builder names for help text.
func builderNames() []string {
	names := make([]string, len(builders))
	for i, b := range builders {
		names[i] = b.Name()
	}
	return names
}
