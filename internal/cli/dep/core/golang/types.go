//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package golang

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
