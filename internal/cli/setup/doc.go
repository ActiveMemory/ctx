//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package setup implements the "ctx setup" command for
// generating AI tool integration configurations.
//
// The setup command outputs configuration snippets and
// step-by-step instructions for integrating ctx with
// various AI coding tools. Each tool gets a tailored
// config that wires ctx's hook system, context loading,
// and MCP server into the tool's native extension
// mechanism.
//
// # Supported Tools
//
//	Claude Code: CLAUDE.md generation and hook config
//	Cursor: .cursor/rules and settings.json snippets
//	Aider: .aider.conf.yml and conventions file
//	GitHub Copilot: .github/copilot-instructions.md
//	Windsurf: .windsurfrules configuration
//
// # Subpackages
//
//	cmd/root: cobra command definition and tool
//	  selection logic
//	core: template rendering and config generation
//	  for each supported tool
package setup
