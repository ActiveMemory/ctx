//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package builder

import (
	"github.com/ActiveMemory/ctx/internal/cli/dep/core/golang"
	"github.com/ActiveMemory/ctx/internal/cli/dep/core/node"
	"github.com/ActiveMemory/ctx/internal/cli/dep/core/python"
	"github.com/ActiveMemory/ctx/internal/cli/dep/core/rust"
)

// Builders is the ordered registry of graph builders.
// Detection walks this list; the first match wins.
var Builders = []GraphBuilder{
	&golang.Builder{},
	&node.Builder{},
	&python.Builder{},
	&rust.Builder{},
}

// Detect returns the first builder whose Detect()
// returns true, or nil if no ecosystem is detected.
//
// Returns:
//   - GraphBuilder: The detected builder, or nil
func Detect() GraphBuilder {
	for _, b := range Builders {
		if b.Detect() {
			return b
		}
	}
	return nil
}

// Find returns the builder matching the given name,
// or nil.
//
// Parameters:
//   - name: Ecosystem name to find (e.g. "go", "node")
//
// Returns:
//   - GraphBuilder: The matching builder, or nil
func Find(name string) GraphBuilder {
	for _, b := range Builders {
		if b.Name() == name {
			return b
		}
	}
	return nil
}

// Names returns all registered builder names for
// help text.
//
// Returns:
//   - []string: Ordered list of ecosystem names
func Names() []string {
	names := make([]string, len(Builders))
	for i, b := range Builders {
		names[i] = b.Name()
	}
	return names
}
