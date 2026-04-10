//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package node

// Builder implements GraphBuilder for Node.js projects.
type Builder struct{}

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
//
// Fields:
//   - Patterns: Workspace glob patterns from either format
type Workspaces struct {
	Patterns []string
}
