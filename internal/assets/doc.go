//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package assets provides the embedded filesystem for ctx.
//
// All templates, skills, hooks, YAML text files, and the Claude
// Code plugin manifest are compiled into the binary via go:embed.
// The single exported variable [FS] is the entry point for all
// embedded asset reads. Subdirectories under assets/read/ provide
// typed accessors grouped by domain (desc, entry, hook, skill, etc.).
package assets
