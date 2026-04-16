//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx why" cobra command.
//
// This command displays project philosophy and design
// documents. It can show a specific document by alias
// or present an interactive menu for browsing.
//
// # Usage
//
//	ctx why [document]
//
// # Arguments
//
// An optional positional argument:
//
//   - document: an alias for the document to display.
//     Valid aliases include "manifesto", "about", and
//     "invariants". When omitted, an interactive menu
//     is shown listing all available documents.
//
// # Behavior
//
// The command operates in two modes:
//
//   - Direct mode: when a document alias is provided,
//     delegates to why/core/show.Doc to render the
//     document content to stdout.
//   - Menu mode: when no argument is given, delegates
//     to why/core/menu.Show to present an interactive
//     selection menu. The user picks a document and
//     its content is displayed.
//
// This command does not require an initialized
// context directory (it is annotated with
// SkipInit) because the documents are embedded
// assets, not project-specific files.
//
// # Output
//
// The full text of the selected design document,
// rendered for terminal display. In menu mode, the
// selection prompt appears first.
//
// # Delegation
//
// Document display uses why/core/show. Interactive
// menu uses why/core/menu. Document aliases are
// defined in config/why.
package root
