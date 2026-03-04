//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package prompt implements the "ctx prompt" command for managing reusable
// prompt templates.
//
// Prompt templates are plain markdown files stored in .context/prompts/.
// They provide a lightweight alternative to full skills for reusable
// prompt patterns like code review checklists or refactoring guard rails.
//
// Subcommands:
//
//   - list:  list available prompt templates (default)
//   - show:  print a prompt template to stdout
//   - add:   create a new prompt from embedded template or stdin
//   - rm:    remove a prompt template
package prompt
