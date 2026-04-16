//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package copilotcli deploys the **GitHub Copilot CLI
// integration**: hook scripts, agent definitions, skills,
// instructions, and MCP configuration that give the Copilot
// CLI feature parity with the Claude Code integration ctx
// ships natively.
//
// Copilot CLI is a different beast from Copilot Chat: it
// runs in the terminal, dispatches via shell scripts under
// `.github/hooks/`, and consumes a different config file
// format. This package handles all of that.
//
// # Public Surface
//
//   - **[DeployHooks](projectRoot)**: the single
//     public entry point called from the setup
//     command. Orchestrates every artifact below.
//
// # What Gets Deployed
//
//   - **`.github/hooks/ctx-hooks.json`**: the hook
//     manifest declaring which `ctx system` command
//     fires on each lifecycle event (sessionStart,
//     preToolUse, postToolUse, sessionEnd). Skipped
//     if a non-ctx version already exists.
//   - **`.github/hooks/scripts/`**: wrapper shell
//     scripts for any non-stdin hooks Copilot CLI
//     expects.
//   - **`.github/copilot/skills/`**: the same skills
//     ctx ships under
//     `internal/assets/integrations/copilot-cli/skills/`.
//   - **`.github/copilot/agents/`**: agent definitions.
//   - **`.github/copilot/INSTRUCTIONS.md`**: the
//     persistent rules Copilot CLI loads on every
//     prompt.
//   - **MCP config**: for Copilot CLI's MCP client so
//     `ctx mcp` is available.
//
// # Idempotency
//
// Each deployment helper checks for an existing
// destination before writing; a present file is left
// alone (preserving user edits) unless the user passes
// `--force`. Skill files are stamped with a content
// hash so `ctx doctor` can detect drift between
// shipped and installed versions.
//
// # Concurrency
//
// Filesystem-bound and stateless.
package copilotcli
