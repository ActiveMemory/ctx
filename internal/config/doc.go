//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package config is the root of the domain-scoped constant
// tree that supplies every magic string, threshold, regex,
// and type definition used across the ctx codebase.
//
// This package itself exports nothing. All constants live
// in sub-packages grouped by domain. Consumers import
// granularly (config/agent, config/ctx, config/hook),
// never this root package directly.
//
// # Organization
//
// Sub-packages fall into several categories:
//
// Core identity and structure:
//
//   - config/ctx: context file names and read order
//   - config/dir: directory path constants
//   - config/file: file extensions and naming
//   - config/entry: entry type constants
//
// AI agent pipeline:
//
//   - config/agent: budget, cooldown, scoring
//   - config/token: token counting constants
//   - config/session: session metadata
//   - config/content: emptiness heuristics
//
// CLI and formatting:
//
//   - config/cli: cobra annotations, XML attributes
//   - config/box: box-drawing for nudges
//   - config/fmt: formatting helpers
//   - config/flag: flag name constants
//   - config/embed: embedded command/flag metadata
//
// Integrations:
//
//   - config/claude: Claude Code model and plugin
//   - config/copilot: Copilot Chat and CLI parsing
//   - config/vscode: VS Code path constants
//
// Lifecycle and hooks:
//
//   - config/hook: hook execution constants
//   - config/ceremony: session ritual detection
//   - config/nudge: nudge content and throttling
//   - config/architecture: map staleness detection
//   - config/drift: context drift detection
//
// Storage and security:
//
//   - config/crypto: AES-256-GCM parameters
//   - config/archive: backup and snapshot naming
//   - config/asset: embedded asset paths
//
// # Design Principles
//
// Every sub-package has zero internal dependencies.
// No config sub-package imports another config sub-
// package. This keeps the dependency graph flat and
// prevents import cycles.
//
// Constants are grouped by the domain they serve, not
// by their Go type. A package like config/archive mixes
// strings, ints, and format templates because they all
// belong to the archival domain.
package config
