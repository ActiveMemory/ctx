//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides shared helpers for the load
// command.
//
// The "ctx load" command reads context files from the
// project directory and presents them to the user or
// to an AI agent. This core package holds the business
// logic that the cmd/load layer delegates to.
//
// # File Name Conversion
//
// The convert sub-package exports [convert.FileNameToTitle],
// which transforms SCREAMING_SNAKE_CASE markdown
// filenames into Title Case strings for display. It
// strips the .md extension, replaces underscores with
// spaces, and capitalizes each word. For example,
// "AGENT_PLAYBOOK.md" becomes "Agent Playbook".
//
// # Read-Order Sorting
//
// The sort sub-package exports [sort.ByReadOrder],
// which arranges context files according to a
// predefined priority list (ctx.ReadOrder). Files not
// in the list receive a fallback priority and appear
// at the end. The function returns a new sorted slice
// without modifying the original.
//
// # Data Flow
//
// The cmd/load layer discovers context files on disk,
// calls ByReadOrder to arrange them, and uses
// FileNameToTitle to generate section headings. The
// sorted, titled output is then rendered by the
// write/load package.
package core
