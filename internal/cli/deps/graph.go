//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"
)

// goPackage represents the subset of `go list -json` output we need.
type goPackage struct {
	ImportPath string   `json:"ImportPath"`
	Name       string   `json:"Name"`
	Imports    []string `json:"Imports"`
	Module     *struct {
		Path string `json:"Path"`
	} `json:"Module"`
}

// detectProjectType checks for known project manifests.
func detectProjectType() string {
	if _, err := os.Stat("go.mod"); err == nil {
		return "go"
	}
	return ""
}

// modulePath reads the module path from go.mod via the first goPackage.
func modulePath(pkgs []goPackage) string {
	for _, p := range pkgs {
		if p.Module != nil && p.Module.Path != "" {
			return p.Module.Path
		}
	}
	return ""
}

// listGoPackages runs `go list -json ./...` and parses the output.
// go list outputs concatenated JSON objects (not an array).
func listGoPackages() ([]goPackage, error) {
	out, err := exec.Command("go", "list", "-json", "./...").Output() //nolint:gosec // fixed args
	if err != nil {
		return nil, err
	}

	var pkgs []goPackage
	dec := json.NewDecoder(strings.NewReader(string(out)))
	for {
		var p goPackage
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

// isStdlib returns true if the import path looks like a Go stdlib package.
// Heuristic: no dot in the first path component.
func isStdlib(path string) bool {
	first := path
	if i := strings.Index(path, "/"); i >= 0 {
		first = path[:i]
	}
	return !strings.Contains(first, ".")
}

// shortPkgName strips the module prefix for readability.
func shortPkgName(importPath, modPath string) string {
	if modPath != "" && strings.HasPrefix(importPath, modPath+"/") {
		return importPath[len(modPath)+1:]
	}
	return importPath
}

// buildInternalGraph returns an adjacency list of internal package dependencies.
// Keys and values use shortened names (module prefix stripped).
func buildInternalGraph() (map[string][]string, error) {
	pkgs, err := listGoPackages()
	if err != nil {
		return nil, err
	}

	modPath := modulePath(pkgs)

	// Build a set of internal packages for filtering.
	internal := make(map[string]bool)
	for _, p := range pkgs {
		internal[p.ImportPath] = true
	}

	graph := make(map[string][]string)
	for _, p := range pkgs {
		short := shortPkgName(p.ImportPath, modPath)
		var deps []string
		for _, imp := range p.Imports {
			if internal[imp] && imp != p.ImportPath {
				deps = append(deps, shortPkgName(imp, modPath))
			}
		}
		if len(deps) > 0 {
			graph[short] = deps
		}
	}
	return graph, nil
}

// buildFullGraph returns an adjacency list including external dependencies.
// Stdlib packages are excluded.
func buildFullGraph() (map[string][]string, error) {
	pkgs, err := listGoPackages()
	if err != nil {
		return nil, err
	}

	modPath := modulePath(pkgs)

	graph := make(map[string][]string)
	for _, p := range pkgs {
		short := shortPkgName(p.ImportPath, modPath)
		var deps []string
		for _, imp := range p.Imports {
			if isStdlib(imp) {
				continue
			}
			if imp == p.ImportPath {
				continue
			}
			deps = append(deps, shortPkgName(imp, modPath))
		}
		if len(deps) > 0 {
			graph[short] = deps
		}
	}
	return graph, nil
}
