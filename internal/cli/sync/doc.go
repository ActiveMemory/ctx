//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync implements the "ctx sync" command for reconciling
// context files with codebase changes.
//
// The sync command scans the project for new directories, package
// manager files, and configuration files that are not yet documented
// in context files. It suggests actions to keep context aligned with
// the actual project structure.
package sync
