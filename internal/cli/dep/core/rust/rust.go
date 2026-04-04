//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rust

import (
	"encoding/json"
	"os"
	"sort"

	cfgDep "github.com/ActiveMemory/ctx/internal/config/dep"
	errDeps "github.com/ActiveMemory/ctx/internal/err/dep"
	execDep "github.com/ActiveMemory/ctx/internal/exec/dep"
)

// Builder implements GraphBuilder for Rust projects.
type Builder struct{}

// Name returns the ecosystem label.
func (r *Builder) Name() string {
	return cfgDep.EcosystemRust
}

// Detect returns true if Cargo.toml exists in the current
// directory.
func (r *Builder) Detect() bool {
	_, err := os.Stat(cfgDep.CargoToml)
	return err == nil
}

// Build produces an adjacency list of Rust dependencies.
//
// Parameters:
//   - external: If true, include all external dependencies
//
// Returns:
//   - map[string][]string: Adjacency list
//   - error: Non-nil if cargo metadata fails
func (r *Builder) Build(
	external bool,
) (map[string][]string, error) {
	if external {
		return FullGraph()
	}
	return InternalGraph()
}

// CargoMetadata represents the subset of `cargo metadata`
// output needed for dependency graph construction.
//
// Fields:
//   - Packages: All packages in the workspace
//   - WorkspaceMembers: Package IDs in the workspace
//   - Resolve: Resolved dependency graph
type CargoMetadata struct {
	Packages         []CargoPackage `json:"packages"`
	WorkspaceMembers []string       `json:"workspace_members"`
	Resolve          *CargoResolve  `json:"resolve"`
}

// CargoPackage represents a package in cargo metadata
// output.
//
// Fields:
//   - ID: Unique package identifier
//   - Name: Package name
//   - Source: Registry source (nil for local)
//   - Dependencies: Declared dependencies
//   - Targets: Build targets
type CargoPackage struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Source       *string       `json:"source"`
	Dependencies []CargoDep    `json:"dependencies"`
	Targets      []CargoTarget `json:"targets"`
}

// CargoDep represents a dependency entry in cargo metadata.
//
// Fields:
//   - Name: Dependency crate name
//   - Kind: Dependency kind (nil=normal, "dev", "build")
type CargoDep struct {
	Name string  `json:"name"`
	Kind *string `json:"kind"`
}

// CargoTarget represents a build target in cargo metadata.
//
// Fields:
//   - Name: Target name
//   - Kind: Target kinds (lib, bin, test, etc.)
type CargoTarget struct {
	Name string   `json:"name"`
	Kind []string `json:"kind"`
}

// CargoResolve represents the resolved dependency graph.
type CargoResolve struct {
	Nodes []CargoNode `json:"nodes"`
}

// CargoNode represents a node in the resolved dependency
// graph.
//
// Fields:
//   - ID: Package identifier
//   - Deps: Resolved dependency IDs
type CargoNode struct {
	ID   string   `json:"id"`
	Deps []string `json:"deps,omitempty"`
}

// RunMetadata runs `cargo metadata` and parses the output.
//
// Returns:
//   - *CargoMetadata: Parsed metadata
//   - error: Non-nil if cargo fails
func RunMetadata() (*CargoMetadata, error) {
	out, cmdErr := execDep.CargoMetadata(true)
	if cmdErr != nil {
		return nil, errDeps.CargoMetadataFailed(cmdErr)
	}

	var meta CargoMetadata
	if unmarshalErr := json.Unmarshal(
		out, &meta,
	); unmarshalErr != nil {
		return nil, errDeps.ParseCargoMetadata(
			unmarshalErr,
		)
	}
	return &meta, nil
}

// RunMetadataFull runs `cargo metadata` with full
// dependency resolution.
//
// Returns:
//   - *CargoMetadata: Parsed metadata with full resolution
//   - error: Non-nil if cargo fails
func RunMetadataFull() (*CargoMetadata, error) {
	out, cmdErr := execDep.CargoMetadata(false)
	if cmdErr != nil {
		return nil, errDeps.CargoMetadataFailed(cmdErr)
	}

	var meta CargoMetadata
	if unmarshalErr := json.Unmarshal(
		out, &meta,
	); unmarshalErr != nil {
		return nil, errDeps.ParseCargoMetadata(
			unmarshalErr,
		)
	}
	return &meta, nil
}

// InternalGraph returns workspace member dependencies on
// each other.
//
// Returns:
//   - map[string][]string: Internal dependency graph
//   - error: Non-nil if cargo metadata fails
func InternalGraph() (map[string][]string, error) {
	meta, metaErr := RunMetadata()
	if metaErr != nil {
		return nil, metaErr
	}

	wsNames := make(map[string]bool)
	for _, pkg := range meta.Packages {
		if pkg.Source == nil {
			wsNames[pkg.Name] = true
		}
	}

	graph := make(map[string][]string)
	for _, pkg := range meta.Packages {
		if pkg.Source != nil {
			continue
		}
		var deps []string
		for _, dep := range pkg.Dependencies {
			if wsNames[dep.Name] &&
				dep.Name != pkg.Name {
				deps = append(deps, dep.Name)
			}
		}
		if len(deps) > 0 {
			sort.Strings(deps)
			graph[pkg.Name] = deps
		}
	}
	return graph, nil
}

// FullGraph returns all dependencies for workspace
// packages.
//
// Returns:
//   - map[string][]string: Full dependency graph
//   - error: Non-nil if cargo metadata fails
func FullGraph() (map[string][]string, error) {
	meta, metaErr := RunMetadataFull()
	if metaErr != nil {
		return nil, metaErr
	}

	graph := make(map[string][]string)
	for _, pkg := range meta.Packages {
		if pkg.Source != nil {
			continue
		}
		var deps []string
		for _, dep := range pkg.Dependencies {
			if dep.Name != pkg.Name {
				deps = append(deps, dep.Name)
			}
		}
		if len(deps) > 0 {
			sort.Strings(deps)
			graph[pkg.Name] = deps
		}
	}
	return graph, nil
}
