//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package skill provides access to embedded skill
// directories, SKILL.md files, and reference
// documents.
//
// # Listing Skills
//
// List returns the names of all bundled skill
// directories deployed by ctx init. Each skill is a
// directory containing a SKILL.md file following the
// Agent Skills specification.
//
//	names, err := skill.List()
//	// => ["ctx-status", "ctx-reflect", ...]
//
// # Reading Skill Content
//
// Content reads a specific skill's SKILL.md file by
// directory name. The returned bytes contain the full
// skill definition including triggers, instructions,
// and examples.
//
//	data, err := skill.Content("ctx-status")
//
// # Skill Structure
//
// Each skill directory under claude/skills/ contains:
//
//   - SKILL.md: the skill definition file with
//     frontmatter (name, triggers) and body
//     (instructions, examples)
//   - references/: optional reference documents
//     that provide additional context
package skill
