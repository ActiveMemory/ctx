//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package codex holds the pure (cobra-free) helpers behind the
// OpenAI Codex integration: locating the Codex home, detecting
// whether the ctx plugin is installed and enabled, merging the
// embedded hooks manifest into a project `.codex/hooks.json`,
// and appending the `[mcp_servers.ctx]` table to a project
// `.codex/config.toml`.
//
// # Detection
//
// [Home] resolves `$CODEX_HOME` (falling back to `~/.codex`).
// [PluginInstalled] and [PluginEnabled] read the plugin cache
// and `config.toml` under that home; [Detect] combines them
// with a PATH lookup for the `codex` binary into a [State]
// that mirrors the Claude Code detector in
// internal/cli/initialize/core/claudecheck.
//
// # File Merging
//
// [MergeHooks] keeps foreign matcher groups, drops ctx-managed
// ones (every handler command starts with the git-root anchor),
// and appends the embedded groups; [EnsureMCPTable] appends the
// MCP table when its header is absent. Neither function touches
// the filesystem: callers read, call, and write, so the deployer
// in internal/cli/setup/core/codex owns every I/O decision and
// the user-facing output.
//
// The TOML helpers deliberately never parse `config.toml`: the
// user owns that file (comments, ordering), so ctx only scans
// for table headers and appends.
package codex
