//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package render

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/dep"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"
)

// MermaidID converts a package path to a valid Mermaid node
// ID.
//
// Parameters:
//   - pkg: Package path to convert
//
// Returns:
//   - string: Safe Mermaid node identifier
func MermaidID(pkg string) string {
	return regex.MermaidUnsafe.ReplaceAllString(
		pkg, token.Underscore,
	)
}

// Mermaid produces a Mermaid graph TD definition.
//
// Parameters:
//   - graph: Adjacency list of package dependencies
//
// Returns:
//   - string: Mermaid graph markup
func Mermaid(graph map[string][]string) string {
	var b strings.Builder
	b.WriteString(dep.MermaidHeader)

	keys := SortedKeys(graph)

	edgeFmt := dep.MermaidEdgeFormat
	for _, pkg := range keys {
		deps := graph[pkg]
		src := MermaidID(pkg)
		for _, d := range deps {
			dst := MermaidID(d)
			io.SafeFprintf(&b, edgeFmt, src, pkg, dst, d)
		}
	}

	return b.String()
}

// Table produces a Package | Imports table.
//
// Parameters:
//   - graph: Adjacency list of package dependencies
//
// Returns:
//   - string: Formatted table output
func Table(graph map[string][]string) string {
	tf := fmt.Sprintf(
		dep.TableRowFormat, dep.TableColPackage,
	)
	var b strings.Builder
	io.SafeFprintf(
		&b, tf,
		dep.TableHeaderPackage, dep.TableHeaderImports,
	)
	io.SafeFprintf(&b, tf,
		strings.Repeat(token.Dash, dep.TableColPackage),
		strings.Repeat(token.Dash, dep.TableColImports),
	)

	keys := SortedKeys(graph)
	for _, pkg := range keys {
		deps := graph[pkg]
		io.SafeFprintf(
			&b, tf, pkg,
			strings.Join(deps, token.CommaSpace),
		)
	}

	return b.String()
}

// JSON produces a machine-readable JSON adjacency list.
//
// Parameters:
//   - graph: Adjacency list of package dependencies
//
// Returns:
//   - string: Pretty-printed JSON
func JSON(graph map[string][]string) string {
	data, _ := json.MarshalIndent(graph, "", "  ")
	return string(data) + token.NewlineLF
}

// SortedKeys returns the keys of a map sorted
// alphabetically.
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
