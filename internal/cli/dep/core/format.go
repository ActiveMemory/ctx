//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// MermaidID converts a package path to a valid Mermaid node ID.
//
// Parameters:
//   - pkg: Package path to convert
//
// Returns:
//   - string: Safe Mermaid node identifier
func MermaidID(pkg string) string {
	r := strings.NewReplacer("/", "_", ".", "_", "-", "_")
	return r.Replace(pkg)
}

// RenderMermaid produces a Mermaid graph TD definition.
//
// Parameters:
//   - graph: Adjacency list of package dependencies
//
// Returns:
//   - string: Mermaid graph markup
func RenderMermaid(graph map[string][]string) string {
	var b strings.Builder
	b.WriteString("graph TD\n")

	// Sort keys for deterministic output.
	keys := SortedKeys(graph)

	for _, pkg := range keys {
		deps := graph[pkg]
		src := MermaidID(pkg)
		for _, dep := range deps {
			dst := MermaidID(dep)
			b.WriteString(fmt.Sprintf("    %s[\"%s\"] --> %s[\"%s\"]\n", src, pkg, dst, dep))
		}
	}

	return b.String()
}

// RenderTable produces a Package | Imports table.
//
// Parameters:
//   - graph: Adjacency list of package dependencies
//
// Returns:
//   - string: Formatted table output
func RenderTable(graph map[string][]string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%-50s %s\n", "Package", "Imports"))
	b.WriteString(fmt.Sprintf("%-50s %s\n", strings.Repeat("-", 50), strings.Repeat("-", 30)))

	keys := SortedKeys(graph)
	for _, pkg := range keys {
		deps := graph[pkg]
		b.WriteString(fmt.Sprintf("%-50s %s\n", pkg, strings.Join(deps, ", ")))
	}

	return b.String()
}

// RenderJSON produces a machine-readable JSON adjacency list.
//
// Parameters:
//   - graph: Adjacency list of package dependencies
//
// Returns:
//   - string: Pretty-printed JSON
func RenderJSON(graph map[string][]string) string {
	data, _ := json.MarshalIndent(graph, "", "  ")
	return string(data) + token.NewlineLF
}

// SortedKeys returns the keys of a map sorted alphabetically.
//
// Parameters:
//   - m: Map to extract keys from
//
// Returns:
//   - []string: Sorted keys
func SortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
