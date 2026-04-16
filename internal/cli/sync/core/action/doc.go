//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package action detects discrepancies between the
// codebase and context documentation, suggesting sync
// actions to reconcile them.
//
// # Detection Algorithm
//
// [Detect] scans the loaded context and runs multiple
// validation checks in sequence:
//
//  1. New directories: identifies top-level directories
//     not mentioned in ARCHITECTURE.md, suggesting they
//     be documented.
//  2. Package manager files: detects package manager
//     files (go.mod, package.json, etc.) that may need
//     dependency documentation.
//  3. Config files: finds common configuration files
//     not mentioned in CONVENTIONS.md, suggesting they
//     be documented.
//
// Each check is delegated to the validate subpackage,
// which returns typed Action structs describing the
// discrepancy and suggested fix.
//
// # Data Flow
//
// The cmd/ layer loads a context entity and passes it
// to [Detect]. The returned actions are displayed to
// the user, who can choose which to apply. This keeps
// context documentation aligned with the actual
// codebase structure.
package action
