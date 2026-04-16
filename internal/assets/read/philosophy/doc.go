//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package philosophy provides access to embedded
// why-documents (manifesto, about, design invariants).
//
// # Reading Documents
//
// WhyDoc reads a philosophy document by name from the
// why/ directory in the embedded filesystem. The .md
// extension is appended automatically.
//
//	doc, err := philosophy.WhyDoc("manifesto")
//	doc, err := philosophy.WhyDoc("about")
//	doc, err := philosophy.WhyDoc("design-invariants")
//
// # Purpose
//
// These documents are the source files for the
// "ctx why" interactive reader. They explain the
// rationale, design philosophy, and invariants behind
// ctx. The documents are synced from docs/ into the
// binary via make sync-why during the build process.
package philosophy
