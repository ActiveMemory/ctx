//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package codex centralizes constants for the OpenAI Codex CLI
// integration: home-directory layout, project-local config,
// plugin/marketplace identity, hook-manifest tokens, and the
// rollout transcript vocabulary the journal parser reads.
//
// # Layout
//
// Codex keeps user state under $CODEX_HOME (default ~/.codex):
//
//   - [DirSessions]: rollout transcripts, one JSONL per thread
//   - [DirPlugins]/[DirPluginCache]: installed plugin copies
//   - [FileConfigTOML]: user config, including enabled plugins
//
// Project-local integration lives under [Dir] (`.codex/`) —
// [FileHooksJSON] and [FileConfigTOML] — plus [DirAgents]
// (`.agents/`) for repo-scoped skills and the marketplace catalog.
//
// # Plugin
//
// ctx ships as a Codex plugin rooted at internal/assets/codex
// (manifest in `.codex-plugin/`, hooks under [DirHooks],
// bundled MCP map in [FileMCPJSON]). The repo marketplace at
// `.agents/plugins/marketplace.json` is named [MarketplaceID];
// the installed identifier is [PluginID].
//
// # Hooks
//
// Codex's lifecycle-hook contract mirrors Claude Code's (same
// event names, same stdin fields, same hookSpecificOutput JSON).
// [HookAnchor] is prepended to every command because Codex runs
// hooks with the session cwd while ctx is CWD-anchored.
//
// # TOML
//
// ctx never parses config.toml; it appends tables and scans for
// header lines ([TOMLHeaderMCPCtx], [TOMLHeaderPluginCtx]) so the
// user's comments and ordering survive untouched.
//
// # Rollouts
//
// The Line*/Item*/Content*/Role* constants name the JSONL
// vocabulary of `rollout-*.jsonl` files; [InjectedUserPrefixes]
// identifies the user-role items Codex injects itself so the
// parser can drop them.
package codex
