//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package golang

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	cfgDep "github.com/ActiveMemory/ctx/internal/config/dep"
	"github.com/ActiveMemory/ctx/internal/config/token"
	execDep "github.com/ActiveMemory/ctx/internal/exec/dep"
)

// GoPackage represents the subset of `go list -json` output
// needed for dependency graph construction.
//
// Fields:
//   - ImportPath: Full import path
//   - Name: Package name
//   - Imports: Direct import paths
//   - Module: Enclosing module (nil for stdlib)
type GoPackage struct {
	ImportPath string   `json:"ImportPath"`
	Name       string   `json:"Name"`
	Imports    []string `json:"Imports"`
	Module     *struct {
		Path string `json:"Path"`
	} `json:"Module"`
}

// Builder implements GraphBuilder for Go projects.
type Builder struct{}

// Name returns the ecosystem label.
func (g *Builder) Name() string { return cfgDep.GoBinary }

// Detect returns true if go.mod exists in the current
// directory.
func (g *Builder) Detect() bool {
	_, err := os.Stat(cfgDep.GoMod)
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
func (g *Builder) Build(
	external bool,
) (map[string][]string, error) {
	if external {
		return FullGraph()
	}
	return InternalGraph()
}

// ModulePath reads the module path from the first GoPackage
// with a Module field.
//
// Parameters:
//   - pkgs: Parsed go list output
//
// Returns:
//   - string: Module path, or empty if not found
func ModulePath(pkgs []GoPackage) string {
	for _, p := range pkgs {
		if p.Module != nil && p.Module.Path != "" {
			return p.Module.Path
		}
	}
	return ""
}

// ListPackages runs `go list -json ./...` and parses the
// output. go list outputs concatenated JSON objects (not an
// array).
//
// Returns:
//   - []GoPackage: Parsed packages
//   - error: Non-nil if go list fails or output is malformed
func ListPackages() ([]GoPackage, error) {
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

// IsStdlib returns true if the import path looks like a Go
// stdlib package. Heuristic: no dot in the first path
// component.
//
// Parameters:
//   - path: Import path to check
//
// Returns:
//   - bool: True if the path is a stdlib package
func IsStdlib(path string) bool {
	first := path
	if i := strings.Index(path, token.Slash); i >= 0 {
		first = path[:i]
	}
	return !strings.Contains(first, token.Dot)
}

// ShortPkgName strips the module prefix for readability.
//
// Parameters:
//   - importPath: Full import path
//   - modPath: Module path prefix to strip
//
// Returns:
//   - string: Shortened path, or original if prefix
//     doesn't match
func ShortPkgName(importPath, modPath string) string {
	if modPath != "" &&
		strings.HasPrefix(importPath, modPath+token.Slash) {
		return importPath[len(modPath)+1:]
	}
	return importPath
}

// InternalGraph returns an adjacency list of internal
// package dependencies. Keys and values use shortened names
// (module prefix stripped).
//
// Returns:
//   - map[string][]string: Internal dependency graph
//   - error: Non-nil if go list fails
func InternalGraph() (map[string][]string, error) {
	pkgs, listErr := ListPackages()
	if listErr != nil {
		return nil, listErr
	}

	modPath := ModulePath(pkgs)

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
				deps = append(
					deps,
					ShortPkgName(imp, modPath),
				)
			}
		}
		if len(deps) > 0 {
			graph[short] = deps
		}
	}
	return graph, nil
}

// FullGraph returns an adjacency list including external
// dependencies. Stdlib packages are excluded.
//
// Returns:
//   - map[string][]string: Full dependency graph
//   - error: Non-nil if go list fails
func FullGraph() (map[string][]string, error) {
	pkgs, listErr := ListPackages()
	if listErr != nil {
		return nil, listErr
	}

	modPath := ModulePath(pkgs)

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
			deps = append(
				deps, ShortPkgName(imp, modPath),
			)
		}
		if len(deps) > 0 {
			graph[short] = deps
		}
	}
	return graph, nil
}
