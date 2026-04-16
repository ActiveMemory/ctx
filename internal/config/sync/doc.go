//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync defines constants for the ctx sync
// command, which detects project changes that need
// documentation updates.
//
// ctx sync scans a project for configuration files,
// dependency manifests, and important directories,
// then reports what has changed since the last sync.
// This package provides the glob patterns, action
// types, directory lists, and package-ecosystem
// mappings the scanner needs.
//
// # Glob Patterns
//
// File patterns like [PatternEslint],
// [PatternPrettier], [PatternTSConfig],
// [PatternEditorConf], [PatternMakefile], and
// [PatternDockerfile] identify config files whose
// changes may require doc updates.
//
// # Action Types
//
//   - [ActionDeps]   — a dependency manifest changed.
//   - [ActionConfig] — a config file changed.
//   - [ActionNewDir] — a new important directory
//     appeared.
//
// # Directory Control
//
//   - [ImportantDirs] — top-level directories that
//     should be documented in ARCHITECTURE.md (api,
//     cmd, internal, src, ...).
//   - [SkipDirs] — directories excluded from scanning
//     (build, dist, node_modules, vendor).
//
// # Package Ecosystems
//
// The [Packages] map links manifest filenames
// (package.json, go.mod, Cargo.toml, ...) to their
// ecosystem descriptions for dependency doc hints.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package sync
