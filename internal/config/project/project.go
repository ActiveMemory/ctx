//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package project

// Project-root file constants (not inside .context/).
const (
	// ImplementationPlan is the high-level project direction file.
	ImplementationPlan = "IMPLEMENTATION_PLAN.md"
	// Makefile is the user's project Makefile.
	Makefile = "Makefile"
	// MakefileCtx is the ctx-owned Makefile include for project root.
	MakefileCtx = "Makefile.ctx"
	// MakefileIncludeDirective is the Make line that pulls in ctx targets.
	// The leading dash suppresses errors when the file is absent.
	MakefileIncludeDirective = "-include Makefile.ctx"
	// GettingStarted is the quick-start reference file written during init.
	GettingStarted = "GETTING_STARTED.md"
)
