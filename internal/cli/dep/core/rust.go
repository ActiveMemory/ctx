//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"os"
	"os/exec"
	"sort"

	ctxerr "github.com/ActiveMemory/ctx/internal/err/deps"
)

// RustEcosystem is the ecosystem label for Rust projects.
const RustEcosystem = "rust"

// RustBuilder implements GraphBuilder for Rust projects.
type RustBuilder struct{}

// Name returns the ecosystem label.
func (r *RustBuilder) Name() string { return RustEcosystem }

// Detect returns true if Cargo.toml exists in the current directory.
func (r *RustBuilder) Detect() bool {
	_, err := os.Stat("Cargo.toml")
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
func (r *RustBuilder) Build(external bool) (map[string][]string, error) {
	if external {
		return BuildRustFullGraph()
	}
	return BuildRustInternalGraph()
}

// CargoMetadata represents the subset of `cargo metadata` output we need.
type CargoMetadata struct {
	Packages         []CargoPackage `json:"packages"`
	WorkspaceMembers []string       `json:"workspace_members"`
	Resolve          *CargoResolve  `json:"resolve"`
}

// CargoPackage represents a package in cargo metadata output.
type CargoPackage struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Source       *string       `json:"source"`
	Dependencies []CargoDep    `json:"dependencies"`
	Targets      []CargoTarget `json:"targets"`
}

// CargoDep represents a dependency entry in cargo metadata.
type CargoDep struct {
	Name string  `json:"name"`
	Kind *string `json:"kind"`
}

// CargoTarget represents a build target in cargo metadata.
type CargoTarget struct {
	Name string   `json:"name"`
	Kind []string `json:"kind"`
}

// CargoResolve represents the resolved dependency graph.
type CargoResolve struct {
	Nodes []CargoNode `json:"nodes"`
}

// CargoNode represents a node in the resolved dependency graph.
type CargoNode struct {
	ID   string   `json:"id"`
	Deps []string `json:"deps,omitempty"`
}

// RunCargoMetadata runs `cargo metadata` and parses the output.
//
// Returns:
//   - *CargoMetadata: Parsed metadata
//   - error: Non-nil if cargo is not found or output is malformed
func RunCargoMetadata() (*CargoMetadata, error) {
	_, lookErr := exec.LookPath("cargo")
	if lookErr != nil {
		return nil, ctxerr.CargoNotFound()
	}

	out, cmdErr := exec.Command("cargo", "metadata", "--format-version", "1", "--no-deps").Output() //nolint:gosec // fixed args
	if cmdErr != nil {
		return nil, ctxerr.CargoMetadataFailed(cmdErr)
	}

	var meta CargoMetadata
	if unmarshalErr := json.Unmarshal(out, &meta); unmarshalErr != nil {
		return nil, ctxerr.ParseCargoMetadata(unmarshalErr)
	}
	return &meta, nil
}

// RunCargoMetadataFull runs `cargo metadata` with full dependency resolution.
//
// Returns:
//   - *CargoMetadata: Parsed metadata with full resolution
//   - error: Non-nil if cargo is not found or output is malformed
func RunCargoMetadataFull() (*CargoMetadata, error) {
	_, lookErr := exec.LookPath("cargo")
	if lookErr != nil {
		return nil, ctxerr.CargoNotFound()
	}

	out, cmdErr := exec.Command("cargo", "metadata", "--format-version", "1").Output() //nolint:gosec // fixed args
	if cmdErr != nil {
		return nil, ctxerr.CargoMetadataFailed(cmdErr)
	}

	var meta CargoMetadata
	if unmarshalErr := json.Unmarshal(out, &meta); unmarshalErr != nil {
		return nil, ctxerr.ParseCargoMetadata(unmarshalErr)
	}
	return &meta, nil
}

// BuildRustInternalGraph returns workspace member dependencies on each other.
//
// Returns:
//   - map[string][]string: Internal dependency graph
//   - error: Non-nil if cargo metadata fails
func BuildRustInternalGraph() (map[string][]string, error) {
	meta, metaErr := RunCargoMetadata()
	if metaErr != nil {
		return nil, metaErr
	}

	// Build a set of workspace member names.
	wsNames := make(map[string]bool)
	for _, pkg := range meta.Packages {
		if pkg.Source == nil { // local packages have no source
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
			if wsNames[dep.Name] && dep.Name != pkg.Name {
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

// BuildRustFullGraph returns all dependencies for workspace packages.
//
// Returns:
//   - map[string][]string: Full dependency graph
//   - error: Non-nil if cargo metadata fails
func BuildRustFullGraph() (map[string][]string, error) {
	meta, metaErr := RunCargoMetadataFull()
	if metaErr != nil {
		return nil, metaErr
	}

	// Identify local packages.
	localPkgs := make(map[string]bool)
	for _, pkg := range meta.Packages {
		if pkg.Source == nil {
			localPkgs[pkg.Name] = true
		}
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
