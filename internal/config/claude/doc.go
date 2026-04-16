//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude centralizes constants for Claude Code
// integration: model identification, API content blocks,
// JSONL parsing, and plugin management.
//
// ctx reads Claude Code session logs (JSONL), detects
// model capabilities, manages plugin installation, and
// parses API content blocks. This package provides all
// the string constants and size limits those operations
// need.
//
// # Model Identification
//
//   - [ModelPrefix] ("claude-") identifies Claude model
//     IDs in session JSONL.
//   - [ModelSuffix1M] ("[1m]") marks models with a 1M
//     context window.
//   - [ModelOpus] detects Opus-family models, which
//     always have 1M context.
//   - [ContextWindow1M] is the numeric window size used
//     for budget calculations.
//
// # API Content Blocks
//
// Claude responses contain typed content blocks. The
// block type constants match the API schema:
//
//   - [BlockText], [BlockThinking], [BlockToolUse],
//     [BlockToolResult]
//
// [ContentBlockArrayPrefix] ("[{") detects the array
// format used by Claude's content block responses.
//
// # JSONL Parsing
//
//   - [FieldUsage] and [FieldInputTokens] are JSON keys
//     for quick-scan token counting.
//   - [MaxTailBytes] (32KB) caps how much of a JSONL file
//     is read from the tail for model detection.
//   - [RoleUser] and [RoleAssistant] identify message
//     roles during conversation parsing.
//
// # Plugin Management
//
// ctx ships as a Claude Code plugin. These constants
// manage plugin identity and permissions:
//
//   - [Binary]: the "claude" executable name
//   - [Md]: the CLAUDE.md project config file
//   - [Settings] / [SettingsGolden]: local settings
//   - [PluginID]: the plugin identifier
//   - [PluginScope] / [PluginScopeWildcard]: permission
//     scope prefixes for skill permissions
//   - [PermSkillPrefix] / [PermSkillSuffix]: tokens
//     that wrap skill names in permission strings
//
// # Why Centralized
//
// The journal importer, the session parser, the setup
// command, and the permission sanitizer all depend on
// these constants. A single package prevents divergence
// between model detection and plugin management code.
package claude
