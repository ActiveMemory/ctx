//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agy implements project-local and plugin deployment logic
// for the Google Antigravity integration.
//
// Deploys .agents/hooks.json (lifecycle hooks),
// ~/.gemini/config/mcp_config.json (MCP server), AGENTS.md, and
// .agents/skills/ (skills and references). When the ctx Antigravity
// plugin is enabled, hooks, MCP, and skills are skipped to avoid
// running them twice.
package agy
