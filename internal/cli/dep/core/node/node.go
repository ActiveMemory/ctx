//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package node

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	cfgDep "github.com/ActiveMemory/ctx/internal/config/dep"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Builder implements GraphBuilder for Node.js projects.
type Builder struct{}

// Name returns the ecosystem label.
func (n *Builder) Name() string {
	return cfgDep.BuilderNode
}

// Detect returns true if package.json exists in the current
// directory.
func (n *Builder) Detect() bool {
	_, err := os.Stat(cfgDep.PackageJSON)
	return err == nil
}

// Build produces an adjacency list of Node.js dependencies.
//
// Parameters:
//   - external: If true, include all dependencies
//
// Returns:
//   - map[string][]string: Adjacency list
//   - error: Non-nil if package.json parsing fails
func (n *Builder) Build(
	external bool,
) (map[string][]string, error) {
	if external {
		return FullGraph()
	}
	return InternalGraph()
}

// PackageJSON represents the fields we need from
// package.json.
//
// Fields:
//   - Name: Package name
//   - Dependencies: Production dependencies
//   - DevDependencies: Development dependencies
//   - Workspaces: Monorepo workspace configuration
type PackageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Workspaces      Workspaces        `json:"workspaces"`
}

// Workspaces handles the two valid package.json workspaces
// formats: array of globs, or object with "packages" array.
type Workspaces struct {
	Patterns []string
}

// UnmarshalJSON handles both array and object formats for
// workspaces.
func (w *Workspaces) UnmarshalJSON(data []byte) error {
	var arr []string
	if unmarshalErr := json.Unmarshal(
		data, &arr,
	); unmarshalErr == nil {
		w.Patterns = arr
		return nil
	}
	var obj struct {
		Packages []string `json:"packages"`
	}
	if unmarshalErr := json.Unmarshal(
		data, &obj,
	); unmarshalErr == nil {
		w.Patterns = obj.Packages
		return nil
	}
	return nil
}

// ReadPackageJSON reads and parses a package.json file.
//
// Parameters:
//   - path: Path to the package.json file
//
// Returns:
//   - PackageJSON: Parsed package data
//   - error: Non-nil if read or parse fails
func ReadPackageJSON(path string) (PackageJSON, error) {
	data, readErr := io.SafeReadUserFile(path)
	if readErr != nil {
		return PackageJSON{}, readErr
	}
	var pkg PackageJSON
	if unmarshalErr := json.Unmarshal(
		data, &pkg,
	); unmarshalErr != nil {
		return PackageJSON{}, unmarshalErr
	}
	return pkg, nil
}

// InternalGraph builds a workspace-to-workspace dependency
// graph. For single-package projects, returns an empty
// graph.
//
// Returns:
//   - map[string][]string: Workspace dependency graph
//   - error: Non-nil if package.json parsing fails
func InternalGraph() (map[string][]string, error) {
	root, readErr := ReadPackageJSON(cfgDep.PackageJSON)
	if readErr != nil {
		return nil, readErr
	}

	if len(root.Workspaces.Patterns) == 0 {
		return map[string][]string{}, nil
	}

	wsPackages, discoverErr := DiscoverWorkspaces(
		root.Workspaces.Patterns,
	)
	if discoverErr != nil {
		return nil, discoverErr
	}

	wsNames := make(map[string]bool)
	for _, ws := range wsPackages {
		wsNames[ws.Name] = true
	}

	graph := make(map[string][]string)
	for _, ws := range wsPackages {
		var deps []string
		for dep := range ws.Dependencies {
			if wsNames[dep] && dep != ws.Name {
				deps = append(deps, dep)
			}
		}
		for dep := range ws.DevDependencies {
			if wsNames[dep] && dep != ws.Name {
				deps = append(deps, dep)
			}
		}
		if len(deps) > 0 {
			sort.Strings(deps)
			graph[ws.Name] = deps
		}
	}
	return graph, nil
}

// FullGraph returns the full dependency list from
// package.json. For workspaces, includes all workspace and
// external deps.
//
// Returns:
//   - map[string][]string: Full dependency graph
//   - error: Non-nil if package.json parsing fails
func FullGraph() (map[string][]string, error) {
	root, readErr := ReadPackageJSON(cfgDep.PackageJSON)
	if readErr != nil {
		return nil, readErr
	}

	if len(root.Workspaces.Patterns) == 0 {
		return SinglePackageGraph(root)
	}

	wsPackages, discoverErr := DiscoverWorkspaces(
		root.Workspaces.Patterns,
	)
	if discoverErr != nil {
		return nil, discoverErr
	}

	graph := make(map[string][]string)
	for _, ws := range wsPackages {
		var deps []string
		for dep := range ws.Dependencies {
			deps = append(deps, dep)
		}
		for dep := range ws.DevDependencies {
			deps = append(deps, dep)
		}
		if len(deps) > 0 {
			sort.Strings(deps)
			graph[ws.Name] = deps
		}
	}
	return graph, nil
}

// SinglePackageGraph returns deps for a single-package
// project.
//
// Parameters:
//   - pkg: Parsed package.json data
//
// Returns:
//   - map[string][]string: Dependency graph
//   - error: Always nil
func SinglePackageGraph(
	pkg PackageJSON,
) (map[string][]string, error) {
	name := pkg.Name
	if name == "" {
		name = cfgDep.WorkspaceRoot
	}

	var deps []string
	for dep := range pkg.Dependencies {
		deps = append(deps, dep)
	}
	for dep := range pkg.DevDependencies {
		deps = append(deps, dep)
	}

	graph := make(map[string][]string)
	if len(deps) > 0 {
		sort.Strings(deps)
		graph[name] = deps
	}
	return graph, nil
}

// DiscoverWorkspaces finds all workspace package.json files
// matching the given glob patterns and returns their parsed
// contents.
//
// Parameters:
//   - patterns: Glob patterns to match workspace directories
//
// Returns:
//   - []PackageJSON: Parsed workspace packages
//   - error: Non-nil if glob matching fails
func DiscoverWorkspaces(
	patterns []string,
) ([]PackageJSON, error) {
	var pkgs []PackageJSON
	seen := make(map[string]bool)

	for _, pattern := range patterns {
		matches, globErr := filepath.Glob(pattern)
		if globErr != nil {
			return nil, globErr
		}
		for _, match := range matches {
			pkgPath := filepath.Join(
				match, cfgDep.PackageJSON,
			)
			info, statErr := os.Stat(pkgPath)
			if statErr != nil || info.IsDir() {
				continue
			}
			if seen[pkgPath] {
				continue
			}
			seen[pkgPath] = true

			pkg, readErr := ReadPackageJSON(pkgPath)
			if readErr != nil {
				continue
			}
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs, nil
}
