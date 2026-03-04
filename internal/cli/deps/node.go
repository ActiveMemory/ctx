//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

// nodeBuilder implements GraphBuilder for Node.js projects.
type nodeBuilder struct{}

func (n *nodeBuilder) Name() string { return "node" }

func (n *nodeBuilder) Detect() bool {
	_, err := os.Stat("package.json")
	return err == nil
}

func (n *nodeBuilder) Build(external bool) (map[string][]string, error) {
	if external {
		return buildNodeFullGraph()
	}
	return buildNodeInternalGraph()
}

// packageJSON represents the fields we need from package.json.
type packageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Workspaces      workspaces        `json:"workspaces"`
}

// workspaces handles the two valid package.json workspaces formats:
// array of globs, or object with "packages" array.
type workspaces struct {
	Patterns []string
}

func (w *workspaces) UnmarshalJSON(data []byte) error {
	// Try array first.
	var arr []string
	if unmarshalErr := json.Unmarshal(data, &arr); unmarshalErr == nil {
		w.Patterns = arr
		return nil
	}
	// Try object with "packages" field.
	var obj struct {
		Packages []string `json:"packages"`
	}
	if unmarshalErr := json.Unmarshal(data, &obj); unmarshalErr == nil {
		w.Patterns = obj.Packages
		return nil
	}
	return nil
}

// readPackageJSON reads and parses a package.json file.
func readPackageJSON(path string) (packageJSON, error) {
	data, readErr := os.ReadFile(path) //nolint:gosec // G304: path is constructed from workspace glob matches
	if readErr != nil {
		return packageJSON{}, readErr
	}
	var pkg packageJSON
	if unmarshalErr := json.Unmarshal(data, &pkg); unmarshalErr != nil {
		return packageJSON{}, unmarshalErr
	}
	return pkg, nil
}

// buildNodeInternalGraph builds a workspace-to-workspace dependency graph.
// For single-package projects, returns an empty graph (no internal deps).
func buildNodeInternalGraph() (map[string][]string, error) {
	root, readErr := readPackageJSON("package.json")
	if readErr != nil {
		return nil, readErr
	}

	if len(root.Workspaces.Patterns) == 0 {
		// Single package — no internal dependency graph to show.
		return map[string][]string{}, nil
	}

	// Discover workspace packages.
	wsPackages, discoverErr := discoverWorkspaces(root.Workspaces.Patterns)
	if discoverErr != nil {
		return nil, discoverErr
	}

	// Build a set of workspace package names for filtering.
	wsNames := make(map[string]bool)
	for _, ws := range wsPackages {
		wsNames[ws.Name] = true
	}

	// Build adjacency list: workspace → workspace deps.
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

// buildNodeFullGraph returns the full dependency list from package.json.
// For workspaces, includes all workspace and external deps.
func buildNodeFullGraph() (map[string][]string, error) {
	root, readErr := readPackageJSON("package.json")
	if readErr != nil {
		return nil, readErr
	}

	if len(root.Workspaces.Patterns) == 0 {
		// Single package — list all deps under package name.
		return buildNodeSinglePackageGraph(root)
	}

	// Workspace project — show each workspace's full deps.
	wsPackages, discoverErr := discoverWorkspaces(root.Workspaces.Patterns)
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

// buildNodeSinglePackageGraph returns deps for a single-package project.
func buildNodeSinglePackageGraph(pkg packageJSON) (map[string][]string, error) {
	name := pkg.Name
	if name == "" {
		name = "root"
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

// discoverWorkspaces finds all workspace package.json files matching the
// given glob patterns and returns their parsed contents.
func discoverWorkspaces(patterns []string) ([]packageJSON, error) {
	var pkgs []packageJSON
	seen := make(map[string]bool)

	for _, pattern := range patterns {
		matches, globErr := filepath.Glob(pattern)
		if globErr != nil {
			return nil, globErr
		}
		for _, match := range matches {
			pkgPath := filepath.Join(match, "package.json")
			info, statErr := os.Stat(pkgPath)
			if statErr != nil || info.IsDir() {
				continue
			}
			if seen[pkgPath] {
				continue
			}
			seen[pkgPath] = true

			pkg, readErr := readPackageJSON(pkgPath)
			if readErr != nil {
				continue
			}
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs, nil
}
