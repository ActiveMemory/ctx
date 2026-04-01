//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package project handles project-root directory and file creation
// during initialization.
//
// [CreateDirs] creates the .context/ directory tree with proper
// permissions. [HandleMakefileCtx] deploys the Makefile.ctx
// template if it does not already exist.
package project
