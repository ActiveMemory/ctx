//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pi generates Pi CLI integration files during project
// setup.
//
// Pi (earendil-works/pi, pi.dev) is a self-extensible coding-agent
// CLI. By design it has no built-in MCP, so the integration surface
// is a thin TypeScript extension plus Agent-Skills-standard skills;
// all real logic stays in the ctx Go binary via ctx system
// subcommands.
//
// # Deployment Steps
//
// [Deploy] performs three operations in sequence:
//  1. Extension deployment: writes .pi/extensions/ctx.ts as a flat
//     top-level file (Pi auto-loads flat files in .pi/extensions/;
//     project-local extensions load only after project trust)
//  2. AGENTS.md: deploys the shared agent instructions template
//     (Pi loads AGENTS.md natively as a context file)
//  3. Skills: copies the bundled ctx skills to .pi/skills/
//
// Extension deployment failure is fatal; AGENTS.md and skills
// failures are warnings that do not halt deployment, matching the
// OpenCode deploy semantics.
//
// # Data Flow
//
// The setup orchestrator calls [Deploy] from the `ctx setup pi
// --write` branch of the setup root command. Deployment reads the
// embedded assets via the agent package (PiExtension, PiSkills) and
// reports progress through the write/setup Pi info helpers.
package pi
