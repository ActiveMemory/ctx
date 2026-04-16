//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package makefile provides access to the embedded
// Makefile.ctx template.
//
// # Template
//
// Ctx returns the Makefile.ctx content deployed by
// ctx init. This file provides common make targets
// (build, test, lint, audit) for projects using ctx.
// It is designed to be included from the project's
// main Makefile via an include directive.
//
//	content, err := makefile.Ctx()
//
// # Deployment
//
// During ctx init, the Makefile.ctx template is
// written to the project root alongside the .context/
// directory. Projects can include it from their main
// Makefile to inherit standard ctx targets.
package makefile
