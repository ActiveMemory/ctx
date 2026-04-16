//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package copilot deploys the **GitHub Copilot integration
// files** that give Copilot Chat in VS Code access to
// project-specific context through the ctx MCP server.
//
// The package is the per-tool deployer called by the
// setup orchestrator when the user opts into Copilot
// integration (`ctx setup copilot`, or as part of
// `ctx init` when Copilot is detected).
//
// # What Gets Deployed
//
// [DeployInstructions] writes two artifacts to the
// project root:
//
//   - **`.github/copilot-instructions.md`** — the
//     prose Copilot Chat reads on session start. Tells
//     it where context lives, which CLI commands it can
//     ask the user to run, and the rule that MCP tool
//     calls beat raw shell commands when both are
//     available. Idempotent: ctx-managed sections are
//     bracketed by markers so user edits outside the
//     markers survive re-deployment.
//   - **`.vscode/mcp.json`** — VS Code Copilot's MCP
//     server registry. Spawns `ctx mcp` over
//     stdin/stdout. Created if missing; merged if
//     present (other MCP servers in the file are
//     preserved).
//
// # Marker Convention
//
// `copilot-instructions.md` uses the same
// `<!-- ctx:copilot --> ... <!-- ctx:copilot:end -->`
// marker pattern as other ctx-managed files so users
// can edit non-ctx prose freely. See
// [internal/cli/initialize/core/merge] for the
// underlying mechanism.
//
// # Concurrency
//
// Filesystem-bound and stateless. Callers serialize
// through process-level execution.
package copilot
