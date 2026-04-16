//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync implements **`ctx sync`**, the command
// that scans the project for new directories, package-
// manager files, and configuration files that are not yet
// represented in the context files, and surfaces them as
// "consider documenting these" suggestions.
//
// Sync is a *suggester*, not a *mutator*: it never edits
// `ARCHITECTURE.md` or `CONVENTIONS.md` on its own. The
// user (or the AI through a skill) sees the report and
// decides what to add. This boundary is intentional:
// auto-population would silently amplify whatever
// scanner mistakes are made.
//
// # The Scan
//
// The scanner looks for:
//
//   - **New top-level directories**: anything in the
//     project root not already mentioned in
//     ARCHITECTURE.md.
//   - **Package-manager files**: `package.json`,
//     `Cargo.toml`, `go.mod`, `pyproject.toml`,
//     `Gemfile`, `requirements.txt`, etc. (full list
//     from [internal/assets/read/lookup.ConfigPatterns]).
//     Their presence often indicates a language or
//     framework that should be documented in
//     CONVENTIONS.md.
//   - **CI / tooling configs**: `.github/workflows/`,
//     `.gitlab-ci.yml`, `.pre-commit-config.yaml`,
//     `.eslintrc*`, `.prettierrc*`, etc.
//   - **Hidden ctx-relevant directories**: `.devbox`,
//     `.vscode`, `.idea`, etc.
//
// # The Output
//
// Each suggestion includes the file/directory path,
// the type of artifact, and a one-line "consider
// documenting in X" pointer. The user runs `ctx sync`
// periodically (or after a major code review) to keep
// `.context/` aligned with reality.
//
// # Sub-Packages
//
//   - **[core/validate]**: the validation predicates
//     used by the scanner (is this a real package
//     manager file or a vendored copy?).
//
// # Concurrency
//
// Filesystem-bound and stateless. Concurrent invocations
// against the same project would each pay the full
// scan cost.
package sync
