//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	execDep "github.com/ActiveMemory/ctx/internal/exec/dep"
)

// GoBuilder implements GraphBuilder for Go projects.
type GoBuilder struct{}

// Name returns the ecosystem label.
func (g *GoBuilder) Name() string { return "go" }

// Detect returns true if go.mod exists in the current directory.
func (g *GoBuilder) Detect() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}

// Build produces an adjacency list of Go dependencies.
//
// Parameters:
//   - external: If true, include third-party dependencies
//
// Returns:
//   - map[string][]string: Adjacency list
//   - error: Non-nil if go list fails
func (g *GoBuilder) Build(external bool) (map[string][]string, error) {
	if external {
		return BuildGoFullGraph()
	}
	return BuildGoInternalGraph()
}

// GoModulePath reads the module path from the first
// GoPackage with a Module field.
//
// Parameters:
//   - pkgs: Parsed go list output
//
// Returns:
//   - string: Module path, or empty if not found
func GoModulePath(pkgs []GoPackage) string {
	for _, p := range pkgs {
		if p.Module != nil && p.Module.Path != "" {
			return p.Module.Path
		}
	}
	return ""
}

// ListGoPackages runs `go list -json ./...` and parses the output.
// go list outputs concatenated JSON objects (not an array).
//
// Returns:
//   - []GoPackage: Parsed packages
//   - error: Non-nil if go list fails or output is malformed
func ListGoPackages() ([]GoPackage, error) {
	out, listErr := execDep.GoListPackages()
	if listErr != nil {
		return nil, listErr
	}

	var pkgs []GoPackage
	dec := json.NewDecoder(strings.NewReader(string(out)))
	for {
		var p GoPackage
		if decErr := dec.Decode(&p); decErr != nil {
			if decErr == io.EOF {
				break
			}
			return nil, decErr
		}
		pkgs = append(pkgs, p)
	}
	return pkgs, nil
}

// IsStdlib returns true if the import path looks like a Go stdlib package.
// Heuristic: no dot in the first path component.
//
// Parameters:
//   - path: Import path to check
//
// Returns:
//   - bool: True if the path is a stdlib package
func IsStdlib(path string) bool {
	first := path
	if i := strings.Index(path, "/"); i >= 0 {
		first = path[:i]
	}
	return !strings.Contains(first, ".")
}

// ShortPkgName strips the module prefix for readability.
//
// Parameters:
//   - importPath: Full import path
//   - modPath: Module path prefix to strip
//
// Returns:
//   - string: Shortened path, or original if prefix doesn't match
func ShortPkgName(importPath, modPath string) string {
	if modPath != "" && strings.HasPrefix(importPath, modPath+"/") {
		return importPath[len(modPath)+1:]
	}
	return importPath
}

// BuildGoInternalGraph returns an adjacency list of
// internal package dependencies.
// Keys and values use shortened names (module prefix stripped).
//
// Returns:
//   - map[string][]string: Internal dependency graph
//   - error: Non-nil if go list fails
func BuildGoInternalGraph() (map[string][]string, error) {
	pkgs, listErr := ListGoPackages()
	if listErr != nil {
		return nil, listErr
	}

	modPath := GoModulePath(pkgs)

	// Build a set of internal packages for filtering.
	internal := make(map[string]bool)
	for _, p := range pkgs {
		internal[p.ImportPath] = true
	}

	graph := make(map[string][]string)
	for _, p := range pkgs {
		short := ShortPkgName(p.ImportPath, modPath)
		var deps []string
		for _, imp := range p.Imports {
			if internal[imp] && imp != p.ImportPath {
				deps = append(deps, ShortPkgName(imp, modPath))
			}
		}
		if len(deps) > 0 {
			graph[short] = deps
		}
	}
	return graph, nil
}

// BuildGoFullGraph returns an adjacency list including external dependencies.
// Stdlib packages are excluded.
//
// Returns:
//   - map[string][]string: Full dependency graph
//   - error: Non-nil if go list fails
func BuildGoFullGraph() (map[string][]string, error) {
	pkgs, listErr := ListGoPackages()
	if listErr != nil {
		return nil, listErr
	}

	modPath := GoModulePath(pkgs)

	graph := make(map[string][]string)
	for _, p := range pkgs {
		short := ShortPkgName(p.ImportPath, modPath)
		var deps []string
		for _, imp := range p.Imports {
			if IsStdlib(imp) {
				continue
			}
			if imp == p.ImportPath {
				continue
			}
			deps = append(deps, ShortPkgName(imp, modPath))
		}
		if len(deps) > 0 {
			graph[short] = deps
		}
	}
	return graph, nil
}
