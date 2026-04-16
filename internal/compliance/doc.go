//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package compliance contains cross-cutting tests that
// verify the entire codebase adheres to project
// standards documented in CONTRIBUTING.md, CLAUDE.md,
// and the lint scripts.
//
// # What It Checks
//
// These tests inspect source files, configs, and build
// artifacts across the whole repository, mirroring the
// checks performed by hack/lint-drift.sh and
// hack/lint-docs.sh so that violations surface in
// go test without requiring bash.
//
// Typical checks include:
//
//   - Copyright headers are present and correctly
//     formatted in every .go file.
//   - Template files are well-formed and parseable.
//   - Go source is gofmt-compliant.
//   - Build tags and file naming conventions are
//     consistent.
//   - YAML asset files and embedded text constants
//     are in sync with their consumers.
//
// # How It Differs from internal/audit
//
// Compliance tests work at the file and string level:
// they read raw source, run regexes, and shell out to
// tools. Audit tests ([internal/audit]) parse full
// ASTs via go/ast and go/packages for type-aware
// checks. Both test suites run in CI; neither produces
// importable symbols.
//
// Every file in this package is a _test.go file except
// this doc.go.
package compliance
