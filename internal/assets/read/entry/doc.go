//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package entry provides access to entry template
// files embedded in the assets filesystem.
//
// Entry templates are Markdown scaffolds used when
// adding new decisions, learnings, tasks, or
// conventions via ctx add. Each template defines the
// structure and required fields for its entry type.
//
// # Listing Templates
//
// List returns the file names of all available entry
// templates from the entry-templates/ asset directory.
//
//	names, err := entry.List()
//	// => ["decision.md", "learning.md", ...]
//
// # Reading Templates
//
// ForName reads a specific template by filename. The
// returned bytes contain the Markdown scaffold ready
// for field substitution.
//
//	content, err := entry.ForName("decision.md")
//
// # Usage
//
// The add command reads the appropriate template,
// substitutes user-provided values into the scaffold,
// and appends the formatted entry to the matching
// context file.
package entry
