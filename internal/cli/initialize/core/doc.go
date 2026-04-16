//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core is the umbrella for the ctx init command's
// business logic.
//
// # Overview
//
// The init command bootstraps a project's .context/
// directory, creates context files from templates,
// configures companion tools, and sets up the encrypted
// scratchpad. This package groups the sub-packages that
// implement each stage of the initialisation pipeline.
//
// # Sub-packages
//
//   - backup -- creates timestamped .bak copies before
//     overwriting existing files.
//   - claude -- creates or merges CLAUDE.md with the
//     ctx-managed section.
//   - claude_check -- validates CLAUDE.md structure and
//     renders diagnostic hints.
//   - entry -- creates context file templates (TASKS.md,
//     DECISIONS.md, etc.) and locates insertion points
//     in existing files.
//   - merge -- marker-delimited section merging for
//     CLAUDE.md and prompt files.
//   - pad -- sets up the encrypted or plaintext
//     scratchpad.
//   - plugin -- detects and enables the ctx companion
//     plugin.
//   - project -- scaffolds project-level files like
//     Makefile and .gitignore.
//   - tpl -- generic template deployment engine used
//     by entry and other sub-packages.
//   - validate -- pre-flight checks (ctx in PATH,
//     essential files present).
//   - vscode -- generates VS Code workspace files
//     (extensions.json, tasks.json, mcp.json).
//
// # Data Flow
//
// The cmd layer's Run function orchestrates the init
// pipeline by calling into these sub-packages in order:
// validate, project scaffolding, entry templates, claude
// handling, pad setup, plugin detection, and vscode
// configuration.
package core
