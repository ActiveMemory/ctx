//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package template provides access to context
// template files from embedded assets.
//
// # Reading Templates
//
// Template reads a context file template by name from
// the context/ asset directory. These templates are
// stamped into .context/ by ctx init and include
// comment-header guidance for when and how to update
// each file.
//
//	content, err := template.Template("TASKS.md")
//	content, err := template.Template("DECISIONS.md")
//
// # Available Templates
//
// The following context files have templates:
//
//   - TASKS.md: work items and progress tracking
//   - DECISIONS.md: architectural decisions
//   - CONVENTIONS.md: code patterns and standards
//   - LEARNINGS.md: gotchas and lessons learned
//   - CONSTITUTION.md: hard rules and invariants
//   - ARCHITECTURE.md: system design overview
//   - GLOSSARY.md: project-specific terminology
//   - AGENT_PLAYBOOK.md: agent instruction guide
package template
