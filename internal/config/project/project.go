//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package project

// Go project directory constants.
const (
	// DirInternal is the conventional Go internal packages directory.
	DirInternal = "internal"
	// DirInternalSlash is DirInternal with a trailing slash for path
	// prefix matching in Markdown content.
	DirInternalSlash = DirInternal + "/"
)

// Project-root file constants (not inside .context/).
const (
	// Makefile is the user's project Makefile.
	Makefile = "Makefile"
	// MakefileCtx is the ctx-owned Makefile include for project root.
	MakefileCtx = "Makefile.ctx"
	// MakefileIncludeDirective is the Make line that pulls in ctx targets.
	// The leading dash suppresses errors when the file is absent.
	MakefileIncludeDirective = "-include Makefile.ctx"
	// GettingStarted is the quick-start reference file written during init.
	GettingStarted = "GETTING_STARTED.md"
	// FallbackName is the project name used when os.Getwd fails.
	FallbackName = "unknown"
)
