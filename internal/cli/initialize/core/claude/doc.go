//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude handles CLAUDE.md creation and merging
// during the ctx init pipeline.
//
// # Overview
//
// CLAUDE.md is the project-level instruction file Claude
// Code reads. The ctx instructions themselves live in
// AGENTS.md (see the sibling agents package), so the
// ctx-managed section of CLAUDE.md is a bridge: an
// "@AGENTS.md" import line that Claude Code expands at
// launch. The bridge keeps Claude Code releases before
// 2.1.277, sessions that cannot read AGENTS.md directly,
// and hand-written CLAUDE.md files (which otherwise
// suppress AGENTS.md) on the same instructions
// (specs/init-agents-md.md).
//
// # Behavior
//
// [HandleMd] creates a new CLAUDE.md from the embedded
// bridge template, or merges the ctx-managed section into
// an existing file using marker-delimited regions.
//
// # Data Flow
//
// When [HandleMd] is called it:
//
//  1. Reads the embedded CLAUDE.md template from the
//     assets package.
//  2. Delegates to the merge sub-package via
//     merge.OrCreate, passing marker boundaries that
//     delimit the ctx-managed section.
//  3. If no CLAUDE.md exists, the template is written
//     as a new file and a creation message is printed.
//  4. If CLAUDE.md exists, the ctx section between
//     markers is replaced (or inserted if absent).
//  5. Force mode overwrites the section without
//     prompting. Auto-merge mode skips interactive
//     confirmation.
package claude
