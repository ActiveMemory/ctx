//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agents handles AGENTS.md creation and merging
// during the ctx init pipeline.
//
// # Overview
//
// AGENTS.md is the vendor-neutral instruction file that
// tells coding agents how to work with a project; Claude
// Code, Codex, Copilot, Cursor, OpenCode and Pi read it.
// During init, this package either creates a new AGENTS.md
// from the embedded template it shares with ctx setup, or
// merges the ctx section into an existing file using
// marker-delimited regions (specs/init-agents-md.md).
//
// # Behavior
//
// [HandleMd] creates a new AGENTS.md from the embedded
// template, or merges the ctx-managed section into an
// existing file using marker-delimited regions.
//
// # Data Flow
//
// When [HandleMd] is called it:
//
//  1. Reads the embedded AGENTS.md template from the
//     assets package (integrations/agents.md).
//  2. Delegates to the merge sub-package via
//     merge.OrCreate, passing the ctx:agents marker
//     boundaries that delimit the ctx-managed section.
//  3. If no AGENTS.md exists, the template is written
//     as a new file and a creation message is printed.
//  4. If AGENTS.md exists without markers, the section
//     is inserted after the H1, after a prompt unless
//     auto-merge is set. With markers present the file
//     is left alone unless force is set, which replaces
//     the section.
package agents
