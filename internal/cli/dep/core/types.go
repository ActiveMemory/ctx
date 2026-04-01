//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

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

// CargoPackage represents a package in cargo metadata output.
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
