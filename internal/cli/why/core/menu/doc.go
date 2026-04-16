//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package menu implements the interactive document
// selection menu for the why command. When the user runs
// "ctx why" without specifying a document alias, this
// package presents a numbered list and reads their
// choice.
//
// # Menu Flow
//
// [Show] performs the following steps:
//
//  1. Print a banner and separator via the write/why
//     package
//  2. Iterate data.DocOrder and print each document as
//     a numbered menu item
//  3. Print a selection prompt
//  4. Read one line from stdin using a buffered reader
//  5. Parse the input as an integer and validate it
//     against the menu range
//  6. Print a separator and delegate to show.Doc to
//     display the selected document
//
// # Error Handling
//
// Show returns an error when stdin cannot be read or
// when the user enters an invalid selection (non-numeric,
// out of range). The cmd/ layer formats these errors
// for the user.
package menu
