//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude provides Claude Code integration types
// and utilities for ctx.
//
// # Configuration Types
//
// The package defines Go structs that mirror Claude
// Code's settings.local.json format:
//
//   - [Settings] is the top-level structure containing
//     hooks and permissions.
//   - [HookConfig] maps lifecycle events (PreToolUse,
//     PostToolUse, UserPromptSubmit, SessionEnd) to
//     lists of [HookMatcher] entries.
//   - [PermissionsConfig] holds allow/deny lists for
//     tool patterns (e.g., "Bash(ctx status:*)").
//
// These types are used by ctx init to generate and
// update the project-level Claude Code configuration
// file.
//
// # Embedded Skills
//
// The package exposes embedded Agent Skill definitions
// via two functions:
//
//   - [SkillList] returns the names of all embedded
//     skill directories (e.g., "ctx-status",
//     "ctx-reflect").
//   - [SkillContent] returns the raw SKILL.md bytes
//     for a given skill name.
//
// Skills follow the Agent Skills specification with
// SKILL.md files containing YAML frontmatter (name,
// description) and autonomy-focused instructions.
// They are installed to .claude/skills/ via ctx init.
//
// # Hook Migration
//
// Hook logic has been moved to internal/cli/system as
// native Go subcommands deployed via the ctx Claude
// Code plugin. The types here remain for reading and
// writing the settings.local.json file.
package claude
