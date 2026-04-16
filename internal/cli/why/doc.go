//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package why provides the "ctx why" command.
//
// Surfaces ctx's philosophy documents -- manifesto,
// about page, and design invariants -- directly in the
// terminal, stripped of MkDocs-specific syntax so they
// read cleanly without a browser.
//
// # Purpose
//
// The why command exists for users and AI agents who
// want to understand ctx's design philosophy and
// guiding principles without leaving the CLI. It reads
// embedded markdown files, strips MkDocs admonitions
// and front matter, and outputs clean terminal text.
//
// # Subpackages
//
//   - cmd/root: cobra command definition and document
//     selection logic
//   - core: MkDocs syntax stripping and terminal
//     formatting
package why
