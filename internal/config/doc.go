//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package config is the root of 60+ sub-packages that provide
// constants, compiled regexes, type definitions, and text keys
// used across the ctx codebase.
//
// Each sub-package groups related constants by domain (agent,
// entry, file, mcp, regex, etc.) with zero internal dependencies.
// Consumers import granularly: config/mcp/tool, config/entry,
// config/regex — never this root package directly.
//
// See README.md in this directory for the full organizational
// guide, package categories, and decision tree for placing new
// constants.
package config
