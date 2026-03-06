//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"encoding/json"
	"fmt"
	"github.com/ActiveMemory/ctx/internal/config"
	"sort"
	"strings"
)

// mermaidID converts a package path to a valid Mermaid node ID.
func mermaidID(pkg string) string {
	r := strings.NewReplacer("/", "_", ".", "_", "-", "_")
	return r.Replace(pkg)
}

// renderMermaid produces a Mermaid graph TD definition.
func renderMermaid(graph map[string][]string) string {
	var b strings.Builder
	b.WriteString("graph TD\n")

	// Sort keys for deterministic output.
	keys := sortedKeys(graph)

	for _, pkg := range keys {
		deps := graph[pkg]
		src := mermaidID(pkg)
		for _, dep := range deps {
			dst := mermaidID(dep)
			b.WriteString(fmt.Sprintf("    %s[\"%s\"] --> %s[\"%s\"]\n", src, pkg, dst, dep))
		}
	}

	return b.String()
}

// renderTable produces a Package | Imports table.
func renderTable(graph map[string][]string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%-50s %s\n", "Package", "Imports"))
	b.WriteString(fmt.Sprintf("%-50s %s\n", strings.Repeat("-", 50), strings.Repeat("-", 30)))

	keys := sortedKeys(graph)
	for _, pkg := range keys {
		deps := graph[pkg]
		b.WriteString(fmt.Sprintf("%-50s %s\n", pkg, strings.Join(deps, ", ")))
	}

	return b.String()
}

// renderJSON produces a machine-readable JSON adjacency list.
func renderJSON(graph map[string][]string) string {
	data, _ := json.MarshalIndent(graph, "", "  ")
	return string(data) + config.NewlineLF
}

// sortedKeys returns the keys of a map sorted alphabetically.
func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
