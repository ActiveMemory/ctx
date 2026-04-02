//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package audit contains AST-based codebase invariant tests.
//
// Unlike internal/compliance (which uses file-level grep and shell tool
// checks), audit tests use go/ast and go/packages to walk parsed syntax
// trees. This gives type-aware, context-sensitive detection that cannot
// be achieved with regex.
//
// Every file in this package is a _test.go file except this doc.go.
// The package produces no binary output and is not importable.
//
// Shared helpers live in helpers_test.go:
//   - [loadPackages] loads and caches parsed packages via sync.Once.
//   - [isTestFile] filters _test.go files.
//   - [posString] formats file:line for error messages.
//
// Each check lives in its own _test.go file, one test function per
// file. See specs/ast-audit-tests.md for the full check catalog.
package audit
