//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/ActiveMemory/ctx/internal/cli/dep/core/golang"
	"github.com/ActiveMemory/ctx/internal/cli/dep/core/node"
	"github.com/ActiveMemory/ctx/internal/cli/dep/core/python"
	"github.com/ActiveMemory/ctx/internal/cli/dep/core/rust"
)

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

// Builders is the ordered registry of graph builders.
// Detection walks this list; the first match wins.
var Builders = []GraphBuilder{
	&golang.Builder{},
	&node.Builder{},
	&python.Builder{},
	&rust.Builder{},
}

// DetectBuilder returns the first builder whose Detect()
// returns true, or nil if no ecosystem is detected.
//
// Returns:
//   - GraphBuilder: The detected builder, or nil
func DetectBuilder() GraphBuilder {
	for _, b := range Builders {
		if b.Detect() {
			return b
		}
	}
	return nil
}

// FindBuilder returns the builder matching the given name,
// or nil.
//
// Parameters:
//   - name: Ecosystem name to find (e.g. "go", "node")
//
// Returns:
//   - GraphBuilder: The matching builder, or nil
func FindBuilder(name string) GraphBuilder {
	for _, b := range Builders {
		if b.Name() == name {
			return b
		}
	}
	return nil
}

// BuilderNames returns all registered builder names for
// help text.
//
// Returns:
//   - []string: Ordered list of ecosystem names
func BuilderNames() []string {
	names := make([]string, len(Builders))
	for i, b := range Builders {
		names[i] = b.Name()
	}
	return names
}
