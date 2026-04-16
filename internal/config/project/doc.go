//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package project defines constants for files and
// directories that live at the project root, outside
// the .context/ directory.
//
// # Go Directory Constants
//
// DirInternal ("internal") and DirInternalSlash
// ("internal/") identify the conventional Go internal
// packages directory. Drift checks use
// DirInternalSlash as a prefix to verify that
// backtick-quoted paths in documentation still exist
// on disk.
//
// # Project-Root Files
//
// Several files are managed at the project root
// during "ctx init" and ongoing operation:
//
//   - Makefile: the user's own Makefile.
//   - MakefileCtx ("Makefile.ctx"): ctx-owned
//     Makefile include with convenience targets.
//   - MakefileIncludeDirective: the Make line
//     "-include Makefile.ctx" that pulls in ctx
//     targets (leading dash suppresses errors when
//     the file is absent).
//   - GettingStarted ("GETTING_STARTED.md"): a
//     quick-start reference written during init.
//   - FallbackName ("unknown"): project name used
//     when os.Getwd fails.
//
// # Why Centralize
//
// File names referenced in multiple commands (init,
// status, drift) belong in one place so renames
// propagate automatically and drift checks can
// validate references.
package project
