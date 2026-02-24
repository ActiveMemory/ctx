//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude provides Claude Code integration types and utilities.
//
// It provides configuration types for reading/writing Claude Code's
// settings.local.json (permissions) and embedded skill definitions.
//
// Hook logic has been moved to the internal/cli/system package as native
// Go subcommands, deployed via the ctx Claude Code plugin.
//
// Embedded assets:
//   - skills/*/SKILL.md: Agent skill definitions for Claude Code
package claude
