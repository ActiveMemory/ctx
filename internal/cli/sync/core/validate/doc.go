//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate is the **measurement layer** behind
// `ctx sync`: walks the project tree looking for files
// and directories that are not yet documented in
// `.context/`, and returns one suggestion per find.
//
// The package is the *predicates*; the CLI
// ([internal/cli/sync]) is the *orchestrator*; neither
// mutates context files — sync is suggestion-only by
// design.
//
// # Public Surface
//
//   - **[CheckPackageFiles](root, ctxFiles)** —
//     walks for package-manager descriptors
//     (`package.json`, `Cargo.toml`, `go.mod`,
//     `pyproject.toml`, …) and returns suggestions
//     when any are unmentioned in CONVENTIONS.md /
//     ARCHITECTURE.md.
//   - **[CheckConfigFiles](root, ctxFiles)** —
//     same shape but for CI / tooling configs
//     (`.github/workflows/*`, `.eslintrc*`,
//     `.pre-commit-config.yaml`, …).
//   - **[CheckNewDirectories](root, ctxFiles)** —
//     same shape but for top-level directories
//     that ARCHITECTURE.md does not mention.
//
// # Pattern Source
//
// The "what counts as a config file" patterns come
// from [internal/assets/read/lookup.ConfigPatterns]
// so a single edit there updates every check.
//
// # Concurrency
//
// Filesystem-bound and stateless. Concurrent
// invocations against the same project each pay
// the full scan cost.
//
// # Related Packages
//
//   - [internal/cli/sync]                  — chief
//     consumer.
//   - [internal/assets/read/lookup]        — supplies
//     the pattern set.
//   - [internal/entity]                    — context
//     file types.
package validate
