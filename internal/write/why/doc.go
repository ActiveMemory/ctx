//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package why provides terminal output for the
// philosophy command (ctx why).
//
// The why command presents embedded project philosophy
// documents through an interactive numbered menu.
// Output functions handle each stage of the menu
// interaction and document display.
//
// # Menu Rendering
//
// [Banner] renders the ctx ASCII art header at the
// top of the menu. [MenuItem] prints a numbered
// choice with its display label. [MenuPrompt]
// prints the selection prompt and waits for input.
//
// # Document Display
//
// [Content] prints the chosen philosophy document
// body to stdout. The content is pre-processed
// before being passed to this function, so it
// contains no formatting logic.
//
// # Visual Structure
//
// [Separator] prints a blank line for visual
// separation between menu sections and between
// the menu and the document body.
package why
