//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package guide provides the "ctx guide" command.
//
// Displays a use-case-oriented cheat sheet covering core
// CLI commands grouped by workflow, available skills, and
// common recipes. The default output fits a single
// terminal screen; full listings are available via flags.
//
// The guide is aimed at new users and AI agents that need
// a quick orientation on what ctx can do. It pulls
// content from embedded asset files and formats it for
// terminal display with section headers and indentation.
//
// # Output Sections
//
//   - Quick start: essential commands for a first session
//   - Context files: what each .context/ file does
//   - Workflows: common multi-command sequences
//   - Skills: available slash-command skills
//   - Recipes: copy-paste command combinations
//
// # Subpackages
//
//   - cmd/root: cobra command definition and flag binding
//   - core: content assembly and terminal formatting
package guide
