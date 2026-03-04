//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
)

// rustBuilder implements GraphBuilder for Rust projects.
type rustBuilder struct{}

func (r *rustBuilder) Name() string { return "rust" }

func (r *rustBuilder) Detect() bool {
	_, err := os.Stat("Cargo.toml")
	return err == nil
}

func (r *rustBuilder) Build(external bool) (map[string][]string, error) {
	if external {
		return buildRustFullGraph()
	}
	return buildRustInternalGraph()
}

// cargoMetadata represents the subset of `cargo metadata` output we need.
type cargoMetadata struct {
	Packages         []cargoPackage `json:"packages"`
	WorkspaceMembers []string       `json:"workspace_members"`
	Resolve          *cargoResolve  `json:"resolve"`
}

// cargoPackage represents a package in cargo metadata output.
type cargoPackage struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Source       *string       `json:"source"`
	Dependencies []cargoDep    `json:"dependencies"`
	Targets      []cargoTarget `json:"targets"`
}

// cargoDep represents a dependency entry in cargo metadata.
type cargoDep struct {
	Name string `json:"name"`
	Kind *string `json:"kind"`
}

// cargoTarget represents a build target in cargo metadata.
type cargoTarget struct {
	Name string   `json:"name"`
	Kind []string `json:"kind"`
}

// cargoResolve represents the resolved dependency graph.
type cargoResolve struct {
	Nodes []cargoNode `json:"nodes"`
}

// cargoNode represents a node in the resolved dependency graph.
type cargoNode struct {
	ID   string   `json:"id"`
	Deps []string `json:"deps,omitempty"`
}

// runCargoMetadata runs `cargo metadata` and parses the output.
func runCargoMetadata() (*cargoMetadata, error) {
	_, lookErr := exec.LookPath("cargo")
	if lookErr != nil {
		return nil, fmt.Errorf("cargo not found in PATH: install Rust toolchain to analyze Cargo projects")
	}

	out, cmdErr := exec.Command("cargo", "metadata", "--format-version", "1", "--no-deps").Output() //nolint:gosec // fixed args
	if cmdErr != nil {
		return nil, fmt.Errorf("cargo metadata failed: %w", cmdErr)
	}

	var meta cargoMetadata
	if unmarshalErr := json.Unmarshal(out, &meta); unmarshalErr != nil {
		return nil, fmt.Errorf("parsing cargo metadata: %w", unmarshalErr)
	}
	return &meta, nil
}

// runCargoMetadataFull runs `cargo metadata` with full dependency resolution.
func runCargoMetadataFull() (*cargoMetadata, error) {
	_, lookErr := exec.LookPath("cargo")
	if lookErr != nil {
		return nil, fmt.Errorf("cargo not found in PATH: install Rust toolchain to analyze Cargo projects")
	}

	out, cmdErr := exec.Command("cargo", "metadata", "--format-version", "1").Output() //nolint:gosec // fixed args
	if cmdErr != nil {
		return nil, fmt.Errorf("cargo metadata failed: %w", cmdErr)
	}

	var meta cargoMetadata
	if unmarshalErr := json.Unmarshal(out, &meta); unmarshalErr != nil {
		return nil, fmt.Errorf("parsing cargo metadata: %w", unmarshalErr)
	}
	return &meta, nil
}

// buildRustInternalGraph returns workspace member dependencies on each other.
func buildRustInternalGraph() (map[string][]string, error) {
	meta, metaErr := runCargoMetadata()
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

// buildRustFullGraph returns all dependencies for workspace packages.
func buildRustFullGraph() (map[string][]string, error) {
	meta, metaErr := runCargoMetadataFull()
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
