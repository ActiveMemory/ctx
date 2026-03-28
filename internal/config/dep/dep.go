//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dep

// Packages is used by sync to detect projects and suggest dependency documentation.
var Packages = map[string]string{
	"package.json":     "Node.js dependencies",
	"go.mod":           "Go module dependencies",
	"Cargo.toml":       "Rust dependencies",
	"requirements.txt": "Python dependencies",
	"Gemfile":          "Ruby dependencies",
}

// Table formatting constants for dependency output.
const (
	// TableColPackage is the column width for package names in table output.
	TableColPackage = 50
	// TableColImports is the column width for import lists in table output.
	TableColImports = 30
	// TableHeaderPackage is the column header for package names.
	TableHeaderPackage = "Package"
	// TableHeaderImports is the column header for import lists.
	TableHeaderImports = "Imports"
	// TableRowFormat is the dynamic-width row format template for table output.
	TableRowFormat = "%%-%ds %%s\n"
	// MermaidEdgeFormat is the Mermaid graph edge format string.
	MermaidEdgeFormat = "    %s[\"%s\"] --> %s[\"%s\"]\n"
)
