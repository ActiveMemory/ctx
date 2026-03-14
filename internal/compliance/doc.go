//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package compliance contains cross-cutting tests that verify the entire
// codebase adheres to project standards.
//
// These tests inspect source files, configs, and build artifacts across the
// whole repository, mirroring the checks performed by the lint-drift and
// lint-docs scripts so that violations surface in go test without requiring
// bash.
package compliance
