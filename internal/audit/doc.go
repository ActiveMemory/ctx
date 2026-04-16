//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package audit contains AST-based codebase invariant
// tests that enforce project conventions at the syntax
// tree level.
//
// Unlike [internal/compliance] (which uses file-level
// grep and shell tool checks), audit tests use go/ast
// and go/packages to walk parsed syntax trees. This
// gives type-aware, context-sensitive detection that
// cannot be achieved with regex.
//
// Every file in this package is a _test.go file except
// this doc.go. The package produces no binary output
// and is not importable.
//
// # Shared Helpers
//
// helpers_test.go provides:
//
//   - [loadPackages] loads and caches parsed packages
//     via sync.Once.
//   - [isTestFile] filters _test.go files.
//   - [posString] formats file:line positions for
//     error messages.
//
// # Check Catalog
//
// Each check lives in its own _test.go file, one test
// function per file. Categories include:
//
//   - **Naming** — stuttery function names, descKey
//     namespace alignment, mixed-visibility files.
//   - **Error handling** — naked errors, errors.As
//     usage, unchecked fmt returns, printf calls.
//   - **Code hygiene** — magic strings, magic values,
//     raw file I/O, raw logging, raw time formats,
//     string-concat paths, literal whitespace.
//   - **Structure** — CLI command structure, core
//     structure, cross-package types, dead exports,
//     type file conventions.
//   - **Documentation** — doc comment alignment, doc
//     comments, doc structure, package doc quality.
//   - **Assets** — YAML content drift, YAML examples
//     registry, YAML linkage.
//   - **Permissions** — flagbind usage, import shadow,
//     variable shadowing.
//
// See specs/ast-audit-tests.md for the full check
// catalog and rationale.
package audit
