//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hook manages git hook installation and removal
// for the context tracing system. It installs two hooks:
// prepare-commit-msg (injects context ref trailers into
// commit messages) and post-commit (records refs to
// trace history).
//
// # Hook Lifecycle
//
// [Enable] installs both hooks by reading their script
// templates from embedded assets and writing them to
// the .git/hooks/ directory. Each hook file is marked
// with a ctx-specific marker so the package can
// distinguish its own hooks from user-installed ones.
//
// [Disable] removes both hooks, but only if they contain
// the ctx marker. User-installed hooks are left
// untouched.
//
// # Safe Installation
//
// [Install] writes a hook script to disk after checking
// for an existing file. If a hook file exists and does
// not contain the ctx marker, Install returns an error
// rather than overwriting a user hook. If the existing
// file is a ctx hook, it is replaced with the new
// version.
//
// [Remove] deletes a hook file only when it contains the
// ctx marker. Non-ctx hooks are silently skipped.
//
// # Path Resolution
//
// [FilePath] resolves the absolute path to a named git
// hook by running git rev-parse --git-dir and appending
// the hooks/ subdirectory.
package hook
