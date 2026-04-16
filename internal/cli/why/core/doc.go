//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core is the umbrella for why subcommand
// business logic. It contains no code of its own;
// all functionality lives in its subpackages.
//
// # Subpackages
//
// The why core layer is split into focused packages
// that each handle one aspect of document display:
//
//   - data -- document alias mappings and display
//     ordering for the interactive menu
//   - menu -- interactive numbered menu that reads
//     user selection from stdin
//   - show -- loads embedded philosophy documents,
//     strips MkDocs syntax, and prints them
//   - strip -- removes MkDocs-specific markup
//     (frontmatter, admonitions, tabs, image refs,
//     relative links) for clean terminal display
//
// # Architecture
//
// Each subpackage exports pure business logic functions
// that the cmd/ layer calls. The data package provides
// static configuration consumed by menu and show. The
// strip package is a pure text transformer with no
// side effects.
package core
