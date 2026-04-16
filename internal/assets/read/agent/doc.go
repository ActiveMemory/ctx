//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agent provides access to agent integration
// files embedded in the assets filesystem.
//
// # Copilot Integration
//
// CopilotInstructions returns the embedded
// copilot-instructions.md template deployed by
// ctx init for GitHub Copilot integration.
//
// CopilotCLIHooksJSON returns the hooks config JSON
// for the Copilot CLI integration.
//
// CopilotCLIScripts returns a map of filename to
// content for all embedded hook scripts in the
// Copilot CLI scripts directory.
//
// CopilotCLISkills returns a map of skill name to
// SKILL.md content for embedded Copilot CLI skills.
//
// # GitHub Agents
//
// AgentsMd returns the AGENTS.md template for
// repository-level agent configuration.
//
// AgentsCtxMd returns the .github/agents/ctx.md
// template for the ctx agent definition.
//
// InstructionsCtxMd returns the path-specific
// instructions template for the context directory.
package agent
