//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package assets provides the embedded filesystem for
// ctx.
//
// All templates, skills, hooks, YAML text files, and
// the Claude Code plugin manifest are compiled into
// the binary via go:embed. The single exported
// variable FS is the entry point for all embedded
// asset reads.
//
// # Embedded Content
//
// The FS variable includes:
//
//   - claude/: CLAUDE.md, plugin.json, and skill
//     definitions with SKILL.md and reference docs
//   - context/: context file templates (TASKS.md,
//     DECISIONS.md, CONVENTIONS.md, etc.)
//   - entry-templates/: Markdown scaffolds for new
//     decisions, learnings, tasks, conventions
//   - hooks/: message templates and trace scripts
//     for lifecycle hooks
//   - integrations/: Copilot instructions, CLI
//     hooks, agent configs, and scripts
//   - project/: README templates for subdirectories
//   - schema/: JSON Schema for .ctxrc validation
//   - why/: philosophy documents (manifesto, about,
//     design invariants)
//   - permissions/: permission text files
//   - commands/: YAML command/flag descriptions
//
// # Accessor Packages
//
// Subdirectories under assets/read/ provide typed
// accessors grouped by domain, so callers do not
// need to know embed paths:
//
//   - read/desc:      command/flag/text descriptions
//   - read/entry:     entry template files
//   - read/hook:      hook message templates
//   - read/skill:     skill SKILL.md files
//   - read/template:  context file templates
package assets
