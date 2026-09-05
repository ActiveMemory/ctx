//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package codex generates the project-local OpenAI Codex
// integration during `ctx setup codex --write`.
//
// Codex CLI runs the same lifecycle hooks as Claude Code, reads
// AGENTS.md natively, and discovers repo skills under
// `.agents/skills/`. This package materializes the embedded
// Codex plugin assets (internal/assets/codex) into the project
// so teams that do not install the user-level plugin still get
// hooks, the MCP server, and skills.
//
// # Deployment Steps
//
// [Deploy] performs these operations in sequence:
//  1. Plugin short-circuit: when the ctx plugin is enabled in
//     `~/.codex/config.toml`, only AGENTS.md is deployed (Codex
//     loads every matching hook from every source, so a project
//     copy would run each hook twice)
//  2. Hooks: create or merge `.codex/hooks.json`, preserving
//     foreign matcher groups and replacing stale ctx groups
//  3. MCP: create `.codex/config.toml` or append the
//     `[mcp_servers.ctx]` table when its header is absent
//  4. AGENTS.md: shared agent instructions via core/agents
//  5. Skills: `.agents/skills/<name>/SKILL.md` for every embedded
//     Codex skill (create / refresh-if-stale / reject foreign)
//
// Every step ends with the /hooks trust reminder: Codex refuses
// to run non-managed hooks until the user reviews them.
//
// The merge and detection logic lives in the cobra-free
// internal/codex package; this package owns file I/O and output.
package codex
